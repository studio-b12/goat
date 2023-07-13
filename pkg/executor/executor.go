package executor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/studio-b12/goat/pkg/advancer"
	"github.com/studio-b12/goat/pkg/clr"
	"github.com/studio-b12/goat/pkg/engine"
	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/goatfile"
	"github.com/studio-b12/goat/pkg/requester"
	"github.com/studio-b12/goat/pkg/util"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

// Executor parses a Goatfiles and executes them.
type Executor struct {
	engineMaker func() engine.Engine
	req         requester.Requester

	Dry     bool
	NoAbort bool
	Skip    []string
	Waiter  advancer.Waiter
}

// New initializes a new instance of Executor using
// the given engineMaker to initialize a new instance
// of Engine for each batch execution. Also, a Requester
// implementation is passed which is used to perform the
// requests.
func New(engineMaker func() engine.Engine, req requester.Requester) *Executor {
	var t Executor

	t.engineMaker = engineMaker
	t.req = req
	t.Waiter = advancer.None{}

	return &t
}

// Execute executes a single or multiple Goatfiles
// from the given file or directory. The given
// initialParams are used as initial state for the
// runtime engine.
func (t *Executor) Execute(pathes []string, initialParams engine.State) (res Result, err error) {
	if len(pathes) == 1 {
		stat, err := os.Stat(pathes[0])
		if err != nil {
			return Result{}, errs.WithPrefix("stat failed:", err)
		}
		if !stat.IsDir() {
			gf, err := t.parseGoatfile(pathes[0])
			if err != nil {
				return Result{}, err
			}

			log.Debug().Msg("Executing goatfile ...")
			return t.ExecuteGoatfile(gf, initialParams)
		}
	}

	return t.executeFromPathes(pathes, initialParams)
}

// ExecuteGoatfile runs the given parsed Goatfile. The given initialParams are
// used as initial state for the runtime engine.
func (t *Executor) ExecuteGoatfile(gf goatfile.Goatfile, initialParams engine.State) (res Result, err error) {
	log := log.Tagged(gf.Path)
	log.Debug().Msg("Parsed Goatfile\n" + gf.String())

	if t.Dry {
		log.Warn().Msg("This is a dry run: no requets will be executed")
		return Result{}, nil
	}

	eng := t.engineMaker()
	eng.SetState(initialParams)

	return t.executeGoatfile(log, gf, eng)
}

func (t *Executor) executeGoatfile(log rogu.Logger, gf goatfile.Goatfile, eng engine.Engine) (res Result, err error) {
	var errsNoAbort errs.Errors

	defer func() {
		// Teardown Procedures

		if t.isSkip("teardown") {
			log.Warn().Msg("skipping teardown steps")
			return
		}

		if len(gf.Teardown) > 0 {
			printSeparator("TEARDOWN")
		}
		for _, act := range gf.Teardown {
			exErr := t.executeAction(log, eng, act, gf)
			res.Teardown.Inc()
			if exErr != nil {
				err = errs.Join(err, exErr)
				if act.Type() == goatfile.ActionRequest {
					log.Error().Err(exErr).Field("req", act).Msg("Teardown step failed")
					res.Teardown.IncFailed()

					// If the returned error comes from the params parsing step, don't
					// cancel the teardown execution. See the following issue for more information.
					// https://github.com/studio-b12/goat/issues/9
					if errs.IsOfType[ParamsParsingError](exErr) {
						continue
					}

					if !t.isAbortOnError(act.(goatfile.Request)) {
						continue
					}
				} else {
					log.Error().Err(exErr).Field("act", act).Msg("Action failed")
				}

				break
			}

			if act.Type() == goatfile.ActionRequest {
				log.Info().Field("req", act).Msg("Teardown step completed")
			}
		}
	}()

	// Setup Procedures

	if t.isSkip("setup") {
		log.Warn().Msg("skipping setup steps")
	} else {
		if len(gf.Setup) > 0 {
			printSeparator("SETUP")
		}
		for _, act := range gf.Setup {
			err := t.executeAction(log, eng, act, gf)
			res.Setup.Inc()
			if err != nil {
				res.Setup.IncFailed()
				if act.Type() == goatfile.ActionRequest {
					log.Error().Err(err).Field("req", act).Msg("Setup step failed")
					if !t.isAbortOnError(act.(goatfile.Request)) {
						errsNoAbort = errsNoAbort.Append(err)
						continue
					}
				}
				return res, err
			}

			if act.Type() == goatfile.ActionRequest {
				log.Info().Field("req", act).Msg("Setup step completed")
			}
		}
	}

	// Test Procedures

	if t.isSkip("tests") {
		log.Warn().Msg("skipping test steps")
	} else {
		if len(gf.Tests) > 0 {
			printSeparator("TESTS")
		}
		for _, act := range gf.Tests {
			err := t.executeTest(act, eng, gf)
			res.Tests.Inc()
			if err != nil {
				res.Tests.IncFailed()
				if act.Type() == goatfile.ActionRequest && !t.isAbortOnError(act.(goatfile.Request)) {
					errsNoAbort = errsNoAbort.Append(err)
					continue
				}
				return res, err
			}
		}
	}

	err = errsNoAbort.Condense()
	return res, err
}

func (t *Executor) executeFromPathes(pathes []string, initialParams engine.State) (finalRes Result, err error) {
	var goatfiles []goatfile.Goatfile

	for _, path := range pathes {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, _ error) error {
			if d.IsDir() && strings.HasPrefix(d.Name(), "_") {
				return fs.SkipDir
			}
			if d.IsDir() ||
				filepath.Ext(d.Name()) != "."+goatfile.FileExtension ||
				strings.HasPrefix(d.Name(), "_") {
				return nil
			}

			gf, err := t.parseGoatfile(path)
			if err != nil {
				return err
			}

			goatfiles = append(goatfiles, gf)
			return nil
		})
		if err != nil {
			return Result{}, err
		}
	}

	if len(goatfiles) == 0 {
		return Result{}, errors.New("no Goatfiles found to execute")
	}

	var mErr errs.Errors

	for _, gf := range goatfiles {
		log.Info().Field("path", gf.Path).Msg(clr.Print(clr.Format("Executing batch ...", clr.ColorFGPurple, clr.FormatBold)))

		res, err := t.ExecuteGoatfile(gf, initialParams)
		finalRes.Merge(res)
		if err != nil {
			entry := log.Error()
			if mErr, ok := err.(errs.Errors); ok {
				errLines := make([]string, 0, len(mErr))
				for _, e := range mErr {
					errLines = append(errLines, clr.Print(clr.Format(e.Error(), clr.ColorFGRed)))
				}
				entry.Field("errors", errLines)
			} else {
				entry.Err(err)
			}
			entry.Msg(clr.Print(clr.Format("Batch execution failed", clr.ColorFGRed, clr.FormatBold)))

			mErr = mErr.Append(wrapBatchExecutionError(err, gf.Path))
			continue
		}

		log.Info().Field("path", gf.Path).Msg(clr.Print(clr.Format("Batch finished successfully", clr.ColorFGPurple, clr.FormatBold)))
	}

	if mErr.HasSome() {
		err = BatchResultError{
			Inner: mErr,
			Total: len(goatfiles),
		}
		return finalRes, err
	}

	return finalRes, nil
}

func (t *Executor) parseGoatfile(path string) (gf goatfile.Goatfile, err error) {
	log.Debug().Field("from", path).Msg("Parsing goatfile ...")

	data, err := os.ReadFile(path)
	if err != nil {
		return goatfile.Goatfile{}, errs.WithPrefix("failed reading file:", err)
	}

	gf, err = goatfile.Unmarshal(string(data), path)
	if err != nil {
		if errs.IsOfType[goatfile.ParseError](err) {
			// TODO: Better wrap this error for visualization and
			//       unwrap-ability.
			return goatfile.Goatfile{}, fmt.Errorf("failed parsing goatfile at %s:%s", path, err.Error())
		}
		return goatfile.Goatfile{}, errs.WithPrefix(fmt.Sprintf("failed parsing goatfile %s:", path), err)
	}

	return gf, nil
}

func (t *Executor) executeTest(
	act goatfile.Action,
	eng engine.Engine,
	gf goatfile.Goatfile,
) (err error) {
	var errsNoAbort errs.Errors
	log := log.Tagged(gf.Path)

	defer func() {
		// Teardown-Each steps

		if t.isSkip("teardown-each") {
			log.Warn().Msg("skipping teardown-each steps")
			return
		}

		for _, postReq := range gf.TeardownEach {
			err := t.executeAction(log, eng, postReq, gf)
			if err != nil {
				if act.Type() == goatfile.ActionRequest {
					log.Error().Err(err).Field("req", act).Msg("Post-Each step failed")

					err = errs.WithPrefix("post-setup-each step failed:", err)

					// If the returned error comes from the params parsing step, don't
					// cancel the teardown-each execution. See the following issue for more information.
					// https://github.com/studio-b12/goat/issues/9
					if errs.IsOfType[ParamsParsingError](err) {
						continue
					}

					if t.isAbortOnError(postReq.(goatfile.Request)) {
						break
					}

					errsNoAbort = errsNoAbort.Append(err)
					continue
				} else {
					log.Error().Err(err).Field("act", act).Msg("Action failed")
					break
				}
			}

			if act.Type() == goatfile.ActionRequest {
				log.Info().Field("req", act).Msg("Teardown-Each step completed")
			}
		}
	}()

	// Setup-Each Steps

	if t.isSkip("setup-each") {
		log.Warn().Msg("skipping setup-each steps")
	} else {
		for _, preAct := range gf.SetupEach {
			err := t.executeAction(log, eng, preAct, gf)
			if err != nil {
				if preAct.Type() == goatfile.ActionRequest {
					log.Error().Err(err).Field("req", act).Msg("Setup-Each step failed")

					err = errs.WithPrefix("Setup-Each step failed:", err)

					if !t.isAbortOnError(preAct.(goatfile.Request)) {
						errsNoAbort = errsNoAbort.Append(err)
						continue
					}
				}

				return err
			}

			if preAct.Type() == goatfile.ActionRequest {
				log.Info().Field("req", act).Msg("Setup-Each step completed")
			}
		}
	}

	// Actual Test Step

	err = t.executeAction(log, eng, act, gf)
	if err != nil {
		if act.Type() == goatfile.ActionRequest {
			log.Error().Err(err).Field("req", act).Msg("Test step failed")

			if !t.isAbortOnError(act.(goatfile.Request)) {
				return err
			}

			errsNoAbort = errsNoAbort.Append(err)
		} else {
			return err
		}
	} else {
		if act.Type() == goatfile.ActionRequest {
			log.Info().Field("req", act).Msg("Test completed")
		}
	}

	return errsNoAbort.Condense()
}

func (t *Executor) executeAction(
	log rogu.Logger,
	eng engine.Engine,
	act goatfile.Action,
	gf goatfile.Goatfile,
) (err error) {
	switch act.Type() {
	case goatfile.ActionRequest:
		req := act.(goatfile.Request)
		err = t.executeRequest(eng, req, gf)
		if err != nil {
			err = errs.WithSuffix(err, fmt.Sprintf("(%s:%d)", req.Path, req.PosLine))
		}
		return err
	case goatfile.ActionLogSection:
		logSection := act.(goatfile.LogSection)
		printSeparator(string(logSection))
		return nil
	case goatfile.ActionExecute:
		execParams := act.(goatfile.Execute)
		_, err := t.executeExecute(path.Dir(gf.Path), execParams, eng)
		if err != nil {
			err = errs.WithSuffix(err, "(imported)")
		}
		return err
	default:
		panic(fmt.Sprintf("An invalid action has been executed: %v\n"+
			"This should actually never happen. If it does though,"+
			"please report this issue to https://github.com/studio-b12/goat.",
			act.Type()))
	}
}

func (t *Executor) executeRequest(eng engine.Engine, req goatfile.Request, gf goatfile.Goatfile) (err error) {
	req.Merge(gf.Defaults)

	preScript, err := util.ReadReaderToString(req.PreScript.Reader())
	if err != nil {
		return errs.WithPrefix("reading preScript failed:", err)
	}

	if preScript != "" {
		err = eng.Run(preScript)
		if err != nil {
			return errs.WithPrefix("preScript failed:", err)
		}
	}

	state := eng.State()
	err = req.ParseWithParams(state)
	if err != nil {
		return errs.WithPrefix("failed infusing request with parameters:",
			ParamsParsingError(err))
	}

	execOpts := ExecOptionsFromMap(req.Options)
	if !execOpts.Condition {
		log.Warn().Field("req", req).Msg("Skipped due to condition")
		return nil
	}

	if execOpts.Delay > 0 {
		log.Info().
			Field("req", req).
			Field("delay", execOpts.Delay).
			Msg(clr.Print(clr.Format("Awaiting delay ...", clr.ColorFGBlack)))
		time.Sleep(execOpts.Delay)
	}

	t.Waiter.Wait()

	httpReq, err := req.ToHttpRequest()
	if err != nil {
		return errs.WithPrefix("failed transforming to http request:", err)
	}

	reqOpts := requester.OptionsFromMap(req.Options)
	httpResp, err := t.req.Do(httpReq, reqOpts)
	if err != nil {
		return errs.WithPrefix("http request failed:", err)
	}

	resp, err := FromHttpResponse(httpResp)
	if err != nil {
		return errs.WithPrefix("response interpretation failed:", err)
	}

	state.Merge(engine.State{"response": resp})
	eng.SetState(state)

	script, err := util.ReadReaderToString(req.Script.Reader())
	if err != nil {
		return errs.WithPrefix("reading script failed:", err)
	}

	if script != "" {
		err = eng.Run(script)
		if err != nil {
			return errs.WithPrefix("script failed:", err)
		}
	}

	return nil
}

func (t *Executor) executeExecute(rootPath string, params goatfile.Execute, eng engine.Engine) (Result, error) {
	pth := goatfile.Extend(path.Join(rootPath, params.File), goatfile.FileExtension)
	gf, err := t.parseGoatfile(pth)
	if err != nil {
		return Result{}, err
	}

	state := eng.State()

	err = goatfile.ApplyTemplateToMap(params.Params, state)
	if err != nil {
		return Result{}, err
	}

	log := log.Tagged(gf.Path)

	isolatedEng := t.engineMaker()
	isolatedEng.SetState(params.Params)

	res, err := t.executeGoatfile(log, gf, isolatedEng)
	if err != nil {
		return res, err
	}

	capturedState := isolatedEng.State()
	for k, v := range capturedState {
		storeAs, ok := params.Returns[k]
		if ok {
			state[storeAs] = v
		}
	}

	eng.SetState(state)

	return res, nil
}

func (t *Executor) isSkip(section string) bool {
	for _, s := range t.Skip {
		if strings.ToLower(s) == section {
			return true
		}
	}
	return false
}

func (t *Executor) isAbortOnError(req goatfile.Request) bool {
	opts := AbortOptionsFromMap(req.Options)

	if opts.AlwaysAbort {
		return true
	}

	if opts.NoAbort || t.NoAbort {
		return false
	}

	return true
}

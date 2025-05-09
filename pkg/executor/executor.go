package executor

import (
	"context"
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

// Executor parses Goatfiles and executes them.
type Executor struct {
	engineMaker func() engine.Engine
	req         requester.Requester

	ctx context.Context

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
//
// The passed context ctx can cancel a stack execution when
// the context is done. All subsequent teardown steps will
// be executed afterward and are not affected by the context
// state.
func New(ctx context.Context, engineMaker func() engine.Engine, req requester.Requester) *Executor {
	var t Executor

	t.ctx = ctx
	t.engineMaker = engineMaker
	t.req = req
	t.Waiter = advancer.None{}

	return &t
}

// Execute executes a single or multiple Goatfiles
// from the given file or directory. The given
// initialParams are used as initial state for the
// runtime engine.
func (t *Executor) Execute(pathes []string, initialParams engine.State, showTeardownParamErrors bool) (res Result, err error) {
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
			return t.ExecuteGoatfile(gf, initialParams, showTeardownParamErrors)
		}
	}

	return t.executeFromPathes(pathes, initialParams, showTeardownParamErrors)
}

// ExecuteGoatfile runs the given parsed Goatfile. The given initialParams are
// used as initial state for the runtime engine.
func (t *Executor) ExecuteGoatfile(gf goatfile.Goatfile, initialParams engine.State, showTeardownParamErrors bool) (res Result, err error) {
	log := log.Tagged(strings.TrimSuffix(gf.Path, ".goat"))

	if t.Dry {
		log.Warn().Msg("This is a dry run: no requets will be executed")
		log.Debug().Msg("Parsed Goatfile\n" + gf.String())
		return Result{}, nil
	}

	eng := t.engineMaker()
	eng.SetState(initialParams)

	return t.executeGoatfile(log, gf, eng, true, showTeardownParamErrors)
}

func (t *Executor) executeGoatfile(
	log rogu.Logger,
	gf goatfile.Goatfile,
	eng engine.Engine,
	printSeperators bool,
	showTeardownParamErrors bool,
) (res Result, err error) {
	var errsNoAbort errs.Errors

	defer func() {
		// Teardown Procedures

		if t.isSkip(goatfile.SectionTeardown) {
			log.Warn().Msg("skipping teardown steps")
			return
		}

		if len(gf.Teardown) > 0 && printSeperators {
			printSeparator("TEARDOWN")
		}
		for _, act := range gf.Teardown {
			sectRes, exErr := t.executeAction(log, eng, act, gf, showTeardownParamErrors)
			res.Teardown.Merge(sectRes)
			if exErr != nil {
				err = errs.Join(err, NewTeardownError(exErr))
				if act.Type() == goatfile.ActionRequest {
					isParamsParseErr := errs.IsOfType[ParamsParsingError](exErr)

					if !isParamsParseErr || showTeardownParamErrors {
						log.Error().Err(exErr).Field("req", act).Msg("Teardown step failed")
					}

					// If the returned error comes from the params parsing step, don't
					// cancel the teardown execution. See the following issue for more information.
					// https://github.com/studio-b12/goat/issues/9
					if isParamsParseErr {
						continue
					}

					if !AbortOptionsFromMap(act.(*goatfile.Request).Options).AlwaysAbort {
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

	if t.isSkip(goatfile.SectionSetup) {
		log.Warn().Msg("skipping setup steps")
	} else {
		if len(gf.Setup) > 0 && printSeperators {
			printSeparator("SETUP")
		}
		for _, act := range gf.Setup {
			select {
			case <-t.ctx.Done():
				return res, ErrCanceled
			default:
				sectRes, err := t.executeAction(log, eng, act, gf, showTeardownParamErrors)
				res.Setup.Merge(sectRes)
				if err != nil {
					if act.Type() == goatfile.ActionRequest {
						log.Error().Err(err).Field("req", act).Msg("Setup step failed")
						if errs.IsOfType[NoAbortError](err) {
							errsNoAbort = errsNoAbort.Append(errors.Unwrap(err))
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
	}

	// Test Procedures

	if t.isSkip(goatfile.SectionTests) {
		log.Warn().Msg("skipping test steps")
	} else {
		if len(gf.Tests) > 0 && printSeperators {
			printSeparator("TESTS")
		}
		for _, act := range gf.Tests {
			select {
			case <-t.ctx.Done():
				return res, ErrCanceled
			default:
				sectRes, err := t.executeTest(act, eng, gf, showTeardownParamErrors)
				res.Tests.Merge(sectRes)
				if err != nil {
					if act.Type() == goatfile.ActionRequest && errs.IsOfType[NoAbortError](err) {
						errsNoAbort = errsNoAbort.Append(errors.Unwrap(err))
						continue
					}
					return res, err
				}
			}
		}
	}

	err = errsNoAbort.Condense()
	return res, err
}

func (t *Executor) executeFromPathes(pathes []string, initialParams engine.State, showTeardownParamErrors bool) (finalRes Result, err error) {
	var goatfiles []goatfile.Goatfile

	for _, path := range pathes {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}

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

		res, err := t.ExecuteGoatfile(gf, initialParams, showTeardownParamErrors)
		finalRes.Merge(res)
		if err != nil {
			entry := log.Error()
			if mErr, ok := err.(errs.Errors); ok {
				errLines := make([]string, 0, len(mErr))
				for _, e := range mErr {
					if !showTeardownParamErrors && errs.IsOfType[TeardownError](e) && errs.IsOfType[ParamsParsingError](err) {
						continue
					}
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
		err = &BatchResultError{
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
	showTeardownParamErrors bool,
) (res ResultSection, err error) {
	var errsNoAbort errs.Errors
	log := log.Tagged(strings.TrimSuffix(gf.Path, ".goat"))

	res, err = t.executeAction(log, eng, act, gf, showTeardownParamErrors)
	if err != nil {
		if act.Type() == goatfile.ActionRequest {
			log.Error().Err(err).Field("req", act).Msg("Test step failed")

			if !errs.IsOfType[NoAbortError](err) {
				return res, err
			}

			errsNoAbort = errsNoAbort.Append(errors.Unwrap(err))
		} else {
			return res, err
		}
	} else {
		if act.Type() == goatfile.ActionRequest {
			log.Info().Field("req", act).Msg("Test completed")
		}
	}

	return res, errsNoAbort.Condense()
}

func (t *Executor) executeAction(
	log rogu.Logger,
	eng engine.Engine,
	act goatfile.Action,
	gf goatfile.Goatfile,
	showTeardownParamErrors bool,
) (res ResultSection, err error) {
	log.Trace().Field("act", act).Msg("Executing action")

	switch act.Type() {

	case goatfile.ActionRequest:
		res.Inc()
		req := act.(*goatfile.Request)
		log.Trace().Fields("options", req.Options).Msg("Request Options")
		err = t.executeRequest(eng, req, gf)
		if err != nil {
			res.IncFailed()
			err = errs.WithSuffix(err, fmt.Sprintf("(%s:%d)", req.Path, req.PosLine))
		}
		return res, err

	case goatfile.ActionLogSection:
		logSection := act.(goatfile.LogSection)
		printSeparator(string(logSection))
		return res, nil

	case goatfile.ActionExecute:
		execParams := act.(goatfile.Execute)
		r, err := t.executeExecute(execParams, eng, showTeardownParamErrors)
		if err != nil {
			err = errs.WithSuffix(err, "(imported)")
		}
		return r.Sum(), err

	default:
		panic(fmt.Sprintf("An invalid action has been executed: %v\n"+
			"This should actually never happen. If it does though,"+
			"please report this issue to https://github.com/studio-b12/goat.",
			act.Type()))
	}
}

func (t *Executor) executeRequest(eng engine.Engine, req *goatfile.Request, gf goatfile.Goatfile) (err error) {
	req.Merge(gf.Defaults)

	if !t.isAbortOnError(req) {
		defer func() {
			if err != nil {
				err = NewNoAbortError(err)
			}
		}()
	}

	state := eng.State()

	err = req.PreSubstituteWithParams(state)
	if err != nil {
		return errs.WithPrefix("failed pre-substituting request with parameters:", err)
	}

	preScript, err := util.ReadReaderToString(req.PreScript.Reader())
	if err != nil {
		return errs.WithPrefix("reading preScript failed:", err)
	}

	if preScript != "" {
		err = eng.Run(preScript)
		if err != nil {
			return errs.WithPrefix("preScript failed:", err)
		}
		state = eng.State()
	}

	err = req.SubstituteWithParams(state)
	if err != nil {
		return errs.WithPrefix("failed substituting request with parameters:",
			NewParamsParsingError(err))
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

	err = req.InsertRawDataIntoBody(state)
	if err != nil {
		return errs.WithPrefix("failed inserting raw variable in body:",
			NewParamsParsingError(err))
	}

	err = req.InsertRawDataIntoFormData(state)
	if err != nil {
		return errs.WithPrefix("failed reading raw variable:",
			NewParamsParsingError(err))
	}

	httpReq, err := req.ToHttpRequest()
	if err != nil {
		return errs.WithPrefix("failed transforming to http request:", err)
	}

	if authOpts, ok := AuthOptionsFromMap(req.Auth); ok {
		httpReq.Header.Set("Authorization", authOpts.HeaderValue())
	}

	reqOpts := requester.OptionsFromMap(req.Options)
	httpResp, err := t.req.Do(httpReq, reqOpts)
	if err != nil {
		return errs.WithPrefix("http request failed:", err)
	}

	resp, err := FromHttpResponse(httpResp, req.Options)
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

func (t *Executor) executeExecute(params goatfile.Execute, eng engine.Engine, showTeardownParamErrors bool) (Result, error) {
	pth := goatfile.Extend(path.Join(path.Dir(params.Path), params.File), goatfile.FileExtension)
	gf, err := t.parseGoatfile(pth)
	if err != nil {
		return Result{}, err
	}

	state := eng.State()

	err = goatfile.ApplyTemplateToMap(params.Params, state)
	if err != nil {
		return Result{}, err
	}

	log := log.Tagged(strings.TrimSuffix(gf.Path, ".goat"))

	isolatedEng := t.engineMaker()
	isolatedEng.SetState(params.Params)

	res, err := t.executeGoatfile(log, gf, isolatedEng, false, showTeardownParamErrors)
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

func (t *Executor) isSkip(section goatfile.SectionName) bool {
	for _, s := range t.Skip {
		if strings.ToLower(s) == string(section) {
			return true
		}
	}
	return false
}

func (t *Executor) isAbortOnError(req *goatfile.Request) bool {
	opts := AbortOptionsFromMap(req.Options)

	if opts.AlwaysAbort {
		return true
	}

	if opts.NoAbort || t.NoAbort {
		return false
	}

	return true
}

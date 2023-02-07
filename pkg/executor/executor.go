package executor

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/studio-b12/goat/pkg/advancer"
	"github.com/studio-b12/goat/pkg/engine"
	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/goatfile"
	"github.com/studio-b12/goat/pkg/requester"
)

// Executor parses a Goatfile and executes it.
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
// from the given directory. The given initialParams are
// used as initial state for the runtime engine.
func (t *Executor) Execute(path string, initialParams engine.State) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat failed: %s", err.Error())
	}

	if stat.IsDir() {
		return t.executeFromDir(path, initialParams)
	}

	gf, err := t.parseGoatfile(path)
	if err != nil {
		return err
	}

	log.Debug().Msg("Executing goatfile ...")
	return t.ExecuteGoatfile(gf, initialParams)
}

// ExecuteGoatfile runs the given parsed Goatfile. The given initialParams are
// used as initial state for the runtime engine.
func (t *Executor) ExecuteGoatfile(gf goatfile.Goatfile, initialParams engine.State) (err error) {
	log.Debug().Msg("Parsed Goatfile\n" + gf.String())

	if t.Dry {
		log.Warn().Msg("This is a dry run: no requets will be executed")
		return nil
	}

	var errsNoAbort errs.Errors

	eng := t.engineMaker()
	eng.SetState(initialParams)

	defer func() {
		// Teardown Procedures

		if t.isSkip("teardown") {
			log.Warn().Msg("skipping teardown steps")
			return
		}

		for _, req := range gf.Teardown {
			err := t.executeRequest(eng, req)
			if err != nil {
				log.Err(err).Str("req", req.String()).Msg("Teardown step failed")

				// If the returned error comes from the params parsing step, don't
				// cancel the teardown execution. See the following issue for more information.
				// https://github.com/studio-b12/goat/issues/9
				if errs.IsOfType[ParamsParsingError](err) {
					continue
				}

				if !t.isAbortOnError(req) {
					errsNoAbort = errsNoAbort.Append(err)
					continue
				}

				break
			}

			log.Info().Str("req", req.String()).Msg("Teardown step completed")
		}
	}()

	// Setup Procedures

	if t.isSkip("setup") {
		log.Warn().Msg("skipping setup steps")
	} else {
		for _, req := range gf.Setup {
			err := t.executeRequest(eng, req)
			if err != nil {
				log.Err(err).Str("req", req.String()).Msg("Setup step failed")
				if !t.isAbortOnError(req) {
					errsNoAbort = errsNoAbort.Append(err)
					continue
				}
				return err
			}

			log.Info().Str("req", req.String()).Msg("Setup step completed")
		}
	}

	// Test Procedures

	if t.isSkip("tests") {
		log.Warn().Msg("skipping test steps")
	} else {
		for _, req := range gf.Tests {
			err := t.executeTest(req, eng, gf)
			if err != nil {
				if !t.isAbortOnError(req) {
					errsNoAbort = errsNoAbort.Append(err)
					continue
				}
				return err
			}
		}
	}

	return errsNoAbort.Condense()
}

func (t *Executor) executeFromDir(path string, initialParams engine.State) error {
	var goatfiles []goatfile.Goatfile

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, _ error) error {
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
		return err
	}

	if len(goatfiles) == 0 {
		return errors.New("no Goatfiles found to execute")
	}

	var errs errs.Errors

	for _, gf := range goatfiles {
		log.Info().Str("path", gf.Path).Msg("Executing batch ...")

		err = t.ExecuteGoatfile(gf, initialParams)
		if err != nil {
			log.Err(err).Msg("Batch execution failed")
			errs = errs.Append(err)
			continue
		}

		log.Info().Str("path", gf.Path).Msg("Batch finished successfully")
	}

	if errs.HasSome() {
		return BatchExecutionError{
			Inner: errs,
			Total: len(goatfiles),
		}
	}

	return nil
}

func (t *Executor) parseGoatfile(path string) (gf goatfile.Goatfile, err error) {
	log.Debug().Str("from", path).Msg("Parsing goatfile ...")

	data, err := os.ReadFile(path)
	if err != nil {
		return goatfile.Goatfile{}, fmt.Errorf("failed reading file: %s", err.Error())
	}

	relCurrDir := filepath.Dir(path)
	gf, err = goatfile.Unmarshal(string(data), relCurrDir)
	if err != nil {
		if errs.IsOfType[goatfile.ParseError](err) {
			return goatfile.Goatfile{}, fmt.Errorf("failed parsing goatfile at %s:%s", path, err.Error())
		}
		return goatfile.Goatfile{}, fmt.Errorf("failed parsing goatfile %s: %s", path, err.Error())
	}

	gf.Path = path

	return gf, nil
}

func (t *Executor) executeTest(req goatfile.Request, eng engine.Engine, gf goatfile.Goatfile) (err error) {
	var errsNoAbort errs.Errors

	defer func() {
		// Teardown-Each steps

		if t.isSkip("teardown-each") {
			log.Warn().Msg("skipping teardown-each steps")
			return
		}

		for _, postReq := range gf.TeardownEach {
			err := t.executeRequest(eng, postReq)
			if err != nil {
				log.Err(err).Str("req", req.String()).Msg("Post-Each step failed")

				err = fmt.Errorf("post-setup-each step failed: %s", err.Error())

				// If the returned error comes from the params parsing step, don't
				// cancel the teardown-each execution. See the following issue for more information.
				// https://github.com/studio-b12/goat/issues/9
				if errs.IsOfType[ParamsParsingError](err) {
					continue
				}

				if t.isAbortOnError(postReq) {
					break
				}

				errsNoAbort = errsNoAbort.Append(err)
				continue
			}

			log.Info().Str("req", req.String()).Msg("Teardown-Each step completed")
		}
	}()

	// Setup-Each Steps

	if t.isSkip("setup-each") {
		log.Warn().Msg("skipping setup-each steps")
	} else {
		for _, preReq := range gf.SetupEach {
			err := t.executeRequest(eng, preReq)
			if err != nil {
				log.Err(err).Str("req", req.String()).Msg("Setup-Each step failed")

				err = fmt.Errorf("Setup-Each step failed: %s", err.Error())

				if !t.isAbortOnError(preReq) {
					errsNoAbort = errsNoAbort.Append(err)
					continue
				}

				return err
			}

			log.Info().Str("req", req.String()).Msg("Setup-Each step completed")
		}
	}

	// Actual Test Step

	err = t.executeRequest(eng, req)
	if err != nil {
		log.Err(err).Str("req", req.String()).Msg("Test step failed")

		if t.isAbortOnError(req) {
			return err
		}

		errsNoAbort = errsNoAbort.Append(err)
	} else {
		log.Info().Str("req", req.String()).Msg("Test completed")
	}

	return errsNoAbort.Condense()
}

func (t *Executor) executeRequest(eng engine.Engine, req goatfile.Request) (err error) {

	state := eng.State()
	err = req.ParseWithParams(state)
	if err != nil {
		return errs.WithPrefix("failed infusing request with parameters:",
			ParamsParsingError(err))
	}

	execOpts := ExecOptionsFromMap(req.Options)
	if !execOpts.Condition {
		log.Warn().Str("req", req.String()).Msg("Skipped due to condition")
		return nil
	}

	t.Waiter.Wait()

	httpReq, err := req.ToHttpRequest()
	if err != nil {
		return fmt.Errorf("failed transforming to http request: %s", err.Error())
	}

	reqOpts := requester.OptionsFromMap(req.Options)
	httpResp, err := t.req.Do(httpReq, reqOpts)
	if err != nil {
		return fmt.Errorf("http request failed: %s", err.Error())
	}

	resp, err := FromHttpResponse(httpResp)
	if err != nil {
		return fmt.Errorf("response interpretation failed: %s", err.Error())
	}

	state.Merge(engine.State{"response": resp})
	eng.SetState(state)

	script, err := readReader(req.Script.Reader())
	if err != nil {
		return errs.WithPrefix("reading script failed:", err)
	}

	err = eng.Run(script)
	if err != nil {
		return fmt.Errorf("script failed: %s", err.Error())
	}

	return nil
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

func readReader(r io.Reader, err error) (string, error) {
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

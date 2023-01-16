package executor

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/studio-b12/gurl/pkg/engine"
	"github.com/studio-b12/gurl/pkg/gurlfile"
	"github.com/studio-b12/gurl/pkg/requester"
)

// Executor parses a Gurlfile and executes it.
type Executor struct {
	engineMaker func() engine.Engine
	req         requester.Requester
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

	return &t
}

// ExecuteFromDir executes a single or multiple Gurlfiles
// from the given directory. The given initialParams are
// used as initial state for the runtime engine.
func (t *Executor) ExecuteFromDir(path string, initialParams engine.State) error {
	log.Debug().Interface("initialParams", initialParams).Send()

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat failed: %s", err.Error())
	}

	if stat.IsDir() {
		return fmt.Errorf("Execution from directories is currently not implemented.")
	}

	log.Debug().Str("from", path).Msg("Reading gurlfile ...")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading file: %s", err.Error())
	}

	log.Debug().Msg("Parsing gurlfile ...")
	gf, err := gurlfile.Unmarshal(string(data), initialParams)
	if err != nil {
		return fmt.Errorf("failed parsing gurlfile: %s", err.Error())
	}

	log.Debug().Msg("Executing gurlfile ...")
	return t.Execute(gf, initialParams)
}

// Execute runs the given parsed Gurlfile. The given initialParams are
// used as initial state for the runtime engine.
func (t *Executor) Execute(gf gurlfile.Gurlfile, initialParams engine.State) (err error) {
	log.Debug().Interface("gf", gf).Send()

	eng := t.engineMaker()
	eng.SetState(initialParams)

	defer func() {
		// Teardown Procedures

		for _, req := range gf.Teardown {
			err := t.executeRequest(eng, req)
			if err != nil {
				log.Error().Str("req", req.String()).Err(err).Msg("teardown step failed")
				break
			}
			log.Info().Str("req", req.String()).Msg("Teardown step completed")
		}
	}()

	// Setup Procedures

	for _, req := range gf.Setup {
		err := t.executeRequest(eng, req)
		if err != nil {
			return err
		}
		log.Info().Str("req", req.String()).Msg("Setup step completed")
	}

	// Test Procedures

	for _, req := range gf.Tests {
		err := t.executeTest(req, eng, gf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Executor) executeTest(req gurlfile.Request, eng engine.Engine, gf gurlfile.Gurlfile) (err error) {
	defer func() {
		for _, postReq := range gf.TeardownEach {
			err := t.executeRequest(eng, postReq)
			if err != nil {
				err = fmt.Errorf("post-setup-each step failed: %s", err.Error())
				break
			}
			log.Info().Str("req", req.String()).Msg("Teardown-Each step completed")
		}
	}()

	for _, preReq := range gf.SetupEach {
		err := t.executeRequest(eng, preReq)
		if err != nil {
			return fmt.Errorf("pre-setup-each step failed: %s", err.Error())
		}
		log.Info().Str("req", req.String()).Msg("Setup-Each step completed")
	}

	err = t.executeRequest(eng, req)
	if err != nil {
		return err
	}
	log.Info().Str("req", req.String()).Msg("Test completed")

	return nil
}

func (t *Executor) executeRequest(eng engine.Engine, req gurlfile.Request) (err error) {
	defer func() {
		// If an error is returned, wrap the error
		// in a ContextError.
		if err != nil {
			err = req.WrapErr(err)
		}
	}()

	state := eng.State()
	parsedReq, err := req.ParseWithParams(state)
	if err != nil {
		return fmt.Errorf("failed infusing request with parameters: %s", err.Error())
	}

	httpReq, err := parsedReq.ToHttpRequest()
	if err != nil {
		return fmt.Errorf("failed transforming to http request: %s", err.Error())
	}

	httpResp, err := t.req.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http request failed: %s", err.Error())
	}

	resp, err := FromHttpResponse(httpResp)
	if err != nil {
		return fmt.Errorf("response interpretation failed: %s", err.Error())
	}

	state.Merge(engine.State{"response": resp})
	eng.SetState(state)

	err = eng.Run(parsedReq.Script)
	if err != nil {
		return fmt.Errorf("script failed: %s", err.Error())
	}

	return nil
}

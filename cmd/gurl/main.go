package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/studio-b12/gurl/internal/version"
	"github.com/studio-b12/gurl/pkg/config"
	"github.com/studio-b12/gurl/pkg/engine"
	"github.com/studio-b12/gurl/pkg/executor"
	"github.com/studio-b12/gurl/pkg/requester"
)

type Args struct {
	Gurlfile string `arg:"positional,required" help:"Gurlfile(s) location"`
	LogLevel int    `arg:"-l,--loglevel" default:"1" help:"Logging level (see https://github.com/rs/zerolog#leveled-logging for reference)"`
	Params   string `arg:"-p,--params" help:"Params file location"`
}

func main() {

	var args Args
	arg.MustParse(&args)

	zerolog.SetGlobalLevel(zerolog.Level(args.LogLevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	state, err := config.Parse[engine.State](args.Params, "GURL_")
	if err != nil {
		log.Fatal().Err(err).Msg("parameter parsing failed")
	}

	engineMaker := engine.NewGoja
	req, err := requester.NewHttpWithCookies()
	if err != nil {
		log.Fatal().Err(err).Msg("requester initialization failed")
	}

	executor := executor.New(engineMaker, req)
	err = executor.ExecuteFromDir(args.Gurlfile, state)
	if err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}

	log.Info().Msg("Execution finished successfully")
}

func (Args) Description() string {
	return "Automation tool for executing and evaluating API requests."
}

func (Args) Version() string {
	return fmt.Sprintf("gurl v%s (%s %s %s)",
		version.Version, version.CommitHash, version.BuildDate, runtime.Version())
}

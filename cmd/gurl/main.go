package main

import (
	"flag"
	"os"
	"time"

	"github.com/dop251/goja"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/studio-b12/gurl/pkg/config"
	"github.com/studio-b12/gurl/pkg/engine"
	"github.com/studio-b12/gurl/pkg/executor"
	"github.com/studio-b12/gurl/pkg/requester"
)

var (
	fLogLevel = flag.Int("level", 2, "Log level - see https://github.com/rs/zerolog#leveled-logging")
	fParams   = flag.String("params", "", "Path to file with initial parameters")
)

var vm = goja.New()

func Assert(v bool) {
	if !v {
		panic(vm.ToValue("test"))
	}
}

func main() {
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.Level(*fLogLevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	path := flag.Arg(0)
	if path == "" {
		log.Fatal().Msg("No path to execute is given")
	}

	var err error
	state := engine.State{}

	if *fParams != "" {
		state, err = config.Parse[engine.State](*fParams, "GURL_")
		if err != nil {
			log.Fatal().Err(err).Msg("parameter parsing failed")
		}
	}

	engineMaker := engine.NewGoja
	req, err := requester.NewHttpWithCookies()
	if err != nil {
		log.Fatal().Err(err).Msg("requester initialization failed")
	}

	executor := executor.New(engineMaker, req)
	err = executor.ExecuteFromDir(path, state)
	if err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}

	log.Info().Msg("Execution finished successfully")
}

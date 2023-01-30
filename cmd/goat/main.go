package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/studio-b12/goat/internal/embedded"
	"github.com/studio-b12/goat/internal/version"
	"github.com/studio-b12/goat/pkg/advancer"
	"github.com/studio-b12/goat/pkg/config"
	"github.com/studio-b12/goat/pkg/engine"
	"github.com/studio-b12/goat/pkg/executor"
	"github.com/studio-b12/goat/pkg/requester"
)

type Args struct {
	Goatfile string `arg:"positional" help:"Goatfile(s) location"`

	Arg      []string      `arg:"-a,--args" help:"Pass params as key value arguments into the execution (format: key=value)"`
	Delay    time.Duration `arg:"-d,--delay" help:"Delay requests by the given duration"`
	Dry      bool          `arg:"--dry" help:"Only parse the goatfile(s) without executing any requests"`
	Gradual  bool          `arg:"-g,--gradual" help:"Advance the requests maually"`
	LogLevel int           `arg:"-l,--loglevel" default:"1" help:"Logging level (see https://github.com/rs/zerolog#leveled-logging for reference)"`
	New      bool          `arg:"--new" help:"Create a new base Goatfile"`
	NoAbort  bool          `arg:"--no-abort" help:"Do not abort batch execution on error."`
	Params   string        `arg:"-p,--params" help:"Params file location"`
	Skip     []string      `arg:"--skip" help:"Section(s) to be skipped during execution"`
}

func main() {

	var args Args
	argParser := arg.MustParse(&args)

	zerolog.SetGlobalLevel(zerolog.Level(args.LogLevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	if args.New {
		createNewGoatfile(args.Goatfile)
		return
	}

	if args.Goatfile == "" {
		argParser.Fail("Goatfile must be specified.")
	}

	state, err := config.Parse(args.Params, "Goat_", engine.State{})
	if err != nil {
		log.Fatal().Err(err).Msg("parameter parsing failed")
	}

	config.ParseKVArgs(args.Arg, state)

	engineMaker := engine.NewGoja
	req := requester.NewHttpWithCookies()

	executor := executor.New(engineMaker, req)
	executor.Dry = args.Dry
	executor.Skip = args.Skip
	executor.NoAbort = args.NoAbort

	if args.Gradual {
		ad := make(advancer.Channel)
		executor.Waiter = ad
		go advanceManually(ad)
	} else if args.Delay != 0 {
		log.Info().Msgf("Delay mode: Advancing every %s", args.Delay.String())
		executor.Waiter = advancer.NewTicker(args.Delay)
	}

	log.Debug().Interface("initialParams", state).Send()

	err = executor.Execute(args.Goatfile, state)
	if err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}

	log.Info().Msg("Execution finished successfully")
}

func (Args) Description() string {
	return "Automation tool for executing and evaluating API requests."
}

func (Args) Version() string {
	return fmt.Sprintf("goat %s (%s %s %s)",
		version.Version, version.CommitHash, version.BuildDate, runtime.Version())
}

func advanceManually(a advancer.Advancer) {
	scanner := bufio.NewScanner(os.Stdin)
	log.Info().Msg("Manual mode: Press [enter] to advance requests")
	for {
		scanner.Scan()
		a.Advance()
	}
}

func createNewGoatfile(name string) {
	if name == "" {
		name = "tests.goat"
	}

	pathTo := filepath.Dir(name)
	if pathTo != "" {
		_, err := os.Stat(pathTo)
		if os.IsNotExist(err) {
			err = os.MkdirAll(pathTo, os.ModePerm)
		}
		if err != nil {
			log.Fatal().
				Err(err).
				Str("at", name).
				Msg("Failed creating new goatfile: Failed creating directory")
		}
	}

	err := os.WriteFile(name, embedded.NewGoatfile, fs.ModePerm)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("at", name).
			Msg("Failed creating new goatfile")
	}

	log.Info().
		Str("at", name).
		Msg("Goatfile created")
}

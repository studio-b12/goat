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
	"github.com/studio-b12/gurl/internal/embedded"
	"github.com/studio-b12/gurl/internal/version"
	"github.com/studio-b12/gurl/pkg/advancer"
	"github.com/studio-b12/gurl/pkg/config"
	"github.com/studio-b12/gurl/pkg/engine"
	"github.com/studio-b12/gurl/pkg/executor"
	"github.com/studio-b12/gurl/pkg/requester"
)

type Args struct {
	Gurlfile string `arg:"positional" help:"Gurlfile(s) location"`

	Delay    time.Duration `arg:"--delay" help:"Delay requests by the given duration."`
	Dry      bool          `arg:"--dry" help:"Only parse the gurlfile(s) without executing any requests"`
	LogLevel int           `arg:"-l,--loglevel" default:"1" help:"Logging level (see https://github.com/rs/zerolog#leveled-logging for reference)"`
	Manual   bool          `arg:"--manual" help:"Advance the requests maually."`
	New      bool          `arg:"--new" help:"Create a new base Gurlfile."`
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
		createNewGurlfile(args.Gurlfile)
		return
	}

	if args.Gurlfile == "" {
		argParser.Fail("Gurlfile must be specified.")
	}

	state, err := config.Parse[engine.State](args.Params, "GURL_")
	if err != nil {
		log.Fatal().Err(err).Msg("parameter parsing failed")
	}

	engineMaker := engine.NewGoja
	req := requester.NewHttpWithCookies()

	executor := executor.New(engineMaker, req)
	executor.Dry = args.Dry
	executor.Skip = args.Skip

	if args.Manual {
		ad := make(advancer.Channel)
		executor.Waiter = ad
		go advanceManually(ad)
	} else if args.Delay != 0 {
		log.Info().Msgf("Delay mode: Advancing every %s", args.Delay.String())
		executor.Waiter = advancer.NewTicker(args.Delay)
	}

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

func advanceManually(a advancer.Advancer) {
	scanner := bufio.NewScanner(os.Stdin)
	log.Info().Msg("Manual mode: Press [enter] to advance requests")
	for {
		scanner.Scan()
		a.Advance()
	}
}

func createNewGurlfile(name string) {
	if name == "" {
		name = "tests.gurl"
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
				Msg("Failed creating new gurlfile: Failed creating directory")
		}
	}

	err := os.WriteFile(name, embedded.NewGurlfile, fs.ModePerm)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("at", name).
			Msg("Failed creating new gurlfile")
	}

	log.Info().
		Str("at", name).
		Msg("Gurlfile created")
}

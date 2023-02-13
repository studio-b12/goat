package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/studio-b12/goat/internal/embedded"
	"github.com/studio-b12/goat/internal/version"
	"github.com/studio-b12/goat/pkg/advancer"
	"github.com/studio-b12/goat/pkg/clr"
	"github.com/studio-b12/goat/pkg/config"
	"github.com/studio-b12/goat/pkg/engine"
	"github.com/studio-b12/goat/pkg/executor"
	"github.com/studio-b12/goat/pkg/requester"
)

type Args struct {
	Goatfile string `arg:"positional" help:"Goatfile(s) location"`

	Arg      []string      `arg:"-a,--args,separate" help:"Pass params as key value arguments into the execution (format: key=value)"`
	Delay    time.Duration `arg:"-d,--delay" help:"Delay requests by the given duration"`
	Dry      bool          `arg:"--dry" help:"Only parse the goatfile(s) without executing any requests"`
	Gradual  bool          `arg:"-g,--gradual" help:"Advance the requests maually"`
	Json     bool          `arg:"--json" help:"Use JSON format instead of pretty console format for logging"`
	LogLevel int           `arg:"-l,--loglevel" default:"1" help:"Logging level (see https://github.com/rs/zerolog#leveled-logging for reference)"`
	New      bool          `arg:"--new" help:"Create a new base Goatfile"`
	NoAbort  bool          `arg:"--no-abort" help:"Do not abort batch execution on error"`
	NoColor  bool          `arg:"--no-color" help:"Supress colored log output"`
	Params   []string      `arg:"-p,--params,separate" help:"Params file location(s)"`
	Skip     []string      `arg:"--skip,separate" help:"Section(s) to be skipped during execution"`
}

func main() {

	var args Args
	argParser := arg.MustParse(&args)

	zerolog.SetGlobalLevel(zerolog.Level(args.LogLevel))
	if !args.Json {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    args.NoColor,
		})
	}

	clr.SetEnable(!args.Json && !args.NoColor)

	if args.New {
		createNewGoatfile(args.Goatfile)
		return
	}

	if args.Goatfile == "" {
		argParser.Fail("Goatfile must be specified.")
	}

	state, err := config.Parse(args.Params, "GOAT_", engine.State{})
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

	log.Debug().Msgf("Initial Params\n%s", state)

	err = executor.Execute(args.Goatfile, state)
	if err != nil {
		log.Fatal().Err(err).Msg(clr.Print(clr.Format("execution failed", clr.ColorFGRed, clr.FormatBold)))
	}

	log.Info().Msg(clr.Print(clr.Format("Execution finished successfully", clr.ColorFGGreen, clr.FormatBold)))
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
	log.Info().Msg(
		"Gradual mode: Press [enter] to advance requests, " +
			"enter 'quit' / 'q' to quit the execution or " +
			"enter 'continue' / 'cont' / 'c' to stop gradual advancement.")

	var skip bool
	for {
		if !skip {
			scanner.Scan()
			txt := scanner.Text()
			switch strings.ToLower(txt) {
			case "quit", "q":
				log.Warn().Msg("Aborted.")
				os.Exit(1)
			case "continue", "cont", "c":
				skip = true
			}
		}
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

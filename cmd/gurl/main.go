package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/studio-b12/gurl/internal/version"
	"github.com/studio-b12/gurl/pkg/advancer"
	"github.com/studio-b12/gurl/pkg/config"
	"github.com/studio-b12/gurl/pkg/engine"
	"github.com/studio-b12/gurl/pkg/executor"
	"github.com/studio-b12/gurl/pkg/requester"
)

type Args struct {
	Gurlfile string        `arg:"positional,required" help:"Gurlfile(s) location"`
	LogLevel int           `arg:"-l,--loglevel" default:"1" help:"Logging level (see https://github.com/rs/zerolog#leveled-logging for reference)"`
	Params   string        `arg:"-p,--params" help:"Params file location"`
	Dry      bool          `arg:"--dry" help:"Only parse the gurlfile(s) without executing any requests"`
	Skip     []string      `arg:"--skip" help:"Section(s) to be skipped during execution"`
	Manual   bool          `arg:"--manual" help:"Advance the requests maually."`
	Delay    time.Duration `arg:"--delay" help:"Delay requests by the given duration."`
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

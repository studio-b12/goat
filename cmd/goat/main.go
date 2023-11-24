package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/studio-b12/goat/internal/embedded"
	"github.com/studio-b12/goat/internal/version"
	"github.com/studio-b12/goat/pkg/advancer"
	"github.com/studio-b12/goat/pkg/clr"
	"github.com/studio-b12/goat/pkg/config"
	"github.com/studio-b12/goat/pkg/engine"
	"github.com/studio-b12/goat/pkg/executor"
	"github.com/studio-b12/goat/pkg/requester"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/level"
	"github.com/zekrotja/rogu/log"
)

type Args struct {
	Goatfile []string `arg:"positional" help:"Goatfile(s) location"`

	Arg      []string      `arg:"-a,--args,separate" help:"Pass params as key value arguments into the execution (format: key=value)"`
	Delay    time.Duration `arg:"-d,--delay" help:"Delay requests by the given duration"`
	Dry      bool          `arg:"--dry" help:"Only parse the goatfile(s) without executing any requests"`
	Gradual  bool          `arg:"-g,--gradual" help:"Advance the requests maually"`
	Json     bool          `arg:"--json" help:"Use JSON format instead of pretty console format for logging"`
	LogLevel string        `arg:"-l,--loglevel" default:"info" help:"Logging level (see https://github.com/zekrotja/rogu#levels for reference)"`
	New      bool          `arg:"--new" help:"Create a new base Goatfile"`
	NoAbort  bool          `arg:"--no-abort" help:"Do not abort batch execution on error"`
	NoColor  bool          `arg:"--no-color" help:"Supress colored log output"`
	Params   []string      `arg:"-p,--params,separate" help:"Params file location(s)"`
	Silent   bool          `arg:"-s,--silent" help:"Disables all logging output"`
	Skip     []string      `arg:"--skip,separate" help:"Section(s) to be skipped during execution"`
	Secure   bool          `arg:"--secure" help:"Validate TLS certificates"`
}

func main() {

	var args Args
	argParser := arg.MustParse(&args)

	if args.Silent {
		log.SetLevel(level.Off)
	} else {
		lvl, ok := level.LevelFromString(args.LogLevel)
		if !ok {
			log.Fatal().Msg("invalid log level; see https://github.com/zekrotja/rogu#levels for reference")
			return
		}
		log.SetLevel(lvl)
	}

	if args.Json {
		w := rogu.NewJsonWriter(os.Stdout)
		log.SetWriter(w)
	} else {
		w := rogu.NewPrettyWriter(os.Stdout)
		w.NoColor = args.NoColor
		w.TimeFormat = time.RFC3339
		w.StyleTag.Width(20)
		log.SetWriter(w)
	}

	clr.SetEnable(!args.Json && !args.NoColor)

	if args.New {
		createNewGoatfile(args.Goatfile)
		return
	}

	if len(args.Goatfile) == 0 {
		argParser.Fail("Goatfile must be specified.")
		return
	}

	state, err := config.Parse(args.Params, "GOAT_", engine.State{})
	if err != nil {
		log.Fatal().Err(err).Msg("parameter parsing failed")
		return
	}

	err = config.ParseKVArgs(args.Arg, state)
	if err != nil {
		log.Fatal().Err(err).Msg("argument parsing failed")
		return
	}

	engineMaker := engine.NewGoja
	req := requester.NewHttpWithCookies(func(client *http.Client) {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !args.Secure,
		}}
	})

	exec := executor.New(engineMaker, req)
	exec.Dry = args.Dry
	exec.Skip = args.Skip
	exec.NoAbort = args.NoAbort

	if args.Gradual {
		ad := make(advancer.Channel)
		exec.Waiter = ad
		go advanceManually(ad)
	} else if args.Delay != 0 {
		log.Info().Msgf("Delay mode: Advancing every %s", args.Delay.String())
		exec.Waiter = advancer.NewTicker(args.Delay)
	}

	log.Debug().Msgf("Initial Params\n%s", state)

	res, err := exec.Execute(args.Goatfile, state)
	res.Log()
	if err != nil {
		entry := log.Fatal().Err(err)

		if batchErr, ok := err.(executor.BatchResultError); ok {
			coloredPathes := batchErr.ErrorMessages()
			for i, p := range coloredPathes {
				coloredPathes[i] = clr.Print(clr.Format(p, clr.ColorFGRed))
			}
			entry.Field("failed_files", coloredPathes)
		}

		entry.Msg(clr.Print(clr.Format("execution failed", clr.ColorFGRed, clr.FormatBold)))
		return
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

func createNewGoatfile(names []string) {
	name := "tests.goat"
	if len(names) > 0 {
		name = names[0]
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
				Field("at", name).
				Msg("Failed creating new goatfile: Failed creating directory")
		}
	}

	err := os.WriteFile(name, embedded.NewGoatfile, fs.ModePerm)
	if err != nil {
		log.Fatal().
			Err(err).
			Field("at", name).
			Msg("Failed creating new goatfile")
		return
	}

	log.Info().
		Field("at", name).
		Msg("Goatfile created")
}

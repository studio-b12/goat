package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/studio-b12/goat/pkg/errs"

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

	Arg           []string      `arg:"-a,--args,separate" help:"Pass params as key value arguments into the execution (format: key=value)"`
	Delay         time.Duration `arg:"-d,--delay,env:GOATARG_DELAY" help:"Delay requests by the given duration"`
	Dry           bool          `arg:"--dry" help:"Only parse the goatfile(s) without executing any requests"`
	Gradual       bool          `arg:"-g,--gradual" help:"Advance the requests maually"`
	Json          bool          `arg:"--json,env:GOATARG_JSON" help:"Use JSON format instead of pretty console format for logging"`
	LogLevel      level.Level   `arg:"-l,--loglevel,env:GOATARG_LOGLEVEL" default:"info" help:"Logging level"`
	New           bool          `arg:"--new" help:"Create a new base Goatfile"`
	NoAbort       bool          `arg:"--no-abort,env:GOATARG_NOABORT" help:"Do not abort batch execution on error"`
	NoColor       bool          `arg:"--no-color,env:GOATARG_NOCOLOR" help:"Supress colored log output"`
	Params        []string      `arg:"-p,--params,separate,env:GOATARG_PARAMS" help:"Params file location(s)"`
	Profile       []string      `arg:"-P,--profile,separate,env:GOATARG_PROFILE" help:"Select a profile from your home config"`
	ReducedErrors bool          `arg:"-R,--reduced-errors,env:GOATARG_REDUCEDERRORS" help:"Hide template errors in teardown steps"`
	Secure        bool          `arg:"--secure,env:GOATARG_SECURE" help:"Validate TLS certificates"`
	Silent        bool          `arg:"-s,--silent,env:GOATARG_SILENT" help:"Disables all logging output"`
	Skip          []string      `arg:"--skip,separate,env:GOATARG_SKIP" help:"Section(s) to be skipped during execution"`
	RetryFailed   bool          `arg:"--retry-failed,env:GOATARG_RETRYFAILED" help:"Retry files which have failed in the previous run"`
}

func main() {

	var args Args
	argParser := arg.MustParse(&args)

	if args.Silent {
		log.SetLevel(level.Off)
	} else {
		log.SetLevel(args.LogLevel)
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

	goatfiles := args.Goatfile

	if args.RetryFailed {
		failed, err := loadLastFailedFiles()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed loading last failed files")
			return
		}
		if len(failed) == 0 {
			log.Fatal().Msg("No failed files have been recorded in previous runs")
			return
		}

		goatfiles = failed
	}

	if len(goatfiles) == 0 {
		argParser.Fail("Goatfile must be specified.")
		return
	}

	state := make(engine.State)

	err := config.LoadProfiles(args.Profile, state)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed loading profiles")
		return
	}

	cfgState, err := config.Parse[engine.State](args.Params, "GOAT_")
	if err != nil {
		log.Fatal().Err(err).Msg("parameter parsing failed")
		return
	}
	state.Merge(cfgState)

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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer cancel()

	exec := executor.New(ctx, engineMaker, req)
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

	res, err := exec.Execute(goatfiles, state, !args.ReducedErrors)
	res.Log()
	if err != nil {
		if args.ReducedErrors {
			err = filterTeardownParamErrors(err)
		}

		entry := log.Fatal().Err(err)

		if batchErr, ok := errs.As[*executor.BatchResultError](err); ok {
			if sErr := storeLastFailedFiles(batchErr.FailedFiles()); err != nil {
				log.Error().Err(sErr).Msg("failed storing latest failed files")
			}

			coloredMessages := batchErr.ErrorMessages()
			for i, p := range coloredMessages {
				coloredMessages[i] = clr.Print(clr.Format(p, clr.ColorFGRed))
			}
			entry.Field("failed_files", coloredMessages)
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

func filterTeardownParamErrors(err error) error {
	switch tErr := err.(type) {
	case errs.Errors:
		newErrors := make(errs.Errors, 0, len(tErr))
		for _, e := range tErr {
			if errs.IsOfType[executor.TeardownError](e) && errs.IsOfType[executor.ParamsParsingError](e) {
				continue
			}
			newErrors = append(newErrors, e)
		}
		return newErrors
	case *executor.BatchResultError:
		tErr.Inner = filterTeardownParamErrors(tErr.Inner).(errs.Errors)
		return tErr
	default:
		return err
	}
}

const lastFailedRunFileName = "goat_last_failed_run"

func storeLastFailedFiles(paths []string) error {
	failedRunPath := path.Join(os.TempDir(), lastFailedRunFileName)
	f, err := os.Create(failedRunPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(strings.Join(paths, "\n"))
	return err
}

func loadLastFailedFiles() (paths []string, err error) {
	failedRunPath := path.Join(os.TempDir(), lastFailedRunFileName)
	f, err := os.Open(failedRunPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return nil, err
		}
		paths = append(paths, scanner.Text())
	}

	return paths, nil
}

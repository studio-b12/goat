# Command Line Tool

The `goat` command is used to create, validate and execute Goatfiles.

You can pass one or multiple Goatfiles or directories as positional argument. For example:
```
goat issue_32.goat tests/integrationtests
```

When passing in a directory, Goat will look for any `*.goat` file recursively. Files and directories prefixed with an underscore (`_`) are ignored. This is especially useful for Goatfiles which are only supposed to be imported or executed in other Goatfiles. If you want to read more about this, take a look into the [Project Structure section](../project-structure.md). 

## Flags

The following sections provide further information about the various flags which can be passed to `goat`.

- **`-a ARGS`, `--args ARGS`**  
  Pass params as key value arguments into the execution. If you want to pass multiple args, specify each pair with its own parameter.  
  *Example: `-a hello=world -a user.name=foo -a "user.password=bar bazz"`*

- **`--delay DELAY`, ` -d DELAY`**  
  Delay all requests by the given duration. The duration is formatted according to the format of Go's [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) function.  
  *Example: `-d 1s30ms`*

- **`--dry`**  
  Only parse the goatfile(s) without executing any requests.

- **`--gradual`, ` -g`**  
  Advance the requests maually.

- **`--json`**  
  Use JSON format instead of pretty console format for logging.

- **`--loglevel LOGLEVEL`, ` -l LOGLEVEL`**  
  Logging level. [Here](https://github.com/zekroTJA/rogu#levels) you can see which values you can use for log levels.  
  *Example: `-l trace`*

- **`--new`**  
  Create a new base Goatfile.

- **`--no-abort`**  
  Do not abort the batch execution on error.

- **`--no-color`**  
  Supress colored log output.

- **`--params PARAMS`, ` -p PARAMS`**  
  Pass parameters in parameter files. These can be either TOML, YAML or JSON files. If you want to pass multiple parameter files, specify each one with its own parameter.  
  *Example: `-p ./local.toml -p ~/credentials.yaml`*

- **`--silent`, ` -s`**  
  Disables all logging output. Only `print` and `println` statements will be printed. This is especially useful if you want to use Goatfiles inwith other scripts.

- **`--skip SKIP`**  
  Section(s) to be skipped during execution.  
  *Example: `--skip teardown`*

- **`--help`, ` -h`**  
  Display the help message.

- **`--version`**  
  Display the installed version.

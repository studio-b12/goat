# Command Line Tool

The `goat` command is used to create, validate and execute Goatfiles.

You can pass one or multiple Goatfiles or directories as positional argument. For example:
```
goat issue_32.goat tests/integrationtests
```

When passing in a directory, Goat will look for any `.goat` file recursively. Files and directories prefixed with an underscore (`_`) are ignored. This is especially useful for Goatfiles which are only supposed to be imported or executed in other Goatfiles. If you want to read more about this, take a look into the [Project Structure section](../project-structure.md). 

## Flags

The following sections provide further information about the various flags which can be passed to `goat`.

- `-a ARGS, --args ARGS` - Pass params as key value arguments into the execution
- `--delay DELAY, -d DELAY` - Delay requests by the given duration
- `--dry` - Only parse the goatfile(s) without executing any requests
- `--gradual, -g` - Advance the requests maually
- `--json` - Use JSON format instead of pretty console format for logging
- `--loglevel LOGLEVEL, -l LOGLEVEL` - Logging level
- `--new` - Create a new base Goatfile
- `--no-abort` - Do not abort batch execution on error
- `--no-color` - Supress colored log output
- `--params PARAMS, -p PARAMS` - Params file location(s)
- `--silent, -s` - Disables all logging output
- `--skip SKIP` - Section(s) to be skipped during execution
- `--help, -h` - display this help and exit
- `--version` - display version and exit

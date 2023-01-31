# Usage

## Parameters

```
❯ goat -h
Automation tool for executing and evaluating API requests.
goat v0.2.0 (ce33e93f 01/30/23 23:21 UTC go1.19.5)
Usage: goat [--args ARGS] [--delay DELAY] [--dry] [--gradual] [--loglevel LOGLEVEL] [--new] [--no-abort] [--params PARAMS] [--skip SKIP] [GOATFILE]

Positional arguments:
  GOATFILE               Goatfile(s) location

Options:
  --args ARGS, -a ARGS   Pass params as key value arguments into the execution (format: key=value)
  --delay DELAY, -d DELAY
                         Delay requests by the given duration
  --dry                  Only parse the goatfile(s) without executing any requests
  --gradual, -g          Advance the requests maually
  --loglevel LOGLEVEL, -l LOGLEVEL
                         Logging level (see https://github.com/rs/zerolog#leveled-logging for reference) [default: 1]
  --new                  Create a new base Goatfile
  --no-abort             Do not abort batch execution on error.
  --params PARAMS, -p PARAMS
                         Params file location
  --skip SKIP            Section(s) to be skipped during execution
  --help, -h             display this help and exit
  --version              display version and exit
```

## Recommended Setup

We are using and recommending the following setup for your Goatfiles.

```
integrationtests/
  tests/
    route-1/
      feature-1/
        _framework.goat
        tests.goat
      feature-2/
        tests.goat
  util/
    _login.goat
  params.yml
```

We are using a directory `integrationtets` at the top of our repositories. In the `util` directory, we put Goatfiles which have a general purpose and are used by a lot of test Goatfiles. For example login and logout procedures. Goatfiles which start with an underscore (`_`) are excluded when you execute goat on a directory. The `tests` directory contains all the actual test cases, split up into routes and features. The `_framework.goat` file contains the setup and teardown procedures for that specific test case and is used by `tests.goat`. Also, we put our `params.yml` in there which contains developer and instance specific parameters which need to be set for the test execution. This file is gitignored and will not be commited to the repository though.

With this setup, the test execution is as simple as using the following command to execute all test suits.
```
goat integrationtests/ -p integrationtests/params.yml
```
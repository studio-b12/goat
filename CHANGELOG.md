[VERSION]

# New Features

- **Added graceful shutdown [#69]**
  When canceling execution of goat, i.e. with <kbd>Ctrl</kbd>+<kbd>C</kbd>, goat will skip all
  subsequent test cases after finishing the currently running one and will still continue on executing teardown steps.
  This ensures integrity and consistency of your test setup, even when canceling the current execution.

# Minor Changes and Bug Fixes

- `--retry-failed` now uses absolute paths so that it can be executed in any directory.
- Fixed a bug that prevented executing Goatfiles via absolute paths.
- Fixed request formatting in the log output.

# ETC

- Goat now has an official logo *(like every cool and fancy open-source project has)*!
  ![](https://raw.githubusercontent.com/studio-b12/goat/refs/heads/main/.github/media/banner.png)
[VERSION]

# New Features

- Added raw data identifier `$` to directly use raw data from a response as request body. See the example in the
  documentation for more information. [#66, #67]

- Added a new flag `--retry-failed` (or `-R` for short). When a batch execution fails, failed files are now saved to a
  temporary file. After that, goat can be executed again with only the flag which will re-run only the failed files.

<!-- # Minor Changes and Bug Fixes -->

<!-- - Fixed an issue when parsing file descriptors on Windows systems. -->
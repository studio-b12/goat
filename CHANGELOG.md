[VERSION]

# New Features

- Profiles have been implemented for a more central control on commonly used parameters across projects.
  See the [documentation](https://studio-b12.github.io/goat/command-line-tool/profiles.html) for all details.

- A new flag `--reduced-errors` (or `-R` for short) has been added to hide template errors in teardown steps.
  This can be useful to reduce the noise of template errors due to missing variables in teardown steps when tests fail.

- Added a new builtin templating function [`formatTimestamp`](https://studio-b12.github.io/goat/templating/builtins.html#formattimestamp) to
  format dates. Also the [`timestamp`](https://studio-b12.github.io/goat/templating/builtins.html#timestamp) builtin
  has been updated to also accept format names like `formatTimestamp`.

# Minor Changes and Bug Fixes

- `response.BodyRaw` is now represented as a UTF-8 encoded string when printed in the `[Script]` section instead of listing
  the list of byte values. It is still an array of bytes though, so you can operate on it as expected.

- Fixed a bug where `noabort` and `alwaysabort` options had no effect on teardown steps.

- The `json` template builtin now accepts an additional parameter for indentation.
  
# Code Base

- The Goatfile parser does now produce an intermediate `AST` structure instead of the `Goatfile` directly. This should allow
  to build tooling around Goatfiles using the provided parser implementation more easily (i.e. like an auto-formatter or LSP). 
  Feel free to discover the new implementation [here](pkg/goatfile/parser.go). [#57]
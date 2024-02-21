[VERSION]

# New Features

- Profiles have been implemented for a more central control on commonly used parameters across projects.
  See the [documentation](https://studio-b12.github.io/goat/command-line-tool/profiles.html) for all details.


# Minor Changes and Bug Fixes

- `response.BodyRaw` is now represented as a UTF-8 encoded string when printed in the `[Script]` section instead of listing
  the list of byte values. It is still an array of bytes though, so you can operate on it as expected.


# Code Base

- The Goatfile parser does now produce an intermediate `AST` structure instead of the `Goatfile` directly. This should allow
  to build tooling around Goatfiles using the provided parser implementation more easily (i.e. like an auto-formatter or LSP). 
  Feel free to discover the new implementation [here](pkg/goatfile/parser.go).
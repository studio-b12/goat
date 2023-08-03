[VERSION]

<!-- # New Features -->

# Minor Changes and Bug Fixes

- Fixed a bug where an `execute`d files location in a `use`d file is looked for in the directory that `use`s the file and not the `use`d file. [#39]
- Fixed a bug in the `errs` package which, when using `Join` with an `error` type as target, overwrites the passed error.

[VERSION]

- Renamed the project to `goat`.
- Added argument `--arg` to set params in the initial execution state.
- Re-importing already imported files now results in an error.
- Fixed a bug where prefixes and suffixes of DetailedError were not shown.
- Parsing errors are now printed to the console so that it is a valid file link in editors like VSCode.
- Fixed position in parsing errors.
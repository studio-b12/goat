[VERSION]

- The `response` object now implements `String()`, which allows proper string printing in the `[Script]` part.
- The argument `--params` can now be passed multiple times to pass multiple parameter files which will then be merged into one initial state for the execution.
- Fixeded import path joining on Windows systems.
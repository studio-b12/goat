# State Management

Goat initializes a state for every Goatfile execution. A state initially consists of the parameters which were passed via a configuration file, environment variables, or the `--args` CLI parameter.

This state is then passed to every request in the executed Goatfile. The request can read and alter the state.
When one request has finished, the state is passed on to the next request and so on until the execution has finished.

Below, you can see a very simple example of what a state lifecycle could look like.

![](../assets/simple-state.excalidraw.svg)

## State with [`use`](../goatfile/import-statement.md)

Because [`use`](../goatfile/import-statement.md) effectively merges the imported Goatfile together with the root Goatfile to a single batch execution, the state is shared between them. So if you define a variable in a Goatfile imported via [`use`](../goatfile/import-statement.md), the variable will be accessible in subsequent imported Goatfiles as well as in the root Goatfile.

## State with [`execute`](../goatfile/execute-statement.md)

Invoking another Goatfile via the [`execute`](../goatfile/execute-statement.md) command will make Goat handle the executed Goatfile like a separate batch with its own, disconnected state.

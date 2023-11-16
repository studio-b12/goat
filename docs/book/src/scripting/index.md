# Scripting

Scripting sections like `[Script]` and `[PreScript]` use a dedicated scripting micro-engine for maximum flexibility in your test setup and procedures.

Goat uses ES5.1 conform JavaScript interpreted by the [goja](https://github.com/dop251/goja) micro-engine.

In each script instance, you have access to the current state values via the global environment variables. Also, you can define global variables using the `var` statement to define values which will be saved in the state after successful script execution.

Example:

> Initial State
> ```toml
> foo = 1
> ```

> Script Execution
> ```js
> var bar = foo * 2;
> ```

> Resulting State
> ```toml
> foo = 1
> bar = 2
> ```

Also, some [built-in functions](./builtins.md) are available in each script instance.

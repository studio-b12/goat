[VERSION]

# New Features

## Logical Section Logging [#28]

The beginning of the processing of a "logical" section (like `Setup`, `Test` and `Teardown`) is now represented in the log output in the same way as context sections are.

Therefore, you don't need to write work around like stuff like the following to visually separate between Goatfile sections.

```
### Tests
##### Tests
// ...
```

An example output would look like the following.
![](https://github.com/studio-b12/goat/assets/16734205/fab842d9-f49c-4ccb-834a-7ddc8e50a8c2)


## More detailed debugg logging [#29]

When debug logging is enabled (log level of `6` or higher), detailed information about the outgoing request as well as the incoming response is provided in the log. Therefore, you don't need to add a `debug(response)` to your `[Script]` field in your Goatfiles anymore.

# Bug Fixes

- In [v0.11.0](https://github.com/studio-b12/goat/releases/tag/v0.11.0), a feature has been added that the file and line of a request that failed is represented in the error message in the console ([#7]). There was a bug that faield reqeusts in imported files show up as being in the file which imports the file. This has been fixed now. [#30]
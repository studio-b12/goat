[VERSION]

# New Features

## External Goatfile Execution [#32]

It is now possible to `execute` a different Goatfile from within a Goatfile.

Here is a quick example.
```
execute ../doStuff (
    username="{{.username}}" 
    password="{{.password}}"
    n=1
) return (
    response as response1
)
```

Executed Goatfiles are ran in a completely separate state with the parameters specified passed in and values captured which are listed in the `return` statement.

More detials to the implementation can be found [**here**](https://github.com/studio-b12/goat/blob/master/docs/implementation.md#execute-directive).

# Minor Changes and Bug Fixes

- It is now possible to substitute file descriptors with parameters (e.g. `@../{{.fileName}}`).


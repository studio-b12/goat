[VERSION]

# New Features

- Added a new request block `[FormData]`, which enables to specify `multipart/form-data` request payloads. [Here](https://studio-b12.github.io/goat/goatfile/requests/formdata.html) you can read the documentation. [#55]

- Added a new request option [`followredirects`](https://studio-b12.github.io/goat/goatfile/requests/options.html#followredirects). This is `true` by default, but can be set to `false` if redirects should not be followed transparently. [#61]

- Added a new script builtin [`assert_eq`](https://studio-b12.github.io/goat/scripting/builtins.html#assert_eq), where you can pass two values which are compared and output for better error clarification.

- Added a new script builtin [`jq`](https://studio-b12.github.io/goat/scripting/builtins.html#jq), to perform JQ commands on any object in script blocks.

# Minor Changes and Bug Fixes

- Fixed an issue when parsing file descriptors on Windows systems.
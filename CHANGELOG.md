[VERSION]

# New Features

- Added a new request block `[FormData]`, which enables to specify `multipart/form-data` request payloads. [Here](https://studio-b12.github.io/goat/goatfile/requests/formdata.html) you can read the documentation. [#55]

- Added a new request option [`followredirects`](https://studio-b12.github.io/goat/goatfile/requests/options.html#followredirects). This is `true` by default, but can be set to `false` if redirects should not be followed transparently. [#61]

# Minor Changes and Bug Fixes

- Fixed an issue when parsing file descriptors on Windows systems.
[VERSION]

# New Features

- A new builtin template function has been added called `isset` which takes a map as first parameter and a map key string as second parameter and returns true if the key is set in the map.

# Minor Changes and Bug Fixes

- Request conditions are now substituted and checked before continuing the request substitution (see #41).

- The sections `Setup-Each` and `Teardown-Each` have been removed (see #37).

- The `PreScript` block will now – as expected – also get substituted as same as the `Script` section (see #35).

- When an array or map value gets substituted, the representation will now be in JSON format (see #36).

- The execution status summary at the end of the end of the program execution does no more count log actions as requests and properly summs up request summaries from executed Goatfiles.
[VERSION]

# New Features

- A new template function `json` has been added which encodes any passed object into a JSON string (e.g. `{{ json .myObject }}`).
- 
- For each in-script logging function, a formatted equivalent has been added (e.g. `infof("Hello, %s!", "world")`).

# Minor Changes and Bug Fixes

- Variables set in the `PreScript` section are now available for substitution in all other parts of the request like URL, options, header, body or script. [#43]

- Section log lines are now hidden when a Goatfile is executed via the `execute` instruction.

- Some more debugging information has been added visible when using the `trace` log level.
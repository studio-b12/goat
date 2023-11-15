# Goatfile

A Goatfile contains the request definitions for a Goat test suite. Goatfiles are written in plain 
UTF-8 encoded text using the Goatfile syntax.

Goatfiles have the file extension `.goat`.

Below, you can see a really simple example Goatfile.

```
### Tests

GET https://api.github.com/repos/{{.repo}}

[Script]
assert(response.StatusCode === 200, `Invalid response code: ${response.StatusCode}`);
assert(response.BodyJson.language === "Go");

---

GET https://api.github.com/repos/{{.repo}}/languages

[Script]
info('Languages:\n' + JSON.stringify(response.BodyJson, null, 2));
assert(response.StatusCode === 200, `Invalid response code: ${response.StatusCode}`);
```

More conclusive examples can be found [here](https://github.com/studio-b12/goat/tree/main/examples).

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
assert(response.Body.language === "Go");

---

GET https://api.github.com/repos/{{.repo}}/languages

[Script]
info('Languages:\n' + JSON.stringify(response.Body, null, 2));
assert(response.StatusCode === 200, `Invalid response code: ${response.StatusCode}`);
```

More conclusive examples can be found [here](https://github.com/studio-b12/goat/tree/main/examples).

## Specification

Below, you can see a simple synopsis of the different building blocks of a Goatfile.

```
         Import  |  use ./setup.goat
                 |
        Comment  |  // This is a comment!
                 |
Section Heading  |  ### Setup
                 |
Context Section  |  ##### Upload Tests
                 |
   Method & URL  |  POST https://example.com
                 |  
        Headers  |  [Header]
                 |  X-Requested-With: XMLHttpRequest
                 |  Content-Type: application/json
                 |  Hash: {{ sha256 .data }}
                 |
           Body  |  [Body]
                 |  ```
                 |  {
                 |    "hello": "world",
                 |    "hash": "{{ sha256 .data }}",
                 |    "dontescapethis": "\{\{ \}\}"
                 |  }
                 |  ```
                 |
  Option Blocks  |  [QueryParams]
                 |  page = 2
                 |  items = 30
                 |  filters = ["time", "name"]
                 |
                 |  [MultipartFormdata]
                 |  image = @myimage.jpeg
                 |
                 |  [Options]
                 |  cookiejar = "foo"
                 |  storecookies = false
                 |
         Script  |  [Script]
                 |  assert(response.status >= 200 && response.status < 400);
                 |  assert(response.body["name"] == "somename");
                 |
                 |  // capture a variable to be used in subsequent responses
                 |  var id = response.body["id"];
                 |
        Request  |  ---
      Separator  |
                 |
        Execute  |  execute ./testFileUpload (
                 |    file="file_a.txt"
                 |    token="{{.auth.token}}"
                 |  ) return (
                 |    fileId as fileId_a
                 |  )
```

In the following sections, you will find a detailed rundown for each component of a Goatfile.

- [Comment](./comments.md)
- [Import Statement](./import-statement.md)
- [Execute Statement](./execute-statement.md)
- [Section](./sections.md)
- [Defaults Section](./defaults-section.md)
- [Log Section](./logsections.md)
- [Request](./requests/index.md)
  - [Method and URL](./requests/method-and-url.md)
  - [Options](./requests/options.md)
  - [Headers](./requests/header.md)
  - [Query Parameters](./requests/query-params.md)
  - [Body](./requests/body.md)
  - [PreScript](./requests/prescript.md)
  - [Script](./requests/script.md)

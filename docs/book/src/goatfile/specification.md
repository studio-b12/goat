# Syntax Specification

Goatfiles are built from the following blocks.

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

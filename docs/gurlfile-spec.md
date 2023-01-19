# Gurlfile Specification

A Gurlfile is an UTF-8 encoded plain text file with the file extension `.gurl`.

```
         Import  |  use ./setup.gurl
                 |
                 |
Section Heading  |  ### Setup
                 |
   Method & URL  |  POST https://example.com
                 |  
        Headers  |  [Headers]
                 |  X-Requested-With: XMLHttpRequest
                 |  Content-Type: application/json
                 |  Hash: {{ sha256 .data }}
                 |
           Body  |  [Body]
                 |  ```
                 |  {
                 |	  "hello": "world",
                 |    "hash": "{{ sha256 .data }}",
                 |    "dontescapethis": "\{\{ \}\}"
                 |  }
                 |  ```
                 |  
        Comment  |  // This is a comment!
                 |
  Option Blocks  |  [QueryParams]
                 |  page = 2
                 |  items = 30
                 |  filters = ["time", "name"]
                 |
                 |  [MultipartFormdata]
                 |  image = @myimage.jpeg
                 |
                 |  [Config]
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
```

## Structure

All requests in all sections of one or more Gurlfiles is called a **batch**.

### Sections

A Gurlfile consists of sections containing one or more requests. Each section has a specific name and function. Though, specifying sections is optional. Defaultly, when no sections are specified, all requests will be assigned to the section `Tests`.

- `Setup`: Requests which will be executed once before all tests are executed. If a request fails, the batch execution is aborted.
- `Setup-Each`: Requests which will be executed before every single request in `Tests`. If a request fails, the batch execution is aborted.
- `Tests`: The actual test requests which will be executed.
- `Teardown`: Requests which will always be executed at the end of the execution of a batch, even if the execution is aborted midways.
- `Teardown-Each`: Requests which will be executed after each single request in `Tests`. Will also be executed if a `Setup-Each` request fails.

When two Gurlfiles are merged (for example on importing one into another), all requests in all sections of one file are appended to the requests in the sections of the other file. The specific order of the sections in the Gurlfile is irrelevant.

### Requests

A request consists of the following parts which must be placed in the right order.

First, the Method and URL is specified separated by one or more spaces or tabs. These fields are mandatory.

*Example:*
```
GET https://example.com
```

In the next lines, request headers and the request body can be defined. There must not be a new line between these.
Request headers are defined as key-value pairs separated by a colon (`:`). The key must not contain any non-word characters
besides numbers, dashes and underscores.

The body can be any type of data which must not contain a newline.

*Example:*
```
Content-Type: application/json
X-Requested-With: XHTMLRequest
{
    "some": "data"
}
```

Followed by one or more newlines, you can specify additional request options. See [options](#options) for more information.
Options are defined as [TOML](https://toml.io/en/) blocks.

*Example:*
```
[QueryParams]
page = 1
count = 30
sortBy = "date"
```

Finally, followed by another one or more newlines, a script section can be specified.
There, assertions and output logic can be specified depending on the implemented scripting language.

Multiple requests are separated by a line containing at least 3 dashes (`-`).

*Example:*
```
// ... Request 1 ...

---

// ... Request 2 ...
```

### Parameters

Parameters can be passed into a request via [Go Templates](https://pkg.go.dev/text/template) placeholders.

A placeholder consists of two opening and closing curly braces enclosing the variable name. Cascading structures are separated by dots. A variable is also always prefixed with a dot.

Lets take the following data structure as an example.
```json
{
    "instance": "https://api.example.com",
    "credentials": {
        "username": "root",
        "password": "rHiZHVs5"
    }
}
```

Then, the defined parameters are referenced as following in the request.

```
POST {{.instance}}/api/auth/login
{
    "username": "{{.credentials.username}}",
    "password": "{{.credentials.password}}"
}
```

On execution of the request, the parameters are injected into the request.

You can also reference captured variables from previous requests.

For example:
```
POST {{.instance}}/api/auth/token
{
    "client_id": "{{.credentials.client_id}}",
    "client_secret": "{{.credentials.client_secret}}"
}

assert(response.StatusCode == 200);

var bearerToken = response.BodyJson.bearer_token;

---

GET {{.instance}}/api/me
Authorization: bearer {{.bearerToken}}
```
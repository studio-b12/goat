# Goatfile Specification

A Goatfile is an UTF-8 encoded plain text file with the file extension `.goat`.

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

All requests in all sections of one or more Goatfiles is called a **batch**.

### Sections

A Goatfile consists of sections containing one or more requests. Each section has a specific name and function. Though, specifying sections is optional. Defaultly, when no sections are specified, all requests will be assigned to the section `Tests`.

- `Setup`: Requests which will be executed once before all tests are executed. If a request fails, the batch execution is aborted.
- `Setup-Each`: Requests which will be executed before every single request in `Tests`. If a request fails, the batch execution is aborted.
- `Tests`: The actual test requests which will be executed.
- `Teardown`: Requests which will always be executed at the end of the execution of a batch, even if the execution is aborted midways.
- `Teardown-Each`: Requests which will be executed after each single request in `Tests`. Will also be executed if a `Setup-Each` request fails.

When two Goatfiles are merged (for example on importing one into another), all requests in all sections of one file are appended to the requests in the sections of the other file. The specific order of the sections in the Goatfile is irrelevant.

A special section also exist called `Defaults`, which can contain different optional blocks. The values in these blocks will be used as default values for all other requests in this batch. For single-value fields like `[Body]`, `[PreBody]` or `[Script]`, the specified value in a request overwrites the default value. Otherwise, the value specified in the defaults will be applied to the request. For multi-value fields like `[Options]`, `[Header]` or `[QueryParams]`, the values of these will be merged where existing values overwrite default values.

### Context Sections

These sections do not alter the execution of the specified requests in a batch and are just there to be used to provide additional context during the execution. For example, these sections can be used to be logged to visualy separate different sections in a batch.

### Requests

A request consists of the following parts.

First, the Method and URL is specified separated by one or more spaces or tabs. These fields are mandatory.

*Example:*
```
GET https://example.com
```

After that, further optional request details are provided in **blocks**. Every block starts with a **block header** followed by the specific block content.


The following blocks are required to be implemented by the specification.

#### `Header`

Define a list of request headers sent with the request.

The values of this block consists of one header per line where the key and value are separated by a colon.

*Example:*
```
[Header]
Accept: */*
Content-Type: application/json
Cookie: token=89034567n08924t
```

#### `Body`

Define a body payload which is sent with the request.

Everything following under this block header is considered to be the content until a new block, request or section is found. Optionally, you can escape the content using three backticks at the start and end of the content.

*Example (unescaped):*
```
[Body]
{
    "username": "root",
    "password": "foo",
}
```

*Example (escaped):*
````
[Body]
```
{
    "username": "root",
    "password": "foo",
}
```
````

It is also possible to import a file as request body. Simply specify the path to the file with a leading `@` to import the file by a relative path from the current Goatfile or an absolute path. The past must be specified in Unix style. Template placeholders in imported body data **will not** be infused on execution!

*Example:*
```
[Body]
@../data/image.png
```

If the path contains spaces, you can wrap the path in quotes.

*Example:*
```
[Body]
@"../body data/image.png"
```


#### `QueryParams`

Allows to specify query parameters in a more human readable way.

Values are defined in a [TOML](https://toml.io/en/)-like key-value representation.

*Example:*
```
[QueryParams]
page = 2
count = 10
orderBy = "date"
filters = ["date", "name"]
```

> For all available blocks in this imlementation, see [implementation.md](implementation.md).

Multiple subsequent requests in a section are separated by a delimiter consisting of at least three dashes (`---`). A request does not need to be terminalized with a delimiter when it is at the end of a batch or at the end of a section. 

*Example:*
```
### Tests

GET https://example1.com

---

GET https://example2.com

### Teardown

POST https://example3.com
```

#### `Script`

Defines an esecutable script which is executed after the request has been performed.

Equal to the `Body` block, everything after the block header is considered to be content until a new block, a delimiter or a new section is detected. The content can also be escaped.

The script parts is getting passed builtin functions as well as the response context. The representation of the script language used and the passed in data is specific to the implementation.

> To find more information about the details of this implementation, see [implementation.md](implementation.md).

*Example (using ECMAScript 5):*
```
[Script]
assert(response.StatusCode == 200);
info(response.Body);
var uid = response.BodyJson.Uid;
```

It is also possible to import a file as script. Simply specify the path to the file with a leading `@` to import the file by a relative path from the current Goatfile or an absolute path. The past must be specified in Unix style. Template placeholders in imported body data **will** be infused on execution!

*Example:*
```
[Script]
@../scripts/check.js
```

If the path contains spaces, you can wrap the path in quotes.

*Example:*
```
[Body]
@"../additional scripts/check.js"
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

[Body]
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

[Body]
{
    "client_id": "{{.credentials.client_id}}",
    "client_secret": "{{.credentials.client_secret}}"
}

[Script]
assert(response.StatusCode == 200);

var bearerToken = response.BodyJson.bearer_token;

---

GET {{.instance}}/api/me

[Header]
Authorization: bearer {{.bearerToken}}
```

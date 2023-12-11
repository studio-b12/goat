# Implementation Details

## Additional Available Request Blocks

### `Options`

Additional options that control the request and response behavior.

Available parameters are:

Name | Type | Default | Description
-|-|-|-
`cookiejar` | `string` or `number` | `"default"` | The cookie jar to be used.
`storecookies` | `boolean` | `true` | Whether or not to store response cookies to the cookiejar.
`sendcookies` | `boolean` | `true` | Whether or not to send stored cookies with the request.
`noabort` | `boolean` | `false` | When set to true, the batch execution will not be arborted when the request fails.
`alwaysabort` | `boolean` | `false` | Forces the batch to abort when the request fails, even when the `--no-abort` flag is set.
`condition` | `boolean` | `true` | Whether or not to execute the request.
`delay` | `string` | `0` | A duration which is awaited before the request is executed. The duration must be formatted in compatibility to Go's [ParseDuration](https://pkg.go.dev/time#ParseDuration) function.

### `PreScript`

The "PreScript" is executed before the actual request is infused with parameters and executed afterwards. This allows setting request parameters via scripting before the execution of a request.

Example:
```
GET https://echo.zekro.de/{{.path}}

[PreScript]
var path = "somepath";
var body = JSON.stringify({"foo": "bar"});

[Body]
{{.body}}

[Script]
assert(response.StatusCode === 200);
assert(response.Body.path === "/somepath");
assert(response.Body.body_string === '{"foo":"bar"}\n');
```

## Script Implementation

This implementation is using ECMAScript 5. Inm detail, it is using the [goja](https://github.com/dop251/goja) engine to perform the script part.

### Builtin Functions

The following functions are available.

#### `assert(condition: boolean, message?: string)`

When the passed condition is `false`, the function will throw an error. Optionally, you can specify a custom message with the `message` parameter. When the thrown error is not catched, the batch execution will abort.

#### `print(message: string)`

Prints a plain `message` to the console output (without a leading newline).

#### `println(message: string)`

Same as `print` but with a leading newline attached.

#### `info(message: string)`

Print an `info` log entry with the passed `message`.

#### `warn(message: string)`

Print a `warn` log entry with the passed `message`.

#### `error(message: string)`

Print an `error` log entry with the passed `message`.

#### `fatal(message: string)`

Print a `fatal` log entry with the passed `message`. This will abort the batch execution.

#### `debug(message: string)`

Print a `debug` log entry with the passed `message`.

### Global Context Variables

Some variables are available in the global context of the execution.

#### `response`

Contains information about the previous response. The Response contains the following fields.

```go
Response {
	StatusCode    int
	Status        string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Header        map[string][]string
	ContentLength int64
	BodyRaw       []byte
	Body		  any
}
```
`Body` is a special field containing the response body content as a JavaScript object which will be populated if the response body can be parsed.
Parsers are currently implemented for `json` and `xml` and are chosen depending on the `responsetype` option or the `Content-Type` header.
If neither are set, the raw response string gets set as `Body`. By setting the `responsetype` to `raw`, implicit body parsing can be prevented.

## Parameters

This implementation is using Go's template implementation for parsing parameters.

[Here](https://pkg.go.dev/text/template) you can find the documentation how to use Go templates.

As you can read in these Docs, you can also perform functions on parameters. As an example:
```go
{{ urlquery "Hello world!" }}
// results in "Hello+world%21"
```

Goat also provides some additional functions which are available in Goatfiles.

- `base64 <string>`: Returns the input strings as base64 encoded string with padding.
- `base64Url <string>`: Returns the input strings as base64 URL encoded string with padding.
- `base64Unpadded <string>`: Returns the input strings as base64 encoded string without padding.
- `base64UrlUnpadded <string>`: Returns the input strings as base64 URL encoded string without padding.
- `md5 <string>`: Returns the HEX encoded MD5 hash of the input string.
- `sha1 <string>`: Returns the HEX encoded SHA1 hash of the input string.
- `sha256 <string>`: Returns the HEX encoded SHA256 hash of the input string.
- `sha512 <string>`: Returns the HEX encoded SHA512 hash of the input string.
- `randomString <integer?>`: Returns a random string with the given length. If no length is passed, the default length is 8 characters.
- `randomInt <integer?>`: Returns a random integer in the range [0, n) where n is given as parameter. If no parameter is passed, n is the max int value.
- `timestamp <format?>`: Returns the current timestamp in the given format. If no format is specified, the timestamp is returned as Unix Seconds string.
- `isset <map[string]any> <string>`: Returns `true` if the given key string is defined and not nil in the given map. Otherwise, `false` is returned.

## `use` Directive

### Syntax

> *UseExpression* :  
> `use` *StringLiteral*

### Example

```
use ../path/to/goatfile
```

### Implementation

The `use` directive allows to "import" a Goatfile *A* into another Goatfile *B*. This behaves like all actions of *A* and *B* were merged together to a single Goatfile where all actions of *B* will be inserted before all actions of *A*. That means that Setups and Teardowns of the imported Goatfile *B* will be executed together with the Setup and Teardown Sections of the base Goatfile *A*.

Imported Goatfiles will be statically checked as same as the executed Goatfile. Circular imports are not allowed.

Because the Goatfile with imported Goatfiles behaves like a single Goatfile, both files share the same execution state and parameters.

Schematical Example:
> `A.goat`
```
use B

### Setup

GET A1

### Tests

GET A2

### Teardown

GET A3
```

> `B.goat`
```
### Setup

GET B1

### Tests

GET B2

### Teardown

GET B3
```

> Result (virtual):
```
### Setup

GET B1
GET A1

### Tests

GET B2
GET A2

### Teardown

GET B3
GET A3
```

## `execute` Directive

### Syntax

> *ExecuteExpression* :  
> `execute` *StringLiteral* *ExecuteSignature?*
>
> *ExecuteSignature* :  
> `(` *ParameterAssignment** `)` *ReturnsValueSignature*?
>
> *ParameterAssignment* :  
> `LF`* *Literal* `=` *Value* `LF`*
>
> *ReturnsValueSignature* :  
> `return` `(` *ReturnValueAssignment** `)`
>
> *ReturnValueAssignment* :  
> `LF`* *Literal* `as` *Literal* `LF`*

### Example
```
execute ../doStuff (
    username="{{.username}}" 
    password="{{.password}}"
    n=1
) return (
    response as response1
)
```

### Implementation

The `execute` directive allows to execute a foreign Goatfile *A* inside another Goatfile *B* with the ability to pass specific parameters and capture return values.

In contrast to the `use` directive, the executed Goatfile *A* is ran like a completely separate Goatfile execution with its own isolated state which does not share any values with the state of file *B*. All parameters which shall be available in *A* must be passed as a list of key-value pairs. Resulting state values of *A* can then be captured to the state of *B* by listing them in the `return` statement with the name of the parameter in the state of *A* and the name the value shall be acessible in *B*.

Goatfiles executed are parsed in place, so they are only statically checked once they are executed within the executing Goatfile.
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
	Body          string
	BodyJson      object
}
```

When the response body is JSON-parsable, `BodyJson` contains the parsed JSON as a JavaScript object. Otherwise, it will be `null`.


## Parameters

This implementation is using Go's template implementation for parsing parameters.

[Here](https://pkg.go.dev/text/template) you can find the documentation how to use Go templates.

As you can read in these Docs, you can also perform functions on parameters. As an example:
```go
{{ urlquery "Hello world!" }}
// results in "Hello+world%21"
```

Goat also provides some additional functions which are available in Goatfiles.

- `base64 <string>`: Returns the input strings as base64 encoded string.
- `base64url <string>`: Returns the input strings as base64 URL encoded string.
- `md5 <string>`: Returns the HEX encoded MD5 hash of the input string.
- `sha1 <string>`: Returns the HEX encoded SHA1 hash of the input string.
- `sha256 <string>`: Returns the HEX encoded SHA256 hash of the input string.
- `sha512 <string>`: Returns the HEX encoded SHA512 hash of the input string.
- `randomString <integer?>`: Returns a random string with the given length. If no length is passed, the default length is 8 characters.
- `randomInt <integer?>`: Returns a random integer in the range [0, n) where n is given as parameter. If no parameter is passed, n is the max int value.
- `timestamp <format?>`: Returns the current timestamp in the given format. If no format is specified, the timestamp is returned as Unix Seconds string.
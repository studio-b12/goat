# Implementation Details

## Additional Available Request Blocks

### `Options`

> **Warning**  
> This is not implemented yet.

Additional options that control the request and response behavior.

Available parameters are:

Name | Type | Default | Description
-|-|-|-
`cookiejar` | `string` or `number` | `"default"` | The cookie jar to be used.
`storecookies` | `boolean` | `true` | Whether or not to store response cookies to the cookiejar.
`sendcookies` | `boolesn` | `true` | Whether or not to send stored cookies with the request.

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
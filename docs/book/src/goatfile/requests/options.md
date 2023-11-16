# Options

> *RequestOptions* :  
> `[Options]` `NL`+ *RequestOptionsContent*
>
> *RequestOptionsContent* :  
> *TomlKeyValues*

## Example

```toml
[Options]
cookiejar = "admin"
condition = {{ isset . "userId" }}
delay = "5s"
```

## Explanation

Define additional options to the requests. The format of the contents of this block is [`TOML`](https://toml.io/).

The option values support template substitution.

Below, all available options are explained.

### `cookiejar`

- **Type**: `string` | `number`
- **Default**: `"default"`

Defines the cookie jar to be used for saving and storing cookies. A cookie jar can be specified by either a number or a string. Every cookie jar contains a separate set of cookies collected from requests performed with that cookie jar specified.

### `storecookies`

- **Type**: `boolean`
- **Default**: `true`

Defines if a cookie set by a request's response shall be stored in the cookie jar.

### `sendcookies`

- **Type**: `boolean`
- **Default**: `true`

Defines if cookies stored in the specified cookie jar shall be sent to the server on request.

### `noabort`

- **Type**: `boolean`
- **Default**: `false`

When enabled, a batch execution will not be canceled when the request execution or assertion failed.

### `alwaysabort`

- **Type**: `boolean`
- **Default**: `false`

Forces a batch request to abort if the request execution or assertion failed, even if the `--no-abort` CLI flag has been passed.

### `condition`

- **Type**: `boolean`
- **Default**: `true`

Defines if a request shall be executed or not. This is useful in combination with template substitution.

> For example, the following request will only be executed when `localAddress` is set in the current state.
> ```
> [Options]
> condition = {{ isset . "localAddress" }}
> ```

### `delay`

- **Type**: `string`
- **Default**: `0`

A duration formatted as a Go [time.ParseDuration](https://pkg.go.dev/time#ParseDuration) compatible string. Execution will pause for this duration before the request is executed.

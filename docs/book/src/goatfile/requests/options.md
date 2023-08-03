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
condition = {{ if .userId }}
delay = "5s"
```

## Explanation

Define additional options to the requests. The format of the contents of this block is [`TOML`](https://toml.io/).

Below, all available otpions are explained.

### `cookiejar`

- **Type**: `string` | `number`
- **Default**: `"default"`

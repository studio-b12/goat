# FormUrlEncoded

> *FormUrlEncoded* :  
> `[FormUrlEncoded]` `NL`+ *FormUrlEncodedContent*
>
> *FormUrlEncodedContent* :  
> *TomlKeyValues*

## Example

```toml
[FormUrlEncoded]
someString = "some string"
someInt = 42
someBool = true
```

## Explanation

Defines entries in a key-value pair format which will be sent in the request body as `application/x-www-form-urlencoded` request.
The format of the contents of this block is [`TOML`](https://toml.io/).

> The example from above results in the following body content.
> ```
> someString=some+string&someInt=42&someBool=true
> ```

Template parameters in parameter values will be substituted.

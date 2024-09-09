# FormData

> *FormData* :  
> `[FormData]` `NL`+ *FormDataContent*
>
> *FormDataContent* :  
> *TomlKeyValues*

## Example

```toml
[FormData]
someString = "some string"
someInt = 42
someFile = @files/goat.png:image/png
someRaw = $someVar:application/zip
```

## Explanation

Defines entries in a key-value pair format which will be sent in the request body as `multipart/form-data` request.
The format of the contents of this block is [`TOML`](https://toml.io/).

The file's content type can be specified after the file descriptor, separated by a colon (`:`). Otherwise,
the content type will default to `application/octet-stream`.

The `$` sign can be used in `FormData` values to directly pass byte array variables. This can e.g. be used when a
previous `response.BodyRaw` should be send to another endpoint as a `FormData` request.

> The example from above results in the following body content.
> ```
> --e8b9253313450dbcf0d09df1a0f3ff36dd00e888339415a239ce167f279c
> Content-Disposition: form-data; name="someInt"
>
> 42
> --e8b9253313450dbcf0d09df1a0f3ff36dd00e888339415a239ce167f279c
> Content-Disposition: form-data; name="someFile"; filename="goat.png"
> Content-Type: image/png
>
> <binary file content>
> --e8b9253313450dbcf0d09df1a0f3ff36dd00e888339415a239ce167f279c
> Content-Disposition: form-data; name="someString"
> 
> some string
> --e8b9253313450dbcf0d09df1a0f3ff36dd00e888339415a239ce167f279c--
> Content-Disposition: form-data; name="someRaw" ; filename="binary-data"
> Content-Type: application/zip
> 
> <raw byte data>
> --e8b9253313450dbcf0d09df1a0f3ff36dd00e888339415a239ce167f279c--
> ```

Template parameters in parameter values will be substituted.

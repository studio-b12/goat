# Auth

> *RequestOptions* :  
> `[Auth]` `NL`+ *RequestAuthContent*
>
> *RequestAuthContent* :  
> *TomlKeyValues*

## Example

```toml
[Auth]
username = "foo"
password = "{{.credentials.password}}"
```

## Explanation

Auth is a utility block for easily defining basic or token authorization. When defined, the `Authorization` header will 
be set accordingly to the request.

## Basic Auth

If you want to add basic auth to your request, simply define a `username` and `password`. If **both** values are set,
the username and password are joined by a `:` and encoded using end-padded base64. The encoded value is then set to the
`Authorization` header using the `basic` auth type (as defined in 
[RFC 7671](https://www.rfc-editor.org/rfc/rfc7617)).

> **Example**
> ```toml
> [Auth]
> username = "foo"
> password = "bar"
> ```
> his input will result in the following header.
> ```
> Authorization: basic Zm9vOmJhcg==
> ```

## Token Auth

You can also specify a `token` set as `Authorization` header. If defined, the token will be prefixed with
a token `type`.

> **Example**
> ```toml
> [Auth]
> type = "bearer"
> token = "foobarbaz"
> ```
> This input will result in the following header.
> ```
> Authorization: bearer foobarbaz
> ```
> 
> ```toml
> [Auth]
> token = "foobarbaz"
> ```
> This input will result in the following header.
> ```
> Authorization: foobarbaz
> ```

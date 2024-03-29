# Header

> *RequestOptions* :  
> `[Header]` `NL`+ *RequestHeaderContent*
>
> *RequestHeaderContent* :  
> *HeaderKV**
>
> *HeaderKV* :  
> `/[A-Za-z\-]+/` `:` `WS`* `/.*/` `NL`

## Example

```toml
[Header]
Content-Type: application/json
Accept: application/json
X-Some-Custom-Header: foo bar bazz
```

## Explanation

Define HTTP headers sent with the request in [HTTP conform header representation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages#headers).

Template parameters in the headers' value fields will be substituted.

# Method and URL

## Syntax

> *RequestHeader* :  
> *StringLiteral* `WS`+ *StringLiteral*

## Example

```
GET https://example.com/api/users
```

## Explanation

The request header defines the method and URL for a request and is the only mandatory element
to define a request.

The method can be any uppercase string.

The URL can either be defined as an unquoted string literal or as a quoted string if spaces are required in the URL. Template substitution is supported.

✅ Valid
```
GET https://example.com/api/users
```

✅ Valid
```
GET https://example.com/api/users/{{.userId}}
```

✅ Valid
```
GET "https://example.com/api/users/some user"
```

✅ Valid
```
GET "https://example.com/api/users/{{ .userId }}"
```

❌ Invalid
```
GET https://example.com/api/users/{{ .userId }}
```

❌ Invalid
```
GET https://example.com/api/users/some user
```

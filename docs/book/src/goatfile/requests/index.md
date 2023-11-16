# Request

## Syntax

> *RequestDefinition* :  
> *RequestHeader* `NL` (*RequestBlock* `NL`)* *RequestDelimiter* `NL`
> 
> *RequestHeader* :  
> *StringLiteral* `WS`+ *StringLiteral*
>
> *RequestBlock* :  
> *RequestBlockHeader* `NL` *RequestBlockContent*
>
> *RequestDelimiter* :  
> `---` `-`*

## Example

```
GET https://example.com/api/users

[Header]
Accept: application/json
Authentication: basic {{ .credentials.token }}

[Script]
assert(response.StatusCode === 200);

---
```

## Explanation

Define a request to be executed. The request definition starts with the request header consisting of the method followed by the URI (separated by one or more spaces). The request header is the only mandatory field for a valid request definition.

After that, you can specify more details about the request in different blocks. In the following documentation sections, all available blocks are listed and explained.

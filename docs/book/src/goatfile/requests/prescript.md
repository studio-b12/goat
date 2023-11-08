# PreScript

> *RequestOptions* :  
> `[PreScript]` `NL`+ *RequestPreScriptContent*
>
> *RequestPreScriptContent* :  
> *BlockDelimitedContent* | *UndelimitedContent*
>
> *BlockDelimitedContent* :  
> *BlockDelimiter* `NL` `/.*/` `NL` *BlockDelimiter*
>
> *UndelimitedContent* :  
> (`/.*/` `NL`)* `NL`
>
> *BlockDelimiter* :  
> `` ``` ``

## Example

````toml
[PreScript]
```
var fileName = requestedFile.Metadata.Name;
```
````

## Explanation

A script section which is evaluated before the actual request is substituted and performed. This can be used to put values from previous responses into variables which can then be later on used in the request parameters.

Example:
```
POST http://example.com/api/user

[Body]
// ...

[Script]
assert(response.StatusCode === 201);
var user = response.BodyJson;

---

GET http://example.com/api/user/{{.userid}}

[PreScript]
var userid = user.id;
```

`PreScript` will always be executed before the request parameters are substituted. So you can use the results in various fields like `[Options]` or `[Body]`, `Header` or `[Script]`.

Scripts are written in ES5.1 conform JavaScript. More on that can be found in the [Script](./script.md) section documentation.

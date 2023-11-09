# Script

> *RequestOptions* :  
> `[Script]` `NL`+ *RequestPreScriptContent*
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
[Script]
```
assert(response.StatusCode === 200, `Response status code was ${response.StatusCode}`);
assert(response.BodyJson.UserName === "Foo Bar");
```
````

## Explanation

A script section which is executed after a request has been performed and responded to. This is generally used to assert response values like status codes, header values or body content.

Scripts are written in ES5.1 conform JavaScript.

The context of the script always contains the current values in the batch state as global variables.

A special variable set in each `[Script]` section is the `response` variable, which contains all information about the request response. The `Response` object contains the following fields.

```go
type Response struct {
	StatusCode    int
	Status        string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Header        map[string][]string
	ContentLength int64
	Body          string
	BodyJson      any
}
```

`BodyJson` is a special field containing the response body content as a JavaScript object which will be populated when the response body can be parsed as a JSON object.

In any script section, you can used built-in functions like `assert`, which are documented [here](../../scripting/builtins.md).

When a script section throws an un-catched exception, the test will be evaluated as *failed*.

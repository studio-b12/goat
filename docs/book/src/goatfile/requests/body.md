# Body

> *RequestOptions* :  
> `[Body]` `NL`+ *RequestBodyContent*
>
> *RequestBodyContent* :  
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
[Body]
```
{
  "user": {
    "name": "{{.userName}}",
    "favorite_programming_langs": [
      "Go",
      "Rust"
    ],
  }
}
```
````

## Explanation

Define HTTP headers sent with the request in the [HTTP conform header representation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages#headers).

If you want to use the parameter braces (`{{`, `}}`) without substitution, you can escape them using a backslash.

Example:
````toml
[Body]
```
{
  "user": {
    "name": "\{\{This will not be substituted\}\}",
    "favorite_programming_langs": [
      "Go",
      "Rust"
    ],
  }
}
```
````

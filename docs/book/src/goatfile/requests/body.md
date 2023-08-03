# Header

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

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

Define the data to be transmitted with the request.

If you want to use the template parameter braces (`{{`, `}}`) without substitution, you can escape them using a backslash.

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

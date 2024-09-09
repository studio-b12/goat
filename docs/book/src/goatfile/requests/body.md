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

## Examples

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

````toml
[Body]
@path/to/some/file
````

````toml
[Body]
$someVar
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

### File Descriptor
With the `@` prefix in front of a file path, it is possible to send a file as a request body.

Example:
````toml
[Body]
@path/to/some/file
````

### Raw Descriptor
With the `$` prefix in front of a variable name it is possible to pass raw byte arrays to the request body. 
This can e.g. be used to send a `response.BodyRaw` to another endpoint.

Example:
````toml
[Body]
$someVar
````

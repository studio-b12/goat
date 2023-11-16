# QueryParams

> *QueryParams* :  
> `[QueryParams]` `NL`+ *QueryParamsContent*
>
> *QueryParamsContent* :  
> *TomlKeyValues*

## Example

```toml
[QueryParams]
page = {{.page}}
count = 100
field = ["username", "age", "id"]
token = "{{.apitoken}}"
```

## Explanation

Define additional query parameters which will be appended to the request URL. The format of the contents of this block is [`TOML`](https://toml.io/).

An array of values will be represented as a repetition of the same query parameter with the different contained values assigned.

> The example from above results in the following query parameters.
> ```
> page=5&count=100&field=username&field=age&field=id
> ```

Template parameters in parameter values will be substituted.

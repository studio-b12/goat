# Defaults Section

## Syntax

> *DefaultsBlock* :  
> `###` `defaults` `NL`+ *PartialRequestProperties*?

## Example

```
### Defaults

[Options]
storecookies = false

[Header]
Authorization: basic {{.token}}

[Script]
assert(response.StatusCode === 200);
```

## Explanation

The `Defaults` section body is structured like a request but without the method and URL section. This includes the blocks `[Options]`, `[Header]`, `[QueryParams]`, `[Body]`, `[PreScript]` and `[Script]`. 

Values specified in `[Header]`, `[Options]` and `[QueryParams]` will be merged with the values set in the requests. Values set in the request will overwrite values set in the defaults, if specified in both.

On the other hand, values specified in `[Body]`, `[PreScript]` or `[Script]` will be used if not specified in the requests. If these blocks are specified in the request, they will overwrite the default values as well (even if they are empty).

Default values will also be applied if imported via a `use` directive. Multiple `Defaults` sections will be merged together in order of specification.

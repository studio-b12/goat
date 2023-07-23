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


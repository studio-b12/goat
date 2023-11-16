# Section

## Syntax

> *SectionDefinition* :  
> *SectionHeader* `NL`+ *SectionContent*
> 
> *SectionHeader* :  
> `###` *StringLiteral*

## Example

```
### Tests

// ...
```

## Explanation

A section defines a discrete segment in a Goatfile. This can either be a specialized section like the `Defaults` section or a request section like `Setup`, `Tests` and `Teardown`. More information on these section types can be found under [Lifecycle](../explanations/lifecycle.md).

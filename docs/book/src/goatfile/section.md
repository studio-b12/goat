# Sections

## Syntax

> *SectionDefinition* :  
> *SectionHeader* `NL`+ *SectionContent*
> 
> *SectionHeader* :  
> `###` *StringLiteral*

## Example

```
### Tests

...
```

## Explanation

A section defines a discrete segment in a Goatfile. This can either be a specialized section like the `Defaults` section or a request section like `setup`, `tests` and `teardown`. More on these sections can be found in the [Test Lifecycle Section](../explanations/test-lifecycle.md).

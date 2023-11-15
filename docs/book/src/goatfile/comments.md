# Comments

## Syntax

> *CommentDefinition* :  
> *LineComment* | *BlockComment*
> 
> *LineComment* :  
> `//` .* `NL`
>
> *BlockComment* :  
> `/*` .* `*/`

## Example

```
// This is a line comment!

/// This is a line comment as well!

/*
  This is a
    multiline
  block comment!
*/
```

## Explanation

Contents of comments are ignored by the parser.

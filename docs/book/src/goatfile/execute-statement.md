# Execute Statement

### Syntax

> *ExecuteExpression* :  
> `execute` *StringLiteral*
>
> *Parameters* :  
> `(` (*KeyValuePair* (`WS`|`NL`)+ )* `)` *ReturnStatement*?
>
> *ReturnStatement* :  
> `(` (*ReturnPair* (`WS`|`NL`)+ )* `)`
>
> *ReturnPair* :  
> *StringLiteral* `as` *StringLiteral*
>
> *KeyValuePair* :  
> *StringLiteral* `WS`* `=` `WS`* *StringLiteral*

### Example

```
execute "../utils/login" (
  username="{{.credentials.username}}"
  password="{{.credentials.password}}"
) return (
  userId as userId
)
```

### Explanation

The `execute` statement allows to run an external Goatfile inside another Goatfile **in its own state context**. Parameters for the executed Goatfile can be passed and values in the resulting state can be taken over into the calling Goatfile's state.

In contrast to the `use` directive, the executed Goatfile `A` is run like a completely separate Goatfile execution with its own isolated state which does not share any values with the state of the executing file `B`. All parameters which shall be available in `A` must be passed as a list of key-value pairs. Resulting state values of `A` can then be captured by the state of `B` by listing them in the `return` statement with the name of the parameter in the state of `A` and the name the value shall be accessible under in `B`.

Executed Goatfiles are parsed in place, so they are only statically checked once they are executed within the executing Goatfile.

# Import Statement

### Syntax

> *UseExpression* :  
> `use` *StringLiteral*

### Example

```
use ../path/to/goatfile
```

### Explanation

Foreign Goatfiles can be imported using the `use` statement.

Imported Goatfile behave like they are merged with the root Goatfile. So when you import a Goatfile `B` into a Goatfile `A`, all actions in all sections of `B` will be inserted **infront** of the actions in the sections of `A`. Meanwhile, the order of the sections `Setup`, `Tests` and `Teardown` stays intact.

Cyclic or repeated imports are not allowed.

Schematical Example:

`A.goat`
```
use B

### Setup

GET A1

### Tests

GET A2

### Teardown

GET A3
```

`B.goat`
```
### Setup

GET B1

### Tests

GET B2

### Teardown

GET B3
```

*Result (virtual):*
```
### Setup

GET B1
GET A1

### Tests

GET B2
GET A2

### Teardown

GET B3
GET A3
```

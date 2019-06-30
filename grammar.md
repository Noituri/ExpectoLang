# Expecto Lang - Grammar
**Table of Contents**
- [Functions](#Functions)

## Functions
Functions are created with `fc` keyword and function's name. Example: `fc add`\
Args are comma separated and specified after function's name inside `()`.
Firstly you have to specify argument's name and then its type. Example: `fc add(x: float, y: float)`.\
Return type is specified after args specification, return type is optional.
Example: `fc add(x: float, y: float): float`.\
To return value use `return <value>` keyword.\
Every function must end with `end` keyword. Example:
```
fc add(x:float, y:float): float
    return x + y
end
```

Calling this function looks like that:
```
add(3.0, 2.0)
```
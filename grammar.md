# Expecto Lang - Grammar
**Table of Contents**
- [Code Style](#Style)
- [Functions](#Functions)

## Style
**TODO write more here**
- snake_case is preferable 

## Functions
Functions are created with `fc` keyword and function's name. Example: `fc add`\
Args are comma separated and specified after function's name inside `()`.
Firstly you have to specify argument's name and then its type. Example: `fc add(x: float, y: float)`.
If function does not have any args it's possible to remove parenthesis. \
Return type is specified after args specification, return type is optional.
Example: `fc add(x: float, y: float): float`.\
To return value use `return <value>` keyword.\
Every function must end with `end` keyword. Example:
```Crystal
fc add(x:float, y:float): float
    return x + y
end
```

Calling this function looks like that:
```Crystal
add(3.0, 2.0)
```

To create private function simply add `_` (underscore character) before function's name f.e
```Crystal
fc _add(x:float, y:float): float
    return x + y
end
```

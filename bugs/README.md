## Bugs

### Examples

1. `let nil = if(true){}; nil == nil` --> runtime error
2. `let nil = if(true){}; nil()` --> runtime error
3. `let add = fn(x,y){x+y}; add(2)` --> runtime error

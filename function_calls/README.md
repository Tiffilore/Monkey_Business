## Function Calls 

There are several places where the author seems to make clear that he only wants to allow *identifiers* and *function literals* 
in the function constituent of function calls. However, the way the parser is implemented any expression is allowed as this constituent.

And we have an obvious problem with function calls:

```
if (true){}()
```

causes a runtime error.

So, one might want to correct the parser and implement it in such a way that it restricts the function constituent of function calls to *identifiers* and *function literals*.

However, that doesn't fix the problem, since 

```
let nil = if (true){}
nil()
```
also causes a runtime error. Moreover, one could come up with perfectly reasonable Monkey code where it makes sense to have expressions besides *identifiers* and *function literals* 
as function constituents of function calls:

```
let func1 = fn...
let func2 = fn...
if(...){func1}{func2}(...)
```

```
let func1 = fn...
let func2 = fn...
let choose_func = fn(x){if (x==1){return func1} if (x==2){return func2}}
choose_func(1)(...)
```

Therefore, it is more reasonable to fix the evaluation of ast nodes.

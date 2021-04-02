## visualizing evaluation

- goal: touch evaluator-code as little as possible
- only works if the evaluator doesn't panic 
- environments are not displyed so far
- evaluation trees really are DAGs
  - possible steps:
    - step 1: display then several times
    - step 2: display them separately


### examples:

```
>> :set logtrace

>> :set evalfile eval1
>> if(true){}

>> :set evalfile eval2
>> if(false){}

>> :set evalfile eval3
>> :expr if(true){let dbl = fn(x){2*x}}

>> :set evalfile eval4
>> :expr dbl(3)

>> :set evalfile eval5
>> :expr dbl(dbl(3))

>> :set evalfile eval6 
>> :expr dbl(1) + dbl(3)

>> :set evalfile eval7
>> fn(){}()
```
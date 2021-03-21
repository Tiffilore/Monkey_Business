# Testing the evaluator

**Table of contents**
1. [Existing Tests](#existing)
2. [Additional Tests](#additional)
    1. [TestRuntimeErrorsNotEnoughArguments](#test1)
    2. [TestArityCallExpressions](#test2)
    3. [TestRuntimeErrorsWithNil1](#test3)
    4. [TestRuntimeErrorsWithNil2](#test4)
    5. [TestRuntimeErrorsWithNull](#test5)
    6. [TestRuntimeErrorsWithInvalidPrograms](#test6)

## Existing Tests: `evaluator_test.go` <a name="existing"></a>

For the evaluator, the majority of tests test only input that is supposed to be valid.

There is one Exception: `func TestErrorHandling(t *testing.T)`, which tests type errors.

With regard to the part of Monkey PL introduced in chapter 3, these concerns the following types of type errors (represented by their errormessages):

```
"type mismatch: INTEGER + BOOLEAN"
"unknown operator: -BOOLEAN"
"unknown operator: BOOLEAN + BOOLEAN"
"unknown operator: STRING - STRING"
"identifier not found: foobar"
```

Additionally, there are two tests with hash literals and index expressions, which will be introduced in chapter 4.

## Additional Tests: `evaluator_add_test.go` <a name="additional"></a>

### `TestRuntimeErrorsNotEnoughArguments` <a name="test1"></a>

- tests whether the evaluation of a call expression with not enough arguments causes a runtime error

### `TestArityCallExpressions` <a name="test2"></a>

- once the problem with runtime errors is fixed, we need to specify how we want to treat call expressions with the wrong number of arguments
    - there is no straightforward answer to that; it needs to be discussed and this test is only a starter

### `TestRuntimeErrorsWithNil1` <a name="test3"></a>

- tests whether expressions containing expressions evaluating to `nil` cause the evaluator to panic
    - definition of `nil`: `let nil=if(true){}` (via empty blockstatement in if expression)
    - tested expressions:
        - nil as value of an operand in an prefix expression
        - nil as value of an operand in an infix expression
        - nil as value of the function of a function call

### `TestRuntimeErrorsWithNil2` <a name="test4"></a>

- exactly like `TestRuntimeErrorsWithNil1`, except:
    - definition of `nil`:
`let nil=fn(){}()` (via empty blockstatement in function of function call)

### `TestRuntimeErrorsWithNull` <a name="test5"></a>

- same test set as before, but this time with expressions evaluating to the **NULL** object
    - definition of `null`: `let null=if(false){}` 
- this test succeeds for the whole test set, but might serve as a starter in discussing how NULL values should be treated and whether the current implementation does it consistently

### `TestRuntimeErrorsWithInvalidPrograms` <a name="test6"></a>

- tests whether the evaluation of defect asts, i.e. asts accompanied by error messages in the parser causes panic
    - spoiler: it does right now
- it might be disputed whether this test does not put too high standards on the evaluator, since usually it will be only used after checking the parser for errors

- examples for defect parse trees (with and without token fields):

  - `@ let`

    <img src="images/ast_wo_tok0.png" width="300" />
    <img src="images/ast_with_tok0.png" width="350" />

  - `let;@;`

    <img src="images/ast_wo_tok1.png" width="600" />

    <img src="images/ast_with_tok1.png" width="600" />



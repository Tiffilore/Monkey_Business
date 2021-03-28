# Testing the evaluator

[_Don't panic!_](https://go-proverbs.github.io/)


**Table of contents**
1. [Existing Tests](#existing)
2. [Additional Tests: `evaluator_add_panic_test.go`](#additional_panic)
    1. [TestRuntimeErrorsWithNil1](#nil1)
    2. [TestRuntimeErrorsWithNil2](#nil2)
    3. [TestRuntimeErrorsWithNull](#null)
    4. [TestRuntimeErrorsWithInvalidPrograms](#invalid)
3. [Additional Tests: `evaluator_add_test.go`](#additional)
    1. [TestArityCallExpressions](#arity_call)
    2. [TestEvalToBoolConsistency and TestEvalToBoolCorrectness](#eval2bool)
    3. [TestDivisionByZero](#div_zero)

## Existing Tests: `evaluator_test.go` <a name="existing"></a>

For the evaluator, the majority of tests test only input that is supposed to be valid.

There is one exception: `func TestErrorHandling(t *testing.T)`, which tests type errors.

With regard to the part of Monkey PL introduced in chapter 3, these concerns the following types of type errors (represented by their errormessages):

```
"type mismatch: INTEGER + BOOLEAN"
"unknown operator: -BOOLEAN"
"unknown operator: BOOLEAN + BOOLEAN"
"unknown operator: STRING - STRING"
"identifier not found: foobar"
```

Additionally, there are two tests with hash literals and index expressions, which will be introduced in chapter 4.


## Additional Tests: `evaluator_add_panic_test.go` <a name="additional_panic"></a>


Look at [this question](https://stackoverflow.com/questions/66841082/why-does-the-monkey-repl-panic-at-my-program) from stackoverflow:

<img src="images/stackoverflow.selection.png" width="1000" />

Can you answer it?

---
_in progress:_

The aim of these tests is primarily to collect all cases that cause the evaluator to panic.

TODO: not enough arguments is now missing 

most of it is due to Monkey in its current state allowing expressions to be evaluated to `nil`.

- when does it happen?
- when the meaning of an expression is derived from a blockstatement:
  - in if expressions (either consequence or for if-else syntax: alternative)
  - in function calls: body of function object that function field evaluates to

- when does a blockstatement evalate to `nil`?
  - when it is empty
  - when its last statement 
    - is a let statement
    - is an expression statement evaluating to `nil`

- examples (types):
```
  - if(<expression evaluating to truthy>)<blockstatement evaluating to nil>
  - if(<expression not evaluating to an error, but to a value that is not truthy)<any blockstatement> else <blockstatement evaluating to nil>
  - <expression evaluating to a function object with a body blockstatement evaluating to nil>(<argument list>)
```

- concrete examples:
    - minimal:
```
if(true){}
if(false){1}{}
fn(){}()
```


### `TestRuntimeErrorsWithNil1` <a name="nil1"></a>

- tests whether expressions containing expressions evaluating to `nil` cause the evaluator to panic
    - definition of `nil`: `let nil=if(true){}` (via empty blockstatement in if expression)
    - tested expressions:
        - nil as value of an operand in an prefix expression
        - nil as value of an operand in an infix expression
        - nil as value of the function of a function call

### `TestRuntimeErrorsWithNil2` <a name="nil2"></a>

- exactly like `TestRuntimeErrorsWithNil1`, except:
    - definition of `nil`:
`let nil=fn(){}()` (via empty blockstatement in function of function call)

### `TestRuntimeErrorsWithNull` <a name="null"></a>

- same test set as before, but this time with expressions evaluating to the **NULL** object
    - definition of `null`: `let null=if(false){}` 
- this test succeeds for the whole test set, but might serve as a starter in discussing how NULL values should be treated and whether the current implementation does it consistently

### `TestRuntimeErrorsWithInvalidPrograms` <a name="invalid"></a>

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




## Additional Tests: `evaluator_add_test.go` <a name="additional"></a>

### `TestArityCallExpressions` <a name="arity_call"></a>

- specifies how to deal with call expressions with not enough / too many arguments
  - in the current implementation, the evaluator panics at the face of not enough arguments
  - matter of specification (thus discussion):
    - do we want to return an error if there are too many arguments?
    - what error messages?


### `TestEvalToBoolConsistency` and `TestEvalToBoolCorrectness` <a name="eval2bool"></a>

In the Monkey PL, every non-erroneous expression can be evaluated to a Boolean value. This can be done in two places: 
- in the Condition field of an if expressions
- (implicitly) in the evaluation of prefix expression with BANG as operator

The first desideratum is **consistency**: we want the evaluation of a condition to a boolean be consistent with its evaluation to a boolean in a prefix expression with BANG as operator.
That means that for any `<expression>`, the evaluation of `if(<expression>){true}else{false}` should yield the same result as the evaluation of `!!<expression>`

The second desideratum is **correctness**: we want the evaluations to be correct. What a correct evaluation is varies from language to language and is a matter of language specification. This test serves as an opportunity to discuss exactly that.

Consistency is being tested in `TestEvalToBoolConsistency`, while correctness is being tested in `TestEvalToBoolCorrectness`. 

In the current implementation, `TestEvalToBoolConsistency` succeeds, but that can easily change if we opt to change the implementation. Evaluation to booleans is implemented twice in the code: for conditions with the help of the function `isTruthy` and for prefix expressions with a BANG operator in `evalBangOperatorExpression`. It thus may serve as a regression test.

`TestEvalToBoolCorrectness` does not suceed, since it needs to be specified correctly first. I wanted to leave the specification open at this stage.

#### Test data

- we want to test the handling of the following types of `object.Object`s:

object type | values
---|---
`Boolean` | true, false
`Integer` | -1, 0, 1
`Null` | the one and only null object
`Error` | any
`Function` | any


- we will skip `ReturnValue` objects, since they can never be values of expressions and for now all object types that are only introduced in chapter 4: `String`, `Builtin`, `Array` and `Hash`.
- we want to add the infamous `nil`, since expressions in MonkeyPL can still evaluate to `nil` given the recent implementation.
- in the given implementation, the only object type for whiches boolean evaluation it matters, what its value is, is `Boolean`. However, in many languages, when numbers are evaluated to Booleans, their evaluation also varies with regard to their value. Since one might opt for such an implementation for Monkey Pl, too, there are three possible values added for 
- we could use expressions in our testdata that evaluate to the desired objects (e.g. `fn(){}()` for `nil`, `if(false){}`for `NULL`), but this has the drawback that any changes in the evaluation of such expressions will undermine our tests. Thus, we opt for creating an environment mapping the name "a" to the respective values (for example, `TRUE` or `&object.Integer{Value: -1}`) and then use the name in the expression. Here is some example code (not from the actual tests) illustrating how this approach works: 

```go
	env := object.NewEnvironment()
	env.Set("a", &object.Integer{Value: -1})
	input := "!!a"
	l := lexer.New(input)
	p := parser.New(l)
	ast := p.ParseProgram()
	result := Eval(ast, env)
```

### `TestDivisionByZero` <a name="div_zero"></a>

In the current implementation, dividing a number by zero causes a runtime exception, as the test shows.

However, maybe we want to opt for an implementation of integer division that differs from Golang's choices.

`TestDivisionByZero` demands that division by zero returns an error. The error message still needs to be specified.





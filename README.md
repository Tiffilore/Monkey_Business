# Extending the interpreter and its interactive environment

 _to be continued..._

## Step 2: Implement a Small Initial Instruction Set for the interactive environment

- new interactive environment in directory `session`
    - inspo from [ghci](https://downloads.haskell.org/~ghc/latest/docs/html/users_guide/ghci.html#ghci-commands) and [gore](https://github.com/motemen/gore) 

| NAME   |                   | USAGE                                                   |
|--------|-------------------|---------------------------------------------------------|
| clear  | ~                 | clear the environment                                   |
| h[elp] | ~                 | list all commands with usage                            |
|        | ~ <cmd>           | print usage command <cmd>                               |
| l[ist] | ~                 | list all identifiers in the environment alphabetically  |
|        |                   |      with types and values                              |
| q[uit] | ~                 | quit the session                                        |
| reset  | ~ prompt          | set prompt to default                                   |
| set    | ~ prompt <prompt> | set prompt string to <prompt>                           |

### `clear`, `h[elp]`, `q[uit]`, `set prompt`
![Demo1](demos/demo1.gif)

### `reset prompt`, `list` 
![Demo2](demos/demo2.gif)

### TODO: Multiline support: `set|unset|reset paste`, `:{...:}`
![Demo3](demos/demo3.gif)


## Step 1: Writing Tests for Bugs in Parser and Evaluator

- `parser/parser_add_test.go`
- `evaluator/evaluator_add_test.go`

## Step 0: Starting Point: Copy the Code

- unaltered code from the book [_Writing an Interpreter in Go_](https://interpreterbook.com/), Version 1.7






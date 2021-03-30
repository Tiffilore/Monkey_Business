# Extending the REPL/ interactive environment

inspired by [ghci](https://downloads.haskell.org/~ghc/latest/docs/html/users_guide/ghci.html#ghci-commands) and [gore](https://github.com/motemen/gore) 


## Initial Command Set

| NAME   |                   | USAGE                                                   |
|--------|-------------------|---------------------------------------------------------|
| clear    | ~                 | clear the environment                                    |
| e[val]   | ~ `<input>`         | print out value of object `<input>` evaluates to           |
| h[elp]   | ~                 | list all commands with usage                             |
|          | ~ `<cmd> `          | print usage command `<cmd>`                                |
| l[ist]   | ~                 | list all identifiers in the environment alphabetically   |
|          |                   |      with types and values                               |
| paste    | ~ `<input>`         | evaluate multiline `<input>` (terminated by blank line)    |
| q[uit]   | ~                 | quit the session                                         |
| reset    | ~ logtype         | set logtype to default                                   |
|          | ~ paste           | set multiline support to default                         |
|          | ~ prompt          | set prompt to default                                    |
| set      | ~ logtype         | when eval, additionally output objecttype                |
|          | ~ paste           | enable multiline support                                 |
|          | ~ prompt `<prompt>` | set prompt string to `<prompt> `                           |
| settings | ~                 | list all settings with their current values and defaults |
| t[ype]   | ~ `<input>`         | show objecttype `<input>` evaluates to                     |
| unset    | ~ logtype         | when eval, don't additionally output objecttype          |
|          | ~ paste           | disable multiline support                                |


#### `clear`, `h[elp]`, `q[uit]`, `set prompt <prompt>`
![Demo1](../demos/demo1.gif)

#### `reset prompt`, `list` 
![Demo2](../demos/demo2.gif)

#### `t[ype]`,  `(set|unset|reset) logtype`, `e[val]`
- eval is default
    - if the user input is not prefixed by `:`, `input` is equivalent to `:eval input`

![Demo3](../demos/demo3.gif)

#### Multiline support: `paste`, `(set|unset|reset) paste` and `settings`

![Demo4](../demos/demo4.gif)

## Adding dimensions: settings `level` and `process` 

| NAME           |                   | USAGE                                                    |
--- | --- | --- |
| clear          | ~                 | clear the environment                                    |
| e[val]         | ~ `<input>`         | print out value of object `<input>` evaluates to           |
| expr[ession]   | ~ `<input>`         | expect `<input>` to be an expression                       |
| h[elp]         | ~                 | list all commands with usage                             |
|                | ~ `<cmd>`           | print usage command `<cmd>`                                |
| l[ist]         | ~                 | list all identifiers in the environment alphabetically   |
|                |                   |      with types and values                               |
| p[arse]        | ~ `<input>`         | parse `<input>`                                            |
| paste          | ~ `<input>`         | evaluate multiline `<input>` (terminated by blank line)    |
| prog[ram]      | ~ `<input>`         | expect `<input>` to be a program                           |
| q[uit]         | ~                 | quit the session                                         |
| reset          | ~ process         | set process to default                                   |
|                | ~ level           | set level to default                                     |
|                | ~ logparse        | set logparse to default                                  |
|                | ~ logtype         | set logtype to default                                   |
|                | ~ paste           | set multiline support to default                         |
|                | ~ prompt          | set prompt to default                                    |
| set            | ~ process `<p>`     | `<p>` must be: [e]val, [p]arse, [t]ype                     |
|                | ~ level `<l>`       | `<l>` must be: [p]rogram, [s]tatement, [e]xpression        |
|                | ~ logparse        | additionally output ast-string                           |
|                | ~ logtype         | additionally output objecttype                           |
|                | ~ paste           | enable multiline support                                 |
|                | ~ prompt `<prompt>` | set prompt string to `<prompt>`                            |
| settings       | ~                 | list all settings with their current values and defaults |
| stmt|statement | ~ `<input>`         | expect `<input>` to be a statement                         |
| t[ype]         | ~ `<input>`         | show objecttype `<input>` evaluates to                     |
| unset          | ~ logparse        | don't additionally output ast-string                     |
|                | ~ logtype         | don't additionally output objecttype                     |
|                | ~ paste           | disable multiline support                                |



#### `(set|reset) level (program|statement|expression)`, `parse`, `logparse`

#### `(set|reset) process  (parse|eval|type)`, `expr[ession]`, `stmt|statement` and `prog[ram]`

![Demo6](../demos/demo6.gif)



## Step 4: Add representation of asts 

The demo shows the output for the current function call in session.go:

The function `process_input_dim` of `session.go` determines the processing of Monkey input, given what the current settings and the current command determine with regard to 
- multiline input 
- input level (expression/statement/program)
- input process (parse/type/eval)

The piece of code, where the new representation is added, is:

```go
	if s.logparse || process == ParseP {
		fmt.Fprintln(s.out, node)
		fmt.Fprintln(s.out, ast.RepresentNodeConsoleTree(node, "|   ", true)) //<---
	}
```
The arrow indicates the added line.

The function RepresentNodeConsoleTree is defined in `ast/ast_add.go` and has the signature 

`func RepresentNodeConsoleTree(node Node, indent string, exclToken bool) string `

The indentation can be set by the second Parameter and the third specifies whether Token-fields are to be displayed or not. Currently, they are not displayed.

![Demo7](../demos/demo7.gif)

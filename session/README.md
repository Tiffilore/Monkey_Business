# Extending the REPL/ interactive environment

inspired by [ghci](https://downloads.haskell.org/~ghc/latest/docs/html/users_guide/ghci.html#ghci-commands) and [gore](https://github.com/motemen/gore) 


## Initial Command Set

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

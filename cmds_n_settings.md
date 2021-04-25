# Commands and Settings 



## Ch-ch-ch-changes

- [X] add c[lear], tr[ace]
- [X] re-arrange commands in help menu
  - before: in alphabetical order 
  - now: in order of registration
- [X] check paths for "clear" and "pdflatex" at beginning
- [ ] decide which setting types are public and which are not !!!
- ask for file if not there?

- write top-level user manual


## Commands new:

```
+----------------+------------------------------+---------------------------------------------------------+
| NAME           |                              | USAGE                                                   |
+----------------+------------------------------+---------------------------------------------------------+
| h[elp]         | ~                            | list all commands with usage                            |
|                | ~ <cmd>                      | print usage command <cmd>                               |
| q[uit]         | ~                            | quit the session                                        |
| cl[earscreen]  | ~                            | clear the terminal screen                               |
| l[ist]         | ~                            | list all identifiers in the environment alphabetically  |
|                |                              |      with types and values                              |
| c[lear]        | ~                            | clear the environment                                   |
| paste          | ~ <input>                    | evaluate multiline <input> (terminated by blank line)   |
| expr[ession]   | ~ <input>                    | expect <input> to be an expression                      |
| stmt|statement | ~ <input>                    | expect <input> to be a statement                        |
| prog[ram]      | ~ <input>                    | expect <input> to be a program                          |
| p[arse]        | ~ <input>                    | print string representation of ast <input> is parsed to |
|                |                              |      --> Node-method: String() string                   |
| p[arse]tree    | ~ <input>                    | print tree representation of <input>' ast               |
|                |                              |      to all set displays                                |
| e[val]         | ~ <input>                    | print value of object <input> evaluates to              |
|                |                              |      --> Object-method: Inspect() string                |
| t[ype]         | ~ <input>                    | show objecttype <input> evaluates to                    |
| tr[ace]        | ~ <input>                    | show evaluation trace interactively step by step        |
| e[val]tree     | ~ <input>                    | print annotated tree representation of <input>'s ast    |
|                |                              |      to all set displays                                |
| settings       | ~                            | list all settings with their current and default values |
| set            | ~ prompt <prompt>            | set prompt string to <prompt>                           |
|                | ~ paste                      | enable multiline support                                |
|                | ~ level <l>                  | <l> must be: p[rogram], s[tatement], e[xpression]       |
|                | ~ process <p>                | <p> must be: p[arse], p[arse]tree, e[val], e[val]tree,  |
|                |                              |      [t]ype, [tr]ace                                    |
|                | ~ logs <+|-l_0...+|-l_n>     | <l_i> must be: p[arse]tree, e[val]tree, [t]ype, [tr]ace |
|                | ~ displays <+|-l_0...+|-l_n> | <l_i> must be: p[arse]tree, e[val]tree, [t]ype, [tr]ace |
|                | ~ verbosity <v>              | <v> must be 0, 1, 2                                     |
|                | ~ inclToken                  | include tokens in representations of asts               |
|                | ~ inclEnv                    | include environments in representations of asts         |
|                | ~ file <f>                   | set file for pdfs to <f>                                |
| reset          | ~                            | reset all settings                                      |
|                | ~ <setting>                  | set <setting> to default value                          |
|                |                              |      for an overview consult :settings and/or :h set    |
| unset          | ~ <setting>                  | set boolean <setting> to false                          |
|                |                              |      for an overview consult :settings and/or :h set    |
+----------------+------------------------------+---------------------------------------------------------+               |                   |      for an overview consult :settings and/or :h set     |

```

## Settings new:

```
+-----------+---------------+---------------+
| SETTING   | CURRENT VALUE | DEFAULT VALUE |
+-----------+---------------+---------------+
| prompt    | >>            | >>            |
| paste     | false         | false         |
| level     | program       | program       |
| process   | eval          | eval          |
| logs      | []            | []            |
| displays  | [console]     | [console]     |
| verbosity | 0             | 0             |
| inclToken | false         | false         |
| inclEnv   | false         | false         |
| file      | tree.pdf      | tree.pdf      |
+-----------+---------------+---------------+
```

## flow 

![Flow](assets/images/flow.png)











paste          | ~ <input>         | evaluate multiline <input> (terminated by blank line)    |
| expr[ession]   | ~ <input>         | expect <input> to be an expression                       |
| stmt|statement | ~ <input>         | expect <input> to be a statement                         |
| prog[ram]      | ~ <input>         | expect <input> to be a program                           |
| p[arse]        | ~ <input>         | parse <input>                                            |
| p[arse]tree    | ~ <input>         | parse <input>                                            |
| e[val]         | ~ <input>         | print out value of object <input> evaluates to           |
| t[ype]         | ~ <input>         | show objecttype <input> evaluates to                     |
| tr[ace]          | ~ <input>         | show evaluation trace step by step                       |
| e[val]tree     | ~ <input>         | print out value of object <input> evaluates to           |
| settings       | ~                 | list all settings with their current values and defaults |
| set            | ~ process <p>     | <p> must be: [e]val, [p]arse, [t]ype                     |
|                | ~ level <l>       | <l> must be: [p]rogram, [s]tatement, [e]xpression        |
|                | ~ logparse        | additionally output ast-string                           |
|                | ~ logtype         | additionally output objecttype                           |
|                | ~ logtrace        | additionally output evaluation trace                     |
|                | ~ incltoken       | include tokens in representations of asts                |
|                | ~ paste           | enable multiline support                                 |
|                | ~ prompt <prompt> | set prompt string to <prompt>                            |
|                | ~ treefile <f>    | set file that outputs ast-pdfs to <f>                    |
|                | ~ evalfile <f>    | set file that outputs eval-pdfs to <f>                   |
| reset          | ~ <setting>       | set <setting> to default                                 |
|                |                   |      for an overview consult :settings and/or :h set     |
| unset          | ~ <setting>       | set boolean <setting> to false                           |
|                |                   |      for an overview consult :settings and/or :h set     |
+-----------
# Commands and Settings

## TODO 
- set filenames => what if we have both, log etree and log ptree?

## Ch-ch-ch-changes

- [X] add c[lear], tr[ace]
- [X] re-arrange commands in help menu
  - before: in alphabetical order
  - now: in order of registration
- [X] check paths for "clear" and "pdflatex" at beginning
- [X] add secondary prompt for paste
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
|                | ~ displays <+|-d_0...+|-d_n> | <d_i> must be: c[ons[ole]], p[df]                       |
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


## check implementations:

- [X] h[elp]
  - session.go: [session.]exec_help     --> commands.go: [commands.]usage(s(string) (string, bool)
  - session.go: [session.]exec_help_all --> commands.go: [commands.]menu() string
- [X] q[uit]
  - session.go: [session.]exec_quit
- [X] cl[earscreen]
  - session.go: returned by f_exec_clearscreen at init
- [-] l[ist]
  - session.go: [session.]exec_list --> visualizer.GetStoreTable(object.env)
- [X] c[lear]
  - session.go: [session.]exec_clear
- [X] paste
  - session.go: [session.]exec_paste 
  + revise [session.]exec_cmd for commands that take argument
- [ ] expr[ession]
- [ ] stmt|statement
- [ ] prog[ram]
- [ ] p[arse]
- [ ] p[arse]tree
- [ ] e[val]
- [ ] t[ype]
- [ ] tr[ace]
- [ ] e[val]tree
- [X] settings
  - settings.go: [session.]exec_settings --> settings.go: menuSettings() string
- [X] set
  - settings.go: [session.]exec_set --> settings.go: set(string) bool
- [X] reset
  - settings.go: [session.]exec_reset --> settings.go: reset(string) bool
  - settings.go: [session.]exec_reset_all 
- [X] unset
  - settings.go: [session.]exec_unset --> settings.go: unset(string) bool

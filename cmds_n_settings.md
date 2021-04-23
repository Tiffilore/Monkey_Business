# Commands and Settings 

- path pdflatex --> set at init
- ask for file if not there?


## Ch-ch-ch-changes

- [X] add c[lear] 
- write top-level user manual
- re-arrange commands in help menu
  - now in alphabetical order 
  - then in chosen order [e.g. order of registration]

:vis -pdf -v -token -eval -env -paste
:vis -cons -vv




## new Commands:

```
+----------------+-------------------+----------------------------------------------------------+
| NAME           |                   | USAGE                                                    |
+----------------+-------------------+----------------------------------------------------------+

GENERAL:
| cl[earscreen]  | ~                 | clear the terminal screen                                |
| h[elp]         | ~                 | list all commands with usage                             |
|                | ~ <cmd>           | print usage command <cmd>                                |
| q[uit]         | ~                 | quit the session                                         |


ENVIRONMENT:
| c[lear]        | ~                 | clear the environment                                    | # c added
| l[ist]         | ~                 | list all identifiers in the environment alphabetically   | -- dependent on verbosity; maybe also on inclenvironment
|                |                   |      with types and values                               |

SETTINGS:
| settings       | ~                 | list all settings with their current values and defaults |
|                |                   |      for an overview consult :settings and/or :h set     |
| set            | ~ ...             |                                                          |
|
| reset          | ~ <setting>       | set <setting> to default                                 |
| unset          | ~ <setting>       | set boolean <setting> to false                           |
|                |                   |      for an overview consult :settings and/or :h set     |

PASTE:
| paste          | ~ <input>         | evaluate multiline <input> (terminated by blank line)    |

LEVEL:
| expr[ession]   | ~ <input>         | expect <input> to be an expression                       |
| stmt|statement | ~ <input>         | expect <input> to be a statement                         |
| prog[ram]      | ~ <input>         | expect <input> to be a program                           |


PROCESS:
| e[val]         | ~ <input>         | print out value of object <input> evaluates to           |
| p[arse]        | ~ <input>         | parse <input>                                            |
| t[ype]         | ~ <input>         | show objecttype <input> evaluates to                     |
| trace          | ~ <input>         | show evaluation trace step by step                       |
| parsetree      | ~ <input>         | show parsetree                                           |
| evaltree       | ~ <input>         | show evaltree                                            |

```

## new Settings:

```
+-----------+---------------+---------------+
| SETTING   | CURRENT VALUE | DEFAULT VALUE |
+-----------+---------------+---------------+
| prompt    | >>            | >>            |--> delete?

| paste     | false         | false         | in {true, false}

| level     | program       | program       | in {program, statement, expression}
| process   | eval          | eval          | in {parse, parsetree, eval, type, trace, evaltree}

| logtree   | false         | false         | in {true, false}
| logtype   | false         | false         | in {true, false}
| logtrace  | false         | false         | in {true, false}

| displays  | cons          | cons          | in {cons, pdf, both}
| verbosity | 0             | 0             | in {0,1,2}
| incltoken | false         | false         | in {true, false}
| inclenv   | false         | false         | in {true, false}
| file      | tree.pdf      | tree.pdf      |FILE
```

## flow 

![Flow](assets/images/flow.png)




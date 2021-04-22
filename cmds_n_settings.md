# Commands and Settings 

## Current Commands and Settings Grouped

### Commands
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
| l[ist]         | ~                 | list all identifiers in the environment alphabetically   |
|                |                   |      with types and values                               |

PASTE:
| paste          | ~ <input>         | evaluate multiline <input> (terminated by blank line)    |

LEVEL:
| expr[ession]   | ~ <input>         | expect <input> to be an expression                       |
| stmt|statement | ~ <input>         | expect <input> to be a statement                         |
| prog[ram]      | ~ <input>         | expect <input> to be a program                           |

TRACE:
| trace          | ~ <input>         | show evaluation trace step by step                       |

PROCESS:
| e[val]         | ~ <input>         | print out value of object <input> evaluates to           |
| p[arse]        | ~ <input>         | parse <input>                                            |
| t[ype]         | ~ <input>         | show objecttype <input> evaluates to                     |


SETTINGS:
| settings       | ~                 | list all settings with their current values and defaults |
|                |                   |      for an overview consult :settings and/or :h set     |
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
| unset          | ~ <setting>       | set boolean <setting> to false                           |
|                |                   |      for an overview consult :settings and/or :h set     |
```

### Settings
```
+-----------+---------------+---------------+
| SETTING   | CURRENT VALUE | DEFAULT VALUE |
+-----------+---------------+---------------+
| paste     | false         | false         |PASTE
| process   | eval          | eval          |PROCESS
| level     | program       | program       |LEVEL
| logparse  | false         | false         |LOG
| logtype   | false         | false         |LOG
| logtrace  | false         | false         |LOG
| prompt    | >>            | >>            |--> delete!
| incltoken | false         | false         |VERBOSITY
| treefile  | tree.pdf      | tree.pdf      |FILE
| evalfile  | eval.pdf      | eval.pdf      |FILE
+-----------+---------------+---------------+
```


### Ch-ch-ch-changes

- add c[lear]?
- delete set prompt
- solutions for logging?
- solutions for pdf-creation?
- write top-level user manual
- re-arrange commands in help menu
  - now in alphabetical order 
  - then in chosen order [e.g. order of registration]

### Options for trees:

- options for tracer: 
  - level --> setting

- verbosity in {V|VV|VVV}
- + eval [not process; there is no equivalent to type]
- + token
- + env 
- level --> setting?
- console | pdf 
- if pdf: filename
- paste not from setting


:vis -pdf -v -token -eval -env -paste
:vis -cons -vv
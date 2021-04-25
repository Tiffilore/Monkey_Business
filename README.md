# Extending the interpreter and its interactive environment

_The purpose of the Monkey Programming Language is to learn different aspects of implementing a language. It was invented by Thorsten Ball, who offers an implementation in his book [Writing An Interpreter In Go](https://interpreterbook.com/).
The main purpose of this repo is to extend the possibilities to explore this language by the help of various visualizations of the chosen abstract syntax tree (ast) and the steps of evaluation. These visualizations can be requested by an extended interactive environment that offers a set of multidimensional options and tools for convenience.
Since Monkey is a language intended for learning purposes, the implementations are designed in such a way that they are open for changes and additions to the parser and evaluator. There is also a set of tests that may give rise to changes in the implementation, because they document certain problems within the interpreter in the state described in the book._

## Project Status

**Currently, this repo is still very much work in progress.**

### What has been done: 
[changelog](changelog.md)

### Next planned steps:

- [X] decide on "final" instruction set
  - what is a setting? what is a command?
  - when should pdfs be created?
- [X] change session.go accordingly
- [ ] decide package structure
- [ ] implement all commands
- [ ] write user manual
- [ ] add workflow 
    - mainly to check installation requirements
- [ ] write discussion doc 
- [ ] add tests for ast: String-methods

## The New Interactive Environment

- implemented in monkey/session
- replaces monkey/repl
- still called in main.go

### Run

#### Prerequisites [not tested yet, TODO]

In addition to go,  the command `pdflatex` needs to be installed for creating pdfs. 
You can check whether `pdflatex` is installed by `which pdflatex`.


For Ubuntu, the installation can be done by:

```sh
sudo apt-get install texlive-base texlive-latex-extra
```


#### Run locally

Cou can use the interactive environment locally by cloning this repo, moving into it and then execute

```sh
go run main.go
```

The interpreter code (i.e. the modules monkey/{token,lexer,ast,parser,object,evaluator}) is the original code from the interpreter book (Version 1.7) with only very few alterations described here (TODO).

You can alter the code or add to it and visualize the differences in the interactive environment.

A starting points for altering might be the additional tests

### How-To Use: TODO 

- [ ] see (yet non-existent) User Manual
- [ ] maybe small demo

## Discussion of the Interpreter: TODO 
- [ ] see (yet non-existent)doc-discussion
- [ ] see tests + their doc (add links)

## Credits 
- original interpreter code: [Writing An Interpreter In Go](https://interpreterbook.com/)
- inspiration for command sets in interactive environments: 
 [swipl](https://www.swi-prolog.org/) and
 [ghci](https://downloads.haskell.org/~ghc/latest/docs/html/users_guide/ghci.html#ghci-commands) 
- inspiration for implementing an interactive environment in go: [gore](https://github.com/motemen/gore) 

## License

[MIT LICENSE](LICENSE)


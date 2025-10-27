# Chapter 1: Lexing

The goal of lexing is to convert raw source code into a list of tokens. Tokens are the smallest way you can seperate source code in a meaningful way.

### Lexical analysis
Scanning source code by character and identifying symbols

- Identifiers (`x`, `add`, `result`)
- Keywords (`let`, `fn`, `if`, `else`, `return`)
- Operators (`+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`)
- Literals (`5`, `"hello"`, `true`, `false`)
- Delimiters (`(`, `)`, `{`, `}`, `,`, `;`)

For example

`let x = 5 + 10;` becomes `LET, IDENT(x), ASSIGN(=), INT(5), PLUS(+), INT(10), SEMICOLON(;)`

Whitespaces are ignored because they're not meaningful in this language.

Each token has a Type and their literal (string value). There are special tokens:  
`ILLEGAL` is for unknown character.  
`EOF` marks the end of the file.

### Structure of the lexer

The lexer reads one character at a time, updating three pointers

```go
type Lexer struct {
    input        string
    position     int //  current position in input
    readPosition int //  next position to read
    ch           byte // current character
}
```

readChar() moves forward in the input string.

NextToken() does most of the work:
- Reads the next character.
- Matche it against known symbols.
- Returns the appropriate token.

Helper function
- isLetter(ch) checks for alphabetic or _.
- isDigit(ch) checks for 0–9.
- readIdentifier() reads names or keywords.
- readNumber() reads integers.
- skipWhitespace() ignores spaces and newlines.

### REPL (Read-Eval-Print-Loop)
Even though we can't evaluate code yet, we can tokenize interactively
```go
>> let x = 5 + 10;
{Type:LET Literal:let}
{Type:IDENT Literal:x}
{Type:= Literal:=}
{Type:INT Literal:5}
{Type:+ Literal:+}
{Type:INT Literal:10}
{Type:; Literal:;}
```

# Chapter 2: Parsing
Convert the token stream into a structured AST (Abstract Tree System)

A parser organizes token according to the grammar rules of the language. It checks for valid syntax and creates a tree representation.

Parser generators like ANTLR, yacc and bison exist, but the goal is to understand how a parser works, not simply use it.  
Writing one manually gives control and clarity over how the language is structured.

### Building the parser

The parser converts tokens into nodes. Those nodes are used to build a tree.  

Example nodes:
- LetStatement
- ReturnStatement
- Identifier
- IntegerLiteral
- PrefixExpression
- InfixExpression
- IfExpression
- FunctionLiteral
- CallExpression

To handle precedence correctly, Pratt Parsing is used.  
Each token has a precedence level (`*` > `+`)  
Each token type has:
- a prefix parse function (for numbers, identifiers and negations)
- an infix parse function (for binary operators)

The parser keeps track of the current token and peek token, building the AST as it advances.  
For example:
```go
let result = add(5, 6);
```
creates:
```
LetStatement
    |--- Name: Identifier("result")
    |--- Value: CallExpression(
        Function=Identifier("add),
        Arguments=[Integer(5), Integer(10)]
    )
```

At the end of parsing, the code is represented as a tree:
```
Program
 |-- LetStatement
 |   |-- Identifier(add)
 |   |-- FunctionLiteral(params=[a, b], body=a+b)
 |-- ExpressionStatement
     |-- CallExpression(add, args=[2, 3])
```

This AST structure will be traversed by the evaluator to run the code.


# Chapter 3: Evaluation
Executing the AST and returning results

The evaluator interprets each node in the AST.  
Each node type defines how it evaluates.

- IntegerLiteral → returns an Integer object.
- InfixExpression → evaluates left and right sides, applies the operator.
- IfExpression → checks condition, evaluates the chosen block.
- FunctionLiteral → creates a Function object storing parameters and body.
- CallExpression → invokes a function

### Object system

To evaluate the code, there needs to be a system to represent runtime values.
```go
type ObjectType string

type Object interface {
    Type()    ObjectType
    Inspect() string
}
```
Built-in Types:

- INTEGER_OBJ
- BOOLEAN_OBJ
- NULL_OBJ
- RETURN_VALUE_OBJ
- ERROR_OBJ
- FUNCTION_OBJ

Each corresponds to a Go struct:
```go
type Integer struct { Value int64 }
type Boolean struct { Value bool }
type Null struct {}
```

### Evaluation
#### Arithmetics
`5 + 10 * 2 + 3` -> Evaluator walks the tree and does the math using go.

#### Return statements
Handled by `ReturnValue` wrapper, stops the execution of a block

#### Environment (variable scope)
An environment map is used
```go
type Environment struct {
    store map[string]Object
    outer *Environment
}
```

It stores variable bindings and supports lexical scoping

#### Functions
When a function is defined:  
```
let add = fn(a, b) { a + b; };
```
It becomes a `Function` object storing its parameters, body and environment at the definition time.

When called
```
add(2, 3);
```
It creates a new environment, binds `a=2` and `b=3` and evaluates the body `a + b`.

### Error handling
If something goes wrong, the evaluator returns an Error object.
```log
ERROR: type mismatch: INTEGER + BOOLEAN
```
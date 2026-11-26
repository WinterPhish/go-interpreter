package session

import (
	"fmt"
	"strings"

	"selenite/evaluator"
	"selenite/lexer"
	"selenite/object"
	"selenite/parser"
	"selenite/token"
)

// Session remembers variables and macros between lines,
// like one REPL window does.
type Session struct {
	env      *object.Environment
	macroEnv *object.Environment
}

func New() *Session {
	return &Session{
		env:      object.NewEnvironment(),
		macroEnv: object.NewEnvironment(),
	}
}

// Explanation is what each phase made of one piece of code.
type Explanation struct {
	Tokens  string // lexer: one token per line, as "TYPE<tab>text"
	Grouped string // parser: the code written back out, with parentheses showing the grouping
	Tree    string // ast: the syntax tree, one node per line
	Result  string // evaluator: what the code came to
	Type    string // object: the type of that result
	Failed  bool   // the parser could not understand the code
}

// Run executes code and returns what the REPL would have printed.
func (s *Session) Run(input string) string {
	return s.Explain(input).Result
}

// Explain executes code, and also keeps what each phase produced along the way.
func (s *Session) Explain(input string) Explanation {
	ex := Explanation{Tokens: tokens(input)}

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		ex.Result = parserErrors(p.Errors())
		ex.Failed = true
		return ex
	}

	ex.Grouped = program.String()
	ex.Tree = Tree(program)

	evaluator.DefineMacros(program, s.macroEnv)
	expanded := evaluator.ExpandMacros(program, s.macroEnv)

	evaluated := evaluator.Eval(expanded, s.env)
	if evaluated != nil {
		ex.Result = evaluated.Inspect()
		ex.Type = fmt.Sprint(evaluated.Type())
	}
	return ex
}

// tokens runs only the lexer over the code and lists what it finds.
func tokens(input string) string {
	l := lexer.New(input)
	var lines []string
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		lines = append(lines, fmt.Sprintf("%s\t%s", tok.Type, tok.Literal))
	}
	return strings.Join(lines, "\n")
}

func parserErrors(errors []string) string {
	var b strings.Builder
	b.WriteString("Woops! We ran into an issue!\n")
	b.WriteString(" parser errors:\n")
	for _, msg := range errors {
		b.WriteString("\t" + msg + "\n")
	}
	return b.String()
}
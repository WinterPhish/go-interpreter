package session

import (
	"strings"

	"selenite/evaluator"
	"selenite/lexer"
	"selenite/object"
	"selenite/parser"
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

// Run executes one line and returns what the REPL would have printed.
func (s *Session) Run(input string) string {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return parserErrors(p.Errors())
	}

	evaluator.DefineMacros(program, s.macroEnv)
	expanded := evaluator.ExpandMacros(program, s.macroEnv)

	evaluated := evaluator.Eval(expanded, s.env)
	if evaluated == nil {
		return ""
	}
	return evaluated.Inspect()
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
package environment

import (
	"fmt"

	"github.com/doeg/golox/golox/token"
)

type Environment struct {
	env map[string]any
}

func New() *Environment {
	return &Environment{
		env: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.env[name] = value
}

func (e *Environment) Get(name *token.Token) (any, error) {
	val, ok := e.env[name.Lexeme]
	if !ok {
		// TODO return a Lox runtime error
		// We return a runtime error instead of a static error since
		// referring to a variable is different from using it; if we make it
		// a static error to _mention_ a variable before it's been declared,
		// it becomes much harder to define recursive functions.
		// So instead, we defer the error to be a runtime error. It's ok
		// to refer to a variable before it's defined as long as you don't
		// **evaluate** the reference.
		return nil, fmt.Errorf("undefined variable %s", name.Lexeme)
	}

	return val, nil
}

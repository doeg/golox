package environment

import (
	"fmt"

	"github.com/doeg/golox/golox/token"
)

type Environment struct {
	// values is a mapping from variable name to value.
	// Note that we use strings as keys instead of tokens, as tokens
	// are entities that exist in a specific place in the source,
	// whereas variables with the same name (in the same scope)
	// should always share the same value.
	values map[string]any
}

func New() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (env *Environment) Define(name string, value any) {
	// Note that we don't check if the variable is already defined
	// when mapping a variable name to a value. In other words,
	// variable declarations can overwrite existing map entries.
	env.values[name] = value
}

func (env *Environment) Get(name *token.Token) (any, error) {
	val, ok := env.values[name.Lexeme]
	if !ok {
		// Note that trying to access an undefined variable name is a runtime
		// error rather than a syntax error. If we made it a static/syntax error
		// to refer to a variable before its used, then that would be pretty inflexible
		// and would make it tough to write recursive functions.
		//
		// TODO explicitly make this a Lox runtime error type.
		return nil, fmt.Errorf("undefined variable %s", name.Lexeme)
	}

	return val, nil
}

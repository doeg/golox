package environment

type Environment struct {
	env map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.env[name] = value
}

func (e )

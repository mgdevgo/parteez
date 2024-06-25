package application

type Environment struct {
	Name      string
	Arguments []string
}

var production = Environment{Name: "production"}
var development = Environment{Name: "development"}
var testing = Environment{}.init("testing")

func (env Environment) init(name string) Environment {
	env.Name = name
	return env
}

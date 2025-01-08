package server

type Environment struct {
	name string
}

var (
	Production  = Environment{name: "production"}
	Development = Environment{name: "development"}
	Testing     = Environment{name: "testing"}
)

func Detect(from []string) Environment {
	var env string
	for i, arg := range from {
		if arg == "-e" || arg == "--env" {
			env = from[i]
			break
		}
	}

	switch env {
	case "prod", "production":
		return Production
	case "dev", "development":
		return Development
	case "test", "testing":
		return Testing
	default:
		return Development
	}
}

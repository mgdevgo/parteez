package location

type Location struct {
	ID        int
	Name      string
	Type      string
	Address   string
	Latitude  string
	Longitude string
	Stages    []string
	Metro     []string
}

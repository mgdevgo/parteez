package scraping

type Event struct {
	Title          string
	Description    string
	AgeRestriction string
	LineUp         []string
	Genres         []string
	StartDate      string
	EndDate        string
	ArtworkURL     string
	TicketsURL     string
}

type Venue struct {
	Name          string
	Address       string
	MetroStations []string
}

type Result struct {
	Event  Event
	Venue  Venue
	Errors []error
}

func NewResult(event Event, errors []error) *Result {
	return &Result{
		Event:  event,
		Errors: errors,
	}
}

func (r *Result) IsFailure() bool {
	return len(r.Errors) > 0
}

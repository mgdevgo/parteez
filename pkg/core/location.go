package core

type LocationType string

const (
	LocationTypeUnknown     LocationType = "UNKNOWN"
	LocationTypeClub        LocationType = "CLUB"
	LocationTypeBar         LocationType = "BAR"
	LocationTypeCafe        LocationType = "CAFE"
	LocationTypeConcertHall LocationType = "CONCERT_HALL"
	LocationTypeSpace       LocationType = "SPACE"
)

type Location struct {
	ID                   int          `json:"id"`
	Name                 string       `json:"name"`
	Type                 LocationType `json:"type"`
	Description          string       `json:"description"`
	ArtworkURL           string       `json:"artworkURL"`
	Stages               []string     `json:"stages"`
	Address              string       `json:"address"`
	NearestMetroStations []string     `json:"nearestMetroStations"`
	// TODO
	// Latitude  string
	// Longitude string
	IsPublic  bool `json:"-"`
	Timestamp `json:"timestamp"`

	// Timestamp
	// CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	// UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

type ListLocationResponse struct {
	Location
	UpcomingEvents []Event `json:"upcomingEvents"`
	PassedEvents   []Event `json:"passedEvents"`
}

var LocationStagesDefault = []string{"main"}

func NewLocation(name string) *Location {
	// if len(stages) <= 0 {
	// 	stages = StagesDefault
	// } else {
	// 	for i := 0; i < len(stages); i++ {
	// 		stages[i] = strings.ToUpper(stages[i])
	// 	}
	// }
	// location := new(Location)
	// location.CreatedAt, location.UpdatedAt = NewTimestamp2()
	return &Location{
		Name:                 name,
		Type:                 LocationTypeUnknown,
		Stages:               []string{"main"},
		NearestMetroStations: make([]string, 0),
	}
}

func (l *Location) IsEmpty() bool {
	if l.ID == 0 {
		return true
	}
	return false
}

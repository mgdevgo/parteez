package types

type LocationType string

const (
	LocationTypeUnknown     LocationType = "UNKNOWN"
	LocationTypeClub        LocationType = "CLUB"
	LocationTypeBar         LocationType = "BAR"
	LocationTypeCafe        LocationType = "CAFE"
	LocationTypeConcertHall LocationType = "CONCERT_HALL"
	LocationTypeSpace       LocationType = "SPACE"
)

type Venue struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Type          LocationType `json:"type"`
	Description   string       `json:"description"`
	ArtworkID     string       `json:"-"`
	Stages        []string     `json:"stages"`
	Address       string       `json:"address"`
	MetroStations []string     `json:"metroStations"`
	// TODO
	// Latitude  string
	// Longitude string
	IsPublic bool `json:"-"`
	Timestamp
}

type ListLocationResponse struct {
	Venue
	UpcomingEvents []Event `json:"upcomingEvents"`
	PassedEvents   []Event `json:"passedEvents"`
}

var LocationStagesDefault = []string{"main"}

func NewLocation(name string) *Venue {
	// if len(stages) <= 0 {
	// 	stages = StagesDefault
	// } else {
	// 	for i := 0; i < len(stages); i++ {
	// 		stages[i] = strings.ToUpper(stages[i])
	// 	}
	// }
	// location := new(Venue)
	// location.CreatedAt, location.UpdatedAt = NewTimestamp2()
	return &Venue{
		Name:          name,
		Type:          LocationTypeUnknown,
		Stages:        []string{"main"},
		MetroStations: make([]string, 0),
	}
}

func (l *Venue) IsEmpty() bool {
	if l.ID == 0 {
		return true
	}
	return false
}

package content

import (
	"time"
)

type Data struct {
	Tittle         string   `json:"title" db:"title"`
	Description    string   `json:"description" db:"description"`
	LineUp         []string `json:"lineUp" db:"-"`
	Genres         []string `json:"genres"`
	AgeRestriction int      `json:"ageRestriction" db:"age_restriction"`
	Promoter       string   `json:"promoter" db:"promoter"`

	Date      time.Time `json:"date" db:"date"`
	StartTime time.Time `json:"startTime" db:"start_time"`
	EndTime   time.Time `json:"endTime" db:"end_time"`

	TicketsURL string `json:"ticketsUrl" db:"tickets_url"`

	ArtworkURL string `json:"artwork_url"`

	LocationName          string
	LocationAddress       string
	LocationMetroStations []string
}

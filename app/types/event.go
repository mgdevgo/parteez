package types

import (
	"time"
)

type EventStatus struct {
	Status string
}

var (
	EventStatusReview    = EventStatus{"REVIEW"}
	EventStatusEditing   = EventStatus{"EDITING"}
	EventStatusPublished = EventStatus{"PUBLISHED"}
)

const (
	EventDefaultAgeRestriction = 18
	EventDefaultLocationID     = -1
)

type Event struct {
	ID             int      `json:"id" db:"id"`
	Tittle         string   `json:"title" db:"title"`
	Description    string   `json:"description" db:"description"`
	LineUp         []LineUp `json:"lineUp" db:"-"`
	GenreIDs       []int    `json:"genreIDs" db:"-"`
	AgeRestriction int      `json:"ageRestriction" db:"age_restriction"`
	Promoter       string   `json:"promoter" db:"promoter"`

	Date      time.Time `json:"date" db:"date"`
	StartTime time.Time `json:"startTime" db:"start_time"`
	EndTime   time.Time `json:"endTime" db:"end_time"`

	TicketsURL string   `json:"ticketsUrl" db:"tickets_url"`
	Tickets    []Ticket `json:"tickets" db:"-"`

	ArtworkID int `json:"-" db:"artwork_id"`
	VenueID   int `json:"-" db:"venue_id"`

	IsPublic bool `json:"isPublic" db:"is_public"`
	Timestamp
}

type Promotion struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type Promotions []Promotion

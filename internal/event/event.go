package event

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Status string

const (
	StatusReview    Status = "review"
	StatusEditing   Status = "editing"
	StatusPublished Status = "published"
)

type Event struct {
	ID          ulid.ULID           `json:"id"`
	Name        string              `json:"name"`
	ImageURL    string              `json:"image_url"`
	Description string              `json:"description"`
	Genres      []string            `json:"genres"`
	LineUp      map[string][]Artist `json:"line_up"`
	StartTime   time.Time           `json:"start_time"`
	EndTime     time.Time           `json:"end_time"`
	MinAge      int                 `json:"min_age"`

	TicketsURL string         `json:"tickets_url"`
	Price      map[string]int `json:"price"`

	Interested int `json:"interested"`

	LocationID int    `json:"location_id"`
	Promoter   string `json:"promoter"`

	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Artist struct {
	ID         int
	Name       string
	ProfileURL string
}

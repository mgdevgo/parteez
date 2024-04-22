package event

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Status string

const (
	StatusReview    Status = "review"
	StatusEditing   Status = "editing"
	StatusPublished Status = "published"
)

const DefaultMinAge = 18

func NewEvent(name string) Event {
	timeNow := time.Now()
	return Event{
		Name:   name,
		Genres: make(Genres, 0),
		LineUp: make(LineUp),
		// TODO: default date
		// StartDate:      time.Time{},
		// EndDate:        time.Time{},
		AgeRestriction: 18,
		Price:          make(Price),
		LocationID:     -1,
		CreatedAt:      timeNow,
		UpdatedAt:      timeNow,
	}
}

type Event struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ArtworkURL     string    `json:"artworkURL"`
	Genres         Genres    `json:"genres"`
	LineUp         LineUp    `json:"lineUp"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AgeRestriction int       `json:"ageRestriction"`
	TicketsURL     string    `json:"ticketsURL"`
	Price          Price     `json:"price,omitempty"`
	LocationID     int       `json:"placeId,omitempty"`
	Promoter       string    `json:"promoter,omitempty"`
	IsPublic       bool      `json:"isPublic,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type Date struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type LineUpNew []struct {
	Stage      string      `json:"stage"`
	Performers []Performer `json:"performers"`
}

type LineUp map[string][]Performer

type Performer struct {
	Name  string `json:"name"`
	Start string `json:"start,omitempty"`
	Live  bool   `json:"live,omitempty"`
}

func (l LineUp) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *LineUp) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("can't scan %T into []byte", value)

	}

	return json.Unmarshal(bytes, &l)
}

type Price map[string]int

func (p Price) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Price) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("price: type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &p)
}

type Genres []string

func (g Genres) Value() (driver.Value, error) {
	return strings.Join(g, ","), nil
}

func (g *Genres) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	*g = strings.Split(s, ",")
	return nil
}

package core

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

func NewEvent(name string) *Event {
	return &Event{
		Name:   name,
		Genres: make(Genres, 0),
		LineUp: make([]LineUpStage, 0),
		// TODO: what is default start/end date?
		AgeRestriction: EventDefaultAgeRestriction,
		LocationID:     EventDefaultLocationID,
		Timestamp:      NewTimestamp(),
	}
}

type Event struct {
	ID             int       `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	ArtworkURL     string    `json:"artworkURL" db:"artwork_url"`
	Genres         Genres    `json:"genres" db:"genres"`
	LineUp         LineUp    `json:"lineUp" db:"line_up"`
	StartDate      time.Time `json:"startDate" db:"start_date"`
	EndDate        time.Time `json:"endDate" db:"end_date"`
	AgeRestriction int       `json:"ageRestriction" db:"age_restriction"`
	Promoter       string    `json:"promoter" db:"promoter"`
	IsPublic       bool      `json:"-" db:"is_public"`
	Tickets        Tickets   `json:"tickets" db:"tickets"`
	LocationID     int       `json:"-" db:"location_id"`
	Location       *Location `json:"location" db:"-"`
	Timestamp
}

type LineUp []LineUpStage

type LineUpStage struct {
	Stage   string         `json:"stage"`
	Artists []LineUpArtist `json:"artists"`
}

type LineUpArtist struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	Live  bool   `json:"live"`
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

type Tickets struct {
	URL   string    `json:"url" db:"tickets_url"`
	Info  string    `json:"info"`
	Price PriceList `json:"price" db:"tickets_price"`
}

type PriceList []Price

type Price struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (p PriceList) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PriceList) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("price list: type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &p)
}

type Promotion struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type Promotions []Promotion

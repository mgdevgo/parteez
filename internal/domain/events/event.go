package events

import (
	"errors"
	"fmt"
	"time"

	"parteez/internal/domain/artwork"
	"parteez/internal/domain/venue"
)

const (
	DEFAULT_AGE_RESTRICTION = 18
)

var (
	ErrEvent              = errors.New("event")
	ErrEventID            = fmt.Errorf("%w: id", ErrEvent)
	ErrEventAlreadyExists = fmt.Errorf("%w: already exists", ErrEvent)
)

type EventID int

func NewEventID(input int) (EventID, error) {
	if input < 0 {
		return 0, ErrEventID
	}

	return EventID(input), nil
}

type Event struct {
	ID             EventID
	Title          string
	Description    string
	AgeRestriction int
	LineUp         LineUp
	Genres         []*Genre
	Promoter       string
	Date           Date
	TicketsURL     string
	Tickets        []Ticket
	ArtworkID      artwork.ArtworkID
	VenueID        venue.VenueID
	Status         Status
	CreatedAt      time.Time
	UpdatedAt      time.Time
	PublishedAt    time.Time
	ArchivedAt     time.Time
}

func NewEvent(id EventID, title, description string, date Date) (*Event, error) {
	if id < 0 {
		return nil, fmt.Errorf("%w: must be positive", ErrEventID)
	}

	now := time.Now()

	return &Event{
		ID:             id,
		Title:          title,
		Description:    description,
		AgeRestriction: DEFAULT_AGE_RESTRICTION,
		LineUp:         LineUp{},
		Genres:         make([]*Genre, 0),
		Promoter:       "",
		Date:           date,
		TicketsURL:     "",
		Tickets:        []Ticket{},
		ArtworkID:      0,
		VenueID:        0,
		Status:         StatusDraft,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (event *Event) AddCover(artworkId artwork.ArtworkID) {
	event.ArtworkID = artworkId
}

func (event *Event) AddVenue(venueId venue.VenueID) {
	event.VenueID = venueId
}

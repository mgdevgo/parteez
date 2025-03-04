package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"parteez/internal/domain/events"
)

type eventRow struct {
	ID             int                     `db:"id"`
	Title          string                  `db:"title"`
	Description    string                  `db:"description"`
	Date           string                  `db:"date"`
	DateRange      pgtype.Range[time.Time] `db:"date"`
	Genres         string                  `db:"genres"`
	LineUp         []byte                  `db:"lineup"`
	ArtworkID      pgtype.Int4             `db:"artwork_id"`
	VenueID        pgtype.Int4             `db:"venue_id"`
	AgeRestriction int                     `db:"age_restriction"`
	Promoter       string                  `db:"promoter"`
	TicketsURL     string                  `db:"tickets_url"`
	Tickets        []byte                  `db:"tickets"`
	IsDraft        bool                    `db:"is_draft"`
	PublishedAt    pgtype.Timestamp        `db:"published_at"`
	CreatedAt      time.Time               `db:"created_at"`
	UpdatedAt      time.Time               `db:"updated_at"`
}

func eventToRow(event *events.Event) (eventRow, error) {
	date := fmt.Sprintf("[%s, %s]",
		event.Date.Start.Format(time.DateTime),
		event.Date.End.Format(time.DateTime),
	)

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("{")
	for i, genre := range event.Genres {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(genre.Name)
	}
	buf.WriteString("}")
	genres := buf.String()

	artworkID := int(event.ArtworkID)
	venueID := int(event.VenueID)

	tickets, err := json.MarshalIndent(event.Tickets, "", "  ")
	if err != nil {
		return eventRow{}, err
	}

	var draft bool
	if event.Status == events.StatusDraft {
		draft = true
	}

	lineup, err := json.MarshalIndent(event.LineUp, "", "  ")
	if err != nil {
		return eventRow{}, err
	}

	return eventRow{
		ID:             int(event.ID),
		Title:          event.Title,
		Description:    event.Description,
		Date:           date,
		Genres:         genres,
		LineUp:         lineup,
		ArtworkID:      pgtype.Int4{Int32: int32(artworkID), Valid: artworkID != 0},
		VenueID:        pgtype.Int4{Int32: int32(venueID), Valid: venueID != 0},
		AgeRestriction: event.AgeRestriction,
		Promoter:       event.Promoter,
		TicketsURL:     event.TicketsURL,
		Tickets:        tickets,
		IsDraft:        draft,
		PublishedAt:    pgtype.Timestamp{Time: event.PublishedAt, Valid: !event.PublishedAt.IsZero()},
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}, nil
}

func rowToEvent(row eventRow) events.Event {
	return events.Event{
		ID:             events.EventID(row.ID),
		Title:          row.Title,
		Description:    row.Description,
		AgeRestriction: 0,
		LineUp:         events.LineUp{},
		Genres:         []*events.Genre{},
		Promoter:       "",
		Date:           events.Date{},
		TicketsURL:     "",
		Tickets:        []events.Ticket{},
		ArtworkID:      0,
		VenueID:        0,
		Status:         "",
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		PublishedAt:    time.Time{},
		ArchivedAt:     time.Time{},
	}
}

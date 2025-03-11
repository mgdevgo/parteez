package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"parteez/internal/domain/artwork"
	"parteez/internal/domain/events"
	"parteez/internal/domain/venue"
)

type eventRow struct {
	ID             int    `db:"id"`
	Title          string `db:"title"`
	Description    string `db:"description"`
	Date           string
	DateRange      pgtype.Range[time.Time] `db:"date"`
	Genres         []string                `db:"genres"`
	LineUp         []byte                  `db:"lineup"`
	ArtworkID      pgtype.Int4             `db:"artwork_id"`
	VenueID        pgtype.Int4             `db:"venue_id"`
	AgeRestriction int                     `db:"age_restriction"`
	Promoter       string                  `db:"promoter"`
	TicketsURL     string                  `db:"tickets_url"`
	Tickets        []byte                  `db:"tickets"`
	IsDraft        bool                    `db:"is_draft"`
	PublishedAt    pgtype.Timestamp        `db:"published_at"`
	// ArchivedAt     pgtype.Timestamp        `db:"archived_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func eventToRow(event *events.Event) (eventRow, error) {
	date := fmt.Sprintf("[%s, %s]",
		event.Date.Start.Format(time.DateTime),
		event.Date.End.Format(time.DateTime),
	)

	genres := make([]string, 0)
	for _, genre := range event.Genres {
		genres = append(genres, string(genre))
	}

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
		// ArchivedAt:     pgtype.Timestamp{Time: event.ArchivedAt, Valid: !event.ArchivedAt.IsZero()},
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}, nil
}

func rowToEvent(row eventRow) (*events.Event, error) {
	var lineup events.LineUp
	if err := json.Unmarshal(row.LineUp, &lineup); err != nil {
		return nil, err
	}

	genres := make([]events.EventGenre, 0)
	for _, genre := range row.Genres {
		genres = append(genres, events.EventGenre(genre))
	}

	var tickets []events.Ticket
	if err := json.Unmarshal(row.Tickets, &tickets); err != nil {
		return nil, err
	}

	var status events.Status
	if row.IsDraft {
		status = events.StatusDraft
	} else if row.PublishedAt.Valid {
		status = events.StatusPublished
	} else {
		status = events.StatusArchived
	}

	return &events.Event{
		ID:             events.EventID(row.ID),
		Title:          row.Title,
		Description:    row.Description,
		AgeRestriction: row.AgeRestriction,
		LineUp:         lineup,
		Genres:         genres,
		Promoter:       row.Promoter,
		Date: events.Date{
			Start: row.DateRange.Lower,
			End:   row.DateRange.Upper,
		},
		TicketsURL:  row.TicketsURL,
		Tickets:     tickets,
		ArtworkID:   artwork.ArtworkID(row.ArtworkID.Int32),
		VenueID:     venue.VenueID(row.VenueID.Int32),
		Status:      status,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		PublishedAt: row.PublishedAt.Time,
		// ArchivedAt:  row.ArchivedAt.Time,
	}, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"iditusi/internal/models"
	"iditusi/internal/repositories"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repositories.EventRepository = (*EventRepository)(nil)

type EventRepository struct {
	connection *pgxpool.Pool
}

func NewPostgresEventRepository(conn *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		connection: conn,
	}
}

type EventRecord struct {
	ID          int
	Name        string
	Description sql.NullString
}

func (r *EventRepository) Save(ctx context.Context, event models.Event) (models.Event, error) {
	isNew := event.ID == 0
	if !isNew {
		// TODO
		return models.Event{}, fmt.Errorf("can't update event - not implemented")
	}

	tx, err := r.connection.Begin(ctx)
	if err != nil {
		return models.Event{}, err
	}

	row := tx.QueryRow(
		ctx, `
INSERT INTO event 
(
	title,
	description,
	date,
	start_time,
	end_time,
	age_restriction,
	tickets_url,
	promoter,
	venue_id,
	artwork_id
)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id`,
		event.Tittle, event.Description,
		event.Date, event.StartTime,
		event.EndTime, event.AgeRestriction,
		event.TicketsURL, event.Promoter,
		event.VenueID, event.ArtworkID,
	)
	if err := row.Scan(event.ID); err != nil {
		// TODO: check
		// if err != nil {
		// 	var pgErr *pgconn.PgError
		// 	if errors.As(err, &pgErr) {
		// 		switch pgErr.Code {
		// 		case "23505":
		// 			return models.Event{}, repositories2.ErrorEventAlreadyExists
		// 		}
		// 	}
		// 	return models.Event{}, err
		// }
		return models.Event{}, err
	}

	if event.LineUp != nil {
		for _, lineup := range event.LineUp {
			var args []any
			var values string

			args = append(args, event.ID, lineup.Stage)

			for _, artist := range lineup.Artists {
				// current arg number
				argIndex := len(args)
				values += fmt.Sprintf("\n($1, $2, $%d, $%d, $%d)", argIndex, argIndex+1, argIndex+2)
				args = append(args, artist.Name, artist.Live, artist.StartAt)
			}

			_, err := tx.Exec(
				ctx, `
INSERT INTO lineup 
	(event_id, stage_name, artist_name, is_live, start_at) 
VALUES
	`+values, args...)

			if err != nil {
				return models.Event{}, fmt.Errorf("can't insert lineup: %w", err)
			}
		}
	}

	var values string
	var args []any

	args = append(args, event.ID)

	for i, genreID := range event.GenreIDs {
		values += fmt.Sprintf("\n($1, $%d)", i+2)
		args = append(args, genreID)
	}

	_, err = tx.Exec(ctx, `
INSERT INTO genres
	(event_id, genre_id)
VALUES 
	`+values, args...,
	)

	if err != nil {
		return models.Event{}, fmt.Errorf("can't insert genre: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func (r *EventRepository) FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]models.Event, error) {
	const query = "SELECT * FROM event WHERE start_date >= $1 AND end_date <= $2"
	// events := make([]models.Event, 0)

	rows, err := r.connection.Query(ctx, query, fromDate, toDate)
	if err != nil {
		return []models.Event{}, err
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Event])
	// err := pgxscan.Select(ctx, s.connection, &events, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []models.Event{}, err
	}
	return events, nil
}

func (r *EventRepository) FindAll(ctx context.Context) ([]models.Event, error) {
	const query = "SELECT * FROM event"
	// events := make([]models.Event, 0)

	rows, err := r.connection.Query(ctx, query)
	if err != nil {
		return []models.Event{}, err
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Event])
	// err := pgxscan.Select(ctx, s.connection, &events, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []models.Event{}, err
	}
	return events, nil
}

// func (r *EventRepository) FindAllSorted(ctx context.Context, fromDate time.Time, filter utils2.Filter) ([]models.Event, error) {
// 	const query = "SELECT * FROM event WHERE start_date >= $1 ORDER BY start_date ASC LIMIT $2 OFFSET $3"
// 	// events := make([]models.Event, 0)
//
// 	rows, err := r.connection.Query(ctx, query, fromDate, filter.Limit, filter.Offset)
// 	if err != nil {
// 		return []models.Event{}, err
// 	}
// 	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Event])
// 	// err := pgxscan.Select(ctx, s.connection, &events, query, params.FromDate, params.Limit, params.Offset)
// 	if err != nil {
// 		return []models.Event{}, err
// 	}
// 	return events, nil
// }

func (r *EventRepository) FindById(ctx context.Context, id int) (models.Event, error) {
	const query = "SELECT * FROM event WHERE id = $1"

	var event models.Event
	err := pgxscan.Get(ctx, r.connection, &event, query, id)
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}

func eventToRecord(event models.Event) EventRecord {
	return EventRecord{}
}

func (r *EventRepository) Delete(ctx context.Context, id int) error {
	_, err := r.connection.Exec(ctx, "DELETE FROM event WHERE id = $1", id)
	return err
}

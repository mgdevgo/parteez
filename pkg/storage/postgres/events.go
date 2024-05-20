package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"iditusi/pkg/core"
	"iditusi/pkg/storage/fileds"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type EventStorage struct {
	driver *pgxpool.Pool
	db     *sqlx.DB
	// .cache .cache.Cache
}

func NewEventStorage(db *pgxpool.Pool) *EventStorage {
	return &EventStorage{
		driver: db,
	}
}

//	type FindOptions struct {
//		limit   int
//		offset  int
//		orderBy string
//		date    time.Time
//	}
//
//	func (o *FindOptions) Fill() {
//		o.limit = 5
//		o.offset = 0
//		o.orderBy = "start_date"
//	}
//
//	func (o *FindOptions) SetLimit(limit int) {
//		if limit <= 0 {
//			return
//		}
//		o.limit = limit
//	}
//
//	func (o *FindOptions) SetOffset(offset int) {
//		if offset <= 0 {
//			return
//		}
//		o.offset = offset
//	}
//
//	func (o *FindOptions) SetDate(date time.Time) {
//		o.date = date
//	}
type EventRecord struct {
	ID          int
	Name        string
	Description sql.NullString
}

func (s *EventStorage) Find(ctx context.Context, params core.EventSearchParams) ([]core.Event, error) {
	const query = "SELECT * FROM event WHERE start_date >= $1 ORDER BY start_date ASC LIMIT $2 OFFSET $3"
	// events := make([]core.Event, 0)

	rows, err := s.driver.Query(ctx, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []core.Event{}, err
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[core.Event])
	// err := pgxscan.Select(ctx, s.driver, &events, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []core.Event{}, err
	}
	return events, nil
}

func (s *EventStorage) FindByID(ctx context.Context, id int) (core.Event, error) {
	const query = "SELECT * FROM event WHERE id = $1"

	var event core.Event
	err := pgxscan.Get(ctx, s.driver, &event, query, id)
	if err != nil {
		return core.Event{}, err
	}

	return event, nil
}

type EventRecord struct {
}

func eventToRecord(event core.Event) EventRecord {
	return EventRecord{}
}

func (s *EventStorage) Save(ctx context.Context, event core.Event) (core.Event, error) {
	isNew := event.ID == 0
	if !isNew {
		// TODO
		return core.Event{}, fmt.Errorf("can't update event - not implemented")
	}

	fields := fileds.Fields{}
	fields.AddField("name", event.Name)
	fields.AddField("description", event.Description)
	fields.AddField("artwork_url", event.ArtworkURL)
	fields.AddField("genre", strings.Join(event.Genres, ","))
	fields.AddField("start_time", event.StartDate)
	fields.AddField("end_time", event.EndDate)
	fields.AddField("age_restriction", event.AgeRestriction)
	fields.AddField("promoter", event.Promoter)
	fields.AddField("location_id", event.LocationID)
	fields.AddField("is_public", event.IsPublic)
	fields.AddField("line_up", event.LineUp)
	fields.AddField("tickets_url", event.Tickets.URL)
	fields.AddField("price", event.Tickets.Price)

	fieldNames, values, args := fields.Build()
	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
	query := fmt.Sprintf(sql, "event", fieldNames, values)

	result := s.driver.QueryRow(ctx, query, args...)
	err := result.Scan(&event.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return core.Event{}, ErrEventAlreadyExists
			}
		}
		return core.Event{}, err
	}
	return event, nil
}

func (s *EventStorage) Delete(ctx context.Context, id int) error {
	_, err := s.driver.Exec(ctx, "DELETE FROM event WHERE id = $1", id)
	return err
}

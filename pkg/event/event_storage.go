package event

import (
	"context"
	"fmt"
	"strings"
	"time"

	"iditusi/pkg/shared/storage"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Reader interface {
	FindByID(id int) (Event, error)
	Find(options FindOptions) ([]Event, error)
}

type Writer interface{}

type Storage interface {
	Reader
	Writer
}

type postgres struct {
	driver *pgxpool.Pool
	// cache cache.Cache
}

var _ Storage = (*postgres)(nil)

func NewStorage(db *pgxpool.Pool) *postgres {
	return &postgres{
		driver: db,
	}
}

type FindOptions struct {
	limit   int
	offset  int
	orderBy string
	date    time.Time
}

func (o *FindOptions) Fill() {
	o.limit = 5
	o.offset = 0
	o.orderBy = "start_date"
}

func (o *FindOptions) SetLimit(limit int) {
	if limit <= 0 {
		return
	}
	o.limit = limit
}

func (o *FindOptions) SetOffset(offset int) {
	if offset <= 0 {
		return
	}
	o.offset = offset
}

func (o *FindOptions) SetDate(date time.Time) {
	o.date = date
}

func (s *postgres) Find(options FindOptions) ([]Event, error) {
	const query = "SELECT * FROM event WHERE start_date >= $1 ORDER BY start_date ASC LIMIT $2 OFFSET $3"
	events := make([]Event, 0)
	err := pgxscan.Select(context.TODO(), s.driver, &events, query, options.date, options.limit, options.offset)
	if err != nil {
		return []Event{}, err
	}
	return events, nil
}

func (s *postgres) FindByID(id int) (Event, error) {
	const query = "SELECT * FROM event WHERE id = $1"

	var event Event
	err := pgxscan.Get(context.TODO(), s.driver, &event, query, id)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func (s *postgres) Save(event Event) (Event, error) {
	// isNew := event.ID
	// if !isNew{
	// 	// TODO: Update
	// }

	fields := storage.Fields{}
	fields.AddField("name", event.Name)
	fields.AddField("description", event.Description)
	fields.AddField("artwork_url", event.ArtworkURL)
	fields.AddField("genre", strings.Join(event.Genres, ","))
	fields.AddField("start_time", event.StartDate)
	fields.AddField("end_time", event.EndDate)
	fields.AddField("tickets_url", event.TicketsURL)
	fields.AddField("age_restriction", event.AgeRestriction)
	fields.AddField("promoter", event.Promoter)
	fields.AddField("location_id", event.LocationID)
	fields.AddField("is_public", event.IsPublic)
	fields.AddField("line_up", event.LineUp)
	fields.AddField("price", event.Price)

	fieldNames, values, args := fields.Build()
	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
	query := fmt.Sprintf(sql, "event", fieldNames, values)

	result := s.driver.QueryRow(context.Background(), query, args...)
	err := result.Scan(&event.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return Event{}, ErrEventAlreadyExists
			}
		}
		return Event{}, err
	}
	return event, nil
}

func (s *postgres) Update(id int, options map[string]any) (Event, error) {
	// TODO implement me
	panic("implement me")
}

func (s *postgres) Delete(id int) error {
	// TODO implement me
	panic("implement me")
}

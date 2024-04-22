package event

import (
	"context"
	"fmt"
	"strings"
	"time"

	"iditusi/pkg/sqlutils"

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

type EventStorage interface {
	Reader
	Writer
}

type storage struct {
	driver *pgxpool.Pool
	// cache cache.Cache
}

var _ EventStorage = (*storage)(nil)

func NewStorage(db *pgxpool.Pool) *storage {
	return &storage{
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

func (s *storage) Find(options FindOptions) ([]Event, error) {
	const query = "SELECT * FROM event WHERE start_date >= $1 ORDER BY start_date ASC LIMIT $2 OFFSET $3"
	events := make([]Event, 0)
	err := pgxscan.Select(context.TODO(), s.driver, &events, query, options.date, options.limit, options.offset)
	if err != nil {
		return []Event{}, err
	}
	return events, nil
}

func (s *storage) FindByID(id int) (Event, error) {
	const query = "SELECT * FROM event WHERE id = $1"
	// rows, err := s.db.Query(ctx, query, ids)
	// if err != nil {
	//	return empty, err
	// }
	// defer rows.Close()
	//
	// for rows.Next() {
	//	var event Event
	//	err = rows.Scan(&event.ID, &event.Name)
	//	if err != nil {
	//		s.log.Error("Row scan", slog.String("Error", err.Error()))
	//	}
	//	events = append(events, event)
	// }

	var err error
	var event Event
	err = pgxscan.Get(context.TODO(), s.driver, &event, query, id)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func (s *storage) Save(event Event) (Event, error) {
	// isNew := event.ID
	// if !isNew{
	// 	// TODO: Update
	// }

	builder := sqlutils.NewFieldsBuilder()
	builder.AddField("name", event.Name)
	builder.AddField("description", event.Description)
	builder.AddField("artwork_url", event.ArtworkURL)
	builder.AddField("genre", strings.Join(event.Genres, ","))
	builder.AddField("start_time", event.StartDate)
	builder.AddField("end_time", event.EndDate)
	builder.AddField("tickets_url", event.TicketsURL)
	builder.AddField("age_restriction", event.AgeRestriction)
	builder.AddField("promoter", event.Promoter)
	builder.AddField("location_id", event.LocationID)
	builder.AddField("is_public", event.IsPublic)
	builder.AddField("line_up", event.LineUp)
	builder.AddField("price", event.Price)

	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
	query := fmt.Sprintf(sql, "event", builder.Fields(), builder.Values())

	result := s.driver.QueryRow(context.Background(), query, builder.Args()...)
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

func (s *storage) Update(id int, options map[string]any) (Event, error) {
	// TODO implement me
	panic("implement me")
}

func (s *storage) Delete(id int) error {
	// TODO implement me
	panic("implement me")
}

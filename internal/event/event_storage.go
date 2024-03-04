package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event exists")
)

type EventStorage interface {
	Create(event Event) (string, error)
	Get(ids []string) ([]Event, error)
	Update(id string, options map[string]any) (Event, error)
	Delete(id string) error
	Transaction() (int, error)
}

type postgresEventStorage struct {
	eventsTable string
	driver      *pgxpool.Pool
}

func NewEventStorage(db *pgxpool.Pool) *postgresEventStorage {
	return &postgresEventStorage{
		eventsTable: "event",
		driver:      db,
	}
}

func (s *postgresEventStorage) Create(event Event) (string, error) {
	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"

	const maxNumberOfFields = 12
	fields := make([]string, 0, maxNumberOfFields)
	values := make([]string, 0, maxNumberOfFields)
	args := make([]any, 0, maxNumberOfFields)
	count := 0

	null := func() bool { return true }
	addField := func(name string, value any, isNull ...bool) {
		count++
		fields = append(fields, name)

		placeholder := "$%d"
		if name == "location_id" {
			placeholder = "(SELECT id FROM location WHERE name = $%d)"
		}

		values = append(values, fmt.Sprintf(placeholder, count))
		if len(isNull) > 0 {
			if v := fmt.Sprintf("%v", value); v == "" || v == "0001-01-01 00:00:00 +0000 UTC" {
				value = nil
			}
		}

		args = append(args, value)
	}

	addField("id", event.ID)
	addField("name", event.Name)
	addField("description", event.Description, null())
	addField("image_url", event.ImageURL, null())
	addField("music_genres", strings.Join(event.MusicGenres, ","), null())

	lbytes, err := json.Marshal(event.LineUp)
	if err != nil {
		return "", err
	}
	addField("line_up", string(lbytes), null())

	addField("start_time", event.StartTime, null())
	addField("end_time", event.EndTime, null())

	pbytes, err := json.Marshal(event.Price)
	if err != nil {
		return "", err
	}
	addField("price", string(pbytes), null())
	addField("tickets_url", event.TicketsURL, null())
	addField("min_age", event.MinAge)

	addField("promoter", event.Promoter, null())
	addField("location_id", event.LocationID, null())

	addField("is_public", event.IsPublic)

	query := fmt.Sprintf(sql, s.eventsTable, strings.Join(fields, ", "), strings.Join(values, ", "))

	fmt.Println(query)
	result := s.driver.QueryRow(context.Background(), query, args...)

	var eventID string
	if err := result.Scan(&eventID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return "", ErrEventExists
			}
		}
		return "", err
	}

	return eventID, nil
}

func (s *postgresEventStorage) Get(ids []string) ([]Event, error) {
	empty := make([]Event, 0)
	const query = "SELECT id, name, line_up, created_at FROM event WHERE id = ANY($1)"
	//rows, err := s.db.Query(ctx, query, ids)
	//if err != nil {
	//	return empty, err
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var event Event
	//	err = rows.Scan(&event.ID, &event.Name)
	//	if err != nil {
	//		s.log.Error("Row scan", slog.String("Error", err.Error()))
	//	}
	//	events = append(events, event)
	//}

	var err error
	events := make([]Event, 0)
	err = pgxscan.Select(context.TODO(), s.driver, &events, query, ids)
	if err != nil {
		return empty, err
	}

	return events, nil
}

func (s *postgresEventStorage) Update(id string, options map[string]any) (Event, error) {
	//TODO implement me
	panic("implement me")
}

func (s *postgresEventStorage) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *postgresEventStorage) Transaction() (int, error) {
	//TODO implement me
	panic("implement me")
}

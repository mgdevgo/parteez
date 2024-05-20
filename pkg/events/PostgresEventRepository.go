package events

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"iditusi/pkg/core"
	"iditusi/pkg/core/repository"
	"iditusi/pkg/storage"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrEventAlreadyExists = errors.New("event already exists")
)

type PostgresEventRepository struct {
	connection *pgxpool.Pool
}

var _ EventRepository = (*PostgresEventRepository)(nil)

func NewPostgresEventStorage(conn *pgxpool.Pool) *PostgresEventRepository {
	return &PostgresEventRepository{
		connection: conn,
	}
}

type EventRecord struct {
	ID          int
	Name        string
	Description sql.NullString
}

func (s *PostgresEventRepository) FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]core.Event, error) {
	const query = "SELECT * FROM event WHERE start_date >= $1 AND end_date <= $2"
	// events := make([]core.Event, 0)

	rows, err := s.connection.Query(ctx, query, fromDate, toDate)
	if err != nil {
		return []core.Event{}, err
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[core.Event])
	// err := pgxscan.Select(ctx, s.connection, &events, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []core.Event{}, err
	}
	return events, nil
}

func (s *PostgresEventRepository) FindAll(ctx context.Context) ([]core.Event, error) {
	const query = "SELECT * FROM event"
	// events := make([]core.Event, 0)

	rows, err := s.connection.Query(ctx, query)
	if err != nil {
		return []core.Event{}, err
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[core.Event])
	// err := pgxscan.Select(ctx, s.connection, &events, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []core.Event{}, err
	}
	return events, nil
}

func (s *PostgresEventRepository) FindAllSorted(ctx context.Context, fromDate time.Time, filter repository.Filter) ([]core.Event, error) {
	const query = "SELECT * FROM event WHERE start_date >= $1 ORDER BY start_date ASC LIMIT $2 OFFSET $3"
	// events := make([]core.Event, 0)

	rows, err := s.connection.Query(ctx, query, fromDate, filter.Limit, filter.Offset)
	if err != nil {
		return []core.Event{}, err
	}
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[core.Event])
	// err := pgxscan.Select(ctx, s.connection, &events, query, params.FromDate, params.Limit, params.Offset)
	if err != nil {
		return []core.Event{}, err
	}
	return events, nil
}

func (s *PostgresEventRepository) FindById(ctx context.Context, id int) (core.Event, error) {
	const query = "SELECT * FROM event WHERE id = $1"

	var event core.Event
	err := pgxscan.Get(ctx, s.connection, &event, query, id)
	if err != nil {
		return core.Event{}, err
	}

	return event, nil
}

func eventToRecord(event core.Event) EventRecord {
	return EventRecord{}
}

func (s *PostgresEventRepository) Save(ctx context.Context, event core.Event) (core.Event, error) {
	isNew := event.ID == 0
	if !isNew {
		// TODO
		return core.Event{}, fmt.Errorf("can't update event - not implemented")
	}

	genres := strings.Join(event.Genres, ",")

	builder := storage.SQLBuilder{}
	builder.WriteLine("INSERT INTO event (name, description, artwork_url, genre, start_time, end_time, age_restriction, promoter, location_id, is_public, line_up, tickets_url, tickets_price)")
	builder.WriteLine("VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", event.Name, event.Description,
		event.ArtworkURL, genres, event.StartDate, event.EndDate, event.AgeRestriction, event.Promoter,
		event.LocationID, event.IsPublic, event.LineUp, event.Tickets.URL, event.Tickets.Price)
	builder.WriteLine("RETURNING id")

	// fields := fileds.Fields{}
	// fields.AddField("name", event.Name)
	// fields.AddField("description", event.Description)
	// fields.AddField("artwork_url", event.ArtworkURL)
	// fields.AddField("genre", strings.Join(event.Genres, ","))
	// fields.AddField("start_time", event.StartDate)
	// fields.AddField("end_time", event.EndDate)
	// fields.AddField("age_restriction", event.AgeRestriction)
	// fields.AddField("promoter", event.Promoter)
	// fields.AddField("location_id", event.LocationID)
	// fields.AddField("is_public", event.IsPublic)
	// fields.AddField("line_up", event.LineUp)
	// fields.AddField("tickets_url", event.Tickets.URL)
	// fields.AddField("price", event.Tickets.Price)

	// fieldNames, values, args := fields.Build()
	// const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
	// query := fmt.Sprintf(sql, "event", fieldNames, values)

	// result := s.connection.QueryRow(ctx, query, args...)
	result := s.connection.QueryRow(ctx, builder.String(), builder.Params()...)
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

func (s *PostgresEventRepository) Delete(ctx context.Context, id int) error {
	_, err := s.connection.Exec(ctx, "DELETE FROM event WHERE id = $1", id)
	return err
}

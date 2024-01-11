package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	tx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event exists")
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
	tx  tx.CtxGetter
}

func EventStorage(db *pgxpool.Pool, log *slog.Logger) *Storage {
	return &Storage{db: db, log: log}
}

func (s *Storage) Save(ctx context.Context, event Event) (string, error) {
	var eventID string
	const op = "event.Storage.Save"

	price, err := json.Marshal(event.Price)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	lineUpBytes, err := json.Marshal(event.LineUp)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	const sql = "INSERT INTO event (id, name, image_url, description, start_time, end_time, line_up, location_id, promoter, tickets_url, price, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id"
	row := s.db.QueryRow(context.Background(), sql,
		event.ID.String(),
		event.Name,
		event.ImageURL,
		event.Description,
		event.StartTime,
		event.EndTime,
		lineUpBytes,
		nil,
		event.Promoter,
		event.TicketsURL,
		price,
		event.Status,
	)
	if err := row.Scan(&eventID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
			switch pgErr.Code {
			case "42601":
				return "", fmt.Errorf("%s: %w", op, ErrEventExists)
			}
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return eventID, nil
}

func (s *Storage) Get(ctx context.Context, ids []string) ([]Event, error) {
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
	err = pgxscan.Select(ctx, s.db, &events, query, ids)
	if err != nil {
		return empty, err
	}

	return events, nil
}

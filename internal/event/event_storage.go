package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	tx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event exists")
)

type storage struct {
	*pgxpool.Pool
	log *slog.Logger
	tx  tx.CtxGetter
}

func NewStorage(db *pgxpool.Pool, log *slog.Logger) *storage {
	return &storage{Pool: db, log: log}
}

type record struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	ImageURL    string    `db:"image_url"`
	Description string    `db:"description"`
	StartsAt    time.Time `db:"starts_at"`
	EndsAt      time.Time `db:"ends_at"`
	LineUp      []byte    `db:"line_up"`
	MinAge      int       `db:"min_age"`

	TicketsURL string `db:"tickets_url"`
	Price      []byte `db:"price"`

	LocationID int `db:"location_id"`

	Status   string    `db:"status"`
	UpdateAt time.Time `db:"updated_at"`
}

func (s *storage) Save(ctx context.Context, event Event) (int, error) {
	var eventID int
	const op = "event.storage.Save"

	price, err := json.Marshal(event.Price)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	lineup, err := json.Marshal(event.LineUp)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	row := s.QueryRow(
		context.Background(),
		"INSERT INTO event (name, image_url, description, starts_at, ends_at, line_up, location_id, promoter, tickets_url, price, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		event.Name,
		event.ArtworkURL,
		event.Description,
		event.StartsAt,
		event.EndsAt,
		lineup,
		event.Location.ID,
		event.Promoter,
		event.TicketsURL,
		price,
		StatusEditing,
	)
	if err := row.Scan(&eventID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
			switch pgErr.Code {
			case "42601":
				return 0, fmt.Errorf("%s: %w", op, ErrEventExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return eventID, nil
}

package postgres

import (
	"parteez/internal/domain/events"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GenreStorage struct {
	connection *pgxpool.Pool
	cache      map[string]*events.Genre
}

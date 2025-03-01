package postgres

import (
	"parteez/internal/domain/event"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GenreStorage struct {
	connection *pgxpool.Pool
	cache      map[string]*event.Genre
}

package postgres

import (
	"context"
	"fmt"
	"iditusi/internal/storage"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool  *pgxpool.Pool
	Users storage.UserStorage
}

func New(URL string) (*pgxpool.Pool, error) {
	const op = "storage.postgres.New"

	ctx := context.TODO()

	db, err := pgxpool.New(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pingContext, cancelPing := context.WithTimeout(ctx, time.Second*2)
	defer cancelPing()

	if err := db.Ping(pingContext); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

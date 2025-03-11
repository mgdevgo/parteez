package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(URL string) (*pgxpool.Pool, error) {
	const op = "pkg.postgres.New"

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

	if err := runMigrations(URL); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

// RunMigrations runs all pending migrations from the specified migrations directory using an existing database connection
func runMigrations(URL string) error {
	const op = "pkg.postgres.runMigrations"
	m, err := migrate.New(
		"file://migrations",
		strings.Replace(URL, "postgres://", "pgx5://", 1),
	)
	if err != nil {
		return fmt.Errorf("%s: failed to create migrate instance: %w", op, err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("%s: failed to run migrations: %w", op, err)
		}
	}

	return nil
}

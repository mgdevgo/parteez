package postgres

import (
	"github.com/jmoiron/sqlx"
)

type Database struct {
}

type Configuration struct {
}

func Connect(URL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", URL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"parteez/internal/event"
)

type Storage struct {
	events *event.Storage
}

func NewStorage(db *pgxpool.Pool, log *slog.Logger) *Storage {
	return &Storage{
		events: event.EventStorage(db, log),
	}
}

//func (s *storage) Events() *event.Storage {
//	return s.events
//}

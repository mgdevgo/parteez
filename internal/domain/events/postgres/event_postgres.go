package postgres

import (
	pgxtrm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventStorage struct {
	pool    *pgxpool.Pool
	context *pgxtrm.CtxGetter
}

func NewEventStorage(pool *pgxpool.Pool, context *pgxtrm.CtxGetter) *EventStorage {
	return &EventStorage{
		pool:    pool,
		context: context,
	}
}

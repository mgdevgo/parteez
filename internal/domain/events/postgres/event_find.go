package postgres

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"parteez/internal/domain/events"
)

func (s *EventStorage) FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*events.Event, error) {
	const query = "SELECT * FROM events WHERE date && tsrange($1, $2, '[]')"
	db := s.context.DefaultTrOrDB(ctx, s.pool)
	result, err := db.Query(ctx, query, fromDate.Format(time.DateTime), toDate.Format(time.DateTime))
	if err != nil {
		return nil, err
	}
	rows, err := pgx.CollectRows(result, pgx.RowToStructByName[eventRow])
	if err != nil {
		return nil, err
	}
	// var rows []eventRow

	// if err := pgxscan.Select(ctx, db, &rows, query,
	// 	fromDate.Format(time.DateTime), toDate.Format(time.DateTime),
	// ); err != nil {
	// 	return nil, err
	// }

	results := make([]*events.Event, len(rows))
	for i, row := range rows {
		var err error
		results[i], err = rowToEvent(row)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (s *EventStorage) FindAll(ctx context.Context) ([]*events.Event, error) {
	panic("not implemented")
}

func (s *EventStorage) FindById(ctx context.Context, id int) (*events.Event, error) {
	const query = "SELECT * FROM event WHERE id = $1"

	var row eventRow
	err := pgxscan.Get(ctx, s.pool, &row, query, id)
	if err != nil {
		return nil, err
	}

	event, err := rowToEvent(row)
	if err != nil {
		return nil, err
	}

	return event, nil
}

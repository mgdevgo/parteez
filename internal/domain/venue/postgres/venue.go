package postgres

import (
	"context"
	"fmt"

	"parteez/internal/domain/venue"

	pgxtrm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type VenueStorage struct {
	pool *pgxpool.Pool
	trm  *pgxtrm.CtxGetter
}

func NewVenueStorage(pool *pgxpool.Pool, trm *pgxtrm.CtxGetter) *VenueStorage {
	return &VenueStorage{
		pool: pool,
		trm:  trm,
	}
}

func (storage *VenueStorage) Save(ctx context.Context, entity *venue.Venue) error {
	const query = `
INSERT INTO venue (id, name, description, type, artwork_id, stages, address, metro_stations, longitude, latitude, visability, updated_at, created_at)
VALUES (COALESCE(NULLIF($1, 0), (SELECT nextval('venue_id_seq'))), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (id) DO UPDATE
	SET name = $2, description = $3, type = $4, artwork_id = $5, stages = $6, address = $7, metro_stations = $8, longitude = $9, latitude = $10, visability = $11, updated_at = $12, created_at = $13
RETURNING id`
	row := venueToRow(entity)

	db := storage.trm.DefaultTrOrDB(ctx, storage.pool)
	result := db.QueryRow(ctx, query, row.ID, row.Name, row.Description, row.Type, row.ArtworkID, row.Stages, row.Address, row.MetroStations, row.Longitude, row.Latitude, row.Visability, row.UpdatedAt, row.CreatedAt)

	if err := result.Scan(&row.ID); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			switch pgerr.Code {
			case "23505":
				if pgerr.ConstraintName == "venue_name_key" {
					return fmt.Errorf("%w: id=%d, name=%s", venue.ErrVenueAlreadyExists, row.ID, row.Name)
				}
			}
		}
		return err
	}

	entity.ID = venue.VenueID(row.ID)

	return nil
}

type filter func(order int) (string, any)

func filterByName(value string) filter {
	return func(order int) (string, any) {
		return fmt.Sprintf("name = $%d", order), value
	}
}

func filterById(value int) filter {
	return func(order int) (string, any) {
		return fmt.Sprintf("id = $%d", order), value
	}
}

func (storage *VenueStorage) findByFilter(ctx context.Context, filters ...filter) ([]*venue.Venue, error) {
	query := "SELECT * FROM venue WHERE 1=1"

	args := make([]any, 0, len(filters))

	for i, filter := range filters {
		q, arg := filter(i + 1)
		query += fmt.Sprintf(" AND %s", q)
		args = append(args, arg)
	}

	db := storage.trm.DefaultTrOrDB(ctx, storage.pool)
	result, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	rows, err := pgx.CollectRows(result, pgx.RowToStructByPos[venueRow])
	if err != nil {
		return nil, err
	}

	venues := make([]*venue.Venue, len(rows))
	for i, row := range rows {
		venue, err := rowToVenue(row)
		if err != nil {
			return nil, err
		}
		venues[i] = venue
	}

	return venues, nil
}

func (storage *VenueStorage) FindById(ctx context.Context, id int) (*venue.Venue, error) {
	venues, err := storage.findByFilter(ctx, filterById(id))
	if err != nil {
		return nil, err
	}
	if len(venues) == 0 {
		return nil, venue.ErrVenueNotFound
	}
	return venues[0], nil
}

func (storage *VenueStorage) FindByName(ctx context.Context, name string) (*venue.Venue, error) {
	venues, err := storage.findByFilter(ctx, filterByName(name))
	if err != nil {
		return nil, err
	}
	if len(venues) == 0 {
		return nil, venue.ErrVenueNotFound
	}
	return venues[0], nil
}

func (storage *VenueStorage) FindAll(ctx context.Context) ([]*venue.Venue, error) {
	return storage.findByFilter(ctx)
}

func (r *VenueStorage) Delete(ctx context.Context, id int) error {
	db := r.trm.DefaultTrOrDB(ctx, r.pool)
	_, err := db.Exec(ctx, "DELETE FROM venue WHERE id = $1", id)
	return err
}

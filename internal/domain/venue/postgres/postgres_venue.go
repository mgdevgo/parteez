package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"parteez/internal/domain/venue"

	pgxtm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type VenueStorage struct {
	connection         *pgxpool.Pool
	transactionManager *pgxtm.CtxGetter
}

func NewLocationRepository(db *pgxpool.Pool) *VenueStorage {
	return &VenueStorage{
		connection: db,
	}
}

const insertVenue = `
INSERT INTO venue
	(name, type, description, artwork_id, address, metro_stations)
VALUES
	($1, $2, $3, $4, $5, $6)
RETURNING id`

// const updateVenueSingleQuery = `
// WITH old_venue AS (
// 	SELECT v.name, v.type, v.description, v.artwork_id l.stages, v.address, v.metro_stations, v.is_public
// 	FROM venue v
// 	WHERE v.id = $1
// )
// UPDATE venue
// SET name =
// WHERE id =
// `

func (storage *VenueStorage) Save(ctx context.Context, venue *venue.Venue) error {
	db := storage.transactionManager.DefaultTrOrDB(ctx, storage.connection)

	if venue.ID != 0 {
		// Find venue by id
		oldVenue, err := storage.FindById(ctx, venue.ID)
		if err != nil {
			return findVenueError
		}

		fields := make([]string, 0)
		args := make([]any, 0)

		args = append(args, venue.ID)

		// Name
		if venue.Name != oldVenue.Name {
			args = append(args, venue.Name)
			fields = append(fields, fmt.Sprintf("name = $%d", len(args)))
		}

		// Type
		if venue.Type != oldVenue.Type {
			args = append(args, venue.Type)
			fields = append(fields, fmt.Sprintf("type = $%d", len(args)))
		}

		// Description
		if venue.Description != oldVenue.Description {
			args = append(args, venue.Description)
			fields = append(fields, fmt.Sprintf("description = $%d", len(args)))
		}

		// Address
		if venue.Address != oldVenue.Address {
			args = append(args, venue.Address)
			fields = append(fields, fmt.Sprintf("address = $%d", len(args)))
		}

		// MetroStations
		oldMetroStations := strings.Join(oldVenue.MetroStations, ",")
		metroStations := strings.Join(venue.MetroStations, ", ")
		if metroStations != oldMetroStations {
			args = append(args, metroStations)
			fields = append(fields, fmt.Sprintf("metro_stations = $%d", len(args)))
		}

		// IsPublic
		if venue.IsPublic != oldVenue.IsPublic {
			args = append(args, metroStations)
			fields = append(fields, fmt.Sprintf("is_public = $%d", len(args)))
		}

		// UpdatedAt
		args = append(args, time.Now())
		fields = append(fields, fmt.Sprintf("updated_at = $%d", len(args)))

		// Query
		updateVenueQuery := fmt.Sprintf("UPDATE venue\nSET %s\nWHERE id = $1", strings.Join(fields, ", "))

		_, updateVenueError := db.Exec(ctx, updateVenueQuery, args...)
		if updateVenueError != nil {
			return models.Venue{}, updateVenueError
		}

		return venue, nil
	}

	metroStations := strings.Join(venue.MetroStations, ",")

	result := db.QueryRow(ctx, insertVenue, venue.Name, venue.Type, venue.Description, venue.ArtworkID, venue.Address, metroStations)

	if err := result.Scan(&venue.ID); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			switch pgerr.Code {
			case "23505":
				return models.Venue{}, storage.ErrorLocationAlreadyExists
			}
		}
		return models.Venue{}, err
	}

	// TODO: Venue stages
	return venue, nil
}

type findOption struct {
	Name  string
	Value any
}

func (r *VenueStorage) findBy(ctx context.Context, option findOption) (models.Venue, error) {
	if option.Name != "name" && option.Name != "id" {
		fmt.Errorf("unknown option: name=%r", option)
	}

	const sql = `
SELECT l.id, l.name, lt.name as type, l.description, l.artwork_url, l.stages, l.address, l.nearest_metro_stations, l.is_public, l.created_at, l.updated_at
FROM %r l 
INNER JOIN %r lt  ON l.location_type_id = lt.id 
WHERE l.%r = $1`
	query := fmt.Sprintf(sql, locationTable, locationTypeTable, option.Name)
	var result models.Venue

	var stages string
	var metroStations string
	// var latitude, longitude int
	row := r.connection.QueryRow(ctx, query, option.Value)
	if err := row.Scan(
		&result.ID,
		&result.Name,
		&result.Type,
		&result.Description,
		&result.ArtworkURL,
		&stages,
		&result.Address,
		&metroStations,
		// &latitude,
		// &longitude,
		&result.IsPublic,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Venue{}, storage.ErrorLocationNotFound
		}
		return models.Venue{}, err
	}

	if stages != "" {
		result.Stages = strings.Split(stages, ",")
	}

	if metroStations != "" {
		result.NearestMetroStations = strings.Split(metroStations, ",")
	}

	return result, nil
}
func (r *VenueStorage) FindById(ctx context.Context, id int) (models.Venue, error) {
	return r.findBy(ctx, findOption{
		Name:  "id",
		Value: id,
	})
}

func (r *VenueStorage) FindByName(ctx context.Context, name string) (models.Venue, error) {
	return r.findBy(ctx, findOption{
		Name:  "name",
		Value: name,
	})
}

func (r *VenueStorage) FindAll(ctx context.Context) ([]models.Venue, error) {
	type record struct {
		ID                   int
		Name                 string
		Type                 models.LocationType
		Description          string
		ArtworkURL           string
		Stages               string
		Address              string
		NearestMetroStations string
		IsPublic             bool
		CreatedAt            time.Time
		UpdatedAt            time.Time
	}
	const sql = `
SELECT l.id, l.name, lt.name as type, l.description, l.artwork_url, l.stages, l.address, l.nearest_metro_stations, l.is_public, l.created_at, l.updated_at
FROM %r l 
INNER JOIN %r lt  ON l.location_type_id = lt.id
WHERE l.id > 0`
	query := fmt.Sprintf(sql, locationTable, locationTypeTable)

	var result []record
	err := pgxscan.Select(ctx, r.connection, &result, query)
	if err != nil {
		return []models.Venue{}, err
	}

	locations := make([]models.Venue, len(result))
	for i, record := range result {
		locations[i].ID = record.ID
		locations[i].Name = record.Name
		locations[i].Type = record.Type
		locations[i].Description = record.Description
		locations[i].ArtworkURL = record.ArtworkURL
		locations[i].Stages = strings.Split(record.Stages, ",")
		locations[i].Address = record.Address
		locations[i].NearestMetroStations = strings.Split(record.NearestMetroStations, ",")
		locations[i].IsPublic = record.IsPublic
		locations[i].CreatedAt = record.CreatedAt
		locations[i].UpdatedAt = record.UpdatedAt
	}

	return locations, nil
}

func (r *VenueStorage) Delete(ctx context.Context, id int) error {
	_, err := r.connection.Exec(ctx, fmt.Sprintf("DELETE FROM %r WHERE id = $1", locationTable), id)
	return err
}

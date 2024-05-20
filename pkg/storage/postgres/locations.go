package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"iditusi/pkg/core"
	"iditusi/pkg/storage/fileds"

	transaction "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type LocationStorage struct {
	driver    *pgxpool.Pool
	txManager *transaction.CtxGetter
}

func NewLocationsStorage(db *pgxpool.Pool) *LocationStorage {
	return &LocationStorage{
		driver: db,
	}
}

func (s *LocationStorage) Save(ctx context.Context, location core.Location) (core.Location, error) {
	fields := fileds.Fields{}
	fields.AddField("name", location.Name)
	fields.AddField("location_type_id", location.Type, "(SELECT id FROM "+locationTypeTable+" WHERE name = $%d)")
	fields.AddField("description", location.Description)
	fields.AddField("artwork_url", location.ArtworkURL)
	fields.AddField("stages", strings.Join(location.Stages, ","))
	fields.AddField("address", location.Address)
	fields.AddField("nearest_metro_stations", strings.Join(location.NearestMetroStations, ","))
	fields.AddField("is_public", location.IsPublic)
	timestamp := time.Now()
	fields.AddField("created_at", timestamp)
	fields.AddField("updated_at", timestamp)

	fieldNames, values, args := fields.Build()
	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
	query := fmt.Sprintf(sql, locationTable, fieldNames, values)

	result := s.driver.QueryRow(ctx, query, args...)

	if err := result.Scan(&location.ID); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			switch pgerr.Code {
			case "23505":
				return core.Location{}, ErrLocationAlreadyExist
			}
		}
		return core.Location{}, err
	}
	return location, nil
}

func (s *LocationStorage) SaveAll(ctx context.Context, locations []core.Location) error {

	// tx := s.txManager.DefaultTrOrDB(ctx, s.driver)

	tx, err := s.driver.Begin(ctx)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fields := fileds.Fields{}
		fields.AddField("id", location.ID)
		fields.AddField("name", location.Name)
		fields.AddField("location_type_id", location.Type, "(SELECT id FROM "+locationTypeTable+" WHERE name = $%d)")
		fields.AddField("description", location.Description)
		fields.AddField("artwork_url", location.ArtworkURL)
		fields.AddField("stages", strings.Join(location.Stages, ","))
		fields.AddField("address", location.Address)
		fields.AddField("nearest_metro_stations", strings.Join(location.NearestMetroStations, ","))
		fields.AddField("is_public", location.IsPublic)
		timestamp := time.Now()
		fields.AddField("created_at", timestamp)
		fields.AddField("updated_at", timestamp)

		fieldNames, values, args := fields.Build()
		const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
		query := fmt.Sprintf(sql, locationTable, fieldNames, values)

		row := tx.QueryRow(ctx, query, args...)

		var id int

		if err := row.Scan(&id); err != nil {
			// var pgerr *pgconn.PgError
			// if errors.As(err, &pgerr) {
			// 	switch pgerr.Code {
			// 	case "23505":
			// 		err := fmt.Errorf()
			// 	}
			// }
			tx.Rollback(ctx)
			break
		}

	}

	return tx.Commit(ctx)
}

type findOption struct {
	Name  string
	Value any
}

func (s *LocationStorage) findBy(ctx context.Context, option findOption) (core.Location, error) {
	if option.Name != "name" && option.Name != "id" {
		fmt.Errorf("unknown option: name=%s", option)
	}

	const sql = `
SELECT l.id, l.name, lt.name as type, l.description, l.artwork_url, l.stages, l.address, l.nearest_metro_stations, l.is_public, l.created_at, l.updated_at
FROM %s l 
INNER JOIN %s lt  ON l.location_type_id = lt.id 
WHERE l.%s = $1`
	query := fmt.Sprintf(sql, locationTable, locationTypeTable, option.Name)
	var result core.Location

	var stages string
	var metroStations string
	// var latitude, longitude int
	row := s.driver.QueryRow(ctx, query, option.Value)
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
			return core.Location{}, ErrLocationNotFound
		}
		return core.Location{}, err
	}

	if stages != "" {
		result.Stages = strings.Split(stages, ",")
	}

	if metroStations != "" {
		result.NearestMetroStations = strings.Split(metroStations, ",")
	}

	return result, nil
}
func (s *LocationStorage) FindByID(ctx context.Context, id int) (core.Location, error) {
	return s.findBy(ctx, findOption{
		Name:  "id",
		Value: id,
	})
}

func (s *LocationStorage) FindByName(ctx context.Context, name string) (core.Location, error) {
	return s.findBy(ctx, findOption{
		Name:  "name",
		Value: name,
	})
}

func (s *LocationStorage) FindAll(ctx context.Context) ([]core.Location, error) {
	type record struct {
		ID                   int
		Name                 string
		Type                 core.LocationType
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
FROM %s l 
INNER JOIN %s lt  ON l.location_type_id = lt.id
WHERE l.id > 0`
	query := fmt.Sprintf(sql, locationTable, locationTypeTable)

	var result []record
	err := pgxscan.Select(ctx, s.driver, &result, query)
	if err != nil {
		return []core.Location{}, err
	}

	locations := make([]core.Location, len(result))
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

func (s *LocationStorage) Delete(ctx context.Context, id int) error {
	_, err := s.driver.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", locationTable), id)
	return err
}

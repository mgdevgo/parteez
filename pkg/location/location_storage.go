package location

import (
	"context"
	"fmt"
	"strings"
	"time"

	"iditusi/pkg/shared/storage"

	transaction "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	_LOCATION_TABLE      = "location"
	_LOCATION_TYPE_TABLE = "location_type"
)

type LocationStorage interface {
	Save(location Location) (Location, error)
	SaveAll(locations []Location) error
	FindByID(id int) (Location, error)
	FindByName(name string) (Location, error)
	FindAll() ([]Location, error)
	Delete(id int) error
}

var _ LocationStorage = (*postgres)(nil)

var ErrAlreadyExist = errors.New("location already exist")
var ErrNotFound = errors.New("location not found")

type postgres struct {
	driver    *pgxpool.Pool
	txManager *transaction.CtxGetter
}

func NewStorage(db *pgxpool.Pool) *postgres {
	return &postgres{
		driver: db,
	}
}

func (s *postgres) Save(location Location) (Location, error) {
	fields := storage.Fields{}
	fields.AddField("name", location.Name)
	fields.AddField("location_type_id", location.Type, "(SELECT id FROM "+_LOCATION_TYPE_TABLE+" WHERE name = $%d)")
	fields.AddField("description", location.Description)
	fields.AddField("artwork_url", location.ArtworkURL)
	fields.AddField("stages", strings.Join(location.Stages, ","))
	fields.AddField("address", location.Address)
	fields.AddField("metro_stations", strings.Join(location.MetroStations, ","))
	fields.AddField("is_public", location.IsPublic)
	timestamp := time.Now()
	fields.AddField("created_at", timestamp)
	fields.AddField("updated_at", timestamp)

	fieldNames, values, args := fields.Build()
	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
	query := fmt.Sprintf(sql, _LOCATION_TABLE, fieldNames, values)

	result := s.driver.QueryRow(context.Background(), query, args...)

	if err := result.Scan(&location.ID); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			switch pgerr.Code {
			case "23505":
				return Location{}, ErrAlreadyExist
			}
		}
		return Location{}, err
	}
	return location, nil
}

func (s *postgres) SaveAll(locations []Location) error {
	ctx := context.TODO()

	// tx := s.txManager.DefaultTrOrDB(ctx, s.driver)

	tx, err := s.driver.Begin(ctx)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fields := storage.Fields{}
		fields.AddField("id", location.ID)
		fields.AddField("name", location.Name)
		fields.AddField("location_type_id", location.Type, "(SELECT id FROM "+_LOCATION_TYPE_TABLE+" WHERE name = $%d)")
		fields.AddField("description", location.Description)
		fields.AddField("artwork_url", location.ArtworkURL)
		fields.AddField("stages", strings.Join(location.Stages, ","))
		fields.AddField("address", location.Address)
		fields.AddField("metro_stations", strings.Join(location.MetroStations, ","))
		fields.AddField("is_public", location.IsPublic)
		timestamp := time.Now()
		fields.AddField("created_at", timestamp)
		fields.AddField("updated_at", timestamp)

		fieldNames, values, args := fields.Build()
		const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"
		query := fmt.Sprintf(sql, _LOCATION_TABLE, fieldNames, values)

		row := tx.QueryRow(context.Background(), query, args...)

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

func (s *postgres) findBy(option findOption) (Location, error) {
	if option.Name != "name" && option.Name != "id" {
		fmt.Errorf("unknown option: name=%s", option)
	}

	const sql = `
SELECT l.id, l.name, lt.name as type, l.description, l.artwork_url, l.stages, l.address, l.metro_stations, l.is_public, l.created_at, l.updated_at
FROM %s l 
INNER JOIN %s lt  ON l.location_type_id = lt.id 
WHERE l.%s = $1`
	query := fmt.Sprintf(sql, _LOCATION_TABLE, _LOCATION_TYPE_TABLE, option.Name)
	var result Location

	var stages string
	var metroStations string
	// var latitude, longitude int
	row := s.driver.QueryRow(context.TODO(), query, option.Value)
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
			return Location{}, ErrNotFound
		}
		return Location{}, err
	}

	if stages != "" {
		result.Stages = strings.Split(stages, ",")
	}

	if metroStations != "" {
		result.MetroStations = strings.Split(metroStations, ",")
	}

	return result, nil
}
func (s *postgres) FindByID(id int) (Location, error) {
	return s.findBy(findOption{
		Name:  "id",
		Value: id,
	})
}

func (s *postgres) FindByName(name string) (Location, error) {
	return s.findBy(findOption{
		Name:  "name",
		Value: name,
	})
}

func (s *postgres) FindAll() ([]Location, error) {
	type record struct {
		ID            int
		Name          string
		Type          Kind
		Description   string
		ArtworkURL    string
		Stages        string
		Address       string
		MetroStations string
		IsPublic      bool
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}
	const sql = `
SELECT l.id, l.name, lt.name as type, l.description, l.artwork_url, l.stages, l.address, l.metro_stations, l.is_public, l.created_at, l.updated_at
FROM %s l 
INNER JOIN %s lt  ON l.location_type_id = lt.id
WHERE l.id > 0`
	query := fmt.Sprintf(sql, _LOCATION_TABLE, _LOCATION_TYPE_TABLE)

	var result []record
	err := pgxscan.Select(context.TODO(), s.driver, &result, query)
	if err != nil {
		return []Location{}, err
	}

	locations := make([]Location, len(result))
	for i, record := range result {
		locations[i].ID = record.ID
		locations[i].Name = record.Name
		locations[i].Type = record.Type
		locations[i].Description = record.Description
		locations[i].ArtworkURL = record.ArtworkURL
		locations[i].Stages = strings.Split(record.Stages, ",")
		locations[i].Address = record.Address
		locations[i].MetroStations = strings.Split(record.MetroStations, ",")
		locations[i].IsPublic = record.IsPublic
		locations[i].CreatedAt = record.CreatedAt
		locations[i].UpdatedAt = record.UpdatedAt
	}

	return locations, nil
}

func (s *postgres) Delete(id int) error {
	_, err := s.driver.Exec(context.TODO(), fmt.Sprintf("DELETE FROM %s WHERE id = $1", _LOCATION_TABLE), id)
	return err
}

package location

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type LocationStorage interface {
	Create(location Location) (string, error)
	Get(id string) (Location, error)
	Update(id string, options map[string]any) (Location, error)
	Delete(id string) error
	Transaction() (int, error)
}

var _ LocationStorage = (*postgresLocationStorage)(nil)

var ErrAlreadyExist = errors.New("location already exist")

type postgresLocationStorage struct {
	tableName string
	driver    *pgxpool.Pool
}

func NewLocationStorage(db *pgxpool.Pool) *postgresLocationStorage {
	return &postgresLocationStorage{
		tableName: "location",
		driver:    db,
	}
}

func (s *postgresLocationStorage) Create(location Location) (string, error) {
	const sql = "INSERT INTO %s (%s) VALUES (%s) RETURNING id"

	const maxNumberOfFields = 12
	fields := make([]string, 0, maxNumberOfFields)
	values := make([]string, 0, maxNumberOfFields)
	args := make([]any, 0, maxNumberOfFields)
	count := 0

	null := func() bool { return true }
	addField := func(name string, value any, isNull ...bool) {
		count++
		fields = append(fields, name)

		placeholder := "$%d"
		if name == "location_type_id" {
			placeholder = "(SELECT id FROM location_type WHERE name = $%d)"
		}

		values = append(values, fmt.Sprintf(placeholder, count))
		if len(isNull) > 0 {
			if fmt.Sprintf("%v", value) == "" {
				value = nil
			}
		}

		args = append(args, value)
	}

	addField("id", location.ID)
	addField("name", location.Name)
	addField("location_type_id", location.Type)
	addField("description", location.Description, null())
	addField("image_url", location.ImageURL, null())
	addField("music_genres", strings.Join(location.MusicGenres, ","), null())
	addField("stages", strings.Join(location.Stages, ","), null())
	addField("address", location.Address, null())
	addField("metro_stations", strings.Join(location.MetroStations, ","), null())
	addField("is_public", location.Public)

	query := fmt.Sprintf(sql, s.tableName, strings.Join(fields, ", "), strings.Join(values, ", "))

	result := s.driver.QueryRow(context.Background(), query, args...)

	var id string
	if err := result.Scan(&id); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			switch pgerr.Code {
			case "23505":
				return "0000000000", ErrAlreadyExist
			}
		}
		return "0000000000", err
	}
	return id, nil
}

func (s *postgresLocationStorage) Get(id string) (Location, error) {
	const sql = `
SELECT l.id, l.name, l.description, l.stages, l.address, t.name as type
FROM %s l 
INNER JOIN %s t  ON l.location_type_id = t.id 
WHERE l.id = $1`
	query := fmt.Sprintf(sql, s.tableName, "location_type")
	var loca Location
	// if err := pgxscan.Select(context.TODO(), s.driver, &locations, query, ids); err != nil {
	//	return []location.Location{}, err
	// }
	var rawStages string
	row := s.driver.QueryRow(context.TODO(), query, id)
	if err := row.Scan(&loca.ID, &loca.Name, &loca.Description, &rawStages, &loca.Address, &loca.Type); err != nil {
		return Location{}, err
	}

	// if err := json.Unmarshal([]byte(stages), &loca.Stages); err != nil {
	//	return location.Location{}, fmt.Errorf("scan: stages: %w", err)
	// }

	loca.Stages = strings.Split(rawStages, ",")

	return loca, nil
}

func (s *postgresLocationStorage) Update(id string, options map[string]any) (Location, error) {
	// TODO implement me
	panic("implement me")
}

func (s *postgresLocationStorage) Delete(id string) error {
	_, err := s.driver.Exec(context.TODO(), fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName), id)
	return err
}

func (s *postgresLocationStorage) Transaction() (int, error) {
	// TODO implement me
	panic("implement me")
}

// massimusz
// hunder

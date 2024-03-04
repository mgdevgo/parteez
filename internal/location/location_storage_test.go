package location

import (
	"context"
	"iditusi/pkg/utils"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/require"
)

func TestPostgresLocationStorageCreate(t *testing.T) {
	type testCase struct {
		PostgresLocationStorage *postgresLocationStorage

		Location Location

		ExpectedID    string
		ExpectedError error
	}

	test := func(name string, t *testing.T, tc *testCase) {
		t.Run(name, func(t *testing.T) {
			actualID, actualError := tc.PostgresLocationStorage.Create(tc.Location)

			require.Equal(t, tc.ExpectedError, actualError)
			require.Equal(t, tc.ExpectedID, actualID)
		})
	}

	PostgresURL := os.Getenv("TEST_POSTGRES_URL")
	db, err := pgxpool.New(context.Background(), PostgresURL)
	if err != nil {

	}
	storage := NewLocationStorage(db)

	id := utils.NewID()
	test("Minimum required fields", t, &testCase{
		PostgresLocationStorage: storage,
		Location: Location{
			ID:     id,
			Name:   "Test_" + nanoid.Must(),
			Type:   LocationUnknown,
			Stages: StagesDefault,
		},
		ExpectedID:    id,
		ExpectedError: nil,
	})

	test("All fields provided", t, &testCase{
		PostgresLocationStorage: storage,
		Location: Location{
			ID:            "6900690069",
			Name:          "Test_" + nanoid.Must(),
			Type:          LocationUnknown,
			Description:   "Description_" + nanoid.Must(),
			ImageURL:      "https://img.test/img/x",
			MusicGenres:   []string{"techno", "electronic"},
			Stages:        StagesDefault,
			Address:       "Address_" + nanoid.Must(),
			MetroStations: []string{"Адмиралтейская"},
			Public:        false,
		},
		ExpectedID:    "6900690069",
		ExpectedError: nil,
	})

	test("Duplicate ID", t, &testCase{
		PostgresLocationStorage: storage,
		Location: Location{
			ID:     "6900690069",
			Name:   "Test_" + nanoid.Must(),
			Type:   LocationUnknown,
			Stages: StagesDefault,
		},
		ExpectedID:    "0000000000",
		ExpectedError: ErrAlreadyExist,
	})
}

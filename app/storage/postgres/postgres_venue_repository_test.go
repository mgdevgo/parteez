package postgres

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"iditusi/pkg/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var POSTGRES_URL = os.Getenv("TEST_POSTGRES_URL")

func Test_storage_Save(t *testing.T) {
	type testCase struct {
		Location Location

		ExpectedID    int
		ExpectedError error
	}

	test := func(name string, t *testing.T, tc *testCase) {
		t.Run(name, func(t *testing.T) {
			db, err := pgxpool.New(context.Background(), POSTGRES_URL)
			if err != nil {
				t.Fatal(err)
			}

			storage := NewStorage(db)
			actualError, _ := storage.Save(tc.Location)

			require.Equal(t, tc.ExpectedError, actualError)
			// require.Equal(t, tc.ExpectedID, actualID)
		})
	}

	test("Minimum required fields", t, &testCase{
		Location: Location{
			ID:     1000000069,
			Name:   "Test_" + nanoid.Must(),
			Type:   Unknown,
			Stages: StagesDefault,
		},
		ExpectedID:    100000069,
		ExpectedError: nil,
	})

	test("All fields provided", t, &testCase{
		Location: Location{
			Name:          "Test_" + nanoid.Must(),
			Type:          Unknown,
			Description:   "Description_" + nanoid.Must(),
			ArtworkURL:    "https://img.test/img/x",
			Stages:        StagesDefault,
			Address:       "Address_" + nanoid.Must(),
			MetroStations: []string{"Адмиралтейская"},
			IsPublic:      false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		ExpectedError: nil,
	})

	test("Insert duplicate", t, &testCase{
		Location: Location{
			ID:     1000000001,
			Name:   "Test_" + nanoid.Must(),
			Type:   Unknown,
			Stages: StagesDefault,
		},
		ExpectedID:    1000000001,
		ExpectedError: repositories.ErrLocationAlreadyExist,
	})
}

func Test_storage_FindByID(t *testing.T) {
	type testCase struct {
		locationID       int
		expectedLocation Location
		expectedError    error
	}

	test := func(name string, t *testing.T, tc *testCase) {
		t.Run(name, func(t *testing.T) {
			db, err := pgxpool.New(context.Background(), POSTGRES_URL)
			if err != nil {
				t.Fatal(err)
			}

			storage := NewStorage(db)

			_, err = storage.Save(tc.expectedLocation)
			require.NoError(t, err)
			defer func() {
				err := storage.Delete(tc.locationID)
				require.NoError(t, err)
			}()

			actualLocation, actualError := storage.FindByID(tc.locationID)
			require.Equal(t, tc.expectedError, actualError)
			require.Equal(t, tc.expectedLocation.ID, actualLocation.ID)
			require.Equal(t, tc.expectedLocation.Name, actualLocation.Name)
			require.Equal(t, tc.expectedLocation.Description, actualLocation.Description)
			require.Equal(t, tc.expectedLocation.Type, actualLocation.Type)
			require.Equal(t, tc.expectedLocation.ArtworkURL, actualLocation.ArtworkURL)
			require.ElementsMatch(t, tc.expectedLocation.Stages, actualLocation.Stages)
			require.Equal(t, tc.expectedLocation.Address, actualLocation.Address)
			require.ElementsMatch(t, tc.expectedLocation.MetroStations, actualLocation.MetroStations)
			require.Equal(t, tc.expectedLocation.IsPublic, actualLocation.IsPublic)
			require.WithinDuration(t, tc.expectedLocation.CreatedAt, actualLocation.CreatedAt, time.Second)
			require.WithinDuration(t, tc.expectedLocation.UpdatedAt, actualLocation.UpdatedAt, time.Second)
		})
	}

	test("Find default location", t, &testCase{
		locationID: 123,
		expectedLocation: Location{
			ID:          123,
			Name:        "odin dva tri",
			Type:        Unknown,
			Description: "very cool places",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		expectedError: nil,
	})
	test("Success", t, &testCase{
		locationID: 69,
		expectedLocation: Location{
			ID:            69,
			Name:          "oh my",
			Type:          Club,
			Description:   "Secret places",
			ArtworkURL:    "htps://img.local/69",
			Stages:        []string{"main", "second"},
			Address:       "Somewhere in the world",
			MetroStations: []string{"Невский проспект"},
			IsPublic:      false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		expectedError: nil,
	})

	test("Success", t, &testCase{
		locationID: 69,
		expectedLocation: Location{
			ID:            69,
			Name:          "oh my",
			Type:          Club,
			Description:   "Secret places",
			ArtworkURL:    "htps://img.local/69",
			Stages:        []string{"main", "second"},
			Address:       "Somewhere in the world",
			MetroStations: []string{"Невский проспект"},
			IsPublic:      false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		expectedError: nil,
	})
}

func Test_storage_BatchSave(t *testing.T) {
	type testCase struct {
		Name string

		Storage *VenueRepository

		Locations []Location

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.Storage.SaveAll(tc.Locations)

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	db, err := pgxpool.New(context.Background(), POSTGRES_URL)
	if err != nil {
		t.Fatal(err)
	}

	storage := NewStorage(db)

	file, err := os.Open("./2024-04-19-591268800.json")
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]map[string]string, 0)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		t.Fatal(err)
	}

	validate(t, &testCase{
		Name:          "OK",
		Storage:       storage,
		Locations:     nil,
		ExpectedError: nil,
	})
}

package event

import (
	"context"
	"iditusi/pkg/utils"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestPostgresEventStorageCreate(t *testing.T) {
	type testCase struct {
		Name string

		PostgresEventStorage *postgresEventStorage

		Event Event

		ExpectedID    string
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualID, actualError := tc.PostgresEventStorage.Create(tc.Event)

			require.Equal(t, tc.ExpectedError, actualError)
			require.Equal(t, tc.ExpectedID, actualID)
		})
	}

	PostgresURL := os.Getenv("TEST_POSTGRES_URL")
	db, err := pgxpool.New(context.Background(), PostgresURL)
	if err != nil {

	}
	storage := NewEventStorage(db)

	validate(t, &testCase{
		Name:                 "Only required",
		PostgresEventStorage: storage,
		Event: Event{
			ID:     "f5W-3xnfz3SETGKA6m3PH",
			Name:   "Test_" + utils.NewID(),
			MinAge: 18,
		},
		ExpectedID:    "f5W-3xnfz3SETGKA6m3PH",
		ExpectedError: err,
	})

	validate(t, &testCase{
		Name:                 "All fields",
		PostgresEventStorage: storage,
		Event: Event{
			ID:          "fQpONkr67Ms3f7_JUBWHW",
			Name:        "Test_" + utils.NewID(),
			ImageURL:    "https://img.test/abc.jpg",
			Description: "Description_fQpONkr67Ms3f7_JUBWHW",
			MusicGenres: []string{"Test", "Test 2"},
			LineUp: map[string][]LineUp{
				"main": {
					{
						Name: "Dj Test",
					},
				},
			},
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(time.Hour * 5),
			MinAge:     18,
			TicketsURL: "https://ticket.test/abc",
			Price:      map[string]int{},
			LocationID: "5487797728",
			Promoter:   "Test label",
			IsPublic:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		ExpectedID:    "fQpONkr67Ms3f7_JUBWHW",
		ExpectedError: err,
	})
}

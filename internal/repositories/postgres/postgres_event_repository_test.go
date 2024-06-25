package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"iditusi/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func Test_EventRepository_Save(t *testing.T) {
	var tests = []struct {
		name    string
		event   models.Event
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			event: models.Event{
				ID:             1,
				Tittle:         "Test Event #1",
				Description:    "The first event",
				LineUp:         []models.LineUp{},
				GenreIDs:       []int{},
				AgeRestriction: 18,
				Promoter:       "",
				Date:           time.Date(2024, time.July, 28, 0, 0, 0, 0, time.UTC),
				StartTime:      time.Date(2024, time.July, 28, 23, 59, 0, 0, time.UTC),
				EndTime:        time.Date(2024, time.July, 29, 6, 0, 0, 0, time.UTC),
				TicketsURL:     "",
				Tickets:        nil,
				ArtworkID:      1,
				VenueID:        1,
				IsPublic:       false,
				Timestamp:      models.NewTimestamp(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			db, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URL"))
			r := NewPostgresEventRepository(db)
			got, err := r.Save(ctx, tt.event)
			if !tt.wantErr(t, err, fmt.Sprintf("Save(%v, %v)", ctx, tt.event)) {
				return
			}
			assert.Equalf(t, tt.event, got, "Save(%v, %v)", ctx, tt.event)
		})
	}
}

package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	pgxtrm "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"parteez/internal/domain/events"
	"parteez/pkg/postgres"
)

func TestEventStorage_FindByDate(t *testing.T) {
	tests := []struct {
		name     string
		fromDate time.Time
		toDate   time.Time
		want     []*events.Event
		wantErr  bool
	}{
		{
			name:     "find by date",
			fromDate: time.Date(2025, 1, 1, 23, 59, 0, 0, time.UTC),
			toDate:   time.Date(2025, 1, 3, 4, 0, 0, 0, time.UTC),
			want: []*events.Event{
				{
					ID:    1,
					Title: "Test 1",
					Date: events.Date{
						Start: time.Date(2025, 1, 1, 23, 0, 0, 0, time.UTC),
						End:   time.Date(2025, 1, 2, 2, 0, 0, 0, time.UTC),
					},
					Status:    events.StatusDraft,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:    2,
					Title: "Test 2",
					Date: events.Date{
						Start: time.Date(2025, 1, 2, 23, 0, 0, 0, time.UTC),
						End:   time.Date(2025, 1, 3, 3, 0, 0, 0, time.UTC),
					},
					Status:    events.StatusDraft,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	pool, err := postgres.New("postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		t.Fatalf("Failed connect to database: %v", err)
	}

	storage := NewEventStorage(pool, pgxtrm.DefaultCtxGetter)

	t.Cleanup(func() {
		_, err := pool.Exec(context.Background(), "TRUNCATE TABLE events RESTART IDENTITY")
		assert.NoError(t, err)
		pool.Close()
	})

	for i := range 5 {
		err := storage.Save(context.Background(), &events.Event{
			Title: fmt.Sprintf("Test %d", i+1),
			Date: events.Date{
				Start: time.Date(2025, 1, 1+i, 23, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 1, 2+i, 2+i, 0, 0, 0, time.UTC),
			},
			Status:    events.StatusDraft,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.FindByDate(context.Background(), tt.fromDate, tt.toDate)
			require.NoError(t, err)
			require.Equal(t, len(tt.want), len(got))
			for i, want := range tt.want {
				assert.Equal(t, want.ID, got[i].ID)
				assert.Equal(t, want.Title, got[i].Title)
				assert.Equal(t, want.Date, got[i].Date)
				assert.Equal(t, want.Status, got[i].Status)
				assert.Equal(t, want.CreatedAt, got[i].CreatedAt.UTC())
				assert.Equal(t, want.UpdatedAt, got[i].UpdatedAt.UTC())
			}
		})
	}
}

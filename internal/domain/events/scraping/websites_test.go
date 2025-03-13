package scraping

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRuporEventsWebsite(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		source := NewRuporEventsWebsite(slog.Default())

		service := NewEventScrapingService([]Source{source}, slog.Default())
		got := service.Scrape(context.Background())

		for result := range got {
			if result.IsFailure() {
				fmt.Println(result.Errors)
				continue
			}

			fmt.Println(result.Event)
		}
	})
}

func TestNewBlankWebsite(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		source := NewBlankWebsite(slog.Default())

		results, err := source.Parse(context.Background())
		require.NoError(t, err)

		var hasValidResults bool
		for result := range results {
			if result.IsFailure() {
				t.Logf("Scraping errors: %v", result.Errors)
				continue
			}

			// Validate required fields are present
			require.NotEmpty(t, result.Event.Title, "Event title should not be empty")
			require.NotEmpty(t, result.Event.AgeRestriction, "Age restriction should not be empty")
			require.NotEmpty(t, result.Event.StartDate, "Start date should not be empty")
			require.NotEmpty(t, result.Event.EndDate, "End date should not be empty")
			require.NotEmpty(t, result.Venue.Name, "Venue name should not be empty")
			require.NotEmpty(t, result.Venue.Address, "Venue address should not be empty")
			require.NotEmpty(t, result.Venue.MetroStations, "Venue metro stations should not be empty")

			// Log successful scrape for debugging
			t.Logf("Successfully scraped event: %s at %s", result.Event.Title, result.Venue.Name)
			hasValidResults = true
		}

		require.True(t, hasValidResults, "Should have at least one valid scraped event")
	})
}

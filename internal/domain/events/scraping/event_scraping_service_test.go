package scraping

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventScrapingService_Scrape(t *testing.T) {
	t.Run("rupor", func(t *testing.T) {
		source, err := ListingEvents(slog.Default())
		require.NoError(t, err)

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

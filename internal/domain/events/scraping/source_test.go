package scraping

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebsite_Parse(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body><h1>Test Page</h1></body></html>`))
	}))
	defer ts.Close()

	logger := slog.Default()
	url, _ := url.Parse(ts.URL)
	website := NewWebsite(url, logger)

	// Test Parse function
	resultChan, err := website.Parse(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, resultChan)

	// Verify channel is closed after scraping
	_, ok := <-resultChan
	assert.False(t, ok, "Channel should be closed after scraping")
}

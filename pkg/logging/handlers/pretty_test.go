package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestPrettyHandlerOptionsNewPrettyHandler(t *testing.T) {
	type testCase struct {
		Name string

		PrettyHandlerOptions PrettyHandlerOptions

		Out io.Writer

		ExpectedPrettyHandler *PrettyHandler
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualPrettyHandler := tc.PrettyHandlerOptions.NewPrettyHandler(tc.Out)

			assert.Equal(t, tc.ExpectedPrettyHandler, actualPrettyHandler)
		})
	}

	validate(t, &testCase{})
}

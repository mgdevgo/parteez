package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorageFind(t *testing.T) {
	type testCase struct {
		Name string

		Storage *postgres

		Options FindOptions

		ExpectedSlice []Event
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.Storage.Find(tc.Options)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "OK",
		Options: FindOptions{
			limit:   5,
			offset:  0,
			orderBy: "start_date",
			date:    time.Date(2024, time.April, 1, 0, 0, 0, 0, nil),
		},
		ExpectedSlice: nil,
		ExpectedError: nil,
	})
}

func Test_storage_Save(t *testing.T) {
	type testCase struct {
		Name string

		Storage *postgres

		Event Event

		ExpectedEvent Event
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualEvent, actualError := tc.Storage.Save(tc.Event)

			assert.Equal(t, tc.ExpectedEvent, actualEvent)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name:          "OK",
		Storage:       nil,
		Event:         Event{},
		ExpectedEvent: Event{},
		ExpectedError: nil,
	})
}

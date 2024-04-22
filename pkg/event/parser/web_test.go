package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"iditusi/pkg/event"
	"iditusi/pkg/location"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func Test_Parse(t *testing.T) {
	p := NewWebParser("https://listing.events/spb")
	p.ParseX()
}

func Test_Import(t *testing.T) {
	file, err := os.Open(os.Getenv("PARSER_DATA"))
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var data []Result
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		t.Fatal(err)
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("TEST_POSTGRES_URL"))
	if err != nil {
		t.Fatal(err)
	}

	eventStorage := event.NewStorage(db)
	locationStorage := location.NewStorage(db)

	for _, i := range data {
		fmt.Printf("Creating a New Event: %s\n", i.Event.Name)

		var err error
		var place location.Location
		place, err = locationStorage.FindByName(i.Location.Name)
		if err != nil {
			if errors.Is(err, location.ErrNotFound) {
				fmt.Printf("Location not found: name=%s\n", i.Location.Name)

				i.Location.Type = location.Unknown
				place, err = locationStorage.Save(i.Location)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("Created a New Location: id=%d\n", place.ID)

			} else {
				t.Fatal(err)
			}
		}

		i.Event.LocationID = place.ID
		i.Event.Price = make(event.Price)

		newEvent, err := eventStorage.Save(i.Event)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("Created a New Event [id:%d, name:%s]\n", newEvent.ID, newEvent.Name)
	}
}

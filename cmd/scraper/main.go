package main

import (
	"context"
	"fmt"
	"iditusi/internal/location"
	"iditusi/pkg/utils"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgxpool"
)

const url = "https://listing.events/spb/bars"

// const url = "https://listing.events/spb/clubs"

func main() {
	clubs := make([]location.Location, 0)
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	// iterating over the list of HTML product elements
	c.OnHTML("article > .item-info", func(e *colly.HTMLElement) {
		// fmt.Println(e.Name, e.Attr("class"))

		venue := location.Location{
			ID:   utils.NewID(),
			Type: location.LocationBar,
		}

		// scraping the data of interest
		venue.Name = e.ChildText("h3.venue-name a")
		div := e.ChildText("div.address")

		if strings.Contains(div, "\n") {
			address := strings.Split(div, "\n")
			metroString := regexp.MustCompile(`\s+`).ReplaceAllString(strings.ReplaceAll(address[1], "м.", ""), " ")
			metroArray := strings.Split(strings.TrimSpace(metroString), ", ")

			venue.MetroStations = metroArray
			venue.Address = address[0]
		}

		venue.Address = strings.TrimSpace(venue.Address)

		e.ForEach("div > span.btn-tag-genres", func(_ int, e *colly.HTMLElement) {
			venue.MusicGenres = append(venue.MusicGenres, strings.TrimSpace(e.Text))
		})

		// adding the product instance with scraped data to the list of products

		clubs = append(clubs, venue)

		fmt.Println(venue.ID, venue.Name, venue.Address, venue.MusicGenres, venue.MetroStations)
	})

	// downloading the target HTML page
	c.Visit(url)

	PostgresURL := os.Getenv("TEST_POSTGRES_URL")
	db, err := pgxpool.New(context.Background(), PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	storage := location.NewLocationStorage(db)

	for _, c := range clubs {
		_, err := storage.Create(c)
		if err != nil {
			fmt.Println(err)
			fmt.Println(c.Name)
		}
	}

}

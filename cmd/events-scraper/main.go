package main

import (
	"context"
	"encoding/json"
	"fmt"
	"iditusi/internal/event"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgxpool"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const url = "https://listing.events/spb"

func main() {
	events := make([]event.Event, 0)
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	// iterating over the list of HTML product elements
	c.OnHTML("article", func(e *colly.HTMLElement) {
		data := parse(e)

		events = append(events, data)
	})
	// downloading the target HTML page
	c.Visit(url)

	// print(events)
	saveToDb(events)

}

func parse(e *colly.HTMLElement) event.Event {
	var data event.Event
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	data.ID = nanoid.Must()
	data.Name = e.ChildText("h3.event-information__title")
	data.ImageURL = e.ChildAttr("div.event-information-banner__afisha > a > img", "src")

	data.StartTime = parseDate(e.ChildText("h2"))

	e.ForEachWithBreak("img", func(i int, e *colly.HTMLElement) bool {
		if e.Attr("src") == "/img/place.svg" {
			data.LocationID = e.Attr("alt")
			return false
		}

		return true
	})

	e.ForEach("div > span.btn-tag-genres", func(_ int, e *colly.HTMLElement) {
		data.MusicGenres = append(data.MusicGenres, strings.TrimSpace(e.Text))
	})

	data.LineUp = make(map[string][]event.LineUp)
	lineup := strings.Split(
		regexp.MustCompile(`\s{2,}`).ReplaceAllString(
			strings.ReplaceAll(e.ChildText("div.container > div.leading-loose"), "\n", ""),
			"",
		),
		",",
	)

	for _, v := range lineup {
		artist := strings.TrimSpace(v)
		if artist != "" {
			data.LineUp["main"] = append(data.LineUp["main"], event.LineUp{
				Name: artist,
			})
		}
	}

	price, err := strconv.Atoi(
		regexp.MustCompile(`\d+`).FindString(
			e.ChildText("div.container > div:nth-child(3) > div:nth-child(2)"),
		),
	)
	if err != nil {
		fmt.Println(err)
	}

	data.Price = map[string]int{}
	data.Price["main"] = price

	data.TicketsURL = "https://listing.events" + e.ChildAttr("div.container > div:nth-child(1) > div:nth-child(1) > a", "href")

	return data
}

func parseDate(date string) time.Time {
	monthsMap := map[string]time.Month{
		"марта":  time.March,
		"апреля": time.April,
		"мая":    time.May,
		"июня":   time.June,
		"июля":   time.July,
	}

	daystring := regexp.MustCompile(`\d+`).FindString(date)
	day, err := strconv.Atoi(daystring)
	if err != nil {
		fmt.Println(err)
	}

	var mounth time.Month
	for i, v := range monthsMap {
		if strings.Contains(date, i) {
			mounth = v
		}
	}

	tz, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
	}

	t := time.Date(2023, mounth, day, 22, 0, 0, 0, tz)
	return t
}

func print(events []event.Event) {
	b, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func saveToDb(events []event.Event) {
	postgresURL := os.Getenv("TEST_POSTGRES_URL")
	db, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {

	}
	storage := event.NewEventStorage(db)

	for _, e := range events {
		id, err := storage.Create(e)
		fmt.Print(id)
		if err != nil {
			fmt.Printf("- Error: %s", err)
		}
		fmt.Print("\n")
	}
}

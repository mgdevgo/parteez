package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"iditusi/pkg/event"
	"iditusi/pkg/location"

	"github.com/gocolly/colly"
)

type WebParser struct {
	rootURL   string
	collector *colly.Collector
}

func NewWebParser(URL string) *WebParser {
	return &WebParser{
		rootURL: URL,
		collector: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0"),
			// colly.Debugger(&debug.LogDebugger{}),
		),
	}
}

type Result struct {
	Event    event.Event
	Location location.Location
}

func (p *WebParser) Parse() {

}

func (p *WebParser) ParseX() []Result {

	results := make([]Result, 0)
	// events := make([]event.Event, 0)
	// locations := make([]location.Location, 0)
	// event[name] -> location[name]
	// eventToLocationMap := make(map[string]string)

	// body > main > article
	p.collector.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
		data := event.Event{}
		place := location.Location{}

		// Age restriction
		age := e.ChildText("h1.event-information__name > span.age-limit")
		data.AgeRestriction, _ = strconv.Atoi(strings.Trim(age, "+"))
		// Name
		data.Name = strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))
		log.Printf("Name: %s\n", data.Name)
		// Description
		data.Description = strings.TrimSpace(e.ChildText("div.event-information__description"))
		// ArtworkURL
		data.ArtworkURL = strings.Trim(e.ChildAttr("div.event-information-banner", "style"), "background-image: url()")
		// Genres
		data.Genres = make([]string, 0)
		e.ForEach("div.bottom-tags-wrapper > span.btn-tag-genres", func(i int, e *colly.HTMLElement) {
			data.Genres = append(data.Genres, strings.TrimSpace(e.Text))
		})
		// LineUp
		data.LineUp = event.LineUp{"main": make([]event.Performer, 0)}
		e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
			data.LineUp["main"] = append(data.LineUp["main"],
				event.Performer{Name: squashSpace(e.Text)},
			)
		})
		// Dates
		loc, _ := time.LoadLocation("Europe/Moscow")
		// Start date and time
		data.StartDate, _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
				strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
			),
			loc,
		)
		// End date. If not presented time.Time zero value
		data.EndDate, _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="endDate"]`, "datetime"),
				strings.TrimSpace(e.ChildText(`p > time[itemprop="endDate"]`)),
			),
			loc,
		)

		// Location name
		place.Name = strings.TrimSpace(e.ChildText("p.venue-name > a"))
		// Location address
		address := strings.Split(e.ChildText("p.venue-adress"), "м.") // typo is left intentionally

		place.Address = squashSpace(address[0])

		// Location metro stations
		if len(address) > 1 {
			metro := strings.Split(address[1], ",")
			for _, i := range metro {
				place.MetroStations = append(place.MetroStations, strings.TrimSpace(i))
			}
		} else {
			place.MetroStations = make([]string, 0)
		}

		// Tickets
		// TODO: get full url instead of shorten
		data.TicketsURL = e.ChildAttr("a.buy-btn", "href")

		// events = append(events, data)
		// locations = append(locations, place)
		// eventToLocationMap[data.Name] = place.Name
		results = append(results, Result{
			data,
			place,
		})
	})

	p.collector.OnHTML("article > div.event-information-banner__afisha > a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		e.Request.Visit(link)
	})

	p.collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting URL: ", r.URL.String())
	})

	p.collector.Visit(p.rootURL)

	p.collector.Wait()

	bytes, _ := json.MarshalIndent(results, "", " ")
	now := time.Now()
	err := os.WriteFile(fmt.Sprintf("./%s-%d.json", strings.Fields(now.String())[0], now.Unix()), bytes, 0644)
	if err != nil {
		log.Println(err)
	}
	return results
}

// squashSpace turn many spaces into one
func squashSpace(text string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " "))
}

func handleParsingError(err error, meta string) {
	log.Printf("parse error [meta: %s]: %s", meta, err.Error())
}

func parseDate(date string, times string) (time.Time, error) {
	if date != "" && times != "" {
		loc, _ := time.LoadLocation("Europe/Moscow")
		localDate, _ := time.ParseInLocation(time.DateOnly, date, loc)
		fmt.Println(strings.Fields(times))
		t, _ := time.Parse(time.TimeOnly, strings.Fields(times)[3]+":00")

		return localDate.Add(time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute())), nil
	}

	return time.Time{}, fmt.Errorf("parseDate: empty input")
}

// Locations
// c.OnHTML("article > .item-info", func(e *colly.HTMLElement) {
// 	// fmt.Println(e.Name, e.Attr("class"))
//
// 	venue := location.Location{
// 		ID:   utils.NewID(),
// 		Type: location.Bar,
// 	}
//
// 	// scraping the data of interest
// 	venue.Name = e.ChildText("h3.venue-name a")
// 	div := e.ChildText("div.address")
//
// 	if strings.Contains(div, "\n") {
// 		address := strings.Split(div, "\n")
// 		metroString := regexp.MustCompile(`\s+`).ReplaceAllString(strings.ReplaceAll(address[1], "м.", ""), " ")
// 		metroArray := strings.Split(strings.TrimSpace(metroString), ", ")
//
// 		venue.MetroStations = metroArray
// 		venue.Address = address[0]
// 	}
//
// 	venue.Address = strings.TrimSpace(venue.Address)
//
// 	e.ForEach("div > span.btn-tag-genres", func(_ int, e *colly.HTMLElement) {
// 		venue.MusicGenres = append(venue.MusicGenres, strings.TrimSpace(e.Text))
// 	})
//
// 	// adding the product instance with scraped data to the list of products
//
//
// 	fmt.Println(venue.ID, venue.Name, venue.Address, venue.MusicGenres, venue.MetroStations)
// })

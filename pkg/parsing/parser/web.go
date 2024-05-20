package parser

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"iditusi/pkg/core"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

/*
	   parsing
	   / 	\ 				\
source_1	source_2  ...source_x
			/ | \
    parser_1  parser_2 ...parser_x
*/
// parser should trigger scrape all sources every 12 hours.
//
type WebPage struct {
	name      string
	url       string
	collector *colly.Collector
	output    chan string
}

func NewWebPage(name string, URL string) *WebPage {
	return &WebPage{
		name: name,
		url:  URL,
		collector: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0"),
			colly.Debugger(&debug.LogDebugger{}),

			// Cache responses to prevent multiple download of pages
			// even if the collector is restarted
			colly.CacheDir("./.cache/"+name),
			colly.Async(true),
		),
		output: make(chan string),
	}
}

func (p *WebPage) Parse() chan string {
	// Create another collector to scrape event details
	detailsCollector := p.collector.Clone()

	// results := make([]Result, 0)
	// events := make([]event.Event, 0)
	// locations := make([]location.Location, 0)
	// event[name] -> location[name]
	// eventToLocationMap := make(map[string]string)
	p.collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting URL: ", r.URL.String())
	})

	p.collector.OnHTML("article > div.event-information-banner__afisha > a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		detailsCollector.Visit(link)
	})

	detailsCollector.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
		var event core.Event
		var place core.Location

		// Age restriction
		age := e.ChildText("h1.event-information__name > span.age-limit")
		event.AgeRestriction, _ = strconv.Atoi(strings.Trim(age, "+"))
		// Name
		event.Name = strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))
		// Description
		event.Description = strings.TrimSpace(e.ChildText("div.event-information__description"))
		// ArtworkURL
		event.ArtworkURL = strings.Trim(e.ChildAttr("div.event-information-banner", "style"), "background-image: url()")
		// Genres
		event.Genres = make([]string, 0)
		e.ForEach("div.bottom-tags-wrapper > span.btn-tag-genres", func(i int, e *colly.HTMLElement) {
			event.Genres = append(event.Genres, strings.TrimSpace(e.Text))
		})
		// LineUp
		artists := make([]core.LineUpArtist, 0)
		e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
			artists = append(artists, core.LineUpArtist{Name: squashSpace(e.Text)})
		})
		event.LineUp = core.LineUp{
			{"main", artists},
		}
		// Dates
		loc, _ := time.LoadLocation("Europe/Moscow")
		// Start date and time
		event.StartDate, _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
				strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
			),
			loc,
		)
		// End date. If not presented time.Time zero value
		event.EndDate, _ = time.ParseInLocation(
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
		stations := make([]string, 0)
		if len(address) > 1 {
			metro := strings.Split(address[1], ",")
			for _, i := range metro {
				stations = append(stations, strings.TrimSpace(i))
			}
		}
		place.NearestMetroStations = stations

		// Tickets
		// TODO: get full url instead of shorten
		event.Tickets.URL = e.ChildAttr("a.buy-btn", "href")

		// events = append(events, event)
		// locations = append(locations, places)
		// eventToLocationMap[event.Name] = places.Name
		// results = append(results, Result{
		// 	event,
		// 	Location: place,
		// })

		p.output <- event.Name
	})

	p.collector.Visit(p.url)

	p.collector.Wait()

	// bytes, _ := json.MarshalIndent(results, "", " ")
	// now := time.Now()
	// err := os.WriteFile(fmt.Sprintf("./%s-%d.json", strings.Fields(now.String())[0], now.Unix()), bytes, 0644)
	// if err != nil {
	// 	log.Println(err)
	// }

	return p.output
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
// 		ID:   id.NewID(),
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

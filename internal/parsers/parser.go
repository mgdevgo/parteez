package parsers

import (
	"sync"

	models2 "iditusi/internal/models"
)

type Parser interface {
	Run() <-chan Result
}

type Source interface {
	Parse() <-chan Result
}

type key struct {
	title string
	date  string
}

type Result struct {
	models2.Event
	Location models2.Venue
}

type parser struct {
	sources      []Source
	errorHandler func(error)
	output       chan Result
	doubles      map[key]struct{}
	doublesMutex sync.RWMutex
}

var Sources = []Source{
	// NewPtichkaWebPage(),
	// web.NewWebPage("ptichka", "https://listing.events/spb"),
	// NewTelegram("tg_blank", -1001331847765),
}

func NewEventParser(sources []Source) *parser {
	return &parser{
		sources:      sources,
		errorHandler: nil,
		output:       make(chan Result),
		doubles:      make(map[key]struct{}),
		doublesMutex: sync.RWMutex{},
	}
}

func (p *parser) Run() <-chan Result {
	wg := sync.WaitGroup{}

	for _, source := range p.sources {
		wg.Add(1)
		output := source.Parse()
		go func() {
			defer wg.Done()
			for item := range output {
				p.doublesMutex.Lock()

				key := key{
					title: item.Name,
					date:  item.StartDate.String(),
				}
				_, ok := p.doubles[key]
				if !ok {
					p.doubles[key] = struct{}{}
					p.output <- item
				}

				p.doublesMutex.Unlock()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(p.output)
	}()

	return p.output
}

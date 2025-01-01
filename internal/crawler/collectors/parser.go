package parsers

import (
	"sync"

	result "iditusi/internal/crawler/content"
)

type Parser interface {
	Run() <-chan result.Data
}

type Source interface {
	Parse() <-chan result.Data
}

type key struct {
	title string
	date  string
}

type parser struct {
	sources      []Source
	errorHandler func(error)
	output       chan result.Data
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
		output:       make(chan result.Data),
		doubles:      make(map[key]struct{}),
		doublesMutex: sync.RWMutex{},
	}
}

func (p *parser) Run() <-chan result.Data {
	wg := sync.WaitGroup{}

	for _, source := range p.sources {
		wg.Add(1)
		output := source.Parse()
		go func() {
			defer wg.Done()
			for item := range output {
				p.doublesMutex.Lock()

				key := key{
					title: item.Tittle,
					date:  item.StartTime.String(),
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

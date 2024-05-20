package parser

import (
	"sync"

	"iditusi/pkg/core"
)

type Parser interface {
	Run() chan string
}

type Source interface {
	Parse() chan string
}

type key struct {
	title string
	date  string
}

type Result struct {
	core.Event
	Location core.Location
}

type parser struct {
	sources      []Source
	errorHandler func(error)
	output       chan string
	doubles      map[key]struct{}
	doublesMutex sync.RWMutex
}

var Sources = []Source{
	NewWebPage("web_ptichka", "https://listing.events/spb"),
	// NewTelegram("tg_blank", -1001331847765),
}

func NewEventParser(sources []Source) *parser {
	return &parser{
		sources:      sources,
		errorHandler: nil,
		output:       make(chan string),
		doubles:      make(map[key]struct{}),
		doublesMutex: sync.RWMutex{},
	}
}

func (p *parser) Run() chan string {
	wg := sync.WaitGroup{}

	for _, source := range p.sources {
		wg.Add(1)
		output := source.Parse()
		go func() {
			defer wg.Done()
			for item := range output {
				p.doublesMutex.Lock()

				key := key{
					title: item,
					date:  "",
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

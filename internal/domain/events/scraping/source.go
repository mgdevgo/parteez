package scraping

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
)

type SourceType int

const (
	SourceTypeTelegram SourceType = iota + 1
	SourceTypeWebsite
)

type Source interface {
	ID() string
	Type() SourceType
	Parse(ctx context.Context) (chan Result, error)
}

type Website struct {
	url       *url.URL
	collector *colly.Collector
	result    chan Result
	logger    *slog.Logger
}

func NewWebsite(URL *url.URL, logger *slog.Logger, optionFuncs ...func(w *Website) colly.CollectorOption) *Website {
	website := &Website{
		url:    URL,
		result: make(chan Result),
		logger: logger,
	}

	options := commonCollyOptions(logger)

	for _, optionFunc := range optionFuncs {
		options = append(options, optionFunc(website))
	}

	website.collector = colly.NewCollector(options...)

	return website
}

func (w *Website) Type() SourceType {
	return SourceTypeWebsite
}

func (w *Website) Options(options ...func(w *Website) colly.CollectorOption) {
	for _, option := range options {
		option(w)
	}
}

func (w *Website) Parse(ctx context.Context) (chan Result, error) {
	if err := w.collector.Visit(w.url.String()); err != nil {
		return nil, err
	}

	go func() {
		w.collector.Wait()
		close(w.result)
	}()

	return w.result, nil
}

func (w *Website) ID() string {
	return w.url.Host
}

type TelegramChannel struct {
	Name          string
	ChannelID     int
	LastVisitedAt time.Time
	output        chan any
}

func NewTelegramChanel(name string, chanelId int) *TelegramChannel {
	return &TelegramChannel{
		Name:      name,
		ChannelID: chanelId,
		output:    make(chan any),
	}
}

func (chanel *TelegramChannel) Parse() chan string {
	panic("not implemented")
}

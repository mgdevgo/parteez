package source

import "context"

type Kind int

const (
	KindTelegram Kind = iota + 1
	KindWebPage
)

type Source interface {
	ID() string
	Name() string
	Parse() chan any
}

type FetchingService interface {
	FetchData(sources []int) chan any
}

type DataSourceService interface {
	ListSources() error
	Add(ctx context.Context, source Source) error
	Update(ctx context.Context, source Source) error
	Disable(ctx context.Context, sourceId int) error
}

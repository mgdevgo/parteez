package scraping

type SourceType int

const (
	SourceTypeTelegram SourceType = iota + 1
	SourceTypeWebsite
)

type Source interface {
	ID() string
	Parse() chan Result
}

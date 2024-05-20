package parsing

import "iditusi/pkg/parsing/parser"

type ParsingLocalService struct {
	parser parser.Parser
}

type ParsingService interface {
	StartProcessing(sourceID int)
	GetErrorReport()
}

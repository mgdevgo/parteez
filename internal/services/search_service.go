package services

import (
	"iditusi/internal/services/local"
)

type SearchService interface {
}

var _ SearchService = (*local.SearchService)(nil)

package search

type SearchService interface {
}

type SearchLocalService struct {
}

func NewSearchService() *SearchLocalService {
	return &SearchLocalService{}
}

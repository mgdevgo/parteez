package crawler

type CrawlerServiceConfig struct{
	CrawlerMaxRetries int
}


type CrawlerService struct {
}

func (service *CrawlerService) AddSource() error {
	panic("not implemented")
}

func (service *CrawlerService) ListSources() error {
	panic("not implemented")
}

func (service *CrawlerService) StartCrawl() error {
	panic("not implemented")
}

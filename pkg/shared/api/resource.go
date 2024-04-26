package api

type Resource struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Href       string `json:"href"`
	Attributes any    `json:"attributes"`
	Relations  any    `json:"relations"`
}

type ResourceCollectionResponse struct {
	Resources []Resource `json:"data"`
}

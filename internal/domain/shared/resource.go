package shared

type Resource struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Href       string `json:"href"`
	Attributes any    `json:"attributes"`
	Relations  any    `json:"relations"`
}

type ResourceCollectionResponse struct {
	Next string     `json:"next,omitempty"`
	Data []Resource `json:"data"`
}

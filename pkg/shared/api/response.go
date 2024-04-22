package api

type APIResponse struct {
	Data []Data `json:"data"`
}

type Data struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Href       string `json:"href"`
	Attributes any    `json:"attributes"`
}

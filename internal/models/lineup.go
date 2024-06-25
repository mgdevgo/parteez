package models

type LineUps []LineUp

type LineUp struct {
	Stage   int      `json:"stage"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name    string `json:"name"`
	StartAt string `json:"startAt"`
	Live    bool   `json:"live"`
}

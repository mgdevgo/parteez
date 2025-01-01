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

type LineupItem struct {
	Name string `json:"name"`
	Live bool   `json:"live"`
}
type Lineup map[string]map[string]LineupItem

// var l Lineup = Lineup{
// 	"main": map[string]LineupItem{
// 		"23:59": LineupItem{
// 			Name: "DJX",
// 			Live: true,
// 		},
// 		"01:00": LineupItem{
// 			Name: "DJ2",
// 			Live: true,
// 		},
// 	},
// }

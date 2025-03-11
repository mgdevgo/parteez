package events

type LineUp struct {
	Artists []string `json:"artists"`
	Stages  []Stage  `json:"stages"`
}

func NewLineUp(artists []string, stages []Stage) LineUp {
	return LineUp{
		Artists: artists,
		Stages:  stages,
	}
}

type Stage struct {
	Name   string       `json:"name"`
	LineUp []LineUpSlot `json:"lineup"`
}

type LineUpSlot struct {
	Time    string   `json:"time"`
	IsLive  bool     `json:"isLive"`
	Artists []string `json:"artists"`
}

func NewLineUpSlot(time string, isLive bool, artist []string) LineUpSlot {
	return LineUpSlot{
		Time:    time,
		IsLive:  isLive,
		Artists: artist,
	}
}

package event

import (
	"parteez/internal/domain/artist"
)

type LineUp struct {
	Stages     []any
	Headliner  artist.ArtistID
	Stage      string
	StageAlias string
	Timetable  []LineUpTimetable
}

func NewLineUp(stage, alias string, timetable []LineUpTimetable) *LineUp {
	return &LineUp{
		Stage:      stage,
		StageAlias: alias,
		Timetable:  timetable,
	}
}

type LineUpTimetable struct {
	Time    string
	IsLive  bool
	Artists []artist.ArtistID
}

func NewLineUpTimetable(time string, isLive bool, artist []artist.ArtistID) LineUpTimetable {
	return LineUpTimetable{
		Time:    time,
		IsLive:  isLive,
		Artists: artist,
	}
}

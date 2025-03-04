package events

type GenreID int

type Genre struct {
	ID   GenreID
	Name string
}

func NewGenre(id GenreID, name string) *Genre {
	return &Genre{
		ID:   id,
		Name: name,
	}
}

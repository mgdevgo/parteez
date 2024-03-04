package list_locations

import "iditusi/internal/location"

type Usecase struct {
	locations location.Storage
}

func NewUsecase(locationStorage location.Storage) *Usecase {
	return &Usecase{locations: locationStorage}
}

type ListEventsSuccessResult struct {
	Locations []location.Location
}

type ListEventsErrorResult struct {
	Code    int
	Error   string
	Message string
}

func (u Usecase) Execute(ids []int) (ListEventsSuccessResult, error) {
	locations, _ := u.locations.Get(ids)

	result := ListEventsSuccessResult{
		Locations: locations,
	}

	return result, nil
}

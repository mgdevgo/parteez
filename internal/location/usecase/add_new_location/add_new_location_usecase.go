package add_new_location

import "iditusi/internal/location"

type UseCase struct {
	locations location.Storage
}

func NewUsecase(locationStorage location.Storage) *UseCase {
	return &UseCase{locations: locationStorage}
}

func (uc UseCase)
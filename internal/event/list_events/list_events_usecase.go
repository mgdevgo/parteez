package list_events

import (
	"iditusi/internal/event"
	"iditusi/internal/location"
)

type ListEventsUsecase struct {
	eventStorage    event.Storage
	locationStorage location.Storage
}

func New(eventStorage event.Storage, locationStorage location.Storage) ListEventsUsecase {
	return ListEventsUsecase{
		eventStorage:    eventStorage,
		locationStorage: locationStorage,
	}
}

type ListEventsSuccessResult struct {
	event.Event
	location.Location
}

func (u ListEventsUsecase) Execute(ids []string) ([]ListEventsSuccessResult, error) {
	events, _ := u.eventStorage.Get(ids)
	locationIDs := make(map[int]location.Location, 0)

	result := make([]ListEventsSuccessResult, 0, len(events))
	for _, _event := range events {
		_location, ok := locationIDs[_event.LocationID]
		if !ok {
			locations, _ := u.locationStorage.Get([]int{_event.LocationID})
			_location = locations[0]
			locationIDs[_location.ID] = _location
		}

		result = append(result, ListEventsSuccessResult{
			Event:    _event,
			Location: _location,
		})
	}

	return result, nil
}

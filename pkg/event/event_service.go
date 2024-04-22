package event

import (
	"fmt"
	"time"

	"iditusi/pkg/location"
)

type service struct {
	eventStorage    EventStorage
	locationStorage location.LocationStorage
}

func NewService(eventStorage EventStorage, locationStorage location.LocationStorage) *service {
	return &service{
		eventStorage:    eventStorage,
		locationStorage: locationStorage,
		// parsers: []parser.Parser
		// workers
	}
}

type CreateEventRequest struct{}

type CreateEventResponse struct {
}

func (s *service) CreateNewEvent(input CreateEventRequest) (CreateEventResponse, error) {

	return CreateEventResponse{}, nil
}

type ListEventsRequest struct {
	Limit  int
	Offset int
}

type ListEventsResponse struct {
	Next string  `json:"next"`
	Data []Event `json:"data"`
}

func (s *service) ListEvents(request ListEventsRequest) (ListEventsResponse, error) {
	// 1. Get date + weekend
	// 2. Fetch events
	// 3. calculate next link
	// 4. retrun data
	now := time.Now()
	loc, _ := time.LoadLocation("Europe/Moscow")
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	options := new(FindOptions)
	options.Fill()
	options.SetLimit(request.Limit)
	options.SetOffset(request.Offset)
	options.SetDate(from)

	events, err := s.eventStorage.Find(*options)
	if err != nil {
		// return api error with message and code
		fmt.Println(err)
		return ListEventsResponse{}, err
	}

	response := ListEventsResponse{
		Data: events,
		Next: "api/v1/+options",
	}

	return response, nil
}

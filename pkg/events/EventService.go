package events

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"iditusi/pkg/core"
)

type EventService interface {
	FindById(ctx context.Context, id int) (core.Event, error)
	FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]core.Event, error)
}

var _ EventService = (*EventLocalService)(nil)

type EventLocalService struct {
	eventStorage    EventRepository
	locationService locations.LocationService
}

type CreateEventRequest struct{}

type CreateEventResponse struct {
}

func (s *EventLocalService) FindById(ctx context.Context, id int) (core.Event, error) {
	return s.eventStorage.FindById(ctx, id)
}

func (s *EventLocalService) FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]core.Event, error) {
	now := time.Now()
	loc, _ := time.LoadLocation("Europe/Moscow")
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	params := core.EventSearchParams{
		FromDate: date,
		Limit:    20,
		Offset:   0,
	}

	events, err := s.eventStorage.Find(ctx, params)
	if err != nil {
		// return api error with message and code
		fmt.Println(err)
		return ListEventsResponse{}, err
	}

	locations := make(map[int]core.Location)
	for _, event := range events {
		_, ok := locations[event.LocationID]
		if !ok {
			location, err := s.locationStorage.FindByID(ctx, event.LocationID)
			if err != nil {
				return ListEventsResponse{}, err
			}
			locations[event.LocationID] = location
		}
	}

	next := fmt.Sprintf("/api/v1/events?offset=%d", request.Offset+request.Limit)
	if limit := request.Limit; limit > 0 && limit != 5 {
		next += "&limit=" + strconv.Itoa(limit)
	}

	response := ListEventsResponse{
		Next: next,
		Data: make([]EventResult, 0),
	}

	for _, event := range events {
		location, _ := locations[event.LocationID]
		response.Data = append(response.Data, EventResult{
			Event:    event,
			Location: &location,
		})
	}

	return response, nil
}

// func (s *EventLocalService) CreateNewEvent(input CreateEventRequest) (CreateEventResponse, error) {
//
// 	return CreateEventResponse{}, nil
// }

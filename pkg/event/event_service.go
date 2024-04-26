package event

import (
	"fmt"
	"strconv"
	"time"

	"iditusi/pkg/location"
)

type service struct {
	eventStorage    Storage
	locationStorage location.Storage
}

func NewService(eventStorage Storage, locationStorage location.Storage) *service {
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
	Next string             `json:"next"`
	Data []ListEventsResult `json:"data"`
}

type ListEventsResult struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	ArtworkURL     string           `json:"artworkURL"`
	Genres         Genres           `json:"genreNames"`
	LineUp         LineUp           `json:"lineUp"`
	StartDate      time.Time        `json:"startDate"`
	EndDate        string           `json:"endDate,omitempty"`
	AgeRestriction string           `json:"ageRestriction"`
	TicketsURL     string           `json:"ticketsURL,omitempty"`
	Price          Price            `json:"price,omitempty"`
	Promoter       string           `json:"promoter,omitempty"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Location       LocationResponse `json:"location,omitempty"`
}

type LocationResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Type          string   `json:"type"`
	ArtworkURL    string   `json:"artworkURL"`
	Stages        []string `json:"stages"`
	Address       string   `json:"address"`
	MetroStations []string `json:"metroStations"`
}

func (s *service) ListEvents(request ListEventsRequest) (ListEventsResponse, error) {
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

	locations := make(map[int]location.Location)
	for _, event := range events {
		_, ok := locations[event.LocationID]
		if !ok {
			location, err := s.locationStorage.FindByID(event.LocationID)
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

	response := ListEventsResponse{Next: next, Data: make([]ListEventsResult, 0)}

	for _, event := range events {
		l, _ := locations[event.LocationID]
		fmt.Println(l)
		response.Data = append(response.Data, ListEventsResult{
			ID:             strconv.Itoa(event.ID),
			Name:           event.Name,
			Description:    event.Description,
			ArtworkURL:     event.ArtworkURL,
			Genres:         event.Genres,
			LineUp:         event.LineUp,
			StartDate:      event.StartDate,
			EndDate:        event.EndDate.String(),
			AgeRestriction: fmt.Sprintf("%d+", event.AgeRestriction),
			TicketsURL:     event.TicketsURL,
			Price:          event.Price,
			Promoter:       event.Promoter,
			UpdatedAt:      event.UpdatedAt,
			Location: LocationResponse{
				ID:            strconv.Itoa(l.ID),
				Name:          l.Name,
				Description:   l.Description,
				Type:          string(l.Type),
				ArtworkURL:    l.ArtworkURL,
				Stages:        l.Stages,
				Address:       l.Address,
				MetroStations: l.MetroStations,
			},
		})
	}

	return response, nil
}

package administration

import (
	"time"
)

type EventManagementService interface {
	GetEvents() ([]string, error)
	CreateEvent(title string, date time.Time) (string, error)
	UpdateEvent(id int, update any) error
	DeleteEvent(id string) error
	MakeEventAvailable(id string) error
}

type EventStorage interface {
}

type EventManagementLocalService struct {
	eventStorage EventStorage
}

func NewEventManagementService(eventStorage EventStorage) *EventManagementLocalService {
	return &EventManagementLocalService{
		eventStorage: eventStorage,
	}
}

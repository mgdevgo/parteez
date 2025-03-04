package events

import "errors"

var ErrStatus = errors.New("status")

type Status string

func NewStatus(input string) (Status, error) {
	status := Status(input)

	if _, ok := statusMap[status]; !ok {
		return "", ErrStatus
	}

	return status, nil
}

const (
	StatusDraft      Status = "DRAFT"
	StatusModeration Status = "MODERATION"
	StatusPublished  Status = "PUBLISHED"
	StatusArchived   Status = "ARCHIVED"
)

var statusMap = map[Status]struct{}{
	StatusDraft:      {},
	StatusModeration: {},
	StatusPublished:  {},
	StatusArchived:   {},
}

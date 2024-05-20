package core

import "time"

type Timestamp struct {
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

func NewTimestamp() Timestamp {
	now := time.Now()
	return Timestamp{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewTimestamp2() (createdAt time.Time, updatedAt time.Time) {
	timeNow := time.Now()
	return timeNow, timeNow
}

// func (timestamp *Timestamp) Init() *Timestamp {
// 	now := time.Now()
// 	timestamp.CreatedAt = now
// 	timestamp.UpdatedAt = now
// 	return timestamp
// }

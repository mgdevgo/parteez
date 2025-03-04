package events

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrDate              = errors.New("date")
	ErrDateStartAfterEnd = fmt.Errorf("%w: start after end", ErrDate)
)

type Date struct {
	Start time.Time
	End   time.Time
}

func NewDate(start, end time.Time) (Date, error) {
	if start.After(end) {
		return Date{}, ErrDateStartAfterEnd
	}
	return Date{
		Start: start,
		End:   end,
	}, nil
}

package storage

import (
	"errors"
	"fmt"
)

var (
	ErrorNotFound      = errors.New("not found")
	ErrorAlreadyExists = errors.New("already exists")
)

var (
	ErrorEventNotFount      = fmt.Errorf("event: %w", ErrorNotFound)
	ErrorEventAlreadyExists = fmt.Errorf("event: %w", ErrorAlreadyExists)

	ErrorLocationNotFound      = fmt.Errorf("location: %w", ErrorNotFound)
	ErrorLocationAlreadyExists = fmt.Errorf("location: %w", ErrorAlreadyExists)

	ErrorUserNotFound      = fmt.Errorf("user: %w", ErrorNotFound)
	ErrorUserAlreadyExists = fmt.Errorf("user: %w", ErrorAlreadyExists)
)

package postgres

import "errors"

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrEventAlreadyExists = errors.New("event already exists")

	ErrLocationNotFound     = errors.New("location not found")
	ErrLocationAlreadyExist = errors.New("location already exist")

	ErrDuplicateEntry = errors.New("duplicate entry")
)

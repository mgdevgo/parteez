package core

import (
	"math/rand"

	"github.com/google/uuid"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/oklog/ulid/v2"
)

type ID interface {
	int | string | uuid.UUID | ulid.ULID
}

const alphabet = "0123456789"
const size = 10

func NewID() string {
	return nanoid.MustGenerate(alphabet, size)
}

func NullID() string { return "0000000000" }

func NewNumericID() int {
	min := 1000000000
	max := 2147483647

	return rand.Intn(max-min+1) + min
}

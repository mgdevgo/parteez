package utils

import nanoid "github.com/matoous/go-nanoid/v2"

const alphabet = "0123456789"
const size = 10

func NewID() string {
	return nanoid.MustGenerate(alphabet, size)
}

func NullID() string { return "0000000000" }

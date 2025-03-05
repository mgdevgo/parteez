package utils

import (
	"encoding/json"
	"fmt"
)

func ToJSON[T any](val T) ([]byte, error) {
	bytes, err := json.Marshal(&val)
	if err != nil {
		return nil, fmt.Errorf("failed to build JSON for %v: %w", val, err)
	}
	return bytes, nil
}

func Keys[T any](val map[string]T) []string {
	var res = make([]string, len(val))

	i := 0

	for k := range val {
		res[i] = k
		i++
	}

	return res
}

package readutil

import (
	"encoding/json"
	"io"
)

func readJSON[T any](r io.Reader) (*T, error) {
	var result T

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

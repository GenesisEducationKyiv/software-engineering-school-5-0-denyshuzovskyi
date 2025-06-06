package readutil

import (
	"encoding/json"
	"fmt"
	"io"
)

//nolint:unused
func readJSON[T any](r io.Reader) (*T, error) {
	var result T

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return &result, nil
}

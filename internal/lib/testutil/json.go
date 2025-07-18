package testutil

import (
	"encoding/json"
	"io"
)

func UnmarshalJSONFromReader[T any](r io.Reader) (T, error) {
	var v T
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}

package cache

import (
	"context"
	"encoding/json"
	"time"

	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
)

type JSONCache[T any] struct {
	underlying Cache
	ttl        time.Duration
}

func NewJSONCache[T any](underlying Cache, ttl time.Duration) *JSONCache[T] {
	return &JSONCache[T]{
		underlying: underlying,
		ttl:        ttl,
	}
}

func (j *JSONCache[T]) Set(ctx context.Context, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return j.underlying.Set(ctx, key, data, j.ttl)
}

func (j *JSONCache[T]) Get(ctx context.Context, key string) (T, error) {
	var zero T

	data, err := j.underlying.Get(ctx, key)
	if err != nil {
		return zero, err
	}
	if len(data) == 0 {
		return zero, commonerrors.ErrEmptyCacheValue
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return zero, err
	}

	return value, nil
}

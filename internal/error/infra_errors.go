package error

import "errors"

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")

	ErrEmptyCacheValue = errors.New("empty cache value")
	ErrCacheMiss       = errors.New("cache miss")
)

package error

import "errors"

var (
	ErrLocationNotFound          = errors.New("no matching location found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrInvalidToken              = errors.New("invalid token")
	ErrTokenNotFound             = errors.New("token not found")
	ErrUnexpectedState           = errors.New("unexpected state")
	ErrValidationFailed          = errors.New("validation failed")

	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

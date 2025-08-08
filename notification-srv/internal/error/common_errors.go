package error

import "errors"

var (
	ErrClosedMessageChannel = errors.New("message channel closed")
	ErrUnsupportedCommand   = errors.New("unsupported command")
)

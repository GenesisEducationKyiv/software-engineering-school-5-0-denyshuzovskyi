package error

import "errors"

var (
	ErrEmailSendingFailed = errors.New("email sending failed")
)

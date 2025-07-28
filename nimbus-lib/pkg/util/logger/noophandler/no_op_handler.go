package noophandler

import (
	"context"
	"log/slog"
)

type NoOpHandler struct{}

func (h *NoOpHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (h *NoOpHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *NoOpHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *NoOpHandler) WithGroup(_ string) slog.Handler {
	return h
}

func NewNoOpHandler() *NoOpHandler {
	return &NoOpHandler{}
}

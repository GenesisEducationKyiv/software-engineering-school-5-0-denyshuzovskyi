package samplinghandler

import (
	"context"
	"log/slog"
	"math/rand"
)

type ProbabilisticSamplingHandler struct {
	handler     slog.Handler
	probability float64 // 0.1 = 10%
	minLevel    slog.Level
}

func NewProbabilisticSamplingHandler(
	handler slog.Handler,
	probability float64,
	minLevel slog.Level,
) slog.Handler {
	return &ProbabilisticSamplingHandler{
		handler:     handler,
		probability: probability,
		minLevel:    minLevel,
	}
}

func (h *ProbabilisticSamplingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ProbabilisticSamplingHandler) Handle(ctx context.Context, record slog.Record) error {
	// Always log if level is >= minLevel
	if record.Level >= h.minLevel {
		return h.handler.Handle(ctx, record)
	}

	if rand.Float64() < h.probability {
		return h.handler.Handle(ctx, record)
	}
	return nil
}

func (h *ProbabilisticSamplingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ProbabilisticSamplingHandler{
		handler:     h.handler.WithAttrs(attrs),
		probability: h.probability,
		minLevel:    h.minLevel,
	}
}

func (h *ProbabilisticSamplingHandler) WithGroup(name string) slog.Handler {
	return &ProbabilisticSamplingHandler{
		handler:     h.handler.WithGroup(name),
		probability: h.probability,
		minLevel:    h.minLevel,
	}
}

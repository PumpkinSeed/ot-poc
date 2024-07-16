package otpoc

import (
	"context"
	"io"
	"log/slog"

	slogcloudlogging "github.com/PumpkinSeed/slog-cloudlogging"
	"go.opentelemetry.io/otel/trace"
)

func handlerWithSpanContext(handler slog.Handler) *spanContextLogHandler {
	return &spanContextLogHandler{Handler: handler}
}

// spanContextLogHandler is an slog.Handler which adds attributes from the
// span context.
type spanContextLogHandler struct {
	slog.Handler
}

// Handle overrides slog.Handler's Handle method. This adds attributes from the
// span context to the slog.Record.
func (t *spanContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	if s := trace.SpanContextFromContext(ctx); s.IsValid() {
		record.AddAttrs(
			slog.Any("logging.googleapis.com/trace", s.TraceID()),
		)
		record.AddAttrs(
			slog.Any("logging.googleapis.com/spanId", s.SpanID()),
		)
		record.AddAttrs(
			slog.Bool("logging.googleapis.com/trace_sampled", s.TraceFlags().IsSampled()),
		)
	}
	return t.Handler.Handle(ctx, record)
}

func setupLogger() {
	googleHandler := slogcloudlogging.NewHandler("guild-xyz-dev", "dev-test-log", &slogcloudlogging.Opts{
		Handler: slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		ForwardHandler: false,
	})
	googleHandler.AutoFlush()

	logger := slog.New(handlerWithSpanContext(googleHandler))
	slog.SetDefault(logger)
}

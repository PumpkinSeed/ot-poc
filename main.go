package otpoc

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

func Run() {
	ctx := context.Background()

	otlploghttp.WithEndpoint("http://localhost:4318")

	if err := setup(ctx); err != nil {
		slog.ErrorContext(ctx, "Setup failed", slog.Any("error", err))
		shutdown()
	}
}

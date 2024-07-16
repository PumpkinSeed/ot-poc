package otpoc

import (
	"context"
	"log/slog"
	"os"
)

func Run() {
	ctx := context.Background()

	//otlploghttp.WithEndpoint("http://localhost:4318")
	//otlploghttp.WithEndpointURL("http://localhost:4318")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
	os.Setenv("OTEL_SERVICE_NAME", "ot-poc")
	os.Setenv("OTEL_GO_X_EXEMPLAR", "true")

	setupLogger()

	if err := setup(ctx); err != nil {
		slog.ErrorContext(ctx, "Setup failed", slog.Any("error", err))
		shutdown()
	}

	s := server{}
	if err := s.run(); err != nil {
		slog.ErrorContext(ctx, "Run failed", slog.Any("error", err))
	}
}

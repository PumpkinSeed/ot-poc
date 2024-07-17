package otpoc

import (
	"io"
	"log/slog"

	slogcloudlogging "github.com/PumpkinSeed/slog-cloudlogging"
)

func setupLogger() {
	googleHandler := slogcloudlogging.NewHandler("guild-xyz-dev", "dev-test-log", &slogcloudlogging.Opts{
		Handler: slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		ForwardHandler: false,
	})
	googleHandler.UseOpenTelemetryTracer = true
	googleHandler.TracePrefix = "projects/guild-xyz-dev/traces/"
	googleHandler.AutoFlush()

	logger := slog.New(googleHandler)
	slog.SetDefault(logger)
}

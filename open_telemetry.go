package otpoc

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

var globalShutdown = map[string]ShutdownInterface{}

type ShutdownInterface interface {
	Shutdown(ctx context.Context) error
}

func setup(ctx context.Context) error {
	// Configure Context Propagation to use the default W3C traceparent format
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	if err := setupTracer(ctx); err != nil {
		return err
	}

	if err := setupMeter(ctx); err != nil {
		return err
	}

	return nil
}

func setupTracer(ctx context.Context) error {
	// Configure Trace Export to send spans as OTLP
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return err
	}
	tp := trace.NewTracerProvider(trace.WithBatcher(spanExporter))
	globalShutdown["tracer"] = tp
	otel.SetTracerProvider(tp)

	return nil
}

func setupMeter(ctx context.Context) error {
	// Configure Metric Export to send metrics as OTLP
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return err
	}
	mp := metric.NewMeterProvider(metric.WithReader(metricReader))
	globalShutdown["meter"] = mp
	otel.SetMeterProvider(mp)

	return nil
}

func shutdown() {
	ctx := context.Background()
	if globalShutdown != nil {
		for service, sd := range globalShutdown {
			if err := sd.Shutdown(ctx); err != nil {
				slog.ErrorContext(ctx, "Shutdown failed",
					slog.String("service", service), slog.Any("error", err))
			}
		}
	}
}

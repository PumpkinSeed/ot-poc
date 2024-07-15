package otpoc

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const scopeName = "github.com/guildxyz/ot-poc"

var (
	meter                = otel.Meter(scopeName)
	tracer               = otel.Tracer(scopeName)
	sleepHistogram       metric.Float64Histogram
	subRequestsHistogram metric.Int64Histogram
)

func init() {
	var err error

	sleepHistogram, err = meter.Float64Histogram("poc.sleep.duration",
		metric.WithDescription("Sample histogram to measure time spent in sleeping"),
		metric.WithExplicitBucketBoundaries(0.05, 0.075, 0.1, 0.125, 0.150, 0.2),
		metric.WithUnit("s"))
	if err != nil {
		panic(err)
	}

	subRequestsHistogram, err = meter.Int64Histogram("poc.subrequests",
		metric.WithDescription("Sample histogram to measure the number of subrequests made"),
		metric.WithExplicitBucketBoundaries(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		metric.WithUnit("{request}"))
	if err != nil {
		panic(err)
	}
}

type server struct{}

func (s server) run() error {
	handleHTTP("/single", s.single)
	handleHTTP("/multi", s.multi)

	return http.ListenAndServe(":8080", nil)
}

// handleHTTP handles the http HandlerFunc on the specified route, and uses
// otelhttp for context propagation, trace instrumentation, and metric
// instrumentation.
func handleHTTP(route string, handleFn http.HandlerFunc) {
	instrumentedHandler := otelhttp.NewHandler(otelhttp.WithRouteTag(route, handleFn), route)

	http.Handle(route, instrumentedHandler)
}

func (s server) single(w http.ResponseWriter, r *http.Request) {
	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)

	hostValue := attribute.String("host.value", r.Host)
	sleepHistogram.Record(r.Context(), sleepTime.Seconds(), metric.WithAttributes(hostValue))

	fmt.Fprintf(w, "work completed in %v\n", sleepTime)
}

func (s server) multi(w http.ResponseWriter, r *http.Request) {
	subRequests := 3 + rand.Intn(4)
	slog.InfoContext(r.Context(), "handle /multi request", slog.Int("subRequests", subRequests))

	ctx, span := tracer.Start(r.Context(), "subrequests")
	defer span.End()

	for i := 0; i < subRequests; i++ {
		if err := s.callSingle(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}

	subRequestsHistogram.Record(ctx, int64(subRequests))

	fmt.Fprintln(w, "ok")
}

func (s server) callSingle(ctx context.Context) error {
	res, err := otelhttp.Get(ctx, "http://localhost:8080/single")
	if err != nil {
		return err
	}

	return res.Body.Close()
}

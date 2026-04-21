package observability

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	httpOnce        sync.Once
	httpTracer      trace.Tracer
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
)

func initHTTP() {
	httpTracer = Tracer("fiber-http")
	m := Meter("fiber-http")
	requestCounter, _ = m.Int64Counter("http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	requestDuration, _ = m.Float64Histogram("http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
	)
}

// Middleware returns a Fiber middleware for OpenTelemetry observability.
// It starts a trace span for every request and records metrics.
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		httpOnce.Do(initHTTP)
		start := time.Now()

		// Tracer and Span
		ctx, span := httpTracer.Start(c.UserContext(), fmt.Sprintf("%s %s", c.Method(), c.Path()),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.path", c.Path()),
				attribute.String("http.host", c.Hostname()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Add Trace ID to response header
		if spanContext := span.SpanContext(); spanContext.HasTraceID() {
			c.Set("X-Trace-Id", spanContext.TraceID().String())
		}

		// Pass the new context with span to Fiber
		c.SetUserContext(ctx)

		// Continue with request chain
		err := c.Next()

		// Record Metrics
		status := c.Response().StatusCode()
		elapsed := time.Since(start).Seconds()

		attrs := metric.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.Int("http.status", status),
		)

		requestCounter.Add(ctx, 1, attrs)
		requestDuration.Record(ctx, elapsed, attrs)

		// Enrich Span with response status
		span.SetAttributes(attribute.Int("http.status_code", status))
		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

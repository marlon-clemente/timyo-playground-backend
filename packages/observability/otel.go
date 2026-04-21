package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	globalTP *sdktrace.TracerProvider
	globalMP *sdkmetric.MeterProvider
)

// Tracer returns a named tracer backed by the initialized provider.
func Tracer(name string) oteltrace.Tracer {
	if globalTP != nil {
		return globalTP.Tracer(name)
	}
	return otel.Tracer(name)
}

// Meter returns a named meter backed by the initialized provider.
func Meter(name string) otelmetric.Meter {
	if globalMP != nil {
		return globalMP.Meter(name)
	}
	return otel.Meter(name)
}

// Init initializes the OpenTelemetry SDK with OTLP HTTP exporters.
// It returns a shutdown function to be called on application exit.
func Init(ctx context.Context, serviceName, endpoint, auth, org string, insecure bool) (func(), error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	// Prepare headers for OpenObserve
	headers := map[string]string{}
	if auth != "" {
		headers["Authorization"] = "Basic " + auth
	}

	// Trace Exporter
	traceOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithHeaders(headers),
		otlptracehttp.WithURLPath(fmt.Sprintf("/api/%s/v1/traces", org)),
	}
	if insecure {
		traceOpts = append(traceOpts, otlptracehttp.WithInsecure())
	}
	traceExporter, err := otlptracehttp.New(ctx, traceOpts...)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	globalTP = tp
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Metric Exporter
	metricOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithHeaders(headers),
		otlpmetrichttp.WithURLPath(fmt.Sprintf("/api/%s/v1/metrics", org)),
	}
	if insecure {
		metricOpts = append(metricOpts, otlpmetrichttp.WithInsecure())
	}
	metricExporter, err := otlpmetrichttp.New(ctx, metricOpts...)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(10*time.Second))),
		sdkmetric.WithResource(res),
	)
	globalMP = mp
	otel.SetMeterProvider(mp)

	shutdown := func() {
		_ = tp.Shutdown(ctx)
		_ = mp.Shutdown(ctx)
	}

	return shutdown, nil
}



package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a new span in the current context.
// The returned function must be called to end the span (typically via defer).
//
//	ctx, end := observability.StartSpan(ctx, "my-operation")
//	defer end()
func StartSpan(ctx context.Context, name string) (context.Context, func()) {
	ctx, span := Tracer("app").Start(ctx, name)
	return ctx, func() { span.End() }
}

// AddEvent records a named event on the active span.
func AddEvent(ctx context.Context, name string) {
	trace.SpanFromContext(ctx).AddEvent(name)
}

// AddError records an error on the active span and marks it as failed.
func AddError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// AddAttribute adds a key-value attribute to the active span.
// Supported value types: string, int, int64, float64, bool. Others are formatted as string.
func AddAttribute(ctx context.Context, key string, value any) {
	span := trace.SpanFromContext(ctx)
	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case int64:
		span.SetAttributes(attribute.Int64(key, v))
	case float64:
		span.SetAttributes(attribute.Float64(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	default:
		span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
	}
}

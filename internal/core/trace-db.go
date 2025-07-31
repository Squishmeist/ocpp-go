package core

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TraceDB(ctx context.Context, tracer trace.Tracer, operationName string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, operationName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.KeyValue{
			Key:   attribute.Key("db.name"),
			Value: attribute.StringValue("turso"),
		}),
	)

	return ctx, span
}

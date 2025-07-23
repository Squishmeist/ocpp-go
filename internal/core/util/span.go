package util

import (
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func JustErrWithSpan(span trace.Span, msg string, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
	span.End()
	return fmt.Errorf("%s: %w", msg, err)
}

// reusable helper for span error handling and error wrapping
func ErrWithSpan[T any](span trace.Span, msg string, err error) (T, error) {
	if err != nil {
		span.RecordError(err)
	}
	span.SetStatus(codes.Error, msg)
	span.End()
	var zero T
	return zero, fmt.Errorf("%s: %w", msg, err)
}

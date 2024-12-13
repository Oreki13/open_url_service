package helper

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func SetOtelError(span trace.Span, err error, stackTrace string) {
	span.RecordError(err)
	span.SetAttributes(
		attribute.String("error.stacktrace", stackTrace),
		attribute.String("error.message", err.Error()),
	)
}

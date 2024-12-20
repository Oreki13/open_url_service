package bootstrap

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"open_url_service/pkg/config"
	"open_url_service/pkg/tracer"
	"open_url_service/pkg/util"
)

type Provider struct {
	provider trace.TracerProvider
}

// RegistryOpenTelemetry NewExporter returns a new `Provider` type. It uses Jaeger exporter and globally sets
// the tracer provider as well as the global tracer for spans.
func RegistryOpenTelemetry(cfg *config.Config) (Provider, error) {
	if !cfg.AppOtelTrace {
		return Provider{
			provider: sdktrace.NewTracerProvider(),
		}, nil
	}

	var exp sdktrace.SpanExporter

	exp, err := tracer.NewExporter(cfg)

	if err != nil {
		return Provider{}, err
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.AppName),
			attribute.String("environment", util.EnvironmentTransform(cfg.AppEnv)),
			// attribute.Int64("ID", id),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return Provider{
		provider: tp,
	}, nil
}

// Close shuts down the tracer provider only if it was not "no operations"
// version.
func (p Provider) Close() error {
	if prv, ok := p.provider.(*sdktrace.TracerProvider); ok {
		return prv.Shutdown(context.TODO())
	}

	return nil
}

package tracer

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"open_url_service/internal/consts"
	"open_url_service/pkg/config"
)

func NewExporter(cfg *config.Config) (trace.SpanExporter, error) {
	switch cfg.AppOtelExporter {
	case consts.TempoExporter:
		return tempoExporter(cfg)
	default:
		return nil, errors.New("unknown otel driver")
	}
}

func tempoExporter(cfg *config.Config) (trace.SpanExporter, error) {
	insecureOpt := otlptracehttp.WithInsecure()
	endpointOpt := otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%v", cfg.TempoHost, cfg.TempoPort))

	return otlptracehttp.New(context.Background(), insecureOpt, endpointOpt)
}

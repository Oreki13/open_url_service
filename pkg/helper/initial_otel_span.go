package helper

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"open_url_service/internal/appctx"
	"open_url_service/pkg/tracer"
)

func InitialOtelSpan(xCtx appctx.Data) (context.Context, trace.Span) {
	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), xCtx.FiberCtx.Path(), nil)

	header, err := json.Marshal(xCtx.FiberCtx.GetReqHeaders())
	if err != nil {
		return ctx, span
	}

	span.SetAttributes(attribute.String("Host", xCtx.FiberCtx.Hostname()))
	span.SetAttributes(attribute.String("IP", xCtx.FiberCtx.IP()))
	span.SetAttributes(attribute.String("Headers", string(header)))

	return ctx, span
}

package trace

import (
	"context"

	oteltrace "go.opentelemetry.io/otel/trace"
)

type spanInfoKey struct{}

type SpanInfo struct {
	TraceID       string `json:"traceID"`
	SpanID        string `json:"spanID"`
	ServiceDomain string `json:"serviceDomain"`
	OperatorName  string `json:"operatorName"`
	StepName      string `json:"stepName"`
	ClientID      string `json:"clientID"`
	SystemID      string `json:"systemID"`
	From          string `json:"from"`
	To            string `json:"to"`
}

func GetSpanInfo(ctx context.Context) (spanInfo SpanInfo) {
	if span, ok := ctx.Value(spanInfoKey{}).(SpanInfo); ok {
		spanInfo = span
	}

	if span := oteltrace.SpanFromContext(ctx); span != nil {
		if span.SpanContext().HasTraceID() {
			spanInfo.TraceID = span.SpanContext().TraceID().String()
		}

		if span.SpanContext().HasSpanID() {
			spanInfo.SpanID = span.SpanContext().SpanID().String()
		}
	}

	return spanInfo
}

func InjectSpanInfo(ctx context.Context, spanInfo SpanInfo) context.Context {
	if span := oteltrace.SpanFromContext(ctx); span != nil {
		if span.SpanContext().HasTraceID() {
			spanInfo.TraceID = span.SpanContext().TraceID().String()
		}

		if span.SpanContext().HasSpanID() {
			spanInfo.SpanID = span.SpanContext().SpanID().String()
		}
	}

	return context.WithValue(ctx, spanInfoKey{}, spanInfo)
}

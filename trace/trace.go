package trace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

var (
	DefaultTracer     = newNoopsTrace()
	traceparentHeader = "traceparent"
)

// @Author: Nghiant5
// This function used for extract the traceparent header
// We can use this for insert to outbox database with field traceparent
// So we can take again the context through kafka message and tracing it on Jaegar
func ExtractTraceparent(ctx context.Context) string {
	m := make(propagation.MapCarrier)
	traceContext := propagation.TraceContext{}

	traceContext.Inject(ctx, m)

	h := m.Get(traceparentHeader)
	return h
}

func InjectTraceparent(ctx context.Context, traceparentId string) context.Context {
	// Create a map carrier to hold the propagation headers
	m := make(propagation.MapCarrier)
	m.Set(traceparentHeader, traceparentId)

	// Instantiate TraceContext propagator
	traceContext := propagation.TraceContext{}

	// Inject the traceparent into the context
	traceContext.Inject(ctx, m)

	// Extract the traceparent back into the context
	newCtx := traceContext.Extract(ctx, m)

	return newCtx
}

func SetDefaultTracer(trace Tracer) {
	DefaultTracer = trace
}

func ExtractSpanInfo(ctx context.Context) SpanInfo {
	return DefaultTracer.ExtractSpanInfo(ctx)
}

func StartTracing(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, SpanFinishFunc) {
	return DefaultTracer.StartTracing(ctx, spanName, opts...)
}

type SpanInfo struct {
	TraceID            string            `json:"traceID"`
	SpanID             string            `json:"spanID"`
	ClientIP           string            `json:"clientIP"`
	Protocol           string            `json:"protocol"`
	Method             string            `json:"method"`
	RequestPath        string            `json:"requestPath"`
	ServiceDomain      string            `json:"serviceDomain"`
	OperatorName       string            `json:"operatorName"`
	StepName           string            `json:"stepName"`
	UserAgent          string            `json:"userAgent"`
	ClientTime         int64             `json:"clientTime"`
	ReceivedTime       int64             `json:"receivedTime"`
	Hostname           string            `json:"hostname"`
	TransactionID      string            `json:"transactionID"`
	ContentLength      int               `json:"contentLength"`
	ClientID           string            `json:"clientID"`
	SystemID           string            `json:"systemID"`
	From               string            `json:"from"`
	To                 string            `json:"to"`
	Username           string            `json:"userName"`
	MessageType        string            `json:"messageType"`
	ReplyTo            []string          `json:"replyTo"`
	RequestHeaders     map[string]string `json:"requestHeaders"`
	RemoteHost         string            `json:"remoteHost"`
	XForwardedFor      string            `json:"xForwardedFor"`
	TransactionTimeout int64             `json:"transactionTimeout"`
}

type SpanFinishFunc func(context.Context, ...SpanFinishOption)

type SpanInfoKey struct{}

type Tracer interface {
	// Extract span info from the context
	ExtractSpanInfo(context.Context) SpanInfo

	// Used for tracing
	StartTracing(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, SpanFinishFunc)
}

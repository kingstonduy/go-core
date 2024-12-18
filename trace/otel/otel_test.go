package otel

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/trace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TestHttpTrace(t *testing.T) {
	var (
		ctx         = context.Background()
		serviceName = "test-service"
		method      = "GET"
		endpoint    = "http://google.com"
		exporter    = "localhost:4318"
	)

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	tracer, err := NewOpenTelemetryTracer(
		ctx,
		trace.WithTraceServiceName(serviceName),
		trace.WithTraceExporterEndpoint(exporter))

	if err != nil {
		t.Fail()
	}

	ctx, finish := tracer.StartTracing(
		ctx,
		"Ping to Google",
		trace.WithTraceRequest("test"),
	)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)

	finish(ctx, trace.WithTraceErrorResponse(err))

	// wait for push tracing
	time.Sleep(10 * time.Second)

}

func TestDefaultTrace(t *testing.T) {
	var (
		ctx         = context.Background()
		serviceName = "test-service"
		exporter    = "localhost:4318"
	)

	tracer, err := NewOpenTelemetryTracer(
		ctx,
		trace.WithTraceServiceName(serviceName),
		trace.WithTraceExporterEndpoint(exporter))

	if err != nil {
		t.Fail()
	}

	trace.SetDefaultTracer(tracer)
	ctx, finish := trace.StartTracing(
		ctx,
		"Ping to Google",
		trace.WithTraceRequest("Data test"),
	)

	finish(ctx, trace.WithTraceErrorResponse(err))

	// wait for push tracing
	time.Sleep(10 * time.Second)
}

package hresty

import (
	"context"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/logger/logrus"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
)

func TestHttpTracing(t *testing.T) {
	client := NewRestyClient()
	client.SetHeader("User-Agent", "test-service")

	logger.SetDefaultLogger(logrus.NewLogrusLogger())

	tracer, err := otel.NewOpenTelemetryTracer(
		context.Background(),
		trace.WithTraceExporterEndpoint("localhost:4318"),
		trace.WithTraceServiceName("test-service"))

	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	ctx, f := tracer.StartTracing(ctx, "Ping to google.com")

	req := client.R().SetContext(ctx)
	req.Get("http://localhost:4318") //nolint

	f(ctx)
	time.Sleep(10 * time.Second) //
}

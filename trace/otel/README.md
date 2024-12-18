#### Usage

```go

var (
		ctx         = context.Background()
		serviceName = "test-service"
		method      = "GET"
		endpoint    = "http://google.com"
		exporter    = "localhost:4318"
	)


	tracer, err := NewOpenTelemetryTracer(
		ctx,
		trace.WithTraceServiceName(serviceName),
		trace.WithTraceExporterEndpoint(exporter))

	if err != nil {
		t.Fail()
	}

	ctx, finish := tracer.StartHttpClientTrace(
		ctx,
		"Ping to Google",
	)
	defer finish(ctx)

```
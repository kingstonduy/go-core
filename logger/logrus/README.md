### Usage

```go
	Format:  Datetime UserAgent 8-traceID OperatorName TransactionID ReceivedTime ClientTime message
```

```go
// without tracing
func GetLogger() logger.Logger {
	levelStr := "LOG_LEVEL"
	level, err := logger.GetLevel(levelStr)
	if err != nil {
		level = logger.InfoLevel
	}
	return logrus.NewLogrusLogger(
		logger.WithLevel(level),
	)
}

// with tracing
func GetLogger(tracer trace.Tracer) logger.Logger {
	levelStr := "LOG_LEVEL"
	level, err := logger.GetLevel(levelStr)
	if err != nil {
		level = logger.InfoLevel
	}
	return logrus.NewLogrusLogger(
		logger.WithLevel(level),
		logger.WithTracer(tracer),
	)
}

// with hooks
broker := kafka.NewKafkaBroker(broker.WithBrokerAddresses("127.0.0.1:9092"))
broker.Connect()

tracer, _ := otel.NewOpenTelemetryTracer(context.Background(),
	trace.WithTraceServiceName("go_logrus_testing"),
)
trace.SetDefaultTracer(tracer)

log := NewLogrusLogger(
	logger.WithLevel(logger.TraceLevel),
	WithHooks(NewKafkaLoggingHook(
		broker,
		"go-core",
		WithKafkaLoggingHookLevels([]logger.Level{logger.InfoLevel, logger.TraceLevel}),
		WithKafkaLoggingHookTopic("mcs.logging.test1"))),
)

ctx := context.Background()
ctx, f := trace.StartTracing(ctx, "test logging kafka events")
defer f(ctx)

log.Fields(map[string]interface{}{
	logger.FIELD_OPERATOR_NAME: "test logging with hooks",
	logger.FIELD_STEP_NAME:     "push event to kafka",
}).Infof(ctx, "test events")

time.Sleep(100 * time.Millisecond)
```

```go
logger := GetLogger()
logger.Error(ctx, "Recovered from panic")
logger.Errorf(ctx, "Recovered from panic: %v", r)
```
package brokerLog

import (
	"context"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/logger/logrus"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
	"github.com/kingstonduy/go-core/transport/broker"
	"github.com/kingstonduy/go-core/transport/broker/kafka"
)

func TestWithDefaultBrokerHook(t *testing.T) {
	broker := kafka.NewKafkaBroker(broker.WithBrokerAddresses("127.0.0.1:9092"))
	if err := broker.Connect(); err != nil {
		t.Error(err)
	}

	tracer, _ := otel.NewOpenTelemetryTracer(context.Background(),
		trace.WithTraceServiceName("go_logrus_testing"),
	)
	trace.SetDefaultTracer(tracer)

	log := logrus.NewLogrusLogger(
		logger.WithLevel(logger.TraceLevel),
		logrus.WithHooks(NewBrokerLogHook(broker, "go-core")),
	)

	ctx := context.Background()
	ctx, f := trace.StartTracing(ctx, "test logging kafka events")
	defer f(ctx)

	log.Fields(map[string]interface{}{
		logger.FIELD_PUBLISHED:     true,
		logger.FIELD_OPERATOR_NAME: "test logging with hooks",
		logger.FIELD_STEP_NAME:     "push event to kafka",
	}).Fatal(ctx, "test events")

	time.Sleep(100 * time.Millisecond)
}

func TestWithBrokerHook(t *testing.T) {
	broker := kafka.NewKafkaBroker(broker.WithBrokerAddresses("127.0.0.1:9092"))
	if err := broker.Connect(); err != nil {
		t.Error(err)
	}

	tracer, _ := otel.NewOpenTelemetryTracer(context.Background(),
		trace.WithTraceServiceName("go_logrus_testing"),
	)
	trace.SetDefaultTracer(tracer)

	log := logrus.NewLogrusLogger(
		logger.WithLevel(logger.TraceLevel),
		logrus.WithHooks(NewBrokerLogHook(
			broker,
			"go-core",
			WithBrokerLogHookLevels([]logger.Level{logger.InfoLevel, logger.TraceLevel}),
			WithBrokerLogHookTopic("mcs.logging.test1"),
			WithBrokerLogEventCondition(func(event BrokerLogEvent) bool {
				published, ok := event[logger.FIELD_PUBLISHED].(bool)
				if !ok {
					return false
				}
				return published
			}),
		)),
	)

	ctx := context.Background()
	ctx, f := trace.StartTracing(ctx, "test logging kafka events")
	defer f(ctx)

	log.Fields(map[string]interface{}{
		logger.FIELD_PUBLISHED:     true,
		logger.FIELD_OPERATOR_NAME: "test logging with hooks",
		logger.FIELD_STEP_NAME:     "push event to kafka",
	}).Infof(ctx, "test events")

	time.Sleep(100 * time.Millisecond)
}

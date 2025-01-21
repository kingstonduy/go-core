package cmd_pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/logger/logrus"
	tracerr "github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
	"github.com/kingstonduy/go-core/transport"
)

func TestHappycase(t *testing.T) {
	tracer, _ := otel.NewOpenTelemetryTracer(
		context.Background(),
		tracerr.WithTraceServiceName("test-service"),
		tracerr.WithServiceVersion("test-version"),
	)

	logger.SetDefaultLogger(logrus.NewLogrusLogger(
		logger.WithTracer(tracer),
	))
	f := func(ctx context.Context, esCommand OutboxWithTrace) error {
		logger.Infof(ctx, "Handling command %s", esCommand.CommandType)
		logger.Info(ctx, esCommand.ToString())
		logger.Info(ctx, "Command handled successfully")
		time.Sleep(time.Second * 5)
		return nil
	}

	cmd := OutboxWithTrace{
		AggregateID: uuid.New().String(),
		CommandID:   uuid.New().String(),
		CommandType: "TestCommand",
		Payload:     "TestPayload",
		Trace: transport.Trace{
			From:     "test-from",
			To:       "test-to",
			Cid:      "test-client-id",
			Cts:      time.Now().UnixMilli(),
			Sid:      "test-sid",
			Username: "test-username",
		},
	}

	// Inject into the context

	dp := NewDispatcherCommandHandler()
	dp.RegisterHandler(cmd.CommandType, f)
	dp.When(context.Background(), cmd)
}

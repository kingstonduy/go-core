package logrus

import (
	"context"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
)

func TestMaskedData(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `{"username": "john_doe", "password": "123"} <password> 1234 </password> <credentials> 12345 </credentials> base64data: ZnNkZnNkZnNkZnNkZnNkZnNkZnNkZnNmc2RzZGZzZGZkc2ZzZGZzZGZzZGZmc2QK`,
			expectedOutput: `{"username": "john_doe", "password": "***"} <password> **** </password> <credentials> ***** </credentials> base64data: ****************************************************************`,
		},
		{
			input:          `{"password": "123"} <password>1234</password> <credentials>12345</credentials> base64data: XYZ123==`,
			expectedOutput: `{"password": "***"} <password>****</password> <credentials>*****</credentials> base64data: XYZ123==`,
		},
		{
			input:          `<![CDATA[ <password> 123 </password> ]]>, <![CDATA[ <credentials> user:1234 </credentials> ]]>, base64data: MNO456==`,
			expectedOutput: `<![CDATA[ <password> *** </password> ]]>, <![CDATA[ <credentials> ********* </credentials> ]]>, base64data: MNO456==`,
		},
		{
			input:          `base64data: ABCDEFGH12345==, <password>123</password>, {"password": "1234"}, <credentials>user:1234</credentials>`,
			expectedOutput: `base64data: ABCDEFGH12345==, <password>***</password>, {"password": "****"}, <credentials>*********</credentials>`,
		},
		{
			input:          `nothing`,
			expectedOutput: `nothing`,
		},
		// Add more test cases as needed
	}

	logger := NewLogrusLogger()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			logger.Info(ctx, tc.input)
		})
	}
}

func TestLogCommon(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "log information 1",
		}, {
			input: "log information 2",
		},
	}

	ctx := context.Background()
	logger := NewLogrusLogger()

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			logger.Info(ctx, tc.input)
		})
	}
}

func TestLogTracing(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "log information 1",
		}, {
			input: "log information 2",
		},
	}

	span := &trace.SpanInfo{
		ClientIP:      "123.123.123.123",
		Protocol:      "Kafka",
		ServiceDomain: "MCS Testing",
		OperatorName:  "Testing API",
		StepName:      "Testing",
		UserAgent:     "Postman",
		ClientTime:    123456789,
		ReceivedTime:  123456789,
		TransactionID: "1234-1234-1234-1234",
		Hostname:      "13.123.12.31",
		Method:        "GET",
		ContentLength: 123,
	}

	tracer, err := otel.NewOpenTelemetryTracer(context.Background(),
		trace.WithTraceServiceName("go_logrus_testing"),
		trace.WithTraceExporterEndpoint("localhost:4318"),
	)

	if err != nil {
		t.Error(err)
	}

	logger := NewLogrusLogger(
		logger.WithTracer(tracer),
	)

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			ctx := context.Background()
			ctx, finish := tracer.StartTracing(ctx, "testing",
				trace.WithTraceSpanInfo(span))

			logger.Info(ctx, tc.input)
			finish(ctx, trace.WithTraceResponse("test"))
		})
	}

	time.Sleep(10 * time.Second)
}

func TestDefaultLogging(t *testing.T) {
	ctx := context.Background()

	logger.SetDefaultLogger(NewLogrusLogger(
		logger.WithLevel(logger.TraceLevel),
	))

	logger.Logf(ctx, logger.InfoLevel, "logging: %s", "data")
	logger.Tracef(ctx, "logging: %s", "data")
	logger.Infof(ctx, "logging: %s", "data")
	logger.Warnf(ctx, "logging: %s", "data")
	logger.Debugf(ctx, "logging: %s", "data")
	logger.Errorf(ctx, "logging: %s", "data")
	logger.Fatalf(ctx, "logging: %s", "data")
}

func TestDefaultLoggingAndTracing(t *testing.T) {
	ctx := context.Background()

	logger.SetDefaultLogger(NewLogrusLogger(
		logger.WithLevel(logger.TraceLevel),
	))

	tracer, _ := otel.NewOpenTelemetryTracer(ctx)
	trace.SetDefaultTracer(tracer)
	ctx, f := trace.StartTracing(ctx, "start tracing")
	defer f(ctx)

	logger.Logf(ctx, logger.InfoLevel, "logging: %s", "data")
	logger.Tracef(ctx, "logging: %s", "data")
	logger.Infof(ctx, "logging: %s", "data")
	logger.Warnf(ctx, "logging: %s", "data")
	logger.Debugf(ctx, "logging: %s", "data")
	logger.Errorf(ctx, "logging: %s", "data")
	logger.Fatalf(ctx, "logging: %s", "data")
}

func TestMaskWithCustomerData(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `{"maskedData": "john_doe"}`,
			expectedOutput: `{"maskedData": "********"}`,
		},
	}

	patterns := []string{
		`\"maskedData\"\s*:\s*\"(.*?)\"`,
	}

	logger.SetDefaultLogger(NewLogrusLogger(
		logger.WithLevel(logger.TraceLevel),
		logger.WithMaskPatterns(patterns...),
	))

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			logger.Info(context.Background(), tc.input)
		})
	}
}

func TestWithFields(t *testing.T) {
	ctx := context.Background()
	tracer, _ := otel.NewOpenTelemetryTracer(ctx)
	trace.SetDefaultTracer(tracer)
	ctx, f := trace.StartTracing(ctx, "start tracing")
	defer f(ctx)

	log := NewLogrusLogger()
	log.Fields(map[string]interface{}{
		logger.FIELD_OPERATOR_NAME:   "operationName",
		logger.FIELD_STEP_NAME:       "stepName",
		logger.FIELD_DURATION:        250,
		logger.FIELD_STATUS_RESPONSE: "01",
	}).Infof(ctx, "test")
	time.Sleep(100 * time.Millisecond)
}

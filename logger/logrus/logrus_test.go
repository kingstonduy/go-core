package logrus

import (
	"context"

	"testing"

	"github.com/kingstonduy/go-core/logger"
	trace "github.com/kingstonduy/go-core/tracer"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()

	// enrich spanInfo
	spanInfo := trace.SpanInfo{
		TraceID:       "traceID",
		SpanID:        "spanID",
		ServiceDomain: "serviceDomain",
		OperatorName:  "operatorName",
		StepName:      "stepName",
		ClientID:      "clientID",
		SystemID:      "systemID",
		From:          "from",
		To:            "to",
	}

	ctx = trace.InjectSpanInfo(ctx, spanInfo)

	logger.InitLogger(NewLogrusLogger())
	for i := 1; i <= 1; i++ {
		logger.Info(ctx, "log by logrus logger ", i)
		logger.Infof(ctx, "log by logrus logger %d", i)
		logger.Debugf(ctx, "log by logrus logger %d", i)
		logger.Debug(ctx, "log by logrus logger ", i)
		logger.Warnf(ctx, "log by logrus logger %d", i)
		logger.Warn(ctx, "log by logrus logger ", i)
		logger.Errorf(ctx, "log by logrus logger %d", i)
		logger.Error(ctx, "log by logrus logger ", i)
	}

}

func TestLoggerWithFields(t *testing.T) {
	ctx := context.Background()
	for i := 1; i <= 1; i++ {
		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Info(ctx, "log by logrus logger ", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Infof(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Debugf(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Debug(ctx, "log by logrus logger ", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Warnf(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Warn(ctx, "log by logrus logger ", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Errorf(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Error(ctx, "log by logrus logger ", i)
	}

}

func TestLoggerWithFieldsAndSpanInfo(t *testing.T) {
	ctx := context.Background()

	// enrich spanInfo
	spanInfo := trace.SpanInfo{
		TraceID:       "traceID",
		SpanID:        "spanID",
		ServiceDomain: "serviceDomain",
		OperatorName:  "operatorName",
		StepName:      "stepName",
		ClientID:      "clientID",
		SystemID:      "systemID",
		From:          "from",
		To:            "to",
	}

	ctx = trace.InjectSpanInfo(ctx, spanInfo)

	logger.InitLogger(NewLogrusLogger())
	for i := 1; i <= 1; i++ {
		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Info(ctx, "log by logrus logger ", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Infof(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Debugf(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Debug(ctx, "log by logrus logger ", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Warnf(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Warn(ctx, "log by logrus logger ", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Errorf(ctx, "log by logrus logger %d", i)

		logger.Fields(map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: "operatorName",
			logger.FIELD_STEP_NAME:     "stepName",
			logger.FIELD_DURATION:      "duration",
		}).Error(ctx, "log by logrus logger ", i)
	}

}

package cmd_pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/metadata"
	"github.com/kingstonduy/go-core/metrics"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/transport"
)

type CommandHandlerFunc func(ctx context.Context, esCommand OutboxWithTrace) error

type DispatcherHandler interface {
	When(ctx context.Context, event OutboxWithTrace) error
	RegisterHandler(commandType string, handler CommandHandlerFunc)
}

var (
	ErrUnknownEventType  = errors.New("unknown command type")
	ErrMonitoringCommand = errors.New("monitoring request failed")

	MetricKeyRequestTotal    = []string{"command", "total"}
	MetricKeyRequestDuration = []string{"command", "duration", "milliseconds"}

	MetricLabelCommandType = "command_type"
	MetricLabelStatusCode  = "status_code"
	MetricLabelErrorCode   = "error_code"

	StepNameDone = "Message-Handled"
)

type DispatcherCommandHandler struct {
	eventHandlers map[string]CommandHandlerFunc
}

func NewDispatcherCommandHandler() DispatcherHandler {
	return &DispatcherCommandHandler{
		eventHandlers: make(map[string]CommandHandlerFunc),
	}
}

// RegisterHandler implements DispatcherHandler.
func (d *DispatcherCommandHandler) RegisterHandler(commandType string, handler CommandHandlerFunc) {
	_, exist := d.eventHandlers[commandType]
	if exist {
		logger.Errorf(context.TODO(), "Can not register handler for commandType: %s", commandType)
		return
	}
	d.eventHandlers[commandType] = handler
}

// When implements DispatcherHandler.
func (d *DispatcherCommandHandler) When(ctx context.Context, esCommand OutboxWithTrace) (err error) {
	ctx, finish := trace.StartTracing(ctx, fmt.Sprintf("Subcription.handleCommand - %s", esCommand.CommandType), trace.WithTraceRequest(esCommand))
	defer func() {
		if pa := recover(); pa != nil {
			logger.Errorf(ctx, "Command is PANIC !! - %+v", pa)
		}
		finish(ctx,
			trace.WithTraceErrorResponse(err),
		)
		d.loggingTimeCommandExcuted(ctx, esCommand, err)
		d.emitMetric(ctx, esCommand, err)
	}()
	ctx = d.monitoringCommand(ctx, esCommand)
	handler, found := d.eventHandlers[string(esCommand.CommandType)]
	if !found {
		logger.Errorf(ctx, "There is no handler for command %s -  AggregateId: %s, CommandType: %s, CommandId: %s", esCommand.AggregateID, esCommand.CommandType, esCommand.CommandID)
		return errors.Wrapf(ErrUnknownEventType, "commandType: %s", esCommand.CommandType)
	}

	return handler(ctx, esCommand)
}

// Setting all trace information to context
func (d *DispatcherCommandHandler) monitoringCommand(ctx context.Context, esCommand OutboxWithTrace) context.Context {
	ctx = transport.MonitorCommand(ctx, transport.MonitorRequestData{
		Protocol:      metadata.ProtocolKafka,
		Method:        "subscribe",
		Hostname:      "",
		ServiceDomain: "",
		UserAgent:     "",
		RequestPath:   esCommand.CommandType,
		// ContentLength: len(payload),
		From:               esCommand.Trace.From,
		To:                 esCommand.Trace.To,
		ClientID:           esCommand.Trace.Cid,
		ClientTime:         esCommand.Trace.Cts,
		Username:           esCommand.Trace.Username,
		Request:            esCommand,
		SystemID:           esCommand.Trace.Sid,
		TransactionTimeout: esCommand.Trace.TransactionTimeout,
		// RequestHeaders: reqHeaders,
	})
	return ctx
}

// It will automatically logging Error if your use case return a error.
// So you just need to return a error in usecase
func (d *DispatcherCommandHandler) loggingTimeCommandExcuted(ctx context.Context, esCommand OutboxWithTrace, err error) {
	spanInfo := trace.ExtractSpanInfo(ctx)
	timeDuration := time.Now().UnixMilli() - spanInfo.ReceivedTime
	codeResponse := errorx.DefaultSuccessResponseCode
	if err != nil {
		var commandError *errorx.Error
		if errors.As(err, &commandError) {
			codeResponse = commandError.Code
		}
	}
	logger.Fields(map[string]interface{}{
		logger.FIELD_OPERATOR_NAME:   spanInfo.RequestPath,
		logger.FIELD_STEP_NAME:       StepNameDone,
		logger.FIELD_DURATION:        timeDuration,
		logger.FIELD_STATUS_RESPONSE: codeResponse,
	}).Infof(ctx, "Executed Command - %s, with Error: %v - Duration: %d", esCommand.ToString(), err, timeDuration)

}

// Write metric of command
func (d *DispatcherCommandHandler) emitMetric(ctx context.Context, esCommand OutboxWithTrace, err error) {
	commandType := esCommand.CommandType
	var statusCode int
	var errorCode string

	if err != nil {
		var commandError *errorx.Error
		if errors.As(err, &commandError) {
			statusCode = commandError.Status
			errorCode = commandError.Code
		}
	} else {
		statusCode = errorx.DefaultSuccessStatusCode
		errorCode = errorx.DefaultSuccessResponseCode
	}

	metrics.IncrCounterWithLabels(
		MetricKeyRequestTotal,
		1,
		[]metrics.Label{
			{
				Name:  MetricLabelCommandType,
				Value: commandType,
			},
			{
				Name:  MetricLabelStatusCode,
				Value: fmt.Sprintf("%v", statusCode),
			},
			{
				Name:  MetricLabelErrorCode,
				Value: fmt.Sprintf("%v", errorCode),
			},
		},
	)
	spanInfo := trace.ExtractSpanInfo(ctx)
	timeReceived := time.UnixMilli(spanInfo.ReceivedTime)
	duration := time.Since(timeReceived)
	metrics.AddSampleWithLabels(
		MetricKeyRequestDuration,
		float32(duration.Milliseconds()),
		[]metrics.Label{
			{
				Name:  MetricLabelCommandType,
				Value: commandType,
			},
		},
	)

}

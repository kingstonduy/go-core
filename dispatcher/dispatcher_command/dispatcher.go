package dispatcherCommand

import (
	"context"
	"fmt"
	"time"

	"github.com/kingstonduy/go-core/dispatcher"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/metadata"
	"github.com/kingstonduy/go-core/metrics"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/transport"
	"github.com/kingstonduy/go-core/transport/broker"
	"github.com/pkg/errors"
)

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

// @Author: Nghiant5
// This is used for dispatching by routing commandType to correct handler.
type DispatcherCommandHandler struct {
	eventHandlers map[string]dispatcher.CommandHandlerFunc
}

func NewDispatcherCommandHandler() dispatcher.DispatcherHandler {
	return &DispatcherCommandHandler{
		eventHandlers: make(map[string]dispatcher.CommandHandlerFunc),
	}
}

// To use dispatcher we should register the handler function for it
// The way to register is in readme of this current folder.
func (d *DispatcherCommandHandler) RegisterHandler(commandType string, handler dispatcher.CommandHandlerFunc) {
	_, exist := d.eventHandlers[commandType]
	if exist {
		logger.Errorf(context.TODO(), "Can not register handler for commandType: %s", commandType)
		return
	}
	d.eventHandlers[commandType] = handler
}

// This function will automatically tracing,logging, routing messageCommand to correct handler
// That we previously registered.
func (d *DispatcherCommandHandler) When(ctx context.Context, esCommand transport.Command, broker broker.Broker) (err error) {

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
		logger.Errorf(ctx, "There is no handler for command %s -  AggregateId: %s, AggregateType: %s ,CommandType: %s, CommandId: %s", esCommand.AggregateID, esCommand.AggregateType, esCommand.CommandType, esCommand.CommandID)
		return errors.Wrapf(ErrUnknownEventType, "commandType: %s", esCommand.CommandType)
	}

	return handler(ctx, esCommand)

}

// Setting all trace information to context
func (d *DispatcherCommandHandler) monitoringCommand(ctx context.Context, esCommand transport.Command) context.Context {
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
func (d *DispatcherCommandHandler) loggingTimeCommandExcuted(ctx context.Context, esCommand transport.Command, err error) {
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
	}).Infof(ctx, "Executed Command - %s, with Error: %v - Duration: %d", esCommand.StringNoPayload(), err, timeDuration)

}

// Write metric of command
func (d *DispatcherCommandHandler) emitMetric(ctx context.Context, esCommand transport.Command, err error) {
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

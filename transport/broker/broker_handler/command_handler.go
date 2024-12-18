package brokerHandler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kingstonduy/go-core/dispatcher"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/transport"
	"github.com/kingstonduy/go-core/transport/broker"
)

// @Author: Nghiant5
// This function is used for handle message with type command
// This uses dispatcher to handle the message
func HandleCommandEvent(ctx context.Context, event broker.Event, dispatcher dispatcher.DispatcherHandler, broker broker.Broker) (err error) {

	if event.Message() == nil || len(event.Message().Body) == 0 {
		logger.Errorf(ctx, "Topic: %s. Empty messsage body", event.Topic())
	}
	cm := transport.NewDefaultCommand()
	if err = json.Unmarshal(event.Message().Body, &cm); err != nil {
		logger.Errorf(ctx, "Topic: %s. Can Not Unmarshal messsage body", event.Topic())
		return err
	}

	if err = dispatcher.When(ctx, cm, broker); err != nil {
		return err
	}

	return nil
}

// This function used for publish response command
// it will automatically define wheather it has an error
// and create a result message type and publish
func PublishResponseCommand(ctx context.Context, broker1 broker.Broker, topic string,
	esCommand transport.Command, opts ...transport.ResponseOption) error {
	options := transport.ResponseOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	esCommand.Forward = ""

	if options.From != nil {
		esCommand.Trace.From = *options.From
	}
	if options.To != nil {
		esCommand.Trace.To = *options.To
	}

	if options.ReceivedTime > 0 {
		esCommand.Trace.Dur = esCommand.Trace.Sts - options.ReceivedTime
	}

	if esCommand.Result == nil {
		res := transport.DefaultSuccessResponse
		esCommand.Result = &res.Result
	}

	//if the error you pass not nil, it will overide the result of esCommand
	// if it have type Errorx it will set corresponding code that you defined
	// otherwise it will regconise internal error and publish code 99
	if options.Error != nil {
		err := options.Error

		failedError := errorx.Failed(err.Error())

		esCommand.Result.StatusCode = failedError.Status
		esCommand.Result.Code = failedError.Code
		esCommand.Result.Message = failedError.Message

		var definedError *errorx.Error
		if errors.As(err, &definedError) {
			esCommand.Result.StatusCode = definedError.Status
			esCommand.Result.Code = definedError.Code
			esCommand.Result.Message = definedError.Message
			esCommand.Result.Details = definedError.Details
		}
	}
	if options.IsResponseEmpty {
		esCommand.Payload = ""
	}

	msg, err := json.Marshal(esCommand)
	if err != nil {
		return errorx.NewError(200, "99", "Internal Error: %v", err)
	}
	kMsg := broker.Message{
		Body: msg,
		Key:  []byte(esCommand.AggregateID),
	}

	err = broker1.Publish(ctx, topic, &kMsg)
	if err != nil {
		jsonBytes, _ := json.Marshal(esCommand)
		logger.Fields(
			map[string]interface{}{
				logger.FIELD_OPERATOR_NAME: topic,
				logger.FIELD_STEP_NAME:     "published-message-successully",
			},
		).Error(ctx, string(jsonBytes))
		return errorx.NewError(200, "99", "Internal Error: %v", err)
	}

	logger.Fields(
		map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: topic,
			logger.FIELD_STEP_NAME:     "published-message-successully",
		},
	).Info(ctx, broker.MakeStringLogsKafka(ctx, kMsg))
	return nil
}

func PublishRequestCommand(ctx context.Context, broker1 broker.Broker, topic string,
	esCommand transport.Command, opts ...transport.ResponseOption) error {

	options := transport.ResponseOptions{}

	for _, opt := range opts {
		opt(&options)
	}
	esCommand.Forward = ""

	if options.From != nil {
		esCommand.Trace.From = *options.From
	}
	if options.To != nil {
		esCommand.Trace.To = *options.To
	}
	esCommand.Trace.Cts = time.Now().UnixMilli()

	msg, err := json.Marshal(esCommand)
	if err != nil {
		return errorx.NewError(200, "99", "Internal Error: %v", err)
	}
	kMsg := broker.Message{
		Body: msg,
		Key:  []byte(esCommand.AggregateID),
	}

	err = broker1.Publish(ctx, topic, &kMsg)
	if err != nil {
		logger.Errorf(ctx, "Can not publish request message - Error: %v - %s ", err, esCommand.String())
		return errorx.NewError(200, "99", "Internal Error: %v", err)
	}
	logger.Fields(
		map[string]interface{}{
			logger.FIELD_OPERATOR_NAME: topic,
			logger.FIELD_STEP_NAME:     "published-message-successully",
		},
	).Info(ctx, broker.MakeStringLogsKafka(ctx, kMsg))
	return nil
}

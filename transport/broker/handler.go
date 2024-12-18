package broker

import (
	"context"
	"encoding/json"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/metadata"
	"github.com/kingstonduy/go-core/pipeline"
	"github.com/kingstonduy/go-core/transport"
)

// Handle request using Pipeline
// Step 1: Parse request from broker event
// Step 2: Send request to Pipeline and receive response
// Step 3: Parse Pipeline response to broker response
// Step 4: If broker response provided and reply header provided, response message to response broker
// error: system error, not API error
func HandleBrokerEvent[TReq any, TRes any](ctx context.Context, e Event, opts ...BrokerEventHandlerOption) error {
	options := NewBrokerEventHandlerOptions(opts...)

	// validate request message
	if e.Message() == nil || len(e.Message().Body) == 0 {
		logger.Infof(ctx, "Topic: %s. Empty message body", e.Topic())
		return EmptyMessageError{}
	}

	// parse the request message
	body, headers := e.Message().Body, e.Message().Headers
	var request transport.Request[TReq]
	err := json.Unmarshal(body, &request)
	if err != nil {
		logger.Infof(ctx, "Topic: %s. Invalid message format", e.Topic())
		return InvalidDataFormatError{}
	}

	trace := request.Trace
	ctx = transport.MonitorRequest(ctx, transport.MonitorRequestData{
		Protocol:       metadata.ProtocolKafka,
		Method:         "subscribe",
		Hostname:       "",
		ServiceDomain:  "",
		UserAgent:      "",
		RequestPath:    e.Topic(),
		ContentLength:  len(body),
		From:           trace.From,
		To:             trace.To,
		ClientID:       trace.Cid,
		ClientTime:     trace.Cts,
		Username:       trace.Username,
		Request:        request,
		MessageType:    headers[metadata.HeaderMessageType],
		RequestHeaders: headers,
		SystemID:       trace.Sid,
	})

	// handle request
	res := handleRequestPipeline[TReq, TRes](ctx, request)

	if options.OnRequestHandledFunc != nil {
		go func() {
			options.OnRequestHandledFunc(ctx, res)
		}()
	}

	return nil
}

func handleRequestPipeline[TReq any, TRes any](ctx context.Context, req transport.Request[TReq]) transport.Response[TRes] {
	res, err := pipeline.Send[TReq, TRes](ctx, req.Data)
	brokerRes := transport.GetResponse[TRes](
		ctx,
		transport.WithError(err),
		transport.WithData(res),
	)

	return brokerRes
}

type BrokerEventHandlerOption func(*BrokerEventHandlerOptions)

type BrokerEventHandlerOptions struct {
	OnRequestHandledFunc func(ctx context.Context, res interface{})
}

func NewBrokerEventHandlerOptions(opts ...BrokerEventHandlerOption) BrokerEventHandlerOptions {
	options := BrokerEventHandlerOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

// Handle after the request handled. res type: transport.Response[T]
func WithOnRequestHandledFunc(f func(ctx context.Context, res interface{})) BrokerEventHandlerOption {
	return func(opts *BrokerEventHandlerOptions) {
		opts.OnRequestHandledFunc = f
	}
}

func Pop[T any](arr []T) ([]T, T) {
	if len(arr) == 0 {
		return arr, *new(T)
	}

	popped := arr[len(arr)-1]
	poppedArr := arr[:len(arr)-1]

	return poppedArr, popped
}

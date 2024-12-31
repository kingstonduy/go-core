package transport

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/trace"
)

type MonitorRequestOptions struct {
	tracer trace.Tracer
	logger logger.Logger
}

type MonitorRequestOption func(*MonitorRequestOptions)

// tracer to inject span info into context. Default: trace.DefaultTrace
func WithTracer(tracer trace.Tracer) MonitorRequestOption {
	return func(o *MonitorRequestOptions) {
		o.tracer = tracer
	}
}

// logging to log request. Default: logger.DefaultLogger
func WithLogger(logger logger.Logger) MonitorRequestOption {
	return func(o *MonitorRequestOptions) {
		o.logger = logger
	}
}

func apply(opts ...MonitorRequestOption) MonitorRequestOptions {
	options := MonitorRequestOptions{
		// default options
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type MonitorRequestData struct {
	// Client IP address
	ClientIP string

	// Request protocol
	// Ex: HTTP, KAFKA, GRPC, IBM-MQ
	Protocol string

	// Request method
	// Ex:
	// 	HTTP: GET, POST, PUT, DELETE, OPTION,...
	Method string

	// Request path
	// Ex:
	// 	HTTP: /api/v1/customers
	// 	KAFKA: topicsName
	RequestPath string

	// Service name
	ServiceDomain string

	// Request user agent
	UserAgent string

	// Request client time
	ClientTime int64

	// Hostname
	Hostname string

	// Transaction ID
	TransactionID string

	// Request content length
	ContentLength int

	// Request body
	Request interface{}

	// Client ID
	ClientID string

	// From service
	From string

	// To service
	To string

	// Request username
	Username string

	// Request message type
	MessageType string

	// Reply to destinations
	ReplyTo []string

	// Request headers
	RequestHeaders map[string]string

	// Remote host
	RemoteHost string

	// X Forwarded For
	XForwardedFor string

	//system id
	SystemID string

	//client transaction timeout for request
	TransactionTimeout int64
}

// Start monitoring a request based on the given options, and return a context with some tracing data
func MonitorRequest(ctx context.Context, data MonitorRequestData, opts ...MonitorRequestOption) context.Context {
	options := apply(opts...)

	var (
		optTracer = options.tracer
		optLogger = options.logger
	)

	// tracing
	if tracer := getMonitoringTracer(optTracer); tracer != nil {
		systemId := data.SystemID
		if len(data.SystemID) == 0 || data.SystemID == "" {
			systemId = uuid.New().String()
		}
		spanInfo := &trace.SpanInfo{
			ClientIP:           data.ClientIP,
			Protocol:           data.Protocol,
			Method:             data.Method,
			RequestPath:        data.RequestPath,
			ServiceDomain:      data.ServiceDomain,
			UserAgent:          data.UserAgent,
			ClientTime:         data.ClientTime,
			Hostname:           data.Hostname,
			TransactionID:      data.TransactionID,
			ContentLength:      data.ContentLength,
			RequestHeaders:     data.RequestHeaders,
			RemoteHost:         data.RemoteHost,
			XForwardedFor:      data.XForwardedFor,
			ClientID:           data.ClientID,
			SystemID:           systemId,
			From:               data.From,
			To:                 data.To,
			Username:           data.Username,
			ReceivedTime:       time.Now().UnixMilli(),
			MessageType:        data.MessageType,
			ReplyTo:            data.ReplyTo,
			TransactionTimeout: data.TransactionTimeout,
		}

		traceCtx, f := tracer.StartTracing(
			ctx, "request",
			trace.WithTraceRequest(data.Request),
			trace.WithTraceSpanInfo(spanInfo),
		)

		ctx = traceCtx
		defer f(ctx)
	}

	// logging
	if log := getMonitoringLogger(optLogger); log != nil {
		b, _ := json.Marshal(data.Request)
		log.Fields(
			map[string]interface{}{
				logger.FIELD_OPERATOR_NAME: data.RequestPath,
				logger.FIELD_STEP_NAME:     "request",
			},
		).Infof(ctx, string(b))
	}

	return ctx
}

// Get transport response. This function will also log the response
func GetResponse[T any](ctx context.Context, opts ...ResponseOption) Response[T] {
	options := ResponseOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	res := DefaultSuccessResponse
	nowMili := time.Now().UnixMilli()

	var spanInfo trace.SpanInfo
	if tracer := getResponseTracer(options.Tracer); tracer != nil {
		spanInfo = tracer.ExtractSpanInfo(ctx)
		res.Trace = Trace{
			From:     spanInfo.To,
			To:       spanInfo.From,
			Cts:      spanInfo.ClientTime,
			Sts:      spanInfo.ReceivedTime,
			Cid:      spanInfo.ClientID,
			Sid:      spanInfo.SystemID,
			Dur:      nowMili - spanInfo.ReceivedTime,
			Username: spanInfo.Username,
		}
	}

	if options.Trace != nil {
		res.Trace = *options.Trace
		res.Trace.Dur = nowMili - res.Trace.Sts
	}

	if options.Error != nil {
		err := options.Error

		failedError := errorx.Failed(err.Error())

		res.Result.StatusCode = failedError.Status
		res.Result.Code = failedError.Code
		res.Result.Message = failedError.Message

		var definedError *errorx.Error
		if errors.As(err, &definedError) {
			res.Result.StatusCode = definedError.Status
			res.Result.Code = definedError.Code
			res.Result.Message = definedError.Message
			res.Result.Details = definedError.Details
		}
	}

	var data T
	if options.Data != nil {
		if d, ok := options.Data.(T); ok {
			data = d
		}
	}

	resp := Response[T]{
		Result: res.Result,
		Trace:  res.Trace,
		Data:   data,
	}

	// logging response
	respBytes, _ := json.Marshal(resp)
	logger.Fields(map[string]interface{}{
		logger.FIELD_OPERATOR_NAME:   spanInfo.RequestPath,
		logger.FIELD_STEP_NAME:       "response",
		logger.FIELD_DURATION:        resp.Trace.Dur,
		logger.FIELD_STATUS_RESPONSE: resp.Result.Code,
	}).Infof(ctx, string(respBytes))

	return resp
}

func getResponseTracer(tracer trace.Tracer) trace.Tracer {
	if tracer != nil {
		return tracer
	}
	return trace.DefaultTracer
}

func getMonitoringLogger(l logger.Logger) logger.Logger {
	if l != nil {
		return l
	}
	return logger.DefaultLogger
}

func getMonitoringTracer(t trace.Tracer) trace.Tracer {
	if t != nil {
		return t
	}
	return trace.DefaultTracer
}

func GetTraceByCtx(ctx context.Context, opts ...ResponseOption) Trace {
	options := ResponseOptions{}

	for _, opt := range opts {
		opt(&options)
	}
	if tracer := getResponseTracer(options.Tracer); tracer != nil {
		spanInfo := tracer.ExtractSpanInfo(ctx)
		nowMili := time.Now().UnixMilli()
		trace := Trace{
			From:     spanInfo.From,
			To:       spanInfo.To,
			Cts:      spanInfo.ClientTime,
			Sts:      spanInfo.ReceivedTime,
			Cid:      spanInfo.ClientID,
			Sid:      spanInfo.SystemID,
			Dur:      nowMili - spanInfo.ReceivedTime,
			Username: spanInfo.Username,
			// MessageType: spanInfo.MessageType,
			// ReplyTo:     spanInfo.ReplyTo,
			TransactionTimeout: spanInfo.TransactionTimeout,
		}

		return trace
	}

	return Trace{}
}

type MessageHandler interface {
	SetTrace(trace Trace)
	GetTrace() Trace
}

// Start monitoring a request based on the given options, and return a context with some tracing data
func MonitorCommand(ctx context.Context, data MonitorRequestData, opts ...MonitorRequestOption) context.Context {
	// tracing

	systemId := data.SystemID
	if len(data.SystemID) == 0 || data.SystemID == "" {
		systemId = uuid.New().String()
	}
	spanInfo := &trace.SpanInfo{
		ClientIP:           data.ClientIP,
		Protocol:           data.Protocol,
		Method:             data.Method,
		RequestPath:        data.RequestPath,
		ServiceDomain:      data.ServiceDomain,
		UserAgent:          data.UserAgent,
		ClientTime:         data.ClientTime,
		Hostname:           data.Hostname,
		TransactionID:      data.TransactionID,
		ContentLength:      data.ContentLength,
		RequestHeaders:     data.RequestHeaders,
		RemoteHost:         data.RemoteHost,
		XForwardedFor:      data.XForwardedFor,
		ClientID:           data.ClientID,
		SystemID:           systemId,
		From:               data.From,
		To:                 data.To,
		Username:           data.Username,
		ReceivedTime:       time.Now().UnixMilli(),
		MessageType:        data.MessageType,
		ReplyTo:            data.ReplyTo,
		TransactionTimeout: data.TransactionTimeout,
	}
	ctx = context.WithValue(ctx, trace.SpanInfoKey{}, *spanInfo)

	return ctx
}

func GetHttpRequestHeaderByCtx(ctx context.Context, opts ...ResponseOption) map[string]string {
	options := ResponseOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	if tracer := getResponseTracer(options.Tracer); tracer != nil {
		spanInfo := tracer.ExtractSpanInfo(ctx)
		return spanInfo.RequestHeaders
	}

	return map[string]string{}
}

package transport

import (
	"github.com/kingstonduy/go-core/trace"
)

type ResponseOptions struct {
	Trace                *Trace
	Data                 interface{}
	Error                error
	Tracer               trace.Tracer
	From                 *string
	To                   *string
	TransactionCreatedAt int64
	ReceivedTime         int64
	ReplyTo              *string
	Result               *Result
	IsResponseEmpty      bool
}

type ResponseOption func(*ResponseOptions)

// The request header mapped to response headers
func WithTrace(trace *Trace) ResponseOption {
	return func(options *ResponseOptions) {
		options.Trace = trace
	}
}

// Data of the response
func WithData(data interface{}) ResponseOption {
	return func(options *ResponseOptions) {
		options.Data = data
	}
}

// With Error mark a response as err response
func WithError(err error) ResponseOption {
	return func(options *ResponseOptions) {
		options.Error = err
	}
}

// With Tracer used for extracting the span info from the context
func WithTraceExtractor(tracer trace.Tracer) ResponseOption {
	return func(options *ResponseOptions) {
		options.Tracer = tracer
	}
}

// This is used for set attribute from in trace
func WithFromPublisher(frm string) ResponseOption {
	return func(options *ResponseOptions) {
		options.From = &frm
	}
}

// This is used for set attribute to in trace
func WithToReceiver(to string) ResponseOption {
	return func(options *ResponseOptions) {
		options.To = &to
	}
}

func WithTransactionCreatedAt(createdAt int64) ResponseOption {
	return func(options *ResponseOptions) {
		options.TransactionCreatedAt = createdAt
	}
}

func WithRequestReceivedAt(receivedAt int64) ResponseOption {
	return func(options *ResponseOptions) {
		options.ReceivedTime = receivedAt
	}
}

func WithReplyTo(replyTo string) ResponseOption {
	return func(options *ResponseOptions) {
		options.ReplyTo = &replyTo
	}
}

func WithResult(result *Result) ResponseOption {
	return func(options *ResponseOptions) {
		options.Result = result
	}
}

// The request header mapped to response headers
func WithIsResponseEmpty(isEmpty bool) ResponseOption {
	return func(options *ResponseOptions) {
		options.IsResponseEmpty = isEmpty
	}
}

package trace

import (
	"crypto/tls"
	"time"

	"github.com/kingstonduy/go-core/metadata"
)

type TraceOptions struct {
	ServiceName       string
	ServiceVersion    string
	ExporterEndpoint  string
	ExporterHeaders   map[string]string
	ExporterTLSConfig *tls.Config
}

type TraceOption func(*TraceOptions)

func WithTraceServiceName(serviceName string) TraceOption {
	return func(options *TraceOptions) {
		options.ServiceName = serviceName
	}
}

func WithServiceVersion(serviceVersion string) TraceOption {
	return func(options *TraceOptions) {
		options.ServiceVersion = serviceVersion
	}
}

func WithTraceExporterEndpoint(ep string) TraceOption {
	return func(options *TraceOptions) {
		options.ExporterEndpoint = ep
	}
}

func WithTraceExporterHeaders(headers map[string]string) TraceOption {
	return func(options *TraceOptions) {
		if options.ExporterHeaders == nil {
			options.ExporterHeaders = make(map[string]string)
		}

		for k, v := range headers {
			options.ExporterHeaders[k] = v
		}
	}
}

func WithTraceExporterHeader(key, value string) TraceOption {
	return func(options *TraceOptions) {
		if options.ExporterHeaders == nil {
			options.ExporterHeaders = make(map[string]string)
		}

		options.ExporterHeaders[key] = value
	}
}

func WithTraceExporterTLSConfig(cfg *tls.Config) TraceOption {
	return func(options *TraceOptions) {
		options.ExporterTLSConfig = cfg
	}
}

func NewTraceOptions(opts ...TraceOption) TraceOptions {
	options := TraceOptions{
		ExporterHeaders: make(map[string]string),
		ServiceName:     metadata.DefaultServiceName,
		ServiceVersion:  metadata.DefaultServiceVersion,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

// span options
type SpanStartOption func(*SpanStartOptions)

type SpanStartOptions struct {
	Request   interface{}
	StartTime time.Time
	SpanInfo  *SpanInfo
}

func WithTraceRequest(request interface{}) SpanStartOption {
	return func(options *SpanStartOptions) {
		options.Request = request
	}
}

func WithTraceStartTime(startTime time.Time) SpanStartOption {
	return func(options *SpanStartOptions) {
		options.StartTime = startTime
	}
}

func WithTraceSpanInfo(span *SpanInfo) SpanStartOption {
	return func(options *SpanStartOptions) {
		options.SpanInfo = span
	}
}

type SpanFinishOption func(*SpanFinishOptions)

type SpanFinishOptions struct {
	Response   interface{}
	Error      error
	FinishTime time.Time
}

func WithTraceResponse(response interface{}) SpanFinishOption {
	return func(options *SpanFinishOptions) {
		options.Response = response
	}
}

func WithTraceErrorResponse(err error) SpanFinishOption {
	return func(options *SpanFinishOptions) {
		options.Error = err
	}
}

func WithTraceFinishTime(finishTime time.Time) SpanFinishOption {
	return func(options *SpanFinishOptions) {
		options.FinishTime = finishTime
	}
}

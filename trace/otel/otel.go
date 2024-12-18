package otel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kingstonduy/go-core/metadata"
	"github.com/kingstonduy/go-core/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	TracerName           = "mcs-tracer"
	DefaultOperationName = "mcs-tracing"
)

type openTelemetryTracer struct {
	options trace.TraceOptions
}

func NewOpenTelemetryTracer(ctx context.Context, opts ...trace.TraceOption) (trace.Tracer, error) {
	options := trace.TraceOptions{
		ServiceName: metadata.DefaultServiceName,
	}

	for _, opt := range opts {
		opt(&options)
	}

	tracer := openTelemetryTracer{
		options: options,
	}

	// set global config trace
	if err := tracer.setGlobalTracer(ctx); err != nil {
		return nil, err
	}

	return &tracer, nil
}

// ExtractSpanInfo implements trace.Tracer.
func (tracer *openTelemetryTracer) ExtractSpanInfo(ctx context.Context) trace.SpanInfo {
	var spanInfo trace.SpanInfo
	if span, ok := ctx.Value(trace.SpanInfoKey{}).(trace.SpanInfo); ok {
		spanInfo = span
	}

	if span := oteltrace.SpanFromContext(ctx); span != nil {
		if span.SpanContext().HasTraceID() {
			spanInfo.TraceID = span.SpanContext().TraceID().String()
		}

		if span.SpanContext().HasSpanID() {
			spanInfo.SpanID = span.SpanContext().SpanID().String()
		}
	}

	if spanInfo.ServiceDomain == "" {
		spanInfo.ServiceDomain = tracer.options.ServiceName
	}

	return spanInfo
}

// StartInternalTrace implements trace.Tracer.
func (*openTelemetryTracer) StartTracing(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.SpanFinishFunc) {
	options := trace.SpanStartOptions{
		StartTime: time.Now(),
	}
	for _, opt := range opts {
		opt(&options)
	}

	tr := otel.Tracer(TracerName)

	if spanName == "" {
		spanName = DefaultOperationName
	}

	ctx, span := tr.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindInternal))

	// start time attribute
	start := options.StartTime
	span.SetAttributes(attribute.String(trace.TracingAttributesStartTime, start.Format(time.RFC1123Z)))

	// request attribute
	if req := options.Request; req != nil {
		reqJson, _ := json.Marshal(req)
		span.SetAttributes(attribute.String(trace.TracingAttributesRequest, string(reqJson)))
	}

	if span := options.SpanInfo; span != nil {
		ctx = context.WithValue(ctx, trace.SpanInfoKey{}, *span)
	}

	return ctx, func(ctx context.Context, opts ...trace.SpanFinishOption) {
		options := trace.SpanFinishOptions{
			FinishTime: time.Now(),
		}
		for _, opt := range opts {
			opt(&options)
		}

		span := oteltrace.SpanFromContext(ctx)
		if span == nil {
			return
		}

		finish := options.FinishTime
		span.SetAttributes(attribute.String(trace.TracingAttributesFinishTime, finish.Format(time.RFC1123Z)))

		if res := options.Response; res != nil {
			resJson, _ := json.Marshal(res)
			span.SetAttributes(attribute.String(trace.TracingAttributesResponse, string(resJson)))
		}

		if err := options.Error; err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}
}

// Setup global tracing configurations
func (o *openTelemetryTracer) setGlobalTracer(ctx context.Context) error {
	traceProviderOptions := make([]tracesdk.TracerProviderOption, 0)

	// default options
	traceProviderOptions = append(
		traceProviderOptions,
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(o.options.ServiceName),
			semconv.ServiceVersionKey.String(o.options.ServiceVersion),
		)))

	// exporter options
	if o.exportEnabled() {
		exporter, err := o.newExporter(ctx)
		if err != nil {
			return err
		}
		traceProviderOptions = append(traceProviderOptions, tracesdk.WithBatcher(exporter))
	}

	tp := tracesdk.NewTracerProvider(traceProviderOptions...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

func (o *openTelemetryTracer) newExporter(ctx context.Context) (tracesdk.SpanExporter, error) {
	if !o.exportEnabled() {
		return nil, fmt.Errorf("exporter is not enabled")
	}

	httpOptions := make([]otlptracehttp.Option, 0)
	var (
		exporterEndpoint  = o.options.ExporterEndpoint
		exporterHeaders   = o.options.ExporterHeaders
		exporterTLSConfig = o.options.ExporterTLSConfig
	)

	httpOptions = append(httpOptions, otlptracehttp.WithInsecure())
	httpOptions = append(httpOptions, otlptracehttp.WithEndpoint(exporterEndpoint))

	if len(exporterHeaders) != 0 {
		httpOptions = append(httpOptions, otlptracehttp.WithHeaders(exporterHeaders))
	}

	if exporterTLSConfig != nil {
		httpOptions = append(httpOptions, otlptracehttp.WithTLSClientConfig(exporterTLSConfig))
	}

	client := otlptracehttp.NewClient(httpOptions...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func (o *openTelemetryTracer) exportEnabled() bool {
	return len(o.options.ExporterEndpoint) != 0
}

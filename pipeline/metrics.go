package pipeline

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/metrics"
	"github.com/kingstonduy/go-core/util"
)

var (
	MetricKeyRequestTotal    = []string{"request", "total"}
	MetricKeyRequestDuration = []string{"request", "duration", "milliseconds"}

	MetricLabelRequestType = "request_type"
	MetricLabelStatusCode  = "status_code"
	MetricLabelErrorCode   = "error_code"
)

type requestMetricsBehavior struct {
	opts RequestMetricBehaviorOptions
}

type RequestMetricBehaviorOptions struct {
	metrics *metrics.Metrics
}

type RequestMetricBehaviorOption func(*RequestMetricBehaviorOptions)

func WithMetrics(metric *metrics.Metrics) RequestMetricBehaviorOption {
	return func(options *RequestMetricBehaviorOptions) {
		options.metrics = metric
	}
}

func NewMetricsBehavior(opts ...RequestMetricBehaviorOption) PipelineBehavior {
	// default options
	options := RequestMetricBehaviorOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return &requestMetricsBehavior{
		opts: options,
	}
}

func (b *requestMetricsBehavior) Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (response interface{}, err error) {
	reqType := util.GetType(request)
	start := time.Now()

	defer func() {
		var requestError *errorx.Error
		if err != nil && !errors.As(err, &requestError) {
			requestError = errorx.Failed(err.Error())
		}

		// Metric request total
		b.metricRequestTotal(reqType, requestError)

		// metric request latency
		b.metricRequestDuration(reqType, time.Since(start))
	}()

	response, err = next(ctx)
	return response, err
}

func (b *requestMetricsBehavior) metricRequestTotal(requestType string, err *errorx.Error) {
	var statusCode int
	var errorCode string

	if err != nil {
		statusCode = err.Status
		errorCode = err.Code
	} else {
		statusCode = errorx.DefaultSuccessStatusCode
		errorCode = errorx.DefaultSuccessResponseCode
	}

	b.getMetrics().IncrCounterWithLabels(
		MetricKeyRequestTotal,
		1,
		[]metrics.Label{
			{
				Name:  MetricLabelRequestType,
				Value: requestType,
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
}

func (b *requestMetricsBehavior) metricRequestDuration(requestType string, duration time.Duration) {
	b.getMetrics().AddSampleWithLabels(
		MetricKeyRequestDuration,
		float32(duration.Milliseconds()),
		[]metrics.Label{
			{
				Name:  MetricLabelRequestType,
				Value: requestType,
			},
		},
	)
}

func (b *requestMetricsBehavior) getMetrics() *metrics.Metrics {
	if b.opts.metrics != nil {
		return b.opts.metrics
	}

	return metrics.Default()
}

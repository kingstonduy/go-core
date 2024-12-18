package pipeline

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/util"
)

// LOGGING
type requestLoggingBehavior struct {
	opts RequestLoggingBehaviorOptions
}

type RequestLoggingBehaviorOptions struct {
	logger logger.Logger
}

type RequestLoggingBehaviorOption func(*RequestLoggingBehaviorOptions)

func WithLogger(log logger.Logger) RequestLoggingBehaviorOption {
	return func(o *RequestLoggingBehaviorOptions) {
		o.logger = log
	}
}

func NewRequestLoggingBehavior(opts ...RequestLoggingBehaviorOption) PipelineBehavior {
	// default options
	options := RequestLoggingBehaviorOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return &requestLoggingBehavior{
		opts: options,
	}
}

func (b *requestLoggingBehavior) Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (response interface{}, err error) {
	start := time.Now()

	defer func() {
		isSuccess := err == nil
		duration := time.Since(start).Milliseconds()

		var (
			requestType     = util.GetType(request)
			requestJson, _  = json.Marshal(request)
			responseJson, _ = json.Marshal(response)
			errJson         = ""
		)
		if err != nil {
			errJson = err.Error()
		}

		b.getLogger().Infof(ctx, "[Request Pipeline] Request Type: %s - Request: %s - Response: %s - Success: %t - Error: %s - Duration: %vms",
			requestType,
			string(requestJson),
			string(responseJson),
			isSuccess,
			errJson,
			duration,
		)
	}()

	response, err = next(ctx)
	return response, err
}

func (b *requestLoggingBehavior) getLogger() logger.Logger {
	if b.opts.logger != nil {
		return b.opts.logger
	}

	return logger.DefaultLogger
}

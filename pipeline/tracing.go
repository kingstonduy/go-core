package pipeline

import (
	"context"
	"fmt"

	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/util"
)

// TRACING
type requestTracingBehavior struct {
	opts TracingBehaviorOptions
}

type TracingBehaviorOptions struct {
	tracer trace.Tracer
}

type TracingBehaviorOption func(*TracingBehaviorOptions)

func WithTracer(tracer trace.Tracer) TracingBehaviorOption {
	return func(options *TracingBehaviorOptions) {
		options.tracer = tracer
	}
}

func NewTracingBehavior(opts ...TracingBehaviorOption) PipelineBehavior {
	// default options
	options := TracingBehaviorOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return &requestTracingBehavior{
		opts: options,
	}
}

func (b *requestTracingBehavior) Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (res interface{}, err error) {
	reqType := util.GetType(request)
	opName := fmt.Sprintf("Request Pipeline - %s", reqType)

	// tracing request
	ctx, finish := b.getTracer().StartTracing(ctx, opName, trace.WithTraceRequest(request))
	defer func() {
		finish(ctx,
			trace.WithTraceResponse(res),
			trace.WithTraceErrorResponse(err),
		)
	}()

	res, err = next(ctx)

	return res, err
}

func (b *requestTracingBehavior) getTracer() trace.Tracer {
	if b.opts.tracer != nil {
		return b.opts.tracer
	}

	return trace.DefaultTracer
}

package pipeline

import (
	"context"

	"github.com/kingstonduy/go-core/errorx"
)

// ERROR HANDLING FOR RECOVERING FROM PANIC
type errorHandlingBehavior struct {
	opts ErrorHandlingOptions
}

type ErrorHandlingOptions struct{}

type ErrorHandlingOption func(*ErrorHandlingOptions)

func NewErrorHandlingBehavior(opts ...ErrorHandlingOption) PipelineBehavior {
	options := ErrorHandlingOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return &errorHandlingBehavior{
		opts: options,
	}
}

func (b *errorHandlingBehavior) Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (res interface{}, err error) {
	// recover from error panic to prevent stop application
	defer func() {
		if r := recover(); r != nil {
			err = errorx.InternalServerError("%v", r)
		}
	}()

	response, err := next(ctx)
	return response, err
}

package pipeline

import (
	"context"
	"reflect"

	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/validation"
)

type validateBehavior struct {
	opts ValidationBehaviorOptions
}

type ValidationBehaviorOptions struct {
	validator validation.Validator
}

type ValidationBehaviorOption func(*ValidationBehaviorOptions)

func WithValidator(validator validation.Validator) ValidationBehaviorOption {
	return func(options *ValidationBehaviorOptions) {
		options.validator = validator
	}
}

func NewValidationBehavior(opts ...ValidationBehaviorOption) PipelineBehavior {
	// Default options
	options := ValidationBehaviorOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return &validateBehavior{
		opts: options,
	}
}

func (v *validateBehavior) Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (interface{}, error) {
	if request == nil || (reflect.ValueOf(request).Kind() == reflect.Ptr && reflect.ValueOf(request).IsNil()) {
		return nil, errorx.BadRequestError("request data is required")
	}

	err := v.getValidator().Validate(request)
	if err != nil {
		return nil, err
	}

	return next(ctx)
}

func (v *validateBehavior) getValidator() validation.Validator {
	if v.opts.validator != nil {
		return v.opts.validator
	}

	return validation.DefaultValidator
}

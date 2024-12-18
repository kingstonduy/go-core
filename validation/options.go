package validation

import "context"

type ValidationOptions struct {
	Context context.Context
}

type ValidationOption func(*ValidationOptions)

func SetOption(k, v interface{}) ValidationOption {
	return func(o *ValidationOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

func NewOptions(opts ...ValidationOption) *ValidationOptions {
	options := ValidationOptions{
		Context: context.Background(),
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &options
}

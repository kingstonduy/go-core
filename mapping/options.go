package mapping

type MapperOptions struct {
	DecodeTimeFormat string
	WeaklyTypedInput bool
}

type MapperOption func(*MapperOptions)

// Option for string to time conversion.
func WithDecodeTimeFormat(fm string) MapperOption {
	return func(c *MapperOptions) {
		c.DecodeTimeFormat = fm
	}
}

// Option for weakly type mapping. Default true
func WithWeaklyTypedInput(enabled bool) MapperOption {
	return func(c *MapperOptions) {
		c.WeaklyTypedInput = enabled
	}
}

func NewMapperOptions(opts ...MapperOption) MapperOptions {
	options := MapperOptions{
		WeaklyTypedInput: true,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

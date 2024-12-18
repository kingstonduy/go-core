package broker

import (
	"context"
)

// SetSubscribeOption returns a function to setup a context with given value
func SetSubscribeOption(k, v interface{}) SubscribeOption {
	return func(o *SubscribeOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetBrokerOption returns a function to setup a context with given value
func SetBrokerOption(k, v interface{}) BrokerOption {
	return func(o *BrokerOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetPublishOption returns a function to setup a context with given value
func SetPublishOption(k, v interface{}) PublishOption {
	return func(o *PublishOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

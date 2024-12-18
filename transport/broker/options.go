package broker

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/google/uuid"
	"github.com/kingstonduy/go-core/logger"
)

type BrokerOption func(*BrokerOptions)

type BrokerOptions struct {
	Context context.Context

	// underlying logger
	Logger logger.Logger

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler

	Addrs []string

	TLSConfig *tls.Config

	Secure bool
}

func NewBrokerOptions(opts ...BrokerOption) BrokerOptions {
	options := BrokerOptions{
		Context: context.Background(),
		Logger:  logger.DefaultLogger,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

func WithBrokerContext(ctx context.Context) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Context = ctx
	}
}

func WithBrokerAddresses(addrs ...string) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Addrs = addrs
	}
}

func WithLogger(log logger.Logger) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Logger = log
	}
}

func WithBrokerErrorHandler(handler Handler) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.ErrorHandler = handler
	}
}

func WithBrokerTLSConfig(t *tls.Config) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.TLSConfig = t
	}
}

func WithBrokerSecure(s bool) BrokerOption {
	return func(opts *BrokerOptions) {
		opts.Secure = s
	}
}

type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Context            context.Context
	Timeout            time.Duration
	ReplyToTopic       string
	ReplyConsumerGroup string
}

func NewPublishOptions(opts ...PublishOption) PublishOptions {
	options := PublishOptions{
		Context: context.Background(),
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

func WithPublishContext(ctx context.Context) PublishOption {
	return func(opts *PublishOptions) {
		opts.Context = ctx
	}
}

func WithPublishTimeout(timeout time.Duration) PublishOption {
	return func(opts *PublishOptions) {
		opts.Timeout = timeout
	}
}

func WithPublishReplyToTopic(replyToTopic string) PublishOption {
	return func(opts *PublishOptions) {
		opts.ReplyToTopic = replyToTopic
	}
}

func WithReplyConsumerGroup(cg string) PublishOption {
	return func(opts *PublishOptions) {
		opts.ReplyConsumerGroup = cg
	}
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Context context.Context

	Group string

	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
}

func WithSubscribeContext(ctx context.Context) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Context = ctx
	}
}

func WithSubscribeGroup(gr string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Group = gr
	}
}

func WithSubscribeAutoAck(autoAck bool) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.AutoAck = autoAck
	}
}

func WithSubscribeQueue(queue string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Queue = queue
	}
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
		Group:   uuid.New().String(),
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

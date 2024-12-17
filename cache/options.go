package cache

import (
	"context"
	"time"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/logger"
	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/logger/logrus"
)

type Options struct {
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context

	// Logger is the be used logger
	Logger logger.Logger

	// Address represents the address or other connection information of the cache service.
	Address string

	Expiration time.Duration
}

// Option manipulates the Options passed.
type Option func(o *Options)

// Expiration sets the duration for items stored in the cache to expire.
func WithDefaultExpiration(d time.Duration) Option {
	return func(o *Options) {
		o.Expiration = d
	}
}

// WithAddress sets the cache service address or connection information.
func WithAddress(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// WithContext sets the cache context, for any extra configuration.
func WithContext(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}

// WithLogger sets underline logger.
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// NewOptions returns a new options struct.
func NewOptions(opts ...Option) Options {
	options := Options{
		Expiration: DefaultExpiration, // no expiration
		Address:    DefaultAddress,
		Logger:     logrus.NewLogrusLogger(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

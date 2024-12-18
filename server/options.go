package server

import (
	"context"

	"github.com/kingstonduy/go-core/logger"
)

type StartHook func(context.Context)

type StopHook func(context.Context)

type Options struct {
	Servers         map[string]Server
	StartHook       StartHook
	StopHook        StopHook
	StopIfOneFailed bool
	Logger          logger.Logger
}

type Option func(*Options)

func WithServer(name string, server Server) Option {
	return func(options *Options) {
		if options.Servers == nil {
			options.Servers = make(map[string]Server)
		}
		options.Servers[name] = server
	}
}

func WithStartHook(start StartHook) Option {
	return func(options *Options) {
		options.StartHook = start
	}
}

func WithStopHook(stop StopHook) Option {
	return func(options *Options) {
		options.StopHook = stop
	}
}

// Default: true
func WithStopIfOneFailed(stop bool) Option {
	return func(options *Options) {
		options.StopIfOneFailed = stop
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(options *Options) {
		options.Logger = logger
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Servers:         make(map[string]Server),
		StopIfOneFailed: true,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

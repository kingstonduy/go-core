package server

import (
	"context"
	"sync"

	"github.com/kingstonduy/go-core/logger"
)

type Server interface {
	// There is some server will block to serve requests
	Start(ctx context.Context) error

	// Stop the server
	Stop(ctx context.Context) error

	// Server is connected or not
	Connected() bool
}

type ServerWrapper struct {
	options Options
}

func NewServerWrapper(opts ...Option) *ServerWrapper {
	options := NewOptions(opts...)
	return &ServerWrapper{
		options: options,
	}
}

func (s *ServerWrapper) Start(ctx context.Context) error {
	if start := s.options.StartHook; start != nil {
		start(ctx)
	}

	wg := &sync.WaitGroup{}

	servers := s.options.Servers
	wg.Add(len(servers))
	for _, server := range servers {
		go func(srv Server) {
			wg.Done() // not defer because srv.Start may be blocking
			err := srv.Start(ctx)
			if err != nil && s.options.StopIfOneFailed {
				s.log(ctx, logger.ErrorLevel, "Stop all server due to one failure: %v", err)
				if err := s.Stop(ctx); err != nil {
					s.log(ctx, logger.ErrorLevel, "Failed to stop server: %v", err)
				}

			}
		}(server)
	}

	wg.Wait()

	return nil
}

func (s *ServerWrapper) Stop(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	servers := s.options.Servers

	wg.Add(len(servers))
	for _, server := range servers {
		go func(srv Server) {
			defer wg.Done()
			if srv.Connected() {
				if err := srv.Stop(ctx); err != nil {
					s.log(ctx, logger.ErrorLevel, "Failed to stop server: %v", err)
				}
			}
		}(server)
	}

	wg.Wait()
	if stop := s.options.StopHook; stop != nil {
		stop(ctx)
	}

	return nil
}

func (s *ServerWrapper) log(ctx context.Context, level logger.Level, message string, args ...interface{}) {
	var log logger.Logger
	if !s.logEnabled() {
		log = logger.DefaultLogger
	}
	log.Logf(ctx, level, message, args...)
}

func (s *ServerWrapper) logEnabled() bool {
	return s.options.Logger != nil
}

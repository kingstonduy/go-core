package database

import (
	"context"
	"time"

	"github.com/kingstonduy/go-core/logger"
)

var (
	DefaultDbLogHooksOptions = dbLogHooksOptions{}
)

type queryBeginTimeKey struct{}

type dbLogHooksOptions struct {
	logger          logger.Logger
	warningDuration time.Duration
}

type DbLogHooksOption func(*dbLogHooksOptions)

func WithLogHooksLogger(logger logger.Logger) DbLogHooksOption {
	return func(options *dbLogHooksOptions) {
		options.logger = logger
	}
}

// only log if query execution duration greater than or equal to the given duration. If 0, log in all cases
func WithLogHooksWarningThreshold(duration time.Duration) DbLogHooksOption {
	return func(options *dbLogHooksOptions) {
		options.warningDuration = duration
	}
}

func newDBLogHooksOptions(opts ...DbLogHooksOption) dbLogHooksOptions {
	options := DefaultDbLogHooksOptions

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

type DbLogHooks struct {
	Logger           logger.Logger
	warningThreshold time.Duration
}

func NewDBLogHooks(opts ...DbLogHooksOption) Hooks {
	options := newDBLogHooksOptions(opts...)

	return &DbLogHooks{
		Logger:           options.logger,
		warningThreshold: options.warningDuration,
	}
}

// After implements Hooks.
func (d *DbLogHooks) After(ctx context.Context, err error, query string, args ...interface{}) (context.Context, error) {
	if begin, ok := ctx.Value(queryBeginTimeKey{}).(time.Time); ok {
		duration := time.Since(begin)
		if duration >= d.warningThreshold {
			d.log(ctx, logger.TraceLevel, "SQL EXECUTION - Query: %s - Durations: %s. Error: %v", query, duration, err)
		}
	}
	return ctx, nil
}

// Before implements Hooks.
func (d *DbLogHooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, queryBeginTimeKey{}, time.Now()), nil
}

func (d *DbLogHooks) log(ctx context.Context, level logger.Level, message string, args ...interface{}) {
	logger := d.getLogger()
	logger.Logf(ctx, level, message, args...)
}

func (d *DbLogHooks) getLogger() logger.Logger {
	if d.logEnabled() {
		return d.Logger
	}
	return logger.DefaultLogger
}

func (d *DbLogHooks) logEnabled() bool {
	return d.Logger != nil
}

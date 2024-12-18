package database

import (
	"database/sql"
	"time"

	"github.com/kingstonduy/go-core/logger"
)

// / Transaction options /////////////////////////////
type TransactionOptions struct {
	// TxOptions holds the transaction options to be used in DB.BeginTx.
	// Isolation is the transaction isolation level.
	// If zero, the driver or database's default level is used.
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

type TransactionOption func(*TransactionOptions)

func WithIsolationLevelOptions(lv sql.IsolationLevel) TransactionOption {
	return func(options *TransactionOptions) {
		options.Isolation = lv
	}
}

func WithReadOnly(readOnly bool) TransactionOption {
	return func(options *TransactionOptions) {
		options.ReadOnly = readOnly
	}
}

func NewTransactionOptions(opts ...TransactionOption) TransactionOptions {
	// Default defaultOptions
	defaultOptions := TransactionOptions{}

	for _, opt := range opts {
		opt(&defaultOptions)
	}

	return defaultOptions
}

// / Pool options /////////////////////////////
type DatabaseOptions struct {
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed
	Logger       logger.Logger
	Hooks        Hooks
}

type DatabaseOption func(*DatabaseOptions)

func WithMaxIdleCount(value int) DatabaseOption {
	return func(options *DatabaseOptions) {
		options.MaxIdleCount = value
	}
}

func WithMaxOpen(value int) DatabaseOption {
	return func(options *DatabaseOptions) {
		options.MaxOpen = value
	}
}

func WithMaxLifetime(time time.Duration) DatabaseOption {
	return func(options *DatabaseOptions) {
		options.MaxLifetime = time
	}
}

func WithMaxIdleTime(time time.Duration) DatabaseOption {
	return func(options *DatabaseOptions) {
		options.MaxIdleTime = time
	}
}

func WithLogger(logger logger.Logger) DatabaseOption {
	return func(options *DatabaseOptions) {
		options.Logger = logger
	}
}

// override before defined hooks
func WithHooks(hooks ...Hooks) DatabaseOption {
	return func(options *DatabaseOptions) {
		options.Hooks = Compose(hooks...)
	}
}

func NewDatabaseOptions(opts ...DatabaseOption) DatabaseOptions {
	defaultOptions := DatabaseOptions{}

	for _, opt := range opts {
		opt(&defaultOptions)
	}

	return defaultOptions
}

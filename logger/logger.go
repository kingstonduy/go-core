package logger

import "context"

type Logger interface {
	// Init initializes options
	Init(options ...Option) error
	// The Logger options
	Options() Options
	// Fields set fields to always be logged
	Fields(fields map[string]interface{}) Logger

	// Log writes a log entry
	Log(ctx context.Context, level Level, args ...interface{})
	// Logf writes a formatted log entry
	Logf(ctx context.Context, level Level, format string, args ...interface{})
	// Log Trace
	Trace(ctx context.Context, args ...interface{})
	// Logf Trace
	Tracef(ctx context.Context, format string, args ...interface{})
	// Log Debug
	Debug(ctx context.Context, args ...interface{})
	// Logf Debug
	Debugf(ctx context.Context, format string, args ...interface{})
	// Log Info
	Info(ctx context.Context, args ...interface{})
	// Logf Info
	Infof(ctx context.Context, format string, args ...interface{})
	// Log Warn
	Warn(ctx context.Context, args ...interface{})
	// Logf Warn
	Warnf(ctx context.Context, format string, args ...interface{})
	// Log Error
	Error(ctx context.Context, args ...interface{})
	// Logf Error
	Errorf(ctx context.Context, format string, args ...interface{})
	// Log Fatal
	Fatal(ctx context.Context, args ...interface{})
	// Logf Fatal
	Fatalf(ctx context.Context, format string, args ...interface{})

	// String returns the name of logger
	String() string
}

func SetDefaultLogger(log Logger) {
	DefaultLogger = log
}

// Default
var (
	DefaultLogger               = newNoopsLogger()
	DefaultEnvLogLevel          = "LOG_LEVEL"
	DefaultLogLevel             = InfoLevel
	DefaultLogLevelString       = "info"
	DefaultMaskSensitiveLogData = true
	DefaultCallerSkipCount      = 7
	DefaultTimestampFormat      = "2006-01-02 15:04:05.000"
)

// Return DefaultLogger with given fields
func Fields(fields map[string]interface{}) Logger {
	return DefaultLogger.Fields(fields)
}

// Log writes a log entry
func Log(ctx context.Context, level Level, args ...interface{}) {
	DefaultLogger.Log(ctx, level, args...)
}

// Logf writes a formatted log entry
func Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	DefaultLogger.Logf(ctx, level, format, args...)
}

// Log Trace
func Trace(ctx context.Context, args ...interface{}) {
	DefaultLogger.Trace(ctx, args...)
}

// Logf Trace
func Tracef(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Tracef(ctx, format, args...)
}

// Log Debug
func Debug(ctx context.Context, args ...interface{}) {
	DefaultLogger.Debug(ctx, args...)
}

// Logf Debug
func Debugf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Debugf(ctx, format, args...)
}

// Log Info
func Info(ctx context.Context, args ...interface{}) {
	DefaultLogger.Info(ctx, args...)
}

// Logf Info
func Infof(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Infof(ctx, format, args...)
}

// Log Warn
func Warn(ctx context.Context, args ...interface{}) {
	DefaultLogger.Warn(ctx, args...)
}

// Logf Warn
func Warnf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Warnf(ctx, format, args...)
}

// Log Error
func Error(ctx context.Context, args ...interface{}) {
	DefaultLogger.Error(ctx, args...)
}

// Logf Error
func Errorf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Errorf(ctx, format, args...)
}

// Log Fatal
func Fatal(ctx context.Context, args ...interface{}) {
	DefaultLogger.Fatal(ctx, args...)
}

// Logf Fatal
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	DefaultLogger.Fatalf(ctx, format, args...)
}

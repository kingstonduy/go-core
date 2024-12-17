package logger

import (
	"context"
	"strings"
)

// Logger interface definition
type ILogger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Fields(fields map[string]interface{}) ILogger
	// Log writes a log entry
	Log(ctx context.Context, level Level, args ...interface{})
	// Logf writes a formatted log entry
	Logf(ctx context.Context, level Level, format string, args ...interface{})
}

var defaultLogger ILogger

func GetDefaulLogger() ILogger {
	if defaultLogger == nil {
		defaultLogger = newDefaultLogger()
	}
	return defaultLogger
}
func InitLogger(logger ILogger) {
	defaultLogger = logger
}

func Fields(fields map[string]interface{}) ILogger {
	return defaultLogger.Fields(fields)
}

// Info logs an informational message
func Info(ctx context.Context, args ...interface{}) {
	defaultLogger.Info(ctx, args...)
}

// Infof logs a formatted informational message
func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Infof(ctx, format, args...)
}

// Debug logs a debug message
func Debug(ctx context.Context, args ...interface{}) {
	defaultLogger.Debug(ctx, args...)
}

// Debugf logs a formatted debug message
func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Debugf(ctx, format, args...)
}

// Warn logs a warning message
func Warn(ctx context.Context, args ...interface{}) {
	defaultLogger.Warn(ctx, args...)
}

// Warnf logs a formatted warning message
func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Warnf(ctx, format, args...)
}

// Error logs an error message
func Error(ctx context.Context, args ...interface{}) {
	defaultLogger.Error(ctx, args...)
}

// Errorf logs a formatted error message
func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Errorf(ctx, format, args...)
}

// Fatal logs a fatal message and exits the program
func Fatal(ctx context.Context, args ...interface{}) {
	defaultLogger.Fatal(ctx, args...)
}

// Fatalf logs a formatted fatal message and exits the program
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Fatalf(ctx, format, args...)
}

func GetColor(level string) string {
	var color string
	switch strings.ToUpper(level) {
	case "DEBUG":
		color = "\033[37m" // Gray
	case "INFO":
		color = "\033[32m" // Green
	case "WARN":
		color = "\033[33m" // Yellow
	case "ERROR":
		color = "\033[31m" // Red
	case "FATAL":
		color = "\033[35m" // Magenta
	default:
		color = "\033[0m" // Default color (reset)
	}
	return color
}

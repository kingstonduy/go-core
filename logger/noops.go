package logger

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type noopsLogger struct{}

// Debug implements Logger.
func (n *noopsLogger) Debug(ctx context.Context, args ...interface{}) {
	n.Log(ctx, DebugLevel, args...)
}

// Debugf implements Logger.
func (n *noopsLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	n.Logf(ctx, DebugLevel, format, args...)
}

// Error implements Logger.
func (n *noopsLogger) Error(ctx context.Context, args ...interface{}) {
	n.Log(ctx, ErrorLevel, args...)
}

// Errorf implements Logger.
func (n *noopsLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	n.Logf(ctx, ErrorLevel, format, args...)
}

// Fatal implements Logger.
func (n *noopsLogger) Fatal(ctx context.Context, args ...interface{}) {
	n.Log(ctx, FatalLevel, args...)
}

// Fatalf implements Logger.
func (n *noopsLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	n.Logf(ctx, FatalLevel, format, args...)
}

// Fields implements Logger.
func (n *noopsLogger) Fields(fields map[string]interface{}) Logger {
	return &noopsLogger{}
}

// Info implements Logger.
func (n *noopsLogger) Info(ctx context.Context, args ...interface{}) {
	n.Log(ctx, InfoLevel, args...)
}

// Infof implements Logger.
func (n *noopsLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	n.Logf(ctx, InfoLevel, format, args...)
}

// Init implements Logger.
func (n *noopsLogger) Init(options ...Option) error {
	return nil
}

// Log implements Logger.
func (n *noopsLogger) Log(ctx context.Context, level Level, args ...interface{}) {
	n.noopsWarning()
	log.Printf("[%s] %v\n", strings.ToUpper(level.String()), fmt.Sprint(args...))
}

// Logf implements Logger.
func (n *noopsLogger) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	n.noopsWarning()
	log.Printf("[%s] %v\n", strings.ToUpper(level.String()), fmt.Sprintf(format, args...))
}

// Options implements Logger.
func (n *noopsLogger) Options() Options {
	return Options{}
}

// String implements Logger.
func (n *noopsLogger) String() string {
	return "noops"
}

// Trace implements Logger.
func (n *noopsLogger) Trace(ctx context.Context, args ...interface{}) {
	n.Log(ctx, TraceLevel, args...)
}

// Tracef implements Logger.
func (n *noopsLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	n.Logf(ctx, TraceLevel, format, args...)
}

// Warn implements Logger.
func (n *noopsLogger) Warn(ctx context.Context, args ...interface{}) {
	n.Log(ctx, WarnLevel, args...)
}

// Warnf implements Logger.
func (n *noopsLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	n.Logf(ctx, WarnLevel, format, args...)
}

func (n *noopsLogger) noopsWarning() {
	log.Print("[WARN] No default logger was set. Using noops logger as default. Set the default logger to do all functions\n")
}

func newNoopsLogger() Logger {
	return &noopsLogger{}
}

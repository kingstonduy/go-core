package logger

import (
	"context"
	"fmt"

	"runtime"
	"strings"

	"github.com/gammazero/workerpool"
	trace "github.com/kingstonduy/go-core/tracer"
	"github.com/sirupsen/logrus"
)

// logrusLogger struct
type defLogger struct {
	entry *logrus.Entry
	wp    *workerpool.WorkerPool
}

// NewLogrusLogger initializes a logrusLogger with a custom formatter
func newDefaultLogger() ILogger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&CustomFormatter{})
	return &defLogger{
		entry: logrus.NewEntry(l),
		wp:    workerpool.New(100),
	}
}

// Log implements I
func (l *defLogger) Log(ctx context.Context, level Level, args ...interface{}) {
	l.wp.Submit(func() {
		l.entry.WithContext(ctx).Log(LoggerToLogrusLevel(level), args...)
	})
}

// Logf implements I
func (l *defLogger) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	l.wp.Submit(func() {
		l.entry.WithContext(ctx).Log(LoggerToLogrusLevel(level), args...)
	})
}

// Fields implements
func (l *defLogger) Fields(fields map[string]interface{}) ILogger {
	return &defLogger{
		entry: l.entry.WithFields(fields),
		wp:    l.wp,
	}
}

func (l *defLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, InfoLevel, args...)
}

func (l *defLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, DebugLevel, args...)
}

func (l *defLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, WarnLevel, args...)
}

func (l *defLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, ErrorLevel, args...)
}

func (l *defLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, FatalLevel, args...)
}

func (l *defLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.InfoLevel, format, args...)
}

func (l *defLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.DebugLevel, format, args...)
}

func (l *defLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.WarnLevel, format, args...)
}

func (l *defLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.ErrorLevel, format, args...)
}

func (l *defLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.FatalLevel, format, args...)
}

// CustomFormatter struct to format logs with color
type CustomFormatter struct {
	logrus.TextFormatter
}

// Format formats the log entry with colors for different log levels
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	ctx := entry.Context
	spanInfo := trace.GetSpanInfo(ctx)
	var (
		time             = entry.Time.Format("2006-01-02 15:04:05.000")
		level            = getLevel(entry)
		numGoroutines    = runtime.NumGoroutine()
		_, file, line, _ = runtime.Caller(7) // Caller depth adjusted for correct file and line number
		fileInfo         = fmt.Sprintf("%s:%d", file, line)
		traceID          = spanInfo.TraceID
		spanID           = spanInfo.SpanID
		service          = spanInfo.ServiceDomain
		operator         = getField(entry, FIELD_OPERATOR_NAME)
		stepName         = getField(entry, FIELD_STEP_NAME)
		clientID         = spanInfo.ClientID
		systemID         = spanInfo.SystemID
		from             = spanInfo.From
		to               = spanInfo.To
		duration         = getField(entry, FIELD_DURATION)
		message          = entry.Message
	)
	// time level [numGoroutines] [file:line] [traceID] [spanID] [service] [operator] [stepName] [clientID] [systemID] [from] [to] [duration] - message
	res := fmt.Sprintf("%s %s [%d] [%s] [%s,%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] - %s\n",
		time,
		level,
		numGoroutines,
		fileInfo,
		traceID,
		spanID,
		service,
		operator,
		stepName,
		clientID,
		systemID,
		from,
		to,
		duration,
		message,
	)

	return []byte(res), nil
}

func getLevel(entry *logrus.Entry) (level string) {
	var color string
	levelName := map[logrus.Level]string{
		logrus.DebugLevel: "DEBUG",
		logrus.InfoLevel:  "INFO",
		logrus.WarnLevel:  "WARN",
		logrus.ErrorLevel: "ERROR",
		logrus.FatalLevel: "FATAL",
		logrus.PanicLevel: "PANIC",
	}[entry.Level]

	if levelName == "" {
		levelName = entry.Level.String() // Fallback to default
	}

	switch entry.Level {
	case logrus.DebugLevel:
		color = GetColor("DEBUG") // Gray
	case logrus.InfoLevel:
		color = GetColor("INFO") // Green
	case logrus.WarnLevel:
		color = GetColor("WARN") // Yellow
	case logrus.ErrorLevel:
		color = GetColor("ERROR") // Red
	case logrus.FatalLevel, logrus.PanicLevel:
		color = GetColor("FATAL") // Magenta
	default:
		color = GetColor("DEFAULT") // Default color (reset)
	}
	level = fmt.Sprintf("%s[%s]\033[0m", color, strings.ToUpper(levelName))

	return level
}

func getField(entry *logrus.Entry, field string) string {
	if value, ok := entry.Data[field]; ok {
		return fmt.Sprintf("%v", value)
	}
	return ""
}

func LoggerToLogrusLevel(level Level) logrus.Level {
	switch level {
	case TraceLevel:
		return logrus.TraceLevel
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

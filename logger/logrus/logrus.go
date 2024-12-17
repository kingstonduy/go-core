package logrus

import (
	"context"
	"fmt"

	"runtime"
	"strings"

	"github.com/gammazero/workerpool"
	"github.com/kingstonduy/go-core/logger"
	trace "github.com/kingstonduy/go-core/tracer"
	"github.com/sirupsen/logrus"
)

// logrusLogger struct
type logrusLogger struct {
	entry *logrus.Entry
	wp    *workerpool.WorkerPool
}

// NewLogrusLogger initializes a logrusLogger with a custom formatter
func NewLogrusLogger() logger.ILogger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&CustomFormatter{})
	return &logrusLogger{
		entry: logrus.NewEntry(l),
		wp:    workerpool.New(100),
	}
}

// Log implements logger.ILogger.
func (l *logrusLogger) Log(ctx context.Context, level logger.Level, args ...interface{}) {
	l.wp.Submit(func() {
		l.entry.WithContext(ctx).Log(LoggerToLogrusLevel(level), args...)
	})
}

// Logf implements logger.ILogger.
func (l *logrusLogger) Logf(ctx context.Context, level logger.Level, format string, args ...interface{}) {
	l.wp.Submit(func() {
		l.entry.WithContext(ctx).Log(LoggerToLogrusLevel(level), args...)
	})
}

// Fields implements logger.Logger.
func (l *logrusLogger) Fields(fields map[string]interface{}) logger.ILogger {
	return &logrusLogger{
		entry: l.entry.WithFields(fields),
		wp:    l.wp,
	}
}

func (l *logrusLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.InfoLevel, args...)
}

func (l *logrusLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.DebugLevel, args...)
}

func (l *logrusLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.WarnLevel, args...)
}

func (l *logrusLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.ErrorLevel, args...)
}

func (l *logrusLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.FatalLevel, args...)
}

func (l *logrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.InfoLevel, format, args...)
}

func (l *logrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.DebugLevel, format, args...)
}

func (l *logrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.WarnLevel, format, args...)
}

func (l *logrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.ErrorLevel, format, args...)
}

func (l *logrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
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
		operator         = getField(entry, logger.FIELD_OPERATOR_NAME)
		stepName         = getField(entry, logger.FIELD_STEP_NAME)
		clientID         = spanInfo.ClientID
		systemID         = spanInfo.SystemID
		from             = spanInfo.From
		to               = spanInfo.To
		duration         = getField(entry, logger.FIELD_DURATION)
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
		color = logger.GetColor("DEBUG") // Gray
	case logrus.InfoLevel:
		color = logger.GetColor("INFO") // Green
	case logrus.WarnLevel:
		color = logger.GetColor("WARN") // Yellow
	case logrus.ErrorLevel:
		color = logger.GetColor("ERROR") // Red
	case logrus.FatalLevel, logrus.PanicLevel:
		color = logger.GetColor("FATAL") // Magenta
	default:
		color = logger.GetColor("DEFAULT") // Default color (reset)
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

func LoggerToLogrusLevel(level logger.Level) logrus.Level {
	switch level {
	case logger.TraceLevel:
		return logrus.TraceLevel
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	case logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

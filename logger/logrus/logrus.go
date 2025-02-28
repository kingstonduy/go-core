package logrus

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/trace"

	"github.com/gammazero/workerpool"
	"github.com/sirupsen/logrus"
)

const (
	LoggerWorkerCount = 100
)

type entryLogger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry

	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
}

type logrusLogger struct {
	Logger entryLogger
	opts   LogrusOptions
	wp     *workerpool.WorkerPool
}

func (l *logrusLogger) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if formatter, ok := l.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		l.opts.Formatter = formatter
	}

	if hs, ok := l.opts.Context.Value(hooksKey{}).(logrus.LevelHooks); ok {
		l.opts.Hooks = hs
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}

	switch ll := l.opts.Context.Value(logrusLoggerKey{}).(type) {
	case *logrus.Logger:
		// overwrite default options
		l.opts.Level = LogrusToLoggerLevel(ll.GetLevel())
		l.opts.Out = ll.Out
		l.opts.Formatter = ll.Formatter
		l.opts.Hooks = ll.Hooks
		l.opts.ReportCaller = ll.ReportCaller
		l.opts.ExitFunc = ll.ExitFunc
		l.Logger = ll
	case *logrus.Entry:
		// overwrite default options
		el := ll.Logger
		l.opts.Level = LogrusToLoggerLevel(el.GetLevel())
		l.opts.Out = el.Out
		l.opts.Formatter = el.Formatter
		l.opts.Hooks = el.Hooks
		l.opts.ReportCaller = el.ReportCaller
		l.opts.ExitFunc = el.ExitFunc
		l.Logger = ll
	case nil:
		log := logrus.New() // defaults
		log.SetLevel(LoggerToLogrusLevel(l.opts.Level))
		log.SetOutput(l.opts.Out)
		log.SetFormatter(l.opts.Formatter)
		log.ReplaceHooks(l.opts.Hooks)
		log.SetReportCaller(l.opts.ReportCaller)
		log.ExitFunc = l.opts.ExitFunc
		l.Logger = log
	default:
		return fmt.Errorf("invalid logrus type: %T", ll)
	}

	return nil
}

func (l *logrusLogger) String() string {
	return "logrus"
}

func (l *logrusLogger) Fields(fields map[string]interface{}) logger.Logger {
	return &logrusLogger{
		Logger: l.Logger.WithFields(fields),
		opts:   l.opts,
		wp:     l.wp,
	}
}

func (l *logrusLogger) Log(ctx context.Context, level logger.Level, args ...interface{}) {
	l.wp.Submit(func() {
		var entry = l.getLogEntry(ctx)
		entry.Log(LoggerToLogrusLevel(level), args...)
	})

}

func (l *logrusLogger) Logf(ctx context.Context, level logger.Level, format string, args ...interface{}) {
	l.wp.Submit(func() {
		var entry = l.getLogEntry(ctx)
		entry.Logf(LoggerToLogrusLevel(level), format, args...)
	})
}

func (l *logrusLogger) getLogEntry(ctx context.Context) entryLogger {
	var entry = l.Logger
	entry = entry.WithContext(ctx)
	return entry
}

func (l *logrusLogger) Options() logger.Options {
	return l.opts.Options
}

// Debug implements logger.Logger.
func (l *logrusLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.DebugLevel, args...)
}

// Debugf implements logger.Logger.
func (l *logrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.DebugLevel, format, args...)
}

// Error implements logger.Logger.
func (l *logrusLogger) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.ErrorLevel, args...)
}

// Errorf implements logger.Logger.
func (l *logrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.ErrorLevel, format, args...)
}

// Fatal implements logger.Logger.
func (l *logrusLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.FatalLevel, args...)
}

// Fatalf implements logger.Logger.
func (l *logrusLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.FatalLevel, format, args...)
}

// Info implements logger.Logger.
func (l *logrusLogger) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.InfoLevel, args...)
}

// Infof implements logger.Logger.
func (l *logrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.InfoLevel, format, args...)
}

// Trace implements logger.Logger.
func (l *logrusLogger) Trace(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.TraceLevel, args...)
}

// Tracef implements logger.Logger.
func (l *logrusLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.TraceLevel, format, args...)
}

// Warn implements logger.Logger.
func (l *logrusLogger) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.WarnLevel, args...)
}

// Warnf implements logger.Logger.
func (l *logrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.Logf(ctx, logger.WarnLevel, format, args...)
}

func (l *logrusLogger) getTracer() trace.Tracer {
	if l.opts.Tracer != nil {
		return l.opts.Tracer
	}
	return trace.DefaultTracer
}

// New builds a new logger based on options.
func NewLogrusLogger(opts ...logger.Option) logger.Logger {
	// Default options
	l := logrusLogger{
		wp: workerpool.New(LoggerWorkerCount),
	}

	loggerOpts := logger.Options{
		MaskSensitiveData: logger.DefaultMaskSensitiveLogData,
		Level:             logger.DefaultLogLevel,
		Fields:            make(map[string]interface{}),
		Out:               os.Stderr,
		Context:           context.Background(),
		CallerSkipCount:   logger.DefaultCallerSkipCount,
	}
	l.opts = LogrusOptions{
		Options: loggerOpts,
		Formatter: &LoggingFormatter{
			logger: &l,
		},
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}
	_ = l.Init(opts...)
	return &l
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

func LogrusToLoggerLevel(level logrus.Level) logger.Level {
	switch level {
	case logrus.TraceLevel:
		return logger.TraceLevel
	case logrus.DebugLevel:
		return logger.DebugLevel
	case logrus.InfoLevel:
		return logger.InfoLevel
	case logrus.WarnLevel:
		return logger.WarnLevel
	case logrus.ErrorLevel:
		return logger.ErrorLevel
	case logrus.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}

type LoggingFormatter struct {
	logger *logrusLogger
}

func (l LoggingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var (
		opts = l.logger.opts
	)

	spanInfo := l.logger.getTracer().ExtractSpanInfo(entry.Context)

	var (
		// these are information which go alongside with each request
		time          = entry.Time.Format(logger.DefaultTimestampFormat)
		level         = strings.ToUpper(entry.Level.String())
		traceID       = spanInfo.TraceID
		spanID        = spanInfo.SpanID
		systemID      = spanInfo.SystemID
		clientIP      = spanInfo.ClientIP
		method        = spanInfo.Method
		serviceDomain = spanInfo.ServiceDomain
		userAgent     = spanInfo.UserAgent
		userName      = spanInfo.Username
		remoteHost    = spanInfo.RemoteHost
		xForwardedFor = spanInfo.XForwardedFor
		contentLength = spanInfo.ContentLength

		// these are information which is specific for each log entry
		operatorName   = l.extractField(entry, logger.FIELD_OPERATOR_NAME)
		stepName       = l.extractField(entry, logger.FIELD_STEP_NAME)
		duration       = l.extractField(entry, logger.FIELD_DURATION)
		statusResponse = l.extractField(entry, logger.FIELD_STATUS_RESPONSE)
	)

	// "DateTime Level [Thread] [X-B3-TraceId,X-B3-SpanId] [systemTraceId] [clientIP] [httpMethod] [serviceDomain] [operatorName] [stepName] [req.userAgent] [user] [processTime] [req.remoteHost] [req.xForwardedFor] [contentLength] [statusResponse] - Message"
	message := fmt.Sprintf("%v %v [] [%v,%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] [%v] - %v\n",
		time,
		level,
		traceID,
		spanID,
		systemID,
		clientIP,
		method,
		serviceDomain,
		operatorName,
		stepName,
		userAgent,
		userName,
		duration,
		remoteHost,
		xForwardedFor,
		contentLength,
		statusResponse,
		entry.Message,
	)

	if opts.MaskSensitiveData {
		message = logger.MaskSensitiveData(message, opts.MaskedPatterns...)
	}

	return []byte(message), nil
}

func (logFormatter LoggingFormatter) extractField(entry *logrus.Entry, field string) string {
	if data, ok := entry.Data[field]; ok {
		return fmt.Sprintf("%v", data)
	}
	return ""
}

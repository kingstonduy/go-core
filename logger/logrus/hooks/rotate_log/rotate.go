package rotateLog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kingstonduy/go-core/metadata"
	strftime "github.com/lestrrat-go/strftime"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_APPLICATION_NAME = "undefined"
	DEFAULT_POD_NAME         = "undefined"
)

// SyslogHook to send logs via syslog.
type RotatelogHook struct {
	rotateLog *RotateLog
	formatter logrus.Formatter
}

// Creates a hook to be added to an instance of logger. This is called with
//
// hook, err := NewRotateLogHook(
//
//	 	WithRotateLogFilePattern("/data/ocb/log/serviceName/podName/podName.%Y%m%d.log"),
//		WithLinkName("./access_log"),
//		WithMaxAge(24 * time.Hour),
//		WithRotationTime(time.Hour),
//		WithClock(UTC))
//
//		if err == nil {
//				log.Hooks.Add(hook.SetFormatter(&logrus.JSONFormatter{}))
//		}`
func NewRotateLogHook(options ...IOption) (*RotatelogHook, error) {
	w, err := newRotateLog(options...)
	if err != nil {
		return nil, err
	}

	return &RotatelogHook{rotateLog: w}, nil
}

func (hook *RotatelogHook) Fire(entry *logrus.Entry) error {
	var line []byte

	if hook.formatter != nil {
		msg, err := hook.formatter.Format(entry)
		if err != nil {
			return err
		}

		line = msg
	} else {
		msg, err := entry.String()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
			return err
		}

		line = []byte(msg)
	}

	_, err := hook.rotateLog.Write(line)

	return err
}

func (hook *RotatelogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *RotatelogHook) SetFormatter(formatter logrus.Formatter) *RotatelogHook {
	hook.formatter = formatter

	return hook
}

func (hook *RotatelogHook) Close() error {
	return hook.rotateLog.Close()
}

// ============= ROTATE LOG==========================
// RotateLog represents a log file that gets
// automatically rotated as you write to it.
type RotateLog struct {
	clock         Clock
	curFn         string
	globPattern   string
	generation    int
	linkName      string
	maxAge        time.Duration
	mutex         sync.RWMutex
	outFh         *os.File
	pattern       *strftime.Strftime
	rotationTime  time.Duration
	rotationCount uint
}

// Clock is the interface used by the RotateLog
// object to determine the current time
type Clock interface {
	Now() time.Time
}
type clockFn func() time.Time

// UTC is an object satisfying the Clock interface, which
// returns the current time in UTC
var UTC = clockFn(func() time.Time { return time.Now().UTC() })

// Local is an object satisfying the Clock interface, which
// returns the current time in the local timezone
var Local = clockFn(time.Now)

func (c clockFn) Now() time.Time {
	return c()
}

// New creates a newRotateLog RotateLog object. A log filename pattern
// must be passed. Optional `Option` parameters may be passed
func newRotateLog(options ...IOption) (*RotateLog, error) {
	var clock Clock = Local
	var rotationCount uint
	var maxAge time.Duration
	var linkName string

	logFilePattern := getLogFilePattern()

	rotationTime := 24 * time.Hour

	for _, o := range options {
		switch o.Name() {
		case OptKeyLogFilePattern:
			logFilePattern = o.Value().(string)
		case OptKeyClock:
			clock = o.Value().(Clock)
		case OptKeyLinkName:
			linkName = o.Value().(string)
		case OptKeyMaxAge:
			maxAge = o.Value().(time.Duration)
			if maxAge < 0 {
				maxAge = 0
			}
		case OptKeyRotationTime:
			rotationTime = o.Value().(time.Duration)
			if rotationTime < 0 {
				rotationTime = 0
			}
		case OptKeyRotationCount:
			rotationCount = o.Value().(uint)
		}
	}

	if maxAge > 0 && rotationCount > 0 {
		return nil, errors.New("options MaxAge and RotationCount cannot be both set")
	}

	if maxAge == 0 && rotationCount == 0 {
		// if both are 0, give maxAge a sane default
		maxAge = 7 * 24 * time.Hour
	}

	globPattern := logFilePattern
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	pattern, err := strftime.New(logFilePattern)
	if err != nil {
		return nil, fmt.Errorf(`invalid strftime pattern. %w`, err)
	}

	return &RotateLog{
		clock:         clock,
		globPattern:   globPattern,
		linkName:      linkName,
		maxAge:        maxAge,
		pattern:       pattern,
		rotationTime:  rotationTime,
		rotationCount: rotationCount,
	}, nil
}

func (rotateLog *RotateLog) genFilename() string {
	now := rotateLog.clock.Now()

	// XXX HACK: Truncate only happens in UTC semantics, apparently.
	// observed values for truncating given time with 86400 secs:
	//
	// before truncation: 2018/06/01 03:54:54 2018-06-01T03:18:00+09:00
	// after  truncation: 2018/06/01 03:54:54 2018-05-31T09:00:00+09:00
	//
	// This is really annoying when we want to truncate in local time
	// so we hack: we take the apparent local time in the local zone,
	// and pretend that it's in UTC. do our math, and put it back to
	// the local zone
	var base time.Time
	if now.Location() != time.UTC {
		base = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
		base = base.Truncate(time.Duration(rotateLog.rotationTime))
		base = time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
	} else {
		base = now.Truncate(time.Duration(rotateLog.rotationTime))
	}
	return rotateLog.pattern.FormatString(base)
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (rotateLog *RotateLog) Write(p []byte) (n int, err error) {
	// Guard against concurrent writes
	rotateLog.mutex.Lock()
	defer rotateLog.mutex.Unlock()

	out, err := rotateLog.getWriter_nolock(false, false)
	if err != nil {
		return 0, fmt.Errorf(`failed to acquite target io.Writer. %w`, err)
	}

	return out.Write(p)
}

// must be locked during this operation
func (rotateLog *RotateLog) getWriter_nolock(bailOnRotateFail, useGenerationalNames bool) (io.Writer, error) {
	generation := rotateLog.generation

	// This filename contains the name of the "NEW" filename
	// to log to, which may be newer than rl.currentFilename
	filename := rotateLog.genFilename()
	if rotateLog.curFn != filename {
		generation = 0
	} else {
		if !useGenerationalNames {
			// nothing to do
			return rotateLog.outFh, nil
		}
		// This is used when we *REALLY* want to rotate a log.
		// instead of just using the regular strftime pattern, we
		// create a new file name using generational names such as
		// "foo.1", "foo.2", "foo.3", etc
		for {
			generation++
			name := fmt.Sprintf("%s.%d", filename, generation)
			if _, err := os.Stat(name); err != nil {
				filename = name
				break
			}
		}
	}

	// if we got here, then we need to create a file
	// create the log directory if not exists
	logPath := filepath.Dir(filename)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %v", logPath, err)
		}
	}

	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", rotateLog.pattern, err)
	}

	if err := rotateLog.rotate_nolock(filename); err != nil {
		err = fmt.Errorf("failed to rotate. %w", err)
		if bailOnRotateFail {
			// Failure to rotate is a problem, but it's really not a great
			// idea to stop your application just because you couldn't rename
			// your log.
			// We only return this error when explicitly needed.
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	rotateLog.outFh.Close()
	rotateLog.outFh = fh
	rotateLog.curFn = filename
	rotateLog.generation = generation

	return fh, nil
}

// CurrentFileName returns the current file name that
// the RotateLog object is writing to
func (rotateLog *RotateLog) CurrentFileName() string {
	rotateLog.mutex.RLock()
	defer rotateLog.mutex.RUnlock()
	return rotateLog.curFn
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type cleanupGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (cleanup *cleanupGuard) Enable() {
	cleanup.mutex.Lock()
	defer cleanup.mutex.Unlock()
	cleanup.enable = true
}
func (g *cleanupGuard) Run() {
	g.fn()
}

// Rotate forcefully rotates the log files. If the generated file name
// clash because file already exists, a numeric suffix of the form
// ".1", ".2", ".3" and so forth are appended to the end of the log file
//
// Thie method can be used in conjunction with a signal handler so to
// emulate servers that generate new log files when they receive a
// SIGHUP
func (rotateLog *RotateLog) Rotate() error {
	rotateLog.mutex.Lock()
	defer rotateLog.mutex.Unlock()
	if _, err := rotateLog.getWriter_nolock(true, true); err != nil {
		return err
	}
	return nil
}

func (rotateLog *RotateLog) rotate_nolock(filename string) error {
	lockfn := filename + `_lock`
	fh, err := os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		// Can't lock, just return
		return err
	}

	var guard cleanupGuard
	guard.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer guard.Run()

	if link := rotateLog.linkName; link != "" {
		// Create the link log path if not exists
		linkPath := filepath.Dir(link)
		if _, err := os.Stat(linkPath); os.IsNotExist(err) {
			err := os.MkdirAll(linkPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create link directory %s: %v", linkPath, err)
			}
		}

		tmpLinkName := filename + `_symlink`
		if err := os.Symlink(filename, tmpLinkName); err != nil {
			return fmt.Errorf(`failed to create new symlink. %w`, err)
		}

		if err := os.Rename(tmpLinkName, rotateLog.linkName); err != nil {
			return fmt.Errorf(`failed to rename new symlink. %w`, err)
		}
	}

	if rotateLog.maxAge <= 0 && rotateLog.rotationCount <= 0 {
		return errors.New("panic: maxAge and rotationCount are both set")
	}

	matches, err := filepath.Glob(rotateLog.globPattern)
	if err != nil {
		return err
	}

	cutoff := rotateLog.clock.Now().Add(-1 * rotateLog.maxAge)
	var toUnlink []string
	for _, path := range matches {
		// Ignore lock files
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		fl, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if rotateLog.maxAge > 0 && fi.ModTime().After(cutoff) {
			continue
		}

		if rotateLog.rotationCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	if rotateLog.rotationCount > 0 {
		// Only delete if we have more than rotationCount
		if rotateLog.rotationCount >= uint(len(toUnlink)) {
			return nil
		}

		toUnlink = toUnlink[:len(toUnlink)-int(rotateLog.rotationCount)]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	guard.Enable()
	go func() {
		// unlink files on a separate goroutine
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}

// Close satisfies the io.Closer interface. You must
// call this method if you performed any writes to
// the object.
func (rotateLog *RotateLog) Close() error {
	rotateLog.mutex.Lock()
	defer rotateLog.mutex.Unlock()

	if rotateLog.outFh == nil {
		return nil
	}

	rotateLog.outFh.Close()
	rotateLog.outFh = nil
	return nil
}

func getLogFilePattern() string {
	applicationName, podName := DEFAULT_APPLICATION_NAME, DEFAULT_POD_NAME

	if value, ok := os.LookupEnv(metadata.EnvApplicationName); ok {
		applicationName = value
	}

	if value, ok := os.LookupEnv(metadata.EnvPodName); ok {
		podName = value
	}

	return fmt.Sprintf("%s/%s/%s.%%Y-%%m-%%d.log", metadata.BASE_PATH_LOG, applicationName, podName)
}

// ============================ OPTIONS =============================
// IOption is used to pass optional arguments to
// the RotateLog constructor
type IOption interface {
	Name() string
	Value() interface{}
}

type Option struct {
	name  string
	value interface{}
}

func NewOption(name string, value interface{}) IOption {
	return &Option{
		name:  name,
		value: value,
	}
}

func (o *Option) Name() string {
	return o.name
}
func (o *Option) Value() interface{} {
	return o.value
}

const (
	OptKeyLogFilePattern = "file-pattern"
	OptKeyClock          = "clock"
	OptKeyLinkName       = "link-name"
	OptKeyMaxAge         = "max-age"
	OptKeyRotationTime   = "rotation-time"
	OptKeyRotationCount  = "rotation-count"
)

// WithRotateLogFilePattern create new Option sets a fileName pattern
// for the log file.You should use patterns
// using the strftime (3) format.
//
// For example: "/var/log/myapp/log.%Y%m%d%H%M%S"
//
// Defaults pattern: ${BASE_PATH}/${APPLICATION_NAME}/${POD_NAME}.yyyy-mm-dd.log
func WithRotateLogFilePattern(filePattern string) IOption {
	return NewOption(OptKeyLogFilePattern, filePattern)
}

// WithRotateLogClock creates a new Option that sets a clock
// that the RotateLogs object will use to determine
// the current time.
//
// By default rotatelogs.Local, which returns the
// current time in the local time zone, is used. If you
// would rather use UTC, use rotatelogs.UTC as the argument
// to this option, and pass it to the constructor.
func WithRotateLogClock(c Clock) IOption {
	return NewOption(OptKeyClock, c)
}

// WithRotateLogLocation creates a new Option that sets up a
// "Clock" interface that the RotateLogs object will use
// to determine the current time.
//
// This option works by always returning the in the given
// location.
func WithRotateLogLocation(loc *time.Location) IOption {
	return NewOption(OptKeyClock, clockFn(func() time.Time {
		return time.Now().In(loc)
	}))
}

// WithRotateLogLinkName creates a new Option that sets the
// symbolic link name that gets linked to the current
// file name being used.
func WithRotateLogLinkName(s string) IOption {
	return NewOption(OptKeyLinkName, s)
}

// WithRotateLogMaxAge creates a new Option that sets the
// max age of a log file before it gets purged from
// the file system.
func WithRotateLogMaxAge(d time.Duration) IOption {
	return NewOption(OptKeyMaxAge, d)
}

// WithRotateLogRotationTime creates a new Option that sets the
// time between rotation.
// Default: 24 * time.Hour
func WithRotateLogRotationTime(d time.Duration) IOption {
	return NewOption(OptKeyRotationTime, d)
}

// WithRotateLogRotationCount creates a new Option that sets the
// number of files should be kept before it gets
// purged from the file system.
func WithRotateLogRotationCount(n uint) IOption {
	return NewOption(OptKeyRotationCount, n)
}

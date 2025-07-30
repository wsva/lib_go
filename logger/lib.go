package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	LogLevelCritical LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelNotice
	LogLevelInfo
	LogLevelDebug
)

// LogLevelNames provides mapping for log levels
var LogLevelNames = map[LogLevel]string{
	LogLevelCritical: "CRITICAL",
	LogLevelError:    "ERROR",
	LogLevelWarning:  "WARNING",
	LogLevelNotice:   "NOTICE",
	LogLevelInfo:     "INFO",
	LogLevelDebug:    "DEBUG",
}

// Logger is the interface for outputing log messages in different levels.
// A new Logger can be created with NewLogger() function.
type Logger interface {
	// SetLevel changes the level of the logger. Default is logging.Info.
	SetLevel(LogLevel)

	// SetCallDepth sets the parameter passed to runtime.Caller().
	// It is used to get the file name from call stack.
	// For example you need to set it to 1 if you are using a wrapper around
	// the Logger. Default value is zero.
	SetCallDepth(int)

	// New creates a new inerhited context logger with given prefixes.
	New(prefixes ...interface{}) Logger

	// Fatal is equivalent to l.Critical followed by a call to os.Exit(1).
	Fatal(format string, args ...interface{})

	// Panic is equivalent to l.Critical followed by a call to panic().
	Panic(format string, args ...interface{})

	// Critical logs a message using CRITICAL as log level.
	Critical(format string, args ...interface{})

	// Error logs a message using ERROR as log level.
	Error(format string, args ...interface{})

	// Warning logs a message using WARNING as log level.
	Warning(format string, args ...interface{})

	// Notice logs a message using NOTICE as log level.
	Notice(format string, args ...interface{})

	// Info logs a message using INFO as log level.
	Info(format string, args ...interface{})

	// Debug logs a message using DEBUG as log level.
	Debug(format string, args ...interface{})
}

// LogRecord contains all of the information about a single log message.
type LogRecord struct {
	Format      string        // Format string
	Args        []interface{} // Arguments to format string
	LoggerName  string        // Name of the logger module
	Level       LogLevel      // Level of the record
	Time        time.Time     // Time of the record (local time)
	Filename    string        // File name of the log call (absolute path)
	Line        int           // Line number in file
	ProcessID   int           // PID
	ProcessName string        // Name of the process
}

type logger struct {
	Name      string
	Level     LogLevel
	Handler   *LogWriterHandler
	calldepth int
}

// NewLogger returns a new Logger implementation. Do not forget to close it at exit.
func NewLogger(name string, l LogLevel) *logger {
	return &logger{
		Name:    name,
		Level:   l,
		Handler: NewWriterHandler(os.Stdout, l),
	}
}

// New creates a new inerhited logger with the given prefixes
func (l *logger) New(prefixes ...interface{}) Logger {
	return newLoggerWithPrefix(*l, "", prefixes...)
}

func (l *logger) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *logger) SetCallDepth(n int) {
	l.calldepth = n
}

// Fatal is equivalent to Critical() followed by a call to os.Exit(1).
func (l *logger) Fatal(format string, args ...interface{}) {
	l.Critical(format, args...)
	l.Handler.Close()
	os.Exit(1)
}

// Panic is equivalent to Critical() followed by a call to panic().
func (l *logger) Panic(format string, args ...interface{}) {
	l.Critical(format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Critical sends a critical level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Critical(format string, args ...interface{}) {
	if l.Level >= LogLevelCritical {
		l.log(LogLevelCritical, format, args...)
	}
}

// Error sends a error level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Error(format string, args ...interface{}) {
	if l.Level >= LogLevelError {
		l.log(LogLevelError, format, args...)
	}
}

// Warning sends a warning level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Warning(format string, args ...interface{}) {
	if l.Level >= LogLevelWarning {
		l.log(LogLevelWarning, format, args...)
	}
}

// Notice sends a notice level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Notice(format string, args ...interface{}) {
	if l.Level >= LogLevelNotice {
		l.log(LogLevelNotice, format, args...)
	}
}

// Info sends a info level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Info(format string, args ...interface{}) {
	if l.Level >= LogLevelInfo {
		l.log(LogLevelInfo, format, args...)
	}
}

// Debug sends a debug level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Debug(format string, args ...interface{}) {
	if l.Level >= LogLevelDebug {
		l.log(LogLevelDebug, format, args...)
	}
}

func (l *logger) log(level LogLevel, format string, args ...interface{}) {
	// Add missing newline at the end.
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	_, file, line, ok := runtime.Caller(l.calldepth + 2)
	if !ok {
		file = "???"
		line = 0
	}

	rec := &LogRecord{
		Format:      format,
		Args:        args,
		LoggerName:  l.Name,
		Level:       level,
		Time:        time.Now(),
		Filename:    file,
		Line:        line,
		ProcessName: filepath.Base(os.Args[0]),
		ProcessID:   os.Getpid(),
	}

	l.Handler.Handle(rec)
}

// LogWriterHandler is a handler implementation that writes the logging output to a io.Writer.
type LogWriterHandler struct {
	Level LogLevel
	w     io.WriteCloser
}

// NewWriterHandler creates a new writer handler with given io.Writer
func NewWriterHandler(w io.WriteCloser, l LogLevel) *LogWriterHandler {
	return &LogWriterHandler{
		Level: l,
		w:     w,
	}
}

// FilterAndFormat filters any record according to loggging level
func (h *LogWriterHandler) FilterAndFormat(rec *LogRecord) string {
	if h.Level >= rec.Level {
		return fmt.Sprintf("%s [%s] %-8s %s",
			fmt.Sprint(rec.Time)[:19],
			rec.LoggerName,
			LogLevelNames[rec.Level],
			fmt.Sprintf(rec.Format, rec.Args...))
	}
	return ""
}

// Handle writes any given Record to the Writer.
func (b *LogWriterHandler) Handle(rec *LogRecord) {
	message := b.FilterAndFormat(rec)
	if message == "" {
		return
	}
	fmt.Fprint(b.w, message)
}

// Close closes WriterHandler
func (b *LogWriterHandler) Close() {
	b.w.Close()
}

type loggerWithPrefix struct {
	prefix string
	logger
}

// Fatal is equivalent to Critical() followed by a call to os.Exit(1).
func (c *loggerWithPrefix) Fatal(format string, args ...interface{}) {
	c.logger.Fatal(c.prefixFormat()+format, args...)
}

// Panic is equivalent to Critical() followed by a call to panic().
func (c *loggerWithPrefix) Panic(format string, args ...interface{}) {
	c.logger.Panic(c.prefixFormat()+format, args...)
}

// Critical sends a critical level log message to the handler. Arguments are
// handled in the manner of fmt.Printf.
func (c *loggerWithPrefix) Critical(format string, args ...interface{}) {
	c.logger.Critical(c.prefixFormat()+format, args...)
}

// Error sends a error level log message to the handler. Arguments are handled
// in the manner of fmt.Printf.
func (c *loggerWithPrefix) Error(format string, args ...interface{}) {
	c.logger.Error(c.prefixFormat()+format, args...)
}

// Warning sends a warning level log message to the handler. Arguments are
// handled in the manner of fmt.Printf.
func (c *loggerWithPrefix) Warning(format string, args ...interface{}) {
	c.logger.Warning(c.prefixFormat()+format, args...)
}

// Notice sends a notice level log message to the handler. Arguments are
// handled in the manner of fmt.Printf.
func (c *loggerWithPrefix) Notice(format string, args ...interface{}) {
	c.logger.Notice(c.prefixFormat()+format, args...)
}

// Info sends a info level log message to the handler. Arguments are handled in
// the manner of fmt.Printf.
func (c *loggerWithPrefix) Info(format string, args ...interface{}) {
	c.logger.Info(c.prefixFormat()+format, args...)
}

// Debug sends a debug level log message to the handler. Arguments are handled
// in the manner of fmt.Printf.
func (c *loggerWithPrefix) Debug(format string, args ...interface{}) {
	c.logger.Debug(c.prefixFormat()+format, args...)
}

// New creates a new Logger from current context
func (c *loggerWithPrefix) New(prefixes ...interface{}) Logger {
	return newLoggerWithPrefix(c.logger, c.prefix, prefixes...)
}

func (c *loggerWithPrefix) prefixFormat() string {
	return c.prefix + " "
}

func newLoggerWithPrefix(l logger, initial string, prefixes ...interface{}) *loggerWithPrefix {
	resultPrefix := "" // resultPrefix holds prefix after initialization
	connector := ""    // connector holds the connector string

	for _, prefix := range prefixes {
		resultPrefix += fmt.Sprintf("%s%+v", connector, prefix)
		switch connector {
		case "=": // if previous is `=` replace with ][
			connector = "]["
		case "][": // if previous is `][` replace with =
			connector = "="
		default:
			connector = "=" // if its first iteration, assing =
		}
	}

	return &loggerWithPrefix{
		prefix: initial + "[" + resultPrefix + "]",
		logger: l,
	}

}

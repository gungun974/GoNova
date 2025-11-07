package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AppLogger struct {
	entry *log.Entry
}

func newAppLogger(name string) AppLogger {
	return AppLogger{
		entry: log.WithField("logger", name),
	}
}

type callerContext string

var callerContextKey callerContext

func (l *AppLogger) setCaller() {
	pc, file, line, ok := runtime.Caller(2)

	if ok {
		l.entry.Context = context.WithValue(context.Background(), callerContextKey, runtime.Frame{
			PC:       pc,
			File:     file,
			Line:     line,
			Function: runtime.FuncForPC(pc).Name(),
		})
	}
}

// Log will log a message at the level given as parameter.
// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
// For this behaviour Entry.Panic or Entry.Fatal should be used instead.
func (l *AppLogger) Log(level log.Level, args ...any) {
	l.setCaller()
	l.entry.Log(level, args...)
}

func (l *AppLogger) Trace(args ...any) {
	l.setCaller()
	l.entry.Trace(args...)
}

func (l *AppLogger) Debug(args ...any) {
	l.setCaller()
	l.entry.Debug(args...)
}

func (l *AppLogger) Print(args ...any) {
	l.setCaller()
	l.entry.Print(args...)
}

func (l *AppLogger) Info(args ...any) {
	l.setCaller()
	l.entry.Info(args...)
}

func (l *AppLogger) Warn(args ...any) {
	l.setCaller()
	l.entry.Warn(args...)
}

func (l *AppLogger) Warning(args ...any) {
	l.setCaller()
	l.entry.Warning(args...)
}

func (l *AppLogger) Error(args ...any) {
	l.setCaller()
	l.entry.Error(args...)
}

func (l *AppLogger) Fatal(args ...any) {
	l.setCaller()
	l.entry.Fatal(args...)
}

func (l *AppLogger) Panic(args ...any) {
	l.setCaller()
	l.entry.Panic(args...)
}

// Entry Printf family functions

func (l *AppLogger) Logf(level log.Level, format string, args ...any) {
	l.setCaller()
	l.entry.Logf(level, format, args...)
}

func (l *AppLogger) Tracef(format string, args ...any) {
	l.setCaller()
	l.entry.Tracef(format, args...)
}

func (l *AppLogger) Debugf(format string, args ...any) {
	l.setCaller()
	l.entry.Debugf(format, args...)
}

func (l *AppLogger) Infof(format string, args ...any) {
	l.setCaller()
	l.entry.Infof(format, args...)
}

func (l *AppLogger) Printf(format string, args ...any) {
	l.setCaller()
	l.entry.Printf(format, args...)
}

func (l *AppLogger) Warnf(format string, args ...any) {
	l.setCaller()
	l.entry.Warnf(format, args...)
}

func (l *AppLogger) Warningf(format string, args ...any) {
	l.setCaller()
	l.entry.Warningf(format, args...)
}

func (l *AppLogger) Errorf(format string, args ...any) {
	l.setCaller()
	l.entry.Errorf(format, args...)
}

func (l *AppLogger) Fatalf(format string, args ...any) {
	l.setCaller()
	l.entry.Fatalf(format, args...)
}

func (l *AppLogger) Panicf(format string, args ...any) {
	l.setCaller()
	l.entry.Panicf(format, args...)
}

type customFormatter struct{}

var projectSourceFolder = ""

func (f *customFormatter) Format(entry *log.Entry) ([]byte, error) {
	levelColor := "\x1b[0m" // Reset
	levelText := strings.ToUpper(entry.Level.String())

	switch entry.Level {
	case log.TraceLevel:
		levelColor = "\x1b[37m"
	case log.DebugLevel:
		levelColor = "\x1b[32m"
	case log.InfoLevel:
		levelColor = "\x1b[34m"
	case log.WarnLevel:
		levelColor = "\x1b[33m"
	case log.ErrorLevel:
		levelColor = "\x1b[31m"
	case log.FatalLevel:
		levelColor = "\x1b[91m"
	case log.PanicLevel:
		levelColor = "\x1b[91m"
	}

	loggerName, ok := entry.Data["logger"]
	if !ok {
		loggerName = "unknown"
	}

	frame := entry.Context.Value(callerContextKey).(runtime.Frame)

	return fmt.Appendf(nil,
		"%s %-20s %-40s %-24s %s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"),
		fmt.Sprintf("%s\x1b[1m%s\x1b[0m", levelColor, levelText),
		fmt.Sprintf(
			"%s:%d",
			strings.TrimPrefix(frame.File, projectSourceFolder),
			frame.Line,
		),
		fmt.Sprintf("\x1b[1m%s\x1b[0m", loggerName),
		entry.Message,
	), nil
}

func init() {
	log.SetFormatter(&customFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)

	_, file, _, ok := runtime.Caller(0)
	if ok {
		projectSourceFolder = path.Dir(path.Dir(path.Dir(file)))
	}
}

var MainLogger = newAppLogger("MainLogger")

var DatabaseLogger = newAppLogger("DatabaseLogger")

var WatcherLogger = newAppLogger("WatcherLogger")

var CommandLogger = newAppLogger("CommandLogger")

var InjectorLogger = newAppLogger("InjectorLogger")

package logutil

import (
	"fmt"
	"golib/caller"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

func StandardLogger() *logrus.Logger {
	return logger
}

type ILogger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}

const sourceKey = "source"

func Writer() io.Writer {
	return logger.Writer()
}

func Muted() {
	logger.SetOutput(io.Discard)
}

func Outer() ILogger {
	return logger.WithField(sourceKey, caller.BriefInfoStr(3))
}

func Trace(args ...interface{}) {
	Outer().Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	Outer().Tracef(format, args...)
}

func Debug(args ...interface{}) {
	Outer().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	Outer().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	Outer().Infof(format, args...)
}

func Info(args ...interface{}) {
	Outer().Info(args...)
}

func Errorf(format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	Outer().Errorf(format, args...)
}

func Error(args ...interface{}) {
	Outer().Error(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	logger.Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	Outer().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	Outer().Fatalf(format, args...)
}

func Log(level logrus.Level, args ...interface{}) {
	logger.Log(level, args...)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	if false {
		_ = fmt.Sprintf(format, args...) // enable printf checker
	}
	logger.Logf(level, format, args...)
}

func Timing() func() {
	t0 := time.Now()
	lg := Outer()
	return func() {
		use := time.Since(t0)
		if use.Milliseconds() > 30 {
			lg.Tracef("function use time %s", use)
		}
	}
}

func DebugTiming() func() {
	t0 := time.Now()
	lg := Outer()
	return func() {
		use := time.Since(t0)
		lg.Debugf("function use time %s", use)
	}
}

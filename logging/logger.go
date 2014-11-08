package logging

import (
	"fmt"
	builtinlog "log"
)

var (
	prefixDebug = "debug"
	prefixInfo  = "info"
	prefixError = "error"
	prefixFatal = "fatal"
)

var (
	debug = false
)

func SetDebug() {
	debug = true
}

func log(prefix string, v ...interface{}) {
	builtinlog.Printf("%s: %s", prefix, fmt.Sprintln(v...))
}

func Debug(v ...interface{}) {
	if debug {
		log(prefixDebug, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	log(prefixInfo, v...)
}

func Infof(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	log(prefixError, v...)
}

func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	log(prefixFatal, v...)
}

func Fatalf(format string, v ...interface{}) {
	Fatal(fmt.Sprintf(format, v...))
}

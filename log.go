package main

import (
	"io"
	"os"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/op/go-logging"
)

// Lg is the Logger instance
var Lg = logging.MustGetLogger("kit")

var lumblog lumberjack.Logger

// CloseLogger and reset logger
func CloseLogger() {
	logging.Reset()
	lumblog.Close()
}

// SetupLogger configures Lg instance.
func SetupLogger(gfile string, glevel string, gformat string) {
	var iolog io.Writer
	// var lumberLog lumberjack.Logger

	if gfile == "" {
		iolog = os.Stderr
	} else {
		lumblog = lumberjack.Logger{
			Filename:   gfile,
			MaxSize:    10, // megabytes
			MaxBackups: 100,
			MaxAge:     100, //days
		}
		iolog = &lumblog
		// iolog = &lumberjack.Logger{
		// 	Filename:   gfile,
		// 	MaxSize:    10, // megabytes
		// 	MaxBackups: 100,
		// 	MaxAge:     100, //days
		// }
	}

	logBk := logging.NewLogBackend(iolog, "", 0)
	var format = logging.MustStringFormatter(gformat)
	// var format = logging.MustStringFormatter(
	// 	// `%{time:01-02 15:04:05} - %{level:.4s} - %{message}`,
	// 	`%{color}%{time:01-02 15:04:05} - %{level:.4s} %{shortfunc} %{color:reset} %{message}`,
	// )

	// `%{color}%{time:01-02 15:04:05} - %{level:.4s} %{shortfunc} %{color:reset} %{message}`,
	logBkFormatted := logging.NewBackendFormatter(logBk, format)
	logBkLeveled := logging.AddModuleLevel(logBkFormatted)

	switch glevel {
	case "DEBUG":
		logBkLeveled.SetLevel(logging.DEBUG, "")
	case "NOTICE":
		logBkLeveled.SetLevel(logging.NOTICE, "")
	case "WARN":
		logBkLeveled.SetLevel(logging.WARNING, "")
	case "ERROR":
		logBkLeveled.SetLevel(logging.ERROR, "")
	default:
		logBkLeveled.SetLevel(logging.INFO, "")
	}
	// logBkLeveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(logBkLeveled)
}

package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger with JSON formatting and appropriate log level
func Init() {
	log = logrus.New()

	// Set JSON formatter
	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	// Set output to stdout
	log.SetOutput(os.Stdout)

	// Set log level based on environment variable, default to info
	level := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn", "warning":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}

// Fields type is used to define structured logging fields
type Fields logrus.Fields

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}

// WithFields creates an entry with specified fields
func WithFields(fields Fields) *logrus.Entry {
	if log == nil {
		Init()
	}
	return log.WithFields(logrus.Fields(fields))
}

// Debug logs a message at the debug level
func Debug(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Debug(args...)
}

// Debugf logs a formatted message at the debug level
func Debugf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Debugf(format, args...)
}

// Info logs a message at the info level
func Info(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Info(args...)
}

// Infof logs a formatted message at the info level
func Infof(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Infof(format, args...)
}

// Warn logs a message at the warn level
func Warn(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Warn(args...)
}

// Warnf logs a formatted message at the warn level
func Warnf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Warnf(format, args...)
}

// Error logs a message at the error level
func Error(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Error(args...)
}

// Errorf logs a formatted message at the error level
func Errorf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Errorf(format, args...)
}

// Fatal logs a message at the fatal level and exits
func Fatal(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Fatal(args...)
}

// Fatalf logs a formatted message at the fatal level and exits
func Fatalf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Fatalf(format, args...)
}

// Panic logs a message at the panic level and panics
func Panic(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Panic(args...)
}

// Panicf logs a formatted message at the panic level and panics
func Panicf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.Panicf(format, args...)
}

package log

import (
	"fmt"
	"io"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

// The global logger
var (
	logger *logrus.Logger
	le     *logrus.Entry
)

// Default initializer
func init() {
	logger = logrus.New()
	logger.SetLevel(logrus.Level(InfoLevel))
	logger.SetOutput(io.Discard)
	le = logger.WithFields(logrus.Fields{})
}

// Initialize the Logger
func Init(logLevel Level, format Format) error {
	// Create new logger
	logger = logrus.New()
	logger.SetLevel(logrus.Level(logLevel))

	switch format {
	case TextFormat:
		formatter := new(logrus.TextFormatter)
		formatter.TimestampFormat = "2006-01-02 15:04:05.000"
		formatter.FullTimestamp = true
		logger.SetFormatter(formatter)
	case FluentdFormat:
		f := joonix.NewFormatter()
		if err := joonix.DisableTimestampFormat(f); err != nil {
			return err
		}
		logger.SetFormatter(f)
	case JsonFormat:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	le = logger.WithFields(logrus.Fields{})
	return nil
}

func SetLevel(logLevel Level) {
	logger.SetLevel(logrus.Level(logLevel))
}

func GetLevel() Level {
	return Level(logger.GetLevel())
}

func withFields(params ...interface{}) *logrus.Entry {
	fields := logrus.Fields{}
	for i := 0; i+1 < len(params); i += 2 {
		val := params[i]
		key, ok := val.(string)
		if !ok {
			key = fmt.Sprintf("%v", val)
		}
		fields[key] = params[i+1]
	}

	return le.WithFields(fields)
}

func Trace(v string, params ...interface{}) {
	if !logger.IsLevelEnabled(logrus.TraceLevel) {
		return
	}

	if len(params) > 1 {
		withFields(params...).Trace(v)
	} else {
		le.Trace(v)
	}
}

func Debug(v string, params ...interface{}) {
	if !logger.IsLevelEnabled(logrus.DebugLevel) {
		return
	}

	if len(params) > 1 {
		withFields(params...).Debug(v)
	} else {
		le.Debug(v)
	}
}

func Info(v string, params ...interface{}) {
	if !logger.IsLevelEnabled(logrus.InfoLevel) {
		return
	}

	if len(params) > 1 {
		withFields(params...).Info(v)
	} else {
		le.Info(v)
	}
}

func Warn(v string, params ...interface{}) {
	if !logger.IsLevelEnabled(logrus.WarnLevel) {
		return
	}

	if len(params) > 1 {
		withFields(params...).Warn(v)
	} else {
		le.Warn(v)
	}
}

func Error(v string, params ...interface{}) {
	if !logger.IsLevelEnabled(logrus.ErrorLevel) {
		return
	}

	if len(params) > 1 {
		withFields(params...).Error(v)
	} else {
		le.Error(v)
	}
}

func Fatal(v string, params ...interface{}) {
	if !logger.IsLevelEnabled(logrus.FatalLevel) {
		return
	}

	if len(params) > 1 {
		withFields(params...).Fatal(v)
	} else {
		le.Fatal(v)
	}
}

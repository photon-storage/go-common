package log

import (
	"context"
	"fmt"
	"io"
	"os"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

// The global logger
var g *log

const chSize int = 64

type settings struct {
	format    Format
	formatter *logrus.TextFormatter
}

type log struct {
	ctx      context.Context
	cancel   context.CancelFunc
	logger   *logrus.Logger
	entry    *logrus.Entry
	settings *settings
	ch       chan func()
}

func (l *log) log(f func()) {
	if l.ch == nil {
		f()
	} else {
		l.ch <- f
	}
}

func (l *log) loop() {
	if l.ch != nil {
		for {
			select {
			case f := <-l.ch:
				f()

			case <-l.ctx.Done():
				// Drain remaining logs in the channel.
				for {
					select {
					case f := <-l.ch:
						f()
					default:
						return
					}
				}
			}
		}
	}
}

func (l *log) stop() {
	l.cancel()
}

// Default initializer
func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.Level(InfoLevel))
	logger.SetOutput(io.Discard)

	ctx, cancel := context.WithCancel(context.Background())
	g = &log{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		entry:  logger.WithFields(logrus.Fields{}),
	}
}

// Initialize the Logger
func Init(logLevel Level, format Format, sync bool) error {
	// Create new logger
	logger := logrus.New()
	logger.SetLevel(logrus.Level(logLevel))
	settings := &settings{
		format: format,
	}

	switch format {
	case TextFormat:
		formatter := new(logrus.TextFormatter)
		formatter.TimestampFormat = "2006-01-02 15:04:05.000"
		formatter.FullTimestamp = true
		logger.SetFormatter(formatter)
		settings.formatter = formatter
	case FluentdFormat:
		f := joonix.NewFormatter()
		if err := joonix.DisableTimestampFormat(f); err != nil {
			return err
		}
		logger.SetFormatter(f)
	case JsonFormat:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logger.SetOutput(os.Stdout)

	if g != nil {
		g.stop()
	}
	ctx, cancel := context.WithCancel(context.Background())
	g = &log{
		ctx:      ctx,
		cancel:   cancel,
		logger:   logger,
		entry:    logger.WithFields(logrus.Fields{}),
		settings: settings,
	}

	if !sync {
		g.ch = make(chan func(), chSize)
		go g.loop()
	}

	return nil
}

func ForceColor() {
	if g.settings == nil {
		return
	}
	switch g.settings.format {
	case TextFormat:
		g.settings.formatter.ForceColors = true
		g.logger.SetFormatter(g.settings.formatter)
	case FluentdFormat:
	case JsonFormat:
	case JournaldFormat:
	}
}

func SetLevel(logLevel Level) {
	g.logger.SetLevel(logrus.Level(logLevel))
}

func GetLevel() Level {
	return Level(g.logger.GetLevel())
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

	return g.logger.WithFields(fields)
}

func Trace(v string, params ...interface{}) {
	if !g.logger.IsLevelEnabled(logrus.TraceLevel) {
		return
	}

	g.log(func() {
		if len(params) > 1 {
			withFields(params...).Trace(v)
		} else {
			g.entry.Trace(v)
		}
	})
}

func Debug(v string, params ...interface{}) {
	if !g.logger.IsLevelEnabled(logrus.DebugLevel) {
		return
	}

	g.log(func() {
		if len(params) > 1 {
			withFields(params...).Debug(v)
		} else {
			g.entry.Debug(v)
		}
	})
}

func Info(v string, params ...interface{}) {
	if !g.logger.IsLevelEnabled(logrus.InfoLevel) {
		return
	}

	g.log(func() {
		if len(params) > 1 {
			withFields(params...).Info(v)
		} else {
			g.entry.Info(v)
		}
	})
}

func Warn(v string, params ...interface{}) {
	if !g.logger.IsLevelEnabled(logrus.WarnLevel) {
		return
	}

	g.log(func() {
		if len(params) > 1 {
			withFields(params...).Warn(v)
		} else {
			g.entry.Warn(v)
		}
	})
}

func Error(v string, params ...interface{}) {
	if !g.logger.IsLevelEnabled(logrus.ErrorLevel) {
		return
	}

	g.log(func() {
		if len(params) > 1 {
			withFields(params...).Error(v)
		} else {
			g.entry.Error(v)
		}
	})
}

func Fatal(v string, params ...interface{}) {
	if !g.logger.IsLevelEnabled(logrus.FatalLevel) {
		return
	}

	g.log(func() {
		if len(params) > 1 {
			withFields(params...).Fatal(v)
		} else {
			g.entry.Fatal(v)
		}
	})
}

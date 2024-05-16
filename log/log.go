package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// The global logger
var g *log

const chSize int = 64

type log struct {
	ctx       context.Context
	cancel    context.CancelFunc
	logger    *logrus.Logger
	entry     *logrus.Entry
	opts      *Options
	formatter *TimeZoneFormatter
	ch        chan func()
	doneCh    chan bool
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
		defer close(l.doneCh)
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

type Options struct {
	Context  context.Context
	LogLevel Level
	Sync     bool
	Location *time.Location
}

// Initialize the Logger
func Init(opts *Options) error {
	if g != nil {
		g.stop()
	}

	formatter := newTimeZoneFormatter(opts.Location)

	// Create new logger
	logger := logrus.New()
	logger.SetLevel(logrus.Level(opts.LogLevel))
	logger.SetFormatter(formatter)
	logger.SetOutput(os.Stdout)

	ctx, cancel := context.WithCancel(opts.Context)
	g = &log{
		ctx:       ctx,
		cancel:    cancel,
		logger:    logger,
		entry:     logger.WithFields(logrus.Fields{}),
		opts:      opts,
		formatter: formatter,
	}

	if !opts.Sync {
		g.ch = make(chan func(), chSize)
		g.doneCh = make(chan bool)
		go g.loop()
	}

	return nil
}

func WaitForDone() {
	if g.doneCh != nil {
		<-g.doneCh
	}
}

func ForceColor() {
	g.formatter.formatter.ForceColors = true
	g.formatter.formatter.DisableColors = false
	g.logger.SetFormatter(g.formatter)
}

func DisableColor() {
	g.formatter.formatter.ForceColors = false
	g.formatter.formatter.DisableColors = true
	g.logger.SetFormatter(g.formatter)
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

func IsDebug() bool {
	return g.logger.IsLevelEnabled(logrus.DebugLevel)
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

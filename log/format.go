package log

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ErrLogFormatInvalid = errors.New("invalid log format name")
)

type Format uint32

const (
	TextFormat Format = iota
	JsonFormat
	FluentdFormat
	JournaldFormat
)

func ParseFormat(format string) (Format, error) {
	switch format {
	case "text":
		return TextFormat, nil
	case "fluentd":
		return FluentdFormat, nil
	case "json":
		return JsonFormat, nil
	case "journald":
		return JournaldFormat, nil
	}

	return TextFormat, ErrLogFormatInvalid
}

type TimeZoneFormatter struct {
	formatter *logrus.TextFormatter
	loc       *time.Location
}

func newTimeZoneFormatter(loc *time.Location) *TimeZoneFormatter {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	formatter.FullTimestamp = true
	return &TimeZoneFormatter{
		formatter: formatter,
		loc:       loc,
	}
}

func (f *TimeZoneFormatter) Format(e *logrus.Entry) ([]byte, error) {
	if f.loc != nil {
		e.Time = e.Time.In(f.loc)
	}
	return f.formatter.Format(e)
}

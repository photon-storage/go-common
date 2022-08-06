package log

import "errors"

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

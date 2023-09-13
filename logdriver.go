package logger

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gomig/utils"
)

// logDriver standard log message
type logDriver struct {
	typ        string
	tags       []string
	timeFormat string
	writers    []io.Writer
	formatter  TimeFormatter
}

func (logDriver) err(format string, args ...any) error {
	return utils.TaggedError([]string{"LogDriver"}, format, args...)
}

func (_log *logDriver) Type(t string) Log {
	_log.typ = t
	return _log
}

func (_log *logDriver) Tags(tags ...string) Log {
	_log.tags = append(_log.tags, tags...)
	return _log
}

func (_log logDriver) Print(format string, params ...any) error {
	for _, writer := range _log.writers {
		// Datetime
		_, err := writer.Write([]byte(_log.formatter(time.Now().UTC(), _log.timeFormat)))
		if err != nil {
			return _log.err(err.Error())
		}

		// Type
		t := []rune(strings.ToUpper(_log.typ))
		if len(t) >= 5 {
			t = t[0:5]
		}
		_, err = writer.Write([]byte(fmt.Sprintf("%6s ", string(t))))
		if err != nil {
			return _log.err(err.Error())
		}

		// Message
		_, err = writer.Write([]byte(fmt.Sprintf(strings.ReplaceAll(format, "\n", ""), params...)))
		if err != nil {
			return _log.err(err.Error())
		}

		// Tags
		for _, tag := range _log.tags {
			_, err = writer.Write([]byte(fmt.Sprintf(" [%s]", tag)))
			if err != nil {
				return _log.err(err.Error())
			}
		}

		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return _log.err(err.Error())
		}
	}

	return nil
}

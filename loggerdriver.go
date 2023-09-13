package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gomig/utils"
)

type loggerDriver struct {
	timeFormat string
	formatter  TimeFormatter
	writers    map[string]io.Writer
}

func (loggerDriver) err(format string, args ...any) error {
	return utils.TaggedError([]string{"LoggerDriver"}, format, args...)
}

func (_logger loggerDriver) log() Log {
	writers := make([]io.Writer, 0)
	for _, w := range _logger.writers {
		writers = append(writers, w)
	}
	return NewLog(_logger.timeFormat, _logger.formatter, "LOG", writers...)
}

func (_logger loggerDriver) Log() Log {
	return _logger.log().Type("LOG")
}

func (_logger loggerDriver) Error() Log {
	return _logger.log().Type("ERROR")
}

func (_logger loggerDriver) Warning() Log {
	return _logger.log().Type("WARN")
}

func (_logger loggerDriver) Divider(divider string, count uint8, title string) error {
	if title != "" {
		title = " " + title + " "
	}
	if len(title)%2 != 0 {
		title = title + " "
	}

	if count%2 != 0 {
		count++
	}
	halfCount := int(count) - len(title)
	if halfCount <= 0 {
		halfCount = 2
	} else {
		halfCount = halfCount / 2
	}
	for _, writer := range _logger.writers {
		_, err := writer.Write([]byte(strings.Repeat(divider, halfCount) + strings.ToUpper(title) + strings.Repeat(divider, halfCount) + "\n"))
		if err != nil {
			return _logger.err(err.Error())
		}
	}
	return nil
}

func (_logger loggerDriver) Raw(format string, params ...any) error {
	for _, writer := range _logger.writers {
		_, err := writer.Write([]byte(fmt.Sprintf(format, params...)))
		if err != nil {
			return _logger.err(err.Error())
		}
	}
	return nil
}

func (_logger loggerDriver) JSON(data any) error {
	if _bytes, err := json.MarshalIndent(data, "", "    "); err != nil {
		return err
	} else {
		for _, writer := range _logger.writers {
			_, err := writer.Write([]byte(string(_bytes)))
			if err != nil {
				return _logger.err(err.Error())
			}
		}
		return nil
	}
}

func (_logger *loggerDriver) AddWriter(name string, writer io.Writer) {
	_logger.writers[name] = writer
}

func (_logger *loggerDriver) RemoveWriter(name string) {
	delete(_logger.writers, name)
}

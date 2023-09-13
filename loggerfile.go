package logger

import (
	"os"
	"path"
	"time"

	"github.com/gomig/utils"
)

type fileLogger struct {
	path       string
	prefix     string
	timeFormat string
	formatter  TimeFormatter
}

func (fileLogger) err(format string, args ...any) error {
	return utils.TaggedError([]string{"FileLogger"}, format, args...)
}

func (_loggerF fileLogger) Write(data []byte) (int, error) {
	utils.CreateDirectory(_loggerF.path)
	filename := _loggerF.prefix + " " + _loggerF.formatter(time.Now().UTC(), _loggerF.timeFormat) + ".log"
	filename = path.Join(_loggerF.path, filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, _loggerF.err(err.Error())
	}
	defer f.Close()
	n, err := f.Write(data)
	if err != nil {
		err = _loggerF.err(err.Error())
	}
	return n, err
}

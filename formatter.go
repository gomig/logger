package logger

import (
	"time"

	"github.com/gomig/jalaali"
)

// TimeFormatter for log date
type TimeFormatter func(t time.Time, format string) string

// GregorianFormatter gregorian date formatter
func GregorianFormatter(t time.Time, format string) string {
	return t.Format(format)
}

// JalaaliFormatter jalaali (tehran) date formatter
func JalaaliFormatter(t time.Time, format string) string {
	return jalaali.NewTehran(time.Now()).Format(format)
}

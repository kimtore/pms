package log

import (
	"fmt"
	"strings"
	"time"
)

type Level int

const (
	ErrorLevel Level = iota
	InfoLevel
	DebugLevel
)

var strLevel = map[Level]string{
	ErrorLevel: "ERROR",
	InfoLevel:  "INFO",
	DebugLevel: "DEBUG",
}

var buffer []string

var since time.Time

func Lines() []string {
	return buffer
}

var maxLevel Level

func SetLevel(level Level) {
	maxLevel = level
}

func ParseLevel(level string) (Level, error) {
	low := strings.ToLower(level)
	for lv, s := range strLevel {
		if low == strings.ToLower(s) {
			return lv, nil
		}
	}
	return ErrorLevel, fmt.Errorf("no such level: %s", level)
}

func appendLine(data string) {
	buffer = append(buffer, data)
}

func init() {
	since = time.Now()
	maxLevel = DebugLevel
	Clear()
}

func Clear() {
	buffer = make([]string, 0)
}

// Printf adds a line to the local buffer, and optionally prints it to an external log writer.
func Printf(format string, args ...interface{}) (int, error) {
	format = strings.Trim(format, " \t\n")
	formatted := fmt.Sprintf(format, args...)
	// TODO: split lines?
	appendLine(formatted)
	return writer.Write([]byte(formatted + "\n"))
}

func Logf(format string, level Level, args ...interface{}) (int, error) {
	if level > maxLevel {
		return 0, nil
	}
	ts := time.Since(since).Seconds()
	prefix := fmt.Sprintf("[%010.3f] [%s] ", ts, strLevel[level])
	return Printf(prefix+format, args...)
}

func Errorf(format string, args ...interface{}) (int, error) {
	return Logf(format, ErrorLevel, args...)
}

func Infof(format string, args ...interface{}) (int, error) {
	return Logf(format, InfoLevel, args...)
}

func Debugf(format string, args ...interface{}) (int, error) {
	return Logf(format, DebugLevel, args...)
}

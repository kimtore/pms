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

var since time.Time

var maxLevel Level

func init() {
	since = time.Now()
	maxLevel = DebugLevel
}

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

// Logf adds a line to the local buffer, and optionally prints it to an external log writer.
func Logf(format string, level Level, args ...interface{}) (int, error) {
	if level > maxLevel {
		return 0, nil
	}
	format = strings.Trim(format, " \t\n")
	formatted := fmt.Sprintf(format, args...)
	msg := Message{
		Timestamp: time.Now(),
		Level:     level,
		Text:      formatted,
	}
	appendMessage(msg)
	return printMsg(msg)
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
func printMsg(msg Message) (int, error) {
	prefix := fmt.Sprintf("[%010.3f] [%s] ", time.Now().Sub(msg.Timestamp).Seconds(), strLevel[msg.Level])
	return writer.Write([]byte(prefix + msg.Text + "\n"))
}


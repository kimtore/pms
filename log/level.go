package log

import (
	"fmt"
	"strings"
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

var maxLevel Level

func (level Level) String() string {
	return strLevel[level]
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

// Package console provides logging functions.
package console

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

var start = time.Now()

// Open opens a log file for writing.
func Open(logfile string) (err error) {
	logFile, err = os.Create(logfile)
	if err != nil {
		return
	}
	return
}

// Close closes an open log file.
func Close() {
	logFile.Close()
}

// Log writes a log line to the log file.
// A timestamp and a newline is automatically added.
// If the log file isn't open, nothing is done.
func Log(format string, args ...interface{}) {
	if logFile == nil {
		return
	}
	since := time.Since(start)
	text := fmt.Sprintf(format, args...)
	text = fmt.Sprintf("[%.5f] %s\n", since.Seconds(), text)
	logFile.WriteString(text)
}

// Package console wraps a console logger into a writer.
// Program-wide logging will go both to a console and a file.

package console

import (
	log "github.com/sirupsen/logrus"
	"io"
)

var LogLines = make([]string, 0)

type writer struct {
	w io.Writer
}

var _ io.Writer = &writer{}

func (w *writer) Write(data []byte) (int, error) {
	LogLines = append(LogLines, string(data))
	return w.w.Write(data)
}

func Writer(w io.Writer) io.Writer {
	return &writer{w: w}
}

// Log writes a log line to the log file.
// A timestamp and a newline is automatically added.
// If the log file isn't open, nothing is done.
func Log(format string, args ...interface{}) {
	log.Printf(format, args...)
}

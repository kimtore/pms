// Package console is a legacy wrapper around the log package.

package console

import (
	"github.com/ambientsound/pms/log"
)

// Log writes a log line to the log file.
// A timestamp and a newline is automatically added.
// If the log file isn't open, nothing is done.
func Log(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

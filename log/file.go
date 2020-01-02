package log

import (
	"io"
	"os"
)

var writer io.Writer = &blackhole{}

type blackhole struct{}

// Write implements Writer. Blackhole writer discards any data written to it.
func (b *blackhole) Write(data []byte) (int, error) {
	return len(data), nil
}

// Open a file for log writing. If successful, set the internal log writer
// to that instance. Returns an error if unsuccessful.
func Configure(filename string, overwrite bool) error {
	w, err := fileWriter(filename, overwrite)
	if err != nil {
		return err
	}
	writer = w
	SetLevel(DebugLevel)
	writeHistory()
	return nil
}

func fileWriter(filename string, overwrite bool) (io.Writer, error) {
	logMode := os.O_WRONLY | os.O_CREATE
	if overwrite {
		logMode |= os.O_TRUNC
	} else {
		logMode |= os.O_APPEND
	}
	return os.OpenFile(filename, logMode, 0666)
}

func writeHistory() {
	messages := Messages(DebugLevel)
	for _, msg := range messages {
		_, _ = msg.Write(writer)
	}
}

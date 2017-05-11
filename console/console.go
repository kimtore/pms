package console

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

var start = time.Now()

func Open(logfile string) (err error) {
	logFile, err = os.Create(logfile)
	if err != nil {
		return
	}
	return
}

func Close() {
	logFile.Close()
}

func Log(format string, args ...interface{}) {
	if logFile == nil {
		return
	}
	since := time.Since(start)
	text := fmt.Sprintf(format, args...)
	text = fmt.Sprintf("[%.5f] %s\n", since.Seconds(), text)
	logFile.WriteString(text)
}

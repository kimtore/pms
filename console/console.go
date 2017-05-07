package console

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

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
	text := fmt.Sprintf(format, args...)
	text = fmt.Sprintf("%s %s\n", time.Now().String(), text)
	logFile.WriteString(text)
}

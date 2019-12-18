package log

import (
	"github.com/ambientsound/pms/config"
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
func Configure(cfg config.Log) error {
	w, err := fileWriter(cfg)
	if err != nil {
		return err
	}
	writer = w
	lvl, err := ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	SetLevel(lvl)
	return nil
}

func fileWriter(cfg config.Log) (io.Writer, error) {
	logMode := os.O_WRONLY | os.O_CREATE
	if cfg.Overwrite {
		logMode |= os.O_TRUNC
	} else {
		logMode |= os.O_APPEND
	}
	return os.OpenFile(cfg.File, logMode, 0666)
}

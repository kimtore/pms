package console

import (
	"log"
)

func Log(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

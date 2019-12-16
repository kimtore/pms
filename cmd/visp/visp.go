package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var buildVersion = "undefined"

func main() {
	exitCode, err := run()
	if exitCode != 0 {
		log.Error(err)
	}
	os.Exit(exitCode)
}

func run() (int, error) {
	return 0, nil
}

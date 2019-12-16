package main

import (
	"fmt"
	"github.com/ambientsound/pms/config"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
)

var buildVersion = "undefined"

const (
	ExitSuccess = iota
	ExitConfiguration
)

func main() {
	exitCode, err := run()
	if exitCode != ExitSuccess {
		log.Error(err)
	}
	os.Exit(exitCode)
}

func run() (int, error) {
	cfg, err := config.Configuration()
	if err != nil {
		flag.Usage()
		return ExitConfiguration, err
	}
	fmt.Println(cfg.Spotify.Username)
	return ExitSuccess, nil
}

package pms

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/ambientsound/pms/options"
)

// SourceDefaultConfig reads, parses, and executes the default config.
func (pms *PMS) SourceDefaultConfig() error {
	reader := strings.NewReader(options.Defaults)
	return pms.SourceConfig(reader)
}

// SourceConfigFile reads, parses, and executes a config file.
func (pms *PMS) SourceConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return pms.SourceConfig(file)
}

// SourceConfig reads, parses, and executes config lines.
func (pms *PMS) SourceConfig(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		err := pms.CLI.Execute(scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}

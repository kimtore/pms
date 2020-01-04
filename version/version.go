// Package version provides access to the program name and compiled version.
package version

import (
	"strconv"
	"time"
)

var (
	Program   = "Visp"
	ShortName = "Visp"
	buildDate = "0"
	Version   = "unknown"
)

func BuildDate() time.Time {
	i, _ := strconv.ParseInt(buildDate, 10, 64)
	return time.Unix(i, 0)
}

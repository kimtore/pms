// Package version provides access to the program name and compiled version.
package version

const shortName string = "PMS"
const longName string = "Practical Music Search"

var version string = "undefined"

func ShortName() string {
	return shortName
}

func LongName() string {
	return longName
}

func Version() string {
	return version
}

func SetVersion(v string) {
	version = v
}

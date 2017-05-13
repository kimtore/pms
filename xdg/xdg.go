// Package xdg provides file paths for cache and configuration, as specified by
// the XDG Base Directory Specification.
//
// See https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
package xdg

import (
	"os"
	"path"
	"strings"
)

// appendPmsDirectory adds "pms" to a directory tree.
func appendPmsDirectory(dir string) string {
	return path.Join(dir, "pms")
}

// ConfigDirectories returns a list of configuration directories. The least
// important directory is listed first.
func ConfigDirectories() []string {
	dirs := make([]string, 0)

	// $XDG_CONFIG_DIRS defines the preference-ordered set of base directories
	// to search for configuration files in addition to the $XDG_CONFIG_HOME base
	// directory. The directories in $XDG_CONFIG_DIRS should be seperated with a
	// colon ':'.
	xdgConfigDirs := os.Getenv("XDG_CONFIG_DIRS")
	if len(xdgConfigDirs) == 0 {
		xdgConfigDirs = "/etc/xdg"
	}

	// Add entries from $XDG_CONFIG_DIRS to directory list.
	configDirs := strings.Split(xdgConfigDirs, ":")
	for i := len(configDirs) - 1; i >= 0; i-- {
		if len(configDirs[i]) > 0 {
			dir := appendPmsDirectory(configDirs[i])
			dirs = append(dirs, dir)
		}
	}

	// $XDG_CONFIG_HOME defines the base directory relative to which user
	// specific configuration files should be stored. If $XDG_CONFIG_HOME is
	// either not set or empty, a default equal to $HOME/.config should be used.
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgConfigHome) == 0 {
		xdgConfigHome = path.Join(os.Getenv("HOME"), ".config")
	}
	dir := appendPmsDirectory(xdgConfigHome)

	// Add $XDG_CONFIG_HOME to directory list.
	dirs = append(dirs, dir)

	return dirs
}

// CacheDirectory returns the cache base directory.
func CacheDirectory() string {
	// $XDG_CACHE_HOME defines the base directory relative to which user
	// specific non-essential data files should be stored. If $XDG_CACHE_HOME is
	// either not set or empty, a default equal to $HOME/.cache should be used.
	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	if len(xdgCacheHome) == 0 {
		xdgCacheHome = path.Join(os.Getenv("HOME"), ".cache")
	}

	return path.Join(xdgCacheHome, "pms")
}

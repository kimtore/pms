package xdg

import (
	"os"
	"path"
)

func CacheDirectory() string {
	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		dir = path.Join(os.Getenv("HOME"), ".cache")
	}
	dir = path.Join(dir, "pms")
	return dir
}

package utils

import (
	"fmt"
)

func TimeString(secs int) string {
	if secs < 0 {
		return "--:--"
	}
	hours := int(secs / 3600)
	secs = secs % 3600
	minutes := int(secs / 60)
	secs = secs % 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

func TimeRunes(secs int) []rune {
	return []rune(TimeString(secs))
}

// ReverseRunes returns a new, reversed rune slice.
func ReverseRunes(src []rune) []rune {
	dest := make([]rune, len(src))
	for i, j := 0, len(src)-1; i <= j; i, j = i+1, j-1 {
		dest[i], dest[j] = src[j], src[i]
	}
	return dest
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

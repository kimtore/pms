// package utils provides simple transformation functions which do not fit anywhere else in particular.
package utils

import (
	"fmt"
	"strings"
)

// TimeString formats length in seconds as H:mm:ss.
func TimeString(secs int) string {
	if secs < 0 {
		return `--:--`
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

// TimeRunes acts as TimeString, but returns a slice of runes.
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

// TokenFilter returns a subset of tokens that match the specified prefix.
func TokenFilter(match string, tokens []string) []string {
	dest := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if strings.HasPrefix(tok, match) {
			dest = append(dest, tok)
		}
	}
	return dest
}

// Min returns the minimum of a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of a and b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

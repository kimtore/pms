package utils

import "fmt"

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

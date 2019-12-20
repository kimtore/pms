package constants

type InputMode int

// Different input modes are handled in different ways. Check
// MultibarWidget.inputMode against these constants.
const (
	MultibarModeNormal InputMode = iota
	MultibarModeInput
	MultibarModeSearch
)

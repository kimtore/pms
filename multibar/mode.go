package multibar

type InputMode int

// Different input modes are handled in different ways. Check
// Multibar.inputMode against these constants.
const (
	ModeNormal InputMode = iota
	ModeInput
	ModeSearch
)

func (m InputMode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInput:
		return "INPUT"
	case ModeSearch:
		return "SEARCH"
	default:
		panic("BUG: unnamed input mode")
	}
}

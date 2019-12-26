package widgets

import (
	"github.com/ambientsound/pms/list"
	"strings"
	"unicode"
)

type column struct {
	col          list.Column
	key          string
	title        string
	rightPadding int
	width        int
}

func ColumnTitle(key string) string {
	var words = make([]string, 0)
	var word = make([]rune, 0)

	split := func() {
		if len(word) > 0 {
			words = append(words, string(word))
			word = make([]rune, 0)
		}
	}

	for _, r := range key {
		if unicode.IsUpper(r) {
			split()
		}
		if len(word) == 0 {
			r = unicode.ToUpper(r)
		}
		word = append(word, r)
	}

	split()

	return strings.Join(words, " ")
}

package actions_test

import (
	"testing"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/actions"
)

// TestSetParser tests
func TestSetParser(t *testing.T) {
	str := "foo=bar baz invfoo nobar"
	parse := actions.NewSetParser()
	token := input.Token{}
	pos := 0
	npos := 0
	_ = parse
	for {
		token, npos = input.NextToken(str[pos:])
		pos += npos
		if token.Class == input.TokenEnd {
			break
		}
	}
}

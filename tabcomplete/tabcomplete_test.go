package tabcomplete_test

import (
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/tabcomplete"
	"github.com/stretchr/testify/assert"
)

var tabCompleteTests = []struct {
	input       string
	success     bool
	completions []string
}{
	{"", true, commands.Keys()},
	{"s", true, []string{
		"se",
		"select",
		"set",
		"sort",
		"stop",
		"style",
	}},
	{"set", true, []string{}},
	{"add ", true, []string{}},
	{"cursor next-of", true, []string{}},
	{"cursor next-of ", true, []string{
		"cursor next-of artist",
		"cursor next-of title",
	}},
	{"foobarbaz", false, []string{}},
	{"foobarbaz ", false, []string{}},
	{"$var", false, []string{}},
	{"{foo", false, []string{}},
	{"# bar", false, []string{}},
}

func TestTabComplete(t *testing.T) {
	for n, test := range tabCompleteTests {

		api := api.NewTestAPI()

		t.Logf("### Test %d: '%s'", n+1, test.input)

		clen := len(test.completions)
		tabComplete := tabcomplete.New(test.input, api)
		sentences := make([]string, clen)
		i := 0

		for i < len(sentences) {
			sentence, err := tabComplete.Scan()
			if test.success {
				assert.Nil(t, err, "Expected success when parsing '%s'", test.input)
			} else {
				assert.NotNil(t, err, "Expected error when parsing '%s'", test.input)
			}
			sentences[i] = sentence
			i++
			if i == clen {
				break
			}
		}

		assert.Equal(t, test.completions, sentences)
		assert.Equal(t, clen, tabComplete.Len())
	}
}

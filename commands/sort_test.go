package commands_test

import (
	"fmt"
	"testing"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/stretchr/testify/assert"
)

var sortTests = []commands.Test{
	// Valid forms
	{``, true, initSort, testSorting, []string{}},
	{`artist title`, true, initSort, testSorting, []string{"title"}},

	// Invalid forms
	{`$`, false, nil, nil, []string{}},
}

func TestSort(t *testing.T) {
	commands.TestVerb(t, "sort", sortTests)
}

func testSorting(data *commands.TestData) {
	// FIXME: test actual sorting
	err := data.Cmd.Exec()
	assert.Nil(data.T, err)
}

func initSort(data *commands.TestData) {
	// Set up the sort option
	// FIXME
	opts := data.Api.Options()
	opts.Add(options.NewStringOption("sort"))
	opts.Get("sort").Set("title")

	list := data.Api.Songlist()
	for i := 0; i < 2; i++ {
		for j := 0; j < 10; j++ {
			s := song.New()
			s.SetTags(mpd.Attrs{
				"artist": fmt.Sprintf("artist %d", 2-i),
				"title":  fmt.Sprintf("title %d", 10-j),
			})
			list.Add(s)
		}
	}
}

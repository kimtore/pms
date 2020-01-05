package commands

import (
	"fmt"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/spotify/aggregator"
	"github.com/google/uuid"
	"strconv"
	"strings"

	"github.com/ambientsound/pms/api"
)

var (
	tagMaps = map[string]string{
		"albumArtist": "artist",
	}
)

// Isolate searches for songs that have similar tags as the selection.
type Isolate struct {
	command
	api  api.API
	tags []string
}

// NewIsolate returns Isolate.
func NewIsolate(api api.API) Command {
	return &Isolate{
		api:  api,
		tags: make([]string, 0),
	}
}

// Parse implements Command.
func (cmd *Isolate) Parse() error {
	var err error
	list := cmd.api.List()
	cmd.tags, err = cmd.ParseTags(list.ColumnNames())
	return err
}

// Exec implements Command.
func (cmd *Isolate) Exec() error {
	list := cmd.api.Tracklist()
	if list == nil {
		return fmt.Errorf("isolate only works within a track list")
	}

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	selection := list.Selection()
	if selection.Len() != 1 {
		return fmt.Errorf("isolate operates on exactly one track")
	}

	row := selection.Row(0)
	queries := make([]string, len(cmd.tags))
	for i, tag := range cmd.tags {
		val := strconv.Quote(row[tag])
		if v, ok := tagMaps[tag]; ok {
			tag = v
		}
		queries[i] = fmt.Sprintf("%s:%s", tag, val)
	}

	query := strings.Join(queries, " AND ")
	log.Debugf("isolate search: %s", query)
	result, err := spotify_aggregator.Search(*client, query, cmd.api.Options().GetInt(options.Limit))

	if err != nil {
		return err
	}

	if result.Len() == 0 {
		return fmt.Errorf("no results found when isolating by %s", strings.Join(cmd.tags, ", "))
	}

	// Post-processing: sort in default order
	sort := cmd.api.Options().GetString(options.Sort)

	err = result.Sort(strings.Split(sort, ","))
	if err != nil {
		log.Errorf("error sorting: %s", err)
	}

	result.SetVisibleColumns(list.VisibleColumns())
	result.SetID(uuid.New().String())

	// Figure out a clever name
	parts := make([]string, len(cmd.tags))
	for i, tag := range cmd.tags {
		parts[i] = tag + ":" + row[tag]
	}
	result.SetName(strings.Join(parts, ", "))

	// Activate results
	cmd.api.SetList(result)

	return nil
}

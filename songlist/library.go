package songlist

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
)

// Library is a Songlist which represents the MPD song library.
type Library struct {
	BaseSonglist
	index           *index.Index
	version         int
	shutdownReIndex chan int
}

func NewLibrary() (s *Library) {
	s = &Library{
		shutdownReIndex: make(chan int, 1),
		version:         -1,
	}
	s.clear()
	return
}

func (s *Library) Name() string {
	return "Library"
}

func (s *Library) SetName(name string) error {
	return fmt.Errorf("The song library name cannot be changed.")
}

func (s *Library) Clear() error {
	return fmt.Errorf("The song library cannot be cleared because it is read-only.")
}

func (s *Library) Delete() error {
	return fmt.Errorf("The song library cannot be deleted using PMS. Try 'rm -rf' in your favorite shell.")
}

func (s *Library) Sort(fields []string) error {
	return fmt.Errorf("The song library is read-only. Please make a copy if you want to sort.")
}

func (s *Library) Remove(index int) error {
	return fmt.Errorf("The song library is read-only.")
}

func (s *Library) RemoveIndices(indices []int) error {
	return fmt.Errorf("The song library is read-only.")
}

// OpenIndex configures the library to use the Bleve search index at the specified path.
func (s *Library) OpenIndex(path string) error {
	var err error

	if s.HasIndex() {
		if err = s.index.Close(); err != nil {
			return err
		}
		s.index = nil
	}

	s.index, err = index.New(path)

	return err
}

// HasIndex returns true if the library has a search index.
func (s *Library) HasIndex() bool {
	return s.index != nil
}

// IndexSynced returns true if the search index is up to date with the MPD version.
func (s *Library) IndexSynced() bool {
	return s.HasIndex() && s.index.Version() == s.version
}

// CloseIndex closes the Bleve search index.
func (s *Library) CloseIndex() error {
	if s.HasIndex() {
		return s.index.Close()
	}
	return nil
}

// SetVersion sets the library version. This is expected to be a UNIX timestamp.
func (s *Library) SetVersion(version int) {
	s.version = version
}

// Version returns the library version.
func (s *Library) Version() int {
	return s.version
}

// ReIndex starts an asynchronous reindexing job. In case this function is
// called again before reindexing is done, ReIndex will abort the old
// reindexing job.
func (s *Library) ReIndex() {
	s.shutdownReIndex <- 0
	s.shutdownReIndex = make(chan int, 1)
	go func() {
		timer := time.Now()
		err := s.index.IndexFull(s.Songs(), s.shutdownReIndex)
		console.Log("Song library index complete, took %s", time.Since(timer).String())

		if err != nil {
			console.Log("Error occurred during library reindex: %s", err)
			return
		}
		s.index.SetVersion(s.Version())
	}()
}

// Search does a search in the Bleve index for a specific natural language
// query string, and returns a new Songlist with the search results.
func (s *Library) Search(q string) (Songlist, error) {
	if !s.HasIndex() {
		return nil, fmt.Errorf("Search index is not open.")
	}

	query := bleve.NewQueryStringQuery(q)
	request := bleve.NewSearchRequest(query)
	request.Size = s.Len()

	ids, _, err := s.index.Query(request)
	if err != nil {
		return nil, err
	}

	list := New()
	list.SetName(q)

	for _, id := range ids {
		song := s.Song(id)
		if song == nil {
			return nil, fmt.Errorf("Search index is corrupt.")
		}
		list.Add(song)
	}

	return list, nil
}

// Isolate takes a songlist and a set of tag keys, and matches the tag values
// of the songlist against the search index.
func (s *Library) Isolate(songs Songlist, tags []string) (Songlist, error) {
	if !s.HasIndex() {
		return nil, fmt.Errorf("Search index is not open.")
	}

	terms := make(map[string]struct{})
	query := bleve.NewBooleanQuery()

	// Create a cartesian join for song values and tag list.
	for _, song := range songs.Songs() {
		subQuery := bleve.NewConjunctionQuery()

		for _, tag := range tags {

			// Ignore empty values
			tagValue := song.StringTags[tag]
			if len(tagValue) == 0 {
				continue
			}

			// Name generation
			terms[tagValue] = struct{}{}

			field := strings.Title(tag)
			query := bleve.NewMatchPhraseQuery(tagValue)
			query.SetField(field)
			subQuery.AddQuery(query)
		}
		query.AddShould(subQuery)
	}

	// Construct a fitting name for this track list
	names := make([]string, 0)
	for k := range terms {
		names = append(names, k)
	}
	name := strings.Join(names, ", ")

	// Make the search
	request := bleve.NewSearchRequest(query)
	request.Size = s.Len()
	r, _, err := s.index.Query(request)
	list := s.Indices(r)

	list.SetName(name)

	return list, err
}

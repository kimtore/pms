package index

import (
	"strings"

	"github.com/ambientsound/pms/console"
	index_song "github.com/ambientsound/pms/index/song"
	"github.com/ambientsound/pms/songlist"

	"github.com/blevesearch/bleve"

	"fmt"
	"strconv"
)

const INDEX_BATCH_SIZE int = 1000

const SEARCH_SCORE_THRESHOLD float64 = 0.5

type Index struct {
	bleveIndex bleve.Index
	Songlist   songlist.Songlist
}

func New(loc string, s songlist.Songlist) (*Index, error) {
	var err error
	i := &Index{}
	i.bleveIndex, err = i.open(loc)
	if err != nil {
		i.bleveIndex, err = i.create(loc)
	}
	i.Songlist = s
	return i, err
}

func (i *Index) create(loc string) (index bleve.Index, err error) {
	mapping, err := buildIndexMapping()
	if err != nil {
		panic(err)
	}
	index, err = bleve.New(loc, mapping)
	if err != nil {
		console.Log("Error while creating index %s: %s", loc, err)
	}
	return
}

func (i *Index) open(loc string) (index bleve.Index, err error) {
	index, err = bleve.Open(loc)
	if err != nil {
		console.Log("Cannot open index %s: %s", loc, err)
	}
	return
}

func (i *Index) Close() error {
	return i.bleveIndex.Close()
}

// Index the entire Songlist.
func (i *Index) IndexFull() error {
	var err error

	// All operations are batched, currently INDEX_BATCH_SIZE are committed each iteration.
	b := i.bleveIndex.NewBatch()

	for pos, s := range i.Songlist.Songs() {
		is := index_song.New(s)
		err = b.Index(strconv.Itoa(pos), is)
		if err != nil {
			return err
		}
		if pos%INDEX_BATCH_SIZE == 0 {
			console.Log("Indexing songs %d/%d...", pos, i.Songlist.Len())
			i.bleveIndex.Batch(b)
			b.Reset()
		}
	}
	console.Log("Indexing last batch...")
	i.bleveIndex.Batch(b)

	console.Log("Finished indexing.")

	return nil
}

// Search takes a natural language query string, matches it against the search
// index, and returns a new Songlist with all matching songs.
func (i *Index) Search(q string) (songlist.Songlist, error) {
	query := bleve.NewQueryStringQuery(q)
	request := bleve.NewSearchRequest(query)

	r, _, err := i.Query(request)
	r.SetName(q)

	return r, err
}

// Isolate takes a songlist and a set of tag keys, and matches the tag values
// of the songlist against the search index.
func (i *Index) Isolate(list songlist.Songlist, tags []string) (songlist.Songlist, error) {
	terms := make(map[string]struct{})
	query := bleve.NewBooleanQuery()
	songs := list.Songs()

	// Create a cartesian join for song values and tag list.
	for _, song := range songs {
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

	request := bleve.NewSearchRequest(query)
	r, _, err := i.Query(request)

	names := make([]string, 0)
	for k, _ := range terms {
		names = append(names, k)
	}
	name := strings.Join(names, ", ")
	r.SetName(name)

	return r, err
}

// Query takes a Bleve search request and returns a songlist with all matching songs.
func (i *Index) Query(request *bleve.SearchRequest) (songlist.Songlist, *bleve.SearchResult, error) {
	r := songlist.New()
	request.Size = i.Songlist.Len()

	sr, err := i.bleveIndex.Search(request)

	if err != nil {
		return r, nil, err
	}

	for _, hit := range sr.Hits {
		if hit.Score < SEARCH_SCORE_THRESHOLD {
			break
		}
		id, err := strconv.Atoi(hit.ID)
		if err != nil {
			panic(fmt.Sprintf("Error when converting index IDs to integer: %s", err))
		}
		song := i.Songlist.Song(id)
		r.Add(song)
		//console.Log("%.2f %s\n", hit.Score, song.Tags["file"])
	}

	console.Log("Query '%s' returned %d results over threshold of %.2f (total %d results) in %s", request, r.Len(), SEARCH_SCORE_THRESHOLD, sr.Total, sr.Took)

	return r, sr, nil
}

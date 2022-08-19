package index

import (
	"bufio"
	"os"
	"path"
	"time"

	"github.com/ambientsound/pms/console"
	index_song "github.com/ambientsound/pms/index/song"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/xdg"

	"github.com/blevesearch/bleve/v2"

	"fmt"
	"strconv"
)

const INDEX_BATCH_SIZE int = 1000

const SEARCH_SCORE_THRESHOLD float64 = 0.5

type Index struct {
	bleveIndex bleve.Index
	path       string
	indexPath  string
	statePath  string
	version    int
}

func createDirectory(dir string) error {
	dirMode := os.ModeDir | 0755
	return os.MkdirAll(dir, dirMode)
}

// New opens a Bleve index and returns Index. In case an index is not found at
// the given path, a new one is created. In case of an error, nil is returned,
// and the error object set accordingly.
func New(basePath string) (*Index, error) {
	var err error

	timer := time.Now()

	err = createDirectory(basePath)
	if err != nil {
		return nil, fmt.Errorf("while creating %s: %s", basePath, err)
	}

	i := &Index{}
	i.path = basePath
	i.indexPath = path.Join(i.path, "index")
	i.statePath = path.Join(i.path, "state")

	// Try to stat the Bleve index path. If it does not exist, create it.
	if _, err := os.Stat(i.indexPath); err != nil {
		if os.IsNotExist(err) {
			i.bleveIndex, err = create(i.indexPath)
			if err != nil {
				return nil, fmt.Errorf("while creating index at %s: %s", i.indexPath, err)
			}

			// After successful creation, reset the MPD library version.
			err = i.SetVersion(0)
			if err != nil {
				return nil, fmt.Errorf("while zeroing out library version at %s: %s", i.statePath, err)
			}

		} else {
			// In case of any other filesystem error, abort operation.
			return nil, fmt.Errorf("while accessing %s: %s", i.indexPath, err)
		}

	} else {

		// If index was statted ok, try to open it.
		i.bleveIndex, err = open(i.indexPath)
		if err != nil {
			return nil, fmt.Errorf("while opening index at %s: %s", i.indexPath, err)
		}
		i.version, err = i.readVersion()
		if err != nil {
			console.Log("index state file is broken: %s", err)
		}
	}

	console.Log("Opened search index in %s", time.Since(timer).String())

	return i, nil
}

// Close closes a Bleve index.
func (i *Index) Close() error {
	return i.bleveIndex.Close()
}

// create creates a Bleve index at the given file system location.
func create(path string) (bleve.Index, error) {
	mapping, err := buildIndexMapping()
	if err != nil {
		return nil, fmt.Errorf("BUG: unable to create search index mapping: %s", err)
	}

	index, err := bleve.New(path, mapping)
	if err != nil {
		return nil, fmt.Errorf("while creating search index %s: %s", path, err)
	}

	return index, nil
}

// open opens a Bleve index at the given file system location.
func open(path string) (bleve.Index, error) {
	index, err := bleve.Open(path)
	if err != nil {
		return nil, fmt.Errorf("while opening search index %s: %s", path, err)
	}

	return index, nil
}

// Path returns the absolute path to where indexes and state for a specific MPD
// server should be stored.
func Path(host, port string) string {
	cacheDir := xdg.CacheDirectory()
	return path.Join(cacheDir, host, port)
}

// SetVersion writes the MPD library version to the state file.
func (i *Index) SetVersion(version int) error {
	file, err := os.Create(i.statePath)
	if err != nil {
		return err
	}
	defer file.Close()
	str := fmt.Sprintf("%d\n", version)
	file.WriteString(str)
	i.version = version
	return nil
}

// readVersion reads the MPD library version from the state file.
func (i *Index) readVersion() (int, error) {
	file, err := os.Open(i.statePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		version, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, err
		}
		return version, nil
	}

	return 0, fmt.Errorf("No data in index mpd library state file")
}

// Version returns the index version. It should correspond to the MPD library version.
func (i *Index) Version() int {
	return i.version
}

// Index the entire Songlist.
func (i *Index) IndexFull(songs []*song.Song, shutdown <-chan int) error {
	songChan := make(chan *song.Song, len(songs))
	console.Log("Feeding all songs into song queue...")
	for _, s := range songs {
		songChan <- s
	}
	console.Log("Done feeding songs.")
	return fullIndex(i.bleveIndex, songChan, shutdown)
}

// fullIndex indexes a stream of songs. This process can be aborted by sending
// a message on the shutdown channel.
func fullIndex(index bleve.Index, songs <-chan *song.Song, shutdown <-chan int) error {
	var err error

	count := 0
	batch := make(chan int, 1)
	size := len(songs)
	console.Log("Start full index.")

	// All operations are batched, currently INDEX_BATCH_SIZE are committed each iteration.
	b := index.NewBatch()

outer:
	for {
		select {
		case n := <-batch:
			console.Log("Indexing songs %d/%d...", count, size)
			index.Batch(b)
			b.Reset()
			if n < 0 {
				break outer
			}
		case s := <-songs:
			is := index_song.New(s)
			err = b.Index(strconv.Itoa(count), is)
			if err != nil {
				return err
			}
			if count%INDEX_BATCH_SIZE == 0 {
				batch <- count
			}
			count += 1
		case _ = <-shutdown:
			return fmt.Errorf("Aborting index batch at position %d", count)
		default:
			batch <- -1
		}
	}

	console.Log("Finished indexing.")

	return nil
}

// Query takes a Bleve search request and returns a songlist with all matching songs.
func (i *Index) Query(request *bleve.SearchRequest) ([]int, *bleve.SearchResult, error) {
	//request.Size = 1000

	sr, err := i.bleveIndex.Search(request)

	if err != nil {
		return make([]int, 0), nil, err
	}

	r := make([]int, 0, len(sr.Hits))

	for _, hit := range sr.Hits {
		if hit.Score < SEARCH_SCORE_THRESHOLD {
			break
		}
		id, err := strconv.Atoi(hit.ID)
		if err != nil {
			return r, nil, fmt.Errorf("Index is corrupt; error when converting index IDs to integer: %s", err)
		}
		r = append(r, id)
	}

	console.Log("Query '%v' returned %d results over threshold of %.2f (total %d results) in %s", request, len(r), SEARCH_SCORE_THRESHOLD, sr.Total, sr.Took)

	return r, sr, nil
}

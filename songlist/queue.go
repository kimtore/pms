package songlist

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/song"
	"github.com/fhs/gompd/v2/mpd"
)

// Queue is a Songlist which represents the MPD play queue.
type Queue struct {
	BaseSonglist
	mpdClient func() *mpd.Client
}

func NewQueue(mpdClient func() *mpd.Client) (s *Queue) {
	s = &Queue{}
	s.mpdClient = mpdClient
	s.clear()
	return
}

func (s *Queue) Name() string {
	return "Queue"
}

// Add adds a song to MPD's queue.
func (s *Queue) Add(song *song.Song) error {
	client := s.mpdClient()
	if client == nil {
		return fmt.Errorf("Cannot communicate with MPD")
	}
	return client.Add(song.StringTags["file"])
}

// AddList appends a songlist to the queue.
func (s *Queue) AddList(songlist Songlist) error {
	client := s.mpdClient()
	if client == nil {
		return fmt.Errorf("Cannot communicate with MPD")
	}
	commandList := client.BeginCommandList()
	if commandList == nil {
		return fmt.Errorf("Cannot begin command list")
	}
	songs := songlist.Songs()
	for _, song := range songs {
		commandList.Add(song.StringTags["file"])
	}
	return commandList.End()
}

// Insert inserts a song at a specified position in the queue.
func (s *Queue) Insert(song *song.Song, position int) error {
	client := s.mpdClient()
	if client == nil {
		return fmt.Errorf("Cannot communicate with MPD")
	}
	_, err := client.AddID(song.StringTags["file"], position)
	return err
}

// InsertList inserts the songs in a songlist into the queue, at a specified position.
func (s *Queue) InsertList(list Songlist, position int) error {
	// Get MPD client
	client := s.mpdClient()
	if client == nil {
		return fmt.Errorf("Cannot communicate with MPD")
	}

	// Create command list
	commandList := client.BeginCommandList()
	if commandList == nil {
		return fmt.Errorf("Cannot begin command list")
	}

	// Insert songs at incrementing positions
	songs := list.Songs()
	for _, song := range songs {
		commandList.AddID(song.StringTags["file"], position)
		position++
	}
	return commandList.End()
}

func (s *Queue) SetName(name string) error {
	return fmt.Errorf("The queue name cannot be changed.")
}

func (s *Queue) Clear() error {
	return fmt.Errorf("Clearing the queue is not implemented yet.")
}

func (s *Queue) Delete() error {
	return fmt.Errorf("The queue cannot be removed.")
}

func (s *Queue) Sort(fields []string) error {
	return fmt.Errorf("Sorting the queue is not implemented yet.")
}

func (s *Queue) Remove(index int) error {
	song := s.Song(index)
	if song == nil {
		return fmt.Errorf("Out of bounds")
	}
	client := s.mpdClient()
	if client == nil {
		return fmt.Errorf("Cannot communicate with MPD")
	}
	console.Log("Telling MPD to delete queue song ID %d", song.ID)
	return client.DeleteID(song.ID)
}

// RemoveIndices removes a selection of songs from MPD's queue.
func (s *Queue) RemoveIndices(indices []int) error {
	client := s.mpdClient()
	if client == nil {
		return fmt.Errorf("Cannot communicate with MPD")
	}
	commandList := client.BeginCommandList()
	if commandList == nil {
		return fmt.Errorf("Cannot begin command list")
	}

	//sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		song := s.Song(i)
		if song != nil {
			commandList.DeleteID(song.ID)
		}
	}

	return commandList.End()
}

// Merge incorporates songs from another songlist, replacing songs that has the same position.
func (q *Queue) Merge(s Songlist) (*Queue, error) {
	newQueue := NewQueue(q.mpdClient)

	oldSongs := q.Songs()
	for i := range oldSongs {
		song := *oldSongs[i]
		newQueue.add(&song)
	}

	newSongs := s.Songs()
	for i := range newSongs {
		song := newSongs[i]
		switch {
		case song.Position < 0:
			return nil, fmt.Errorf("Song number %d does not have a position", i)
		case song.Position == newQueue.Len():
			newQueue.add(song)
		case song.Position > newQueue.Len():
			return nil, fmt.Errorf("Song number %d has position greater than list length, there are parts missing", i)
		default:
			if err := newQueue.Replace(song.Position, song); err != nil {
				return nil, fmt.Errorf("Cannot replace song %d: %s", i, err)
			}
		}
	}

	return newQueue, nil
}

// IndexAtSong returns true if the song at the specified index is at a song with the same ID.
func (q *Queue) IndexAtSong(i int, song *song.Song) bool {
	check := q.Song(i)
	return song != nil && check != nil && check.ID == song.ID
}

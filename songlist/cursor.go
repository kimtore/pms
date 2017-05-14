package songlist

import (
	"github.com/ambientsound/pms/song"
)

// CursorSong returns the song currently selected by the cursor.
func (s *BaseSonglist) CursorSong() *song.Song {
	return s.Song(s.Cursor())
}

// MoveCursorUp moves the cursor up by the specified offset.
func (s *BaseSonglist) MoveCursorUp(i int) {
	s.MoveCursor(-i)
}

// MoveCursorUp moves the cursor down by the specified offset.
func (s *BaseSonglist) MoveCursorDown(i int) {
	s.MoveCursor(i)
}

// MoveCursorUp moves the cursor by the specified offset.
func (s *BaseSonglist) MoveCursor(i int) {
	s.SetCursor(s.Cursor() + i)
}

// SetCursor sets the cursor to an absolute position.
func (s *BaseSonglist) SetCursor(i int) {
	s.cursor = i
	s.ValidateCursor(0, s.Len()-1)
	s.expandVisualSelection()
}

// Cursor returns the cursor position.
func (s *BaseSonglist) Cursor() int {
	return s.cursor
}

func (s *BaseSonglist) CursorToSong(song *song.Song) error {
	index, err := s.Locate(song)
	if err != nil {
		return err
	}
	//console.Log("Located %s at position %d, id %d", song.StringTags["file"], index, song.ID)
	s.SetCursor(index)
	return nil
}

// ValidateCursor makes sure the cursor is within minimum and maximum boundaries.
func (s *BaseSonglist) ValidateCursor(ymin, ymax int) {
	if s.Cursor() < ymin {
		s.cursor = ymin
	}
	if s.Cursor() > ymax {
		s.cursor = ymax
	}
}

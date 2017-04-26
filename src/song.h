/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2015  Kim Tore Jensen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * song.h
 * 	contains what info is stored about a song
 */

#ifndef _SONG_H_
#define _SONG_H_

#include <string>
#include <mpd/client.h>

typedef signed long song_t;

#define MPD_SONG_NO_TIME -1
#define MPD_SONG_NO_ID -1
#define MPD_SONG_NO_NUM -1

using namespace std;

/*
 * Remember to update this as libmpd changes.
 */
class Song
{
public:
			Song(const mpd_song *);
			Song(const Song *);
			Song(const string);
			~Song();
	
	/* Common function to initialize special fields that MPD don't return */

	void		init();
	string		strip_leading_zeroes(string * src);
	string		dirname();

	/**
	 * Match this song against a search string and criteria flags.
	 *
	 * Returns true if song matches, false otherwise.
	 */
	bool		match(string term, long flags);

	/* Custom parameters only used by PMS */
	
	bool		selected;
	string		trackshort;
	string		discshort;

	/* Standard parameters imported from libmpdclient.h */

	string		file;
	string		artist;
	string		albumartist;
	string		albumartistsort;
	string		artistsort;
	string		title;
	string		album;
	string		track;
	string		name;
	string		date;
	string		originaldate;
	string		year;
	string		originalyear;

	string		genre;
	string		composer;
	string		performer;
	string		disc;
	string		comment;

	int		time;
	song_t		pos;
	song_t		id;
};

#endif /* _SONG_H_ */

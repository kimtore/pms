/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2010  Kim Tore Jensen
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
#include "libmpdclient.h"

using namespace std;

/*
 * Remember to update this as libmpd changes.
 */
class Song
{
public:
			Song(mpd_Song *);
			Song(Song *);
			Song(string);
			~Song();
	
	/* Common function to initialize special fields that MPD don't return */

	void		init();
	string		dirname();

	/* Custom parameters only used by PMS */
	
	bool		selected;
	string		trackshort;

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
	string		year;

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

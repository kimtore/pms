/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS	<<Practical Music Search>>
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
 * song.cpp
 *	contains what info is stored about a song
 */


#include "song.h"
#include "songlist.h"
#include "pms.h"
#include <string>
#include <vector>


Song::Song(const mpd_song * song)
{
	selected	= false;

	assert(mpd_song_get_uri(song) != NULL);

	file			= Pms::tostring(mpd_song_get_uri(song));
	artist			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_ARTIST, 0));
	albumartist		= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_ALBUM_ARTIST, 0));

	/* If the AlbumArtist tag is missing, and songs are sorted by that tag,
	 * the sort will be messed up. This is a hack which may misrepresent
	 * the songs, but at least sorting and displaying will look nice. */
	if (!albumartist.size()) {
		albumartist = artist;
	}

/* The ability to process the ArtistSort and AlbumArtistSort tags were added to
 * libmpdclient in version 2.11. If these tags are not provided, a rudimentary
 * sort rewrite will be performed in init(). */
#if LIBMPDCLIENT_CHECK_VERSION(2, 11, 0)
	artistsort		= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_ARTIST_SORT, 0));
	albumartistsort		= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_ALBUM_ARTIST_SORT, 0));
#else
	artistsort		= "";
	albumartistsort		= "";
#endif

	title			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_TITLE, 0));
	album			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_ALBUM, 0));
	track			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_TRACK, 0));
	name			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_NAME, 0));
	date			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_DATE, 0));

	genre			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_GENRE, 0));
	composer		= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_COMPOSER, 0));
	performer		= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_PERFORMER, 0));
	disc			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_DISC, 0));
	comment			= Pms::tostring(mpd_song_get_tag(song, MPD_TAG_COMMENT, 0));

	time			= mpd_song_get_duration(song);
	pos			= mpd_song_get_pos(song);
	id			= mpd_song_get_id(song);

	init();
}

Song::Song(const Song * song)
{
	selected		= false;

	file			= song->file;
	artist			= song->artist;
	albumartist		= song->albumartist;
	artistsort		= song->artistsort;
	albumartistsort		= song->albumartistsort;
	title			= song->title;
	album			= song->album;
	track			= song->track;
	trackshort		= song->trackshort;
	name			= song->name;
	date			= song->date;

	genre			= song->genre;
	composer		= song->composer;
	performer		= song->performer;
	disc			= song->disc;
	comment			= song->comment;

	time			= song->time;
	pos			= song->pos;
	id			= song->id;

	init();
}

Song::Song(const string uri)
{
	selected		= false;

	file			= uri;
	artist			= "";
	albumartist		= "";
	artistsort		= "";
	albumartistsort		= "";
	title			= "";
	album			= "";
	track			= "";
	trackshort		= "";
	name			= "";
	date			= "";
	year			= "";

	genre			= "";
	composer		= "";
	performer		= "";
	disc			= "";
	comment			= "";

	time			= MPD_SONG_NO_TIME;
	pos			= MPD_SONG_NO_NUM;
	id			= MPD_SONG_NO_ID;
}

Song::~Song()
{
}

/*
 * Initialize custom parameters
 */
void		Song::init()
{
	const string			the = "the";
	string				tmp;
	vector<string *>		original;
	vector<string *>		rewritten;
	vector<string *>::iterator	src;
	vector<string *>::iterator	dest;

	/* year from date */
	if (date.size() >= 4) {
		year = date.substr(0, 4);
	} else {
		year = "";
	}

	/* strip zeros and total tracks from the 'track' tag,
	 * and store it in 'trackshort'. */
	trackshort = strip_leading_zeroes(&track);

	/* Generate rudimentary sort names if none available, by
	 * rewriting 'The Artist' to 'Artist, The'. */
	if (artistsort.size() == 0) {
		original.push_back(&artist);
		rewritten.push_back(&artistsort);
	}
	if (albumartistsort.size() == 0) {
		original.push_back(&albumartist);
		rewritten.push_back(&albumartistsort);
	}

	src = original.begin();
	dest = rewritten.begin();

	while (src != original.end()) {

		/* String is too short, skip. */
		if ((*src)->size() <= the.size()) {
			**dest = **src;
			goto next;
		}

		tmp = (*src)->substr(0, 4);
		if (!lcstrcmp(tmp, the)) {
			**dest = **src;
			goto next;
		}

		/* If artist name consists of "the ...", place it at the end of the string. */
		**dest = (*src)->substr(4) + ", " + (*src)->substr(0, 3);

next:
		++src;
		++dest;
	}
}

string
Song::strip_leading_zeroes(string * src)
{
	bool zero = true;
	string s;
	string::const_iterator iter;

	iter = src->begin();
	while (iter != src->end() && *iter != '/') {
		if (!zero || *iter != '0') {
			zero = false;
			s += *iter;
		}
		++iter;
	}

	return s;
}

/*
 * Return directory name
 */
string		Song::dirname()
{
	string		ret = "";
	size_t		p;

	if (file.size() == 0)
		return ret;

	p = file.find_last_of("/\\");
	if (p == string::npos)
		return ret;

	return file.substr(0, p);
}

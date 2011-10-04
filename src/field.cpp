/* vi:set ts=8 sts=8 sw=8:
 *
 * Practical Music Search
 * Copyright (c) 2006-2011  Kim Tore Jensen
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
 */

#include "field.h"
#include "song.h"
#include "mpd.h"
#include <vector>
#include <string>
using namespace std;

extern MPD mpd;

Field::Field(field_t nfield, string name, string mpd_name, string tit, unsigned int minl, unsigned int maxl)
{
	type = nfield;
	str = name;
	cstr = mpd_name;
	title = tit;
	minlen = minl;
	maxlen = maxl;
}

string Field::format(Song * song)
{
	string tmp;

	if (!song)
		return "";

	if (type < FIELD_COLUMN_VALUES)
		return song->f[type];

	switch(type)
	{
		case FIELD_ELAPSED:
			return time_format((int)mpd.status.elapsed);

		case FIELD_MODES:
			tmp = "----";
			if (mpd.status.repeat)
				tmp[0] = 'r';
			if (mpd.status.random)
				tmp[1] = 'z';
			if (mpd.status.single)
				tmp[2] = 's';
			if (mpd.status.consume)
				tmp[3] = 'c';
			return tmp;

		case FIELD_STATE:
			if (mpd.status.state == MPD_STATE_PLAY)
				return "Playing";
			else if (mpd.status.state == MPD_STATE_STOP)
				return "Stopped";
			else if (mpd.status.state == MPD_STATE_PAUSE)
				return "Paused";
			break;

		case FIELD_QUEUESIZE:
			return tostring(mpd.playlist.size());

		case FIELD_QUEUELENGTH:
			return tostring(mpd.playlist.size());

		default:
			break;
	}

	return "";
}

Fieldtypes::Fieldtypes()
{
	fields.push_back(new Field(FIELD_POS, "pos", "Pos", "Pos", 0, 0));
	fields.push_back(new Field(FIELD_ID, "id", "Id", "ID", 0, 0));
	fields.push_back(new Field(FIELD_TIME, "length", "Time", "Length", 5, 7));
	fields.push_back(new Field(FIELD_DIRECTORY, "directory", "directory", "Directory", 0, 0));
	fields.push_back(new Field(FIELD_FILE, "file", "file", "Filename", 0, 0));
	fields.push_back(new Field(FIELD_NAME, "name", "Name", "Name", 0, 0));
	fields.push_back(new Field(FIELD_ARTIST, "artist", "Artist", "Artist", 0, 0));
	fields.push_back(new Field(FIELD_ARTISTSORT, "artist", "ArtistSort", "Artist", 0, 0));
	fields.push_back(new Field(FIELD_ALBUM, "album", "Album", "Album", 0, 0));
	fields.push_back(new Field(FIELD_TITLE, "title", "Title", "Title", 0, 0));
	fields.push_back(new Field(FIELD_TRACK, "track", "Track", "Track", 5, 5));
	fields.push_back(new Field(FIELD_DATE, "date", "Date", "Date", 4, 10));
	fields.push_back(new Field(FIELD_GENRE, "genre", "Genre", "Genre", 0, 0));
	fields.push_back(new Field(FIELD_DISC, "disc", "Disc", "Disc", 4, 4));
	fields.push_back(new Field(FIELD_ALBUMARTIST, "albumartist", "AlbumArtist", "Album artist", 0, 0));
	fields.push_back(new Field(FIELD_ALBUMARTISTSORT, "albumartistsort", "AlbumArtistSort", "Album artist", 0, 0));

	fields.push_back(new Field(FIELD_YEAR, "year", "", "Year", 4, 4));
	fields.push_back(new Field(FIELD_TRACKSHORT, "trackshort", "", "#", 2, 2));

	/* Topbar fields */
	fields.push_back(new Field(FIELD_ELAPSED, "elapsed", "", "", 0, 0));
	fields.push_back(new Field(FIELD_MODES, "modes", "", "", 0, 0));
	fields.push_back(new Field(FIELD_STATE, "state", "", "", 0, 0));
	fields.push_back(new Field(FIELD_QUEUESIZE, "queuesize", "", "", 0, 0));
	fields.push_back(new Field(FIELD_QUEUELENGTH, "queuelength", "", "", 0, 0));
}

Fieldtypes::~Fieldtypes()
{
	vector<Field *>::iterator i;
	for (i = fields.begin(); i != fields.end(); ++i)
		delete *i;
	fields.clear();
}

Field *	Fieldtypes::find_mpd(string &value)
{
	vector<Field *>::iterator i;
	for (i = fields.begin(); i != fields.end(); ++i)
		if ((*i)->cstr == value)
			return *i;
	return NULL;
}

Field *	Fieldtypes::find(string &value)
{
	vector<Field *>::iterator i;
	for (i = fields.begin(); i != fields.end(); ++i)
		if ((*i)->str == value)
			return *i;
	return NULL;
}

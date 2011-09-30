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
#include <vector>
#include <string>
using namespace std;

Field::Field(field_t nfield, string name, string mpd_name, unsigned int minl, unsigned int maxl)
{
	type = nfield;
	str = name;
	cstr = mpd_name;
	minlen = minl;
	maxlen = maxl;
}

Fieldtypes::Fieldtypes()
{
	fields.push_back(new Field(FIELD_POS, "pos", "Pos", 0, 0));
	fields.push_back(new Field(FIELD_ID, "id", "Id", 0, 0));
	fields.push_back(new Field(FIELD_TIME, "length", "Time", 5, 7));
	fields.push_back(new Field(FIELD_DIRECTORY, "directory", "directory", 0, 0));
	fields.push_back(new Field(FIELD_FILE, "file", "file", 0, 0));
	fields.push_back(new Field(FIELD_NAME, "name", "Name", 0, 0));
	fields.push_back(new Field(FIELD_ARTIST, "artist", "Artist", 0, 0));
	fields.push_back(new Field(FIELD_ARTISTSORT, "artist", "ArtistSort", 0, 0));
	fields.push_back(new Field(FIELD_ALBUM, "album", "Album", 0, 0));
	fields.push_back(new Field(FIELD_TITLE, "title", "Title", 0, 0));
	fields.push_back(new Field(FIELD_TRACK, "track", "Track", 5, 5));
	fields.push_back(new Field(FIELD_DATE, "date", "Date", 4, 10));
	fields.push_back(new Field(FIELD_GENRE, "genre", "Genre", 0, 0));
	fields.push_back(new Field(FIELD_DISC, "disc", "Disc", 4, 4));
	fields.push_back(new Field(FIELD_ALBUMARTIST, "albumartist", "AlbumArtist", 0, 0));
	fields.push_back(new Field(FIELD_ALBUMARTISTSORT, "albumartistsort", "AlbumArtistSort", 0, 0));

	fields.push_back(new Field(FIELD_YEAR, "year", "", 4, 4));
	fields.push_back(new Field(FIELD_TRACKSHORT, "trackshort", "", 2, 2));
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

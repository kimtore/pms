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
#include "config.h"
#include "curses.h"
#include "window.h"
#include <vector>
#include <string>
using namespace std;

extern MPD mpd;
extern Config config;
extern Curses curses;
extern Windowmanager wm;

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
	Wmain * win;
	Wsonglist * ws;
	string tmp;
	int i;

	if (type < FIELD_COLUMN_VALUES)
	{
		if (!song)
			return "";
		return song->f[type];
	}

	switch(type)
	{
		case FIELD_ELAPSED:
			return time_format((int)mpd.status.elapsed);

		case FIELD_REMAINING:
			return time_format((int)(mpd.status.length - mpd.status.elapsed));

		case FIELD_PROGRESSBAR:
			tmp.clear();
			if (mpd.status.length == -1 || mpd.status.elapsed == -1)
				return tmp;
			i = mpd.status.elapsed * (curses.topbar.right - curses.topbar.left) / mpd.status.length;
			while (i-- >= 0)
				tmp += '=';
			tmp += '>';
			return tmp;

		case FIELD_VOLUME:
			if (mpd.status.volume == 0 && config.mute)
				tmp = "Muted (" + tostring(config.volume) + "%%)";
			else
				tmp = tostring(mpd.status.volume) + "%%";
			return tmp;

		case FIELD_MODES:
			tmp = "----";
			if (mpd.status.repeat)
				tmp[0] = 'r';
			if (config.random)
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
			i = mpd.playlist.size();
			if (mpd.currentsong)
				i -= mpd.playlist.spos(mpd.currentsong->pos);
			tmp = tostring(i);
			if (mpd.status.repeat && mpd.active_songlist == &mpd.playlist)
				tmp += '+';
			return tmp;

		case FIELD_QUEUELENGTH:
			if (mpd.currentsong)
				tmp = time_format(mpd.playlist.length(mpd.currentsong->pos) - (mpd.currentsong->time == -1 ? 0 : (int)mpd.status.elapsed));
			else
				tmp = time_format(mpd.playlist.length());
			if (mpd.status.repeat && mpd.active_songlist == &mpd.playlist)
				tmp += '+';
			return tmp;

		case FIELD_LISTSIZE:
			if ((win = WMAIN(wm.active)) != NULL)
				return tostring(win->content_size());
			else
				return "0";

		case FIELD_LISTLENGTH:
			if ((ws = WSONGLIST(wm.active)) != NULL)
				return time_format(ws->songlist->length());
			else
				return time_format(-1);

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
	fields.push_back(new Field(FIELD_ARTISTSORT, "artistsort", "ArtistSort", "Artist", 0, 0));
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
	fields.push_back(new Field(FIELD_REMAINING, "remaining", "", "", 0, 0));
	fields.push_back(new Field(FIELD_VOLUME, "volume", "", "", 0, 0));
	fields.push_back(new Field(FIELD_MODES, "modes", "", "", 0, 0));
	fields.push_back(new Field(FIELD_STATE, "state", "", "", 0, 0));
	fields.push_back(new Field(FIELD_PROGRESSBAR, "progressbar", "", "", 0, 0));
	fields.push_back(new Field(FIELD_QUEUESIZE, "queuesize", "", "", 0, 0));
	fields.push_back(new Field(FIELD_QUEUELENGTH, "queuelength", "", "", 0, 0));
	fields.push_back(new Field(FIELD_LISTSIZE, "listsize", "", "", 0, 0));
	fields.push_back(new Field(FIELD_LISTLENGTH, "listlength", "", "", 0, 0));
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

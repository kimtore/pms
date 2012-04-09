/* vi:set ts=8 sts=8 sw=8 noet:
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

#include "window.h"
#include "songlist.h"
#include "curses.h"
#include "config.h"
#include "mpd.h"
#include "field.h"
#include "song.h"
#include <string>
#include <vector>

using namespace std;

extern Curses * curses;
extern Config * config;
extern MPD * mpd;
extern Windowmanager * wm;

void Wsonglist::draw()
{
	unsigned int x = 0, i, it;
	string wtitle;

	if (!rect || !visible())
		return;

	if (config->show_column_headers)
	{
		i = config->show_window_title ? 1 : 0;
		curses->clearline(rect, i, config->colors.columnheader);
		for (it = 0; it < column_len.size(); ++it)
		{
			curses->print(rect, config->colors.columnheader, i, x, config->songlist_columns[it]->title.c_str());
			x += column_len[it] + 1;
		}
	}

	if (config->show_window_title)
	{
		wtitle = title;
		if (songlist->searchresult)
			wtitle += "  <<search results>>";
		curses->clearline(rect, 0, config->colors.windowtitle);
		curses->print(rect, config->colors.windowtitle, 0, 0, wtitle.c_str());
	}

	Window::draw();
	wm->readout->draw();
}

void Wsonglist::drawline(int rely)
{
	unsigned int it;
	Song * song;
	Color * color;
	unsigned int linepos = rely + position;
	int x = 0;
	size_t vstart;
	size_t vstop;

	if (config->show_window_title)
		++rely;
	if (config->show_column_headers)
		++rely;

	if (!songlist || rely + rect->top > rect->bottom || linepos >= songlist->size())
	{
		curses->clearline(rect, rely, config->colors.standard);
		return;
	}

	songlist->visual_pos(&vstart, &vstop);

	song = songlist->at(linepos);

	if (linepos == cursor)
		color = config->colors.cursor;
	else if (linepos >= vstart && linepos <= vstop)
		color = config->colors.selection;
	else if (song->pos == mpd->status.song)
		color = config->colors.playing;
	else if (song->pos == -1 && mpd->currentsong && song->fhash == mpd->currentsong->fhash)
		color = config->colors.playing;
	else
		color = NULL;

	curses->clearline(rect, rely, color ? color : config->colors.standard);

	for (it = 0; it < column_len.size(); ++it)
	{
		curses->print(rect, color ? color : config->colors.field[config->songlist_columns[it]->type], rely, x, song->f[config->songlist_columns[it]->type].c_str());
		x += column_len[it] + 1;
	}
}

Song * Wsonglist::cursorsong()
{
	if (songlist->size() == 0)
		return NULL;
	
	move_cursor(0);
	return songlist->at(cursor);
}

selection_t Wsonglist::get_selection(long multiplier)
{
	vector<Song *> * sel;
	unsigned int s;
	Song * song;

	sel = songlist->get_selection();

	/* Append cursor song (plus multiplier) if no selection */
	if (sel->empty() && songlist->size())
	{
		s = cursor;
		while (--multiplier >= 0)
		{
			if ((song = songlist->at(s)) == NULL)
				break;
			sel->push_back(song);
			++s;
		}
	}

	return sel;
}

unsigned int Wsonglist::height()
{
	if (!rect) return 0;
	return Wmain::height() - (config->show_column_headers ? 1 : 0);
}

unsigned int Wsonglist::content_size()
{
	return songlist ? songlist->size() : 0;
}

void Wsonglist::move_cursor(int offset)
{
	Wmain::move_cursor(offset);
	if (songlist && songlist->visual_stop != -1)
		songlist->visual_stop = cursor;
}

void Wsonglist::update_column_length()
{
	vector<Field *>::iterator column;
	unsigned int it;
	unsigned int max;
	unsigned int len = 0;
	unsigned int oldlen = 0;

	column_len.clear();

	if (!rect || !songlist)
		return;

	max = rect->right - rect->left - config->songlist_columns.size() + 1;

	for (column = config->songlist_columns.begin(); column != config->songlist_columns.end(); ++column)
	{
		column_len.push_back((*column)->minlen);
		len += (*column)->minlen;
	}

	while (len <= max)
	{
		oldlen = len;
		for (it = 0; it < column_len.size(); ++it)
		{
			if (len > max)
				break;
			if (config->songlist_columns[it]->maxlen > 0 && column_len[it] >= config->songlist_columns[it]->maxlen)
				continue;

			column_len[it]++;
			++len;
		}

		/* break out of infinite loops if there are no expanding columns. */
		if (len == oldlen)
			break;
	}
}

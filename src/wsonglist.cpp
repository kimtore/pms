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

#include "window.h"
#include "songlist.h"
#include "curses.h"
#include "config.h"
#include "mpd.h"
#include "field.h"
#include <string>
#include <vector>

using namespace std;

extern Curses curses;
extern Config config;
extern MPD mpd;

void Wsonglist::drawline(int rely)
{
	vector<Field *>::iterator column;
	Song * song;
	Color * color;
	unsigned int linepos = rely + position;
	int x = 0;

	if (!songlist || rely + rect->top > rect->bottom || linepos >= songlist->size())
	{
		curses.clearline(rect, rely, config.colors.standard);
		return;
	}

	song = songlist->songs[linepos];
	if (linepos == cursor)
		color = config.colors.cursor;
	else if (song->pos == mpd.status.song)
		color = config.colors.playing;
	else
		color = NULL;

	curses.clearline(rect, rely, color ? color : config.colors.standard);

	for (column = config.songlist_columns.begin(); column != config.songlist_columns.end(); ++column)
	{
		curses.print(rect, color ? color : config.colors.field[(*column)->type], rely, x, song->f[(*column)->type].c_str());
		x += 30;
	}
}

unsigned int Wsonglist::content_size()
{
	return songlist ? songlist->size() : 0;
}

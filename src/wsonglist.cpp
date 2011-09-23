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
#include <string>
#include <vector>

using namespace std;

extern Curses curses;
extern Config config;
extern MPD mpd;

void Wsonglist::drawline(int rely)
{
	Song * song;
	Color * color;
	unsigned int linepos = rely + position;

	curses.clearline(rect, rely);
	if (!songlist || rely + rect->top > rect->bottom || linepos >= songlist->size())
		return;

	song = songlist->songs[linepos];
	if (linepos == cursor)
		color = config.colors.cursor;
	else if (song->pos == mpd.state.song)
		color = config.colors.playing;
	else
		color = config.colors.standard;

	curses.print(rect, color, rely, 0, song->file.c_str());
}

unsigned int Wsonglist::content_size()
{
	return songlist ? songlist->size() : 0;
}

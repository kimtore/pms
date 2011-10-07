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

#include "color.h"
#include "curses.h"
#include "field.h"

short Color::color_count = 0;

Colortable::Colortable()
{
	pair_content(-1, &dfront, &dback);

	standard = new Color(dfront, dback, 0);
	topbar = new Color(COLOR_WHITE, -1, 0);
	statusbar = new Color(COLOR_WHITE, -1, 0);
	windowtitle = new Color(COLOR_CYAN, -1, A_BOLD);
	columnheader = new Color(COLOR_WHITE, -1, 0);
	console = new Color(COLOR_WHITE, -1, 0);
	error = new Color(COLOR_WHITE, COLOR_RED, A_BOLD);
	readout = new Color(COLOR_WHITE, -1, 0);

	cursor = new Color(COLOR_BLACK, COLOR_WHITE, 0);
	playing = new Color(COLOR_BLACK, COLOR_YELLOW, 0);

	field[FIELD_DIRECTORY] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_FILE] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_POS] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_ID] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_TIME] = new Color(COLOR_MAGENTA, -1, 0);
	field[FIELD_NAME] = new Color(COLOR_WHITE, -1, A_BOLD);
	field[FIELD_ARTIST] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_ARTISTSORT] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_ALBUM] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_TITLE] = new Color(COLOR_WHITE, -1, A_BOLD);
	field[FIELD_TRACK] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_DATE] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_DISC] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_GENRE] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_ALBUMARTIST] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_ALBUMARTISTSORT] = new Color(COLOR_YELLOW, -1, 0);

	field[FIELD_YEAR] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_TRACKSHORT] = new Color(COLOR_CYAN, -1, 0);

	field[FIELD_ELAPSED] = new Color(COLOR_GREEN, -1, 0);
	field[FIELD_REMAINING] = new Color(COLOR_MAGENTA, -1, 0);
	field[FIELD_MODES] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_STATE] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_QUEUESIZE] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_QUEUELENGTH] = new Color(COLOR_WHITE, -1, 0);
}

Colortable::~Colortable()
{
}

Color::Color(short nfront, short nback, int nattr)
{
	id = Color::color_count;
	set(nfront, nback, nattr);
	Color::color_count++;
}

void Color::set(short nfront, short nback, int nattr)
{
	front = nfront;
	back = nback;
	attr = nattr;
	init_pair(id, front, back);
	pair = COLOR_PAIR(id) | attr;
}

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

short Color::color_count = 0;

Colortable::Colortable()
{
	pair_content(-1, &dfront, &dback);

	standard = new Color(dfront, dback, 0);
	statusbar = new Color(COLOR_WHITE, -1, 0);
	console = new Color(COLOR_WHITE, -1, 0);
	error = new Color(COLOR_WHITE, COLOR_RED, A_BOLD);
	readout = new Color(COLOR_WHITE, -1, 0);

	cursor = new Color(COLOR_BLACK, COLOR_WHITE, 0);
	playing = new Color(COLOR_BLACK, COLOR_YELLOW, 0);
}

Colortable::~Colortable()
{
	delete standard;
	delete statusbar;
	delete console;
	delete error;
	delete readout;

	delete cursor;
	delete playing;
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

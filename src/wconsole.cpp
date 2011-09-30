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
#include "console.h"
#include "curses.h"
#include "config.h"
#include <string>
#include <vector>

using namespace std;

extern vector<Logline *> logbuffer;
extern Curses curses;
extern Config config;

void Wconsole::drawline(int rely)
{
	unsigned int linepos = rely + position;

	curses.clearline(rect, rely);
	if (rely + rect->top > rect->bottom || linepos >= logbuffer.size())
		return;

	curses.print(rect, logbuffer[linepos]->level == MSG_LEVEL_ERR ? config.colors.error : config.colors.console, rely, 0, logbuffer[linepos]->line.c_str());
}

unsigned int Wconsole::content_size()
{
	return logbuffer.size();
}

void Wconsole::move_cursor(int offset)
{
	return scroll_window(offset);
}

void Wconsole::set_cursor(unsigned int absolute)
{
	return set_position(absolute);
}

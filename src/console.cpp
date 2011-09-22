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

#include "console.h"
#include "window.h"
#include "curses.h"
#include <stdarg.h>
#include <stdio.h>
#include <string>
#include <vector>

using namespace std;

vector<Logline *> logbuffer;
extern Curses curses;
extern Windowmanager wm;

Logline::Logline(int lvl, const char * ln)
{
	level = lvl;
	line = ln;
}

void console_log(int level, const char * format, ...)
{
	va_list		ap;
	char		buffer[1024];

	va_start(ap, format);
	vsprintf(buffer, format, ap);
	va_end(ap);

	logbuffer.push_back(new Logline(level, buffer));
	if (level <= MSG_LEVEL_INFO || wm.console->visible())
	{
		wm.statusbar->draw();
		wm.console->scroll_window(1);
		curses.flush();
	}
}

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

#include <ncurses.h>

bool init_curses()
{
	if ((initscr()) == NULL)
		return false;

	raw();
	noecho();
	halfdelay(5);
	keypad(stdscr, true);
	curs_set(0);

	if (has_colors())
	{
		start_color();
		use_default_colors();
	}

	clear();
	refresh();

	return true;
}

void shutdown_curses()
{
	endwin();
}

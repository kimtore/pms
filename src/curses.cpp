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

#include "curses.h"
#include <cstring>

Curses::Curses()
{
	if ((initscr()) == NULL)
	{
		ready = false;
		return;
	}

	raw();
	noecho();
	halfdelay(10);
	keypad(stdscr, true);
	curs_set(0);

	if (has_colors() && start_color())
	{
		use_default_colors();
		hascolors = true;
	}

	detect_dimensions();
	clear();
	refresh();

	ready = true;
}

Curses::~Curses()
{
	endwin();
}

void Curses::detect_dimensions()
{
	memset(&self, 0, sizeof self);
	memset(&topbar, 0, sizeof topbar);
	memset(&main, 0, sizeof main);
	memset(&statusbar, 0, sizeof statusbar);

	self.right = COLS - 1;
	self.bottom = LINES - 1;

	topbar.top = 1;
	topbar.bottom = topbar.top;
	topbar.right = self.right;

	main.top = topbar.bottom + 1;
	main.bottom = self.bottom - 1;
	main.right = self.right;

	statusbar.top = self.bottom;
	statusbar.bottom = self.bottom;
	statusbar.right = self.right;
}

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
#include "curses.h"
#include "command.h"
#include <vector>

extern Curses curses;

Windowmanager::Windowmanager()
{
	/* Setup static windows that are not in the window list */
	topbar = new Wtopbar;
	topbar->set_rect(&curses.topbar);
	statusbar = new Wstatusbar;
	statusbar->set_rect(&curses.statusbar);
	readout = new Wreadout;
	readout->set_rect(&curses.readout);

	/* Setup static windows that appear in the window list */
	console = new Wconsole;
	console->set_rect(&curses.main);
	playlist = new Wsonglist;
	playlist->set_rect(&curses.main);
	library = new Wsonglist;
	library->set_rect(&curses.main);
	windows.push_back(WWINDOW(console));
	windows.push_back(WWINDOW(playlist));
	windows.push_back(WWINDOW(library));

	/* Activate playlist window */
	activate(WWINDOW(console));
	context = CONTEXT_CONSOLE;
}

void Windowmanager::draw()
{
	topbar->draw();
	statusbar->draw();
	readout->draw();
	active->draw();
}

bool Windowmanager::activate(Window * nactive)
{
	unsigned int i;

	for (i = 0; i < windows.size(); ++i)
	{
		if (windows[i] == nactive)
		{
			active_index = i;
			active = nactive;
			active->clear();
			active->draw();
			curses.flush();
			return true;
		}
	}

	return false;
}

void Windowmanager::cycle(int offset)
{
	if (offset >= 0)
		offset %= windows.size();
	else
		offset %= -windows.size();

	offset = active_index + offset;
	if (offset < 0)
		offset = windows.size() - offset;
	else if (offset >= (int)windows.size())
		offset -= windows.size();

	active_index = (unsigned int)offset;
	active = windows[active_index];
	active->draw();
	curses.flush();
}

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
#include "mpd.h"
#include "config.h"
#include <vector>

extern Curses curses;
extern MPD mpd;
extern Config config;

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
	console->title = "Console";
	playlist = new Wsonglist;
	playlist->set_rect(&curses.main);
	playlist->title = "Playlist";
	library = new Wsonglist;
	library->set_rect(&curses.main);
	library->title = "Library";
	windows.push_back(WMAIN(console));
	windows.push_back(WMAIN(playlist));
	windows.push_back(WMAIN(library));

	/* Activate playlist window */
	activate(WMAIN(console));
}

void Windowmanager::draw()
{
	topbar->draw();
	statusbar->draw();
	readout->draw();
	active->draw();
}

void Windowmanager::flush()
{
	curses.flush();
}

bool Windowmanager::activate(Wmain * nactive)
{
	Wsonglist * ws;
	unsigned int i;

	for (i = 0; i < windows.size(); ++i)
	{
		if (windows[i] == nactive)
		{
			if ((ws = WSONGLIST(nactive)) != NULL && config.playback_follows_window)
				mpd.activate_songlist(ws->songlist);

			active_index = i;
			active = nactive;
			context = active->context;
			active->clear();
			active->draw();
			readout->draw();
			curses.flush();
			return true;
		}
	}

	return false;
}

void Windowmanager::cycle(int offset)
{
	Wsonglist * ws;

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
	context = active->context;

	if ((ws = WSONGLIST(active)) != NULL && config.playback_follows_window)
		mpd.activate_songlist(ws->songlist);

	active->draw();
	readout->draw();
	curses.flush();
}

void Windowmanager::update_column_length()
{
	Wsonglist * w;
	vector<Wmain *>::iterator i;
	for (i = windows.begin(); i != windows.end(); ++i)
	{
		if ((w = WSONGLIST(*i)) != NULL)
			w->update_column_length();
	}
}

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

#include "build.h"
#include "console.h"
#include "curses.h"
#include "config.h"
#include "window.h"
#include "mpd.h"
#include "input.h"
#include "pms.h"
#include "field.h"
#include <glib.h>
#include <stdio.h>
#include <stdlib.h>
#include <cstring>
#include <sys/time.h>

Fieldtypes	fieldtypes;
Curses		curses;
Config		config;
MPD		mpd;
Windowmanager	wm;
Input		input;
PMS		pms;

int main(int argc, char *argv[])
{
	struct timeval cl;
	struct timeval conn;

	if (!curses.ready)
	{
		perror("Fatal: failed to initialise ncurses.\n");
		return 1;
	}

	memset(&conn, 0, sizeof conn);

	curses.detect_dimensions();
	wm.playlist->songlist = &mpd.playlist;
	wm.library->songlist = &mpd.library;
	wm.draw();
	stinfo("%s %d.%d", PMS_APP_NAME, PMS_VERSION_MAJOR, PMS_VERSION_MINOR);

	while(!config.quit)
	{
		gettimeofday(&cl, NULL);
		if (!mpd.is_connected())
		{
			if (cl.tv_sec - conn.tv_sec >= (int)config.reconnect_delay)
			{
				mpd.mpd_connect(config.host, config.port);
				mpd.set_password(config.password);
				mpd.get_status();
				mpd.get_playlist();
				mpd.get_library();
				mpd.read_opts();
				mpd.update_playstring();
				wm.draw();
				curses.flush();
				gettimeofday(&conn, NULL);
			}
		}

		/* Get updates from MPD, run clock, do updates */
		mpd.poll();

		/* Draw topbar after poll() */
		wm.topbar->draw();

		/* Statusbar needs redraw with playstring? */
		wm.statusbar->draw();

		/* Refresh screen every tick */
		curses.flush();

		/* Check for any input events and run them */
		pms.run_event(input.next());
	}
}

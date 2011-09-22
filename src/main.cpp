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
#include <glib.h>
#include <stdio.h>

Config		config;
MPD		mpd;
Curses		curses;
Windowmanager	wm;
Input		input;
PMS		pms;

int main(int argc, char *argv[])
{
	if (!curses.ready)
	{
		perror("Fatal: failed to initialise ncurses.\n");
		return 1;
	}

	wm.draw();
	stinfo("%s %d.%d", PMS_APP_NAME, PMS_VERSION_MAJOR, PMS_VERSION_MINOR);

	while(!config.quit)
	{
		if (!mpd.is_connected())
		{
			mpd.mpd_connect(config.host, config.port);
			mpd.set_password(config.password);
			mpd.get_status();
		}
		mpd.poll();
		pms.run_event(input.next());
	}
}

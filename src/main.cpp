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
#include "debug.h"
#include "curses.h"
#include "config.h"
#include "mpd.h"
#include <glib.h>
#include <stdio.h>

Config		config;
MPD		mpd;

int main(int argc, char *argv[])
{
	printf("%s %d.%d\n", PMS_APP_NAME, PMS_VERSION_MAJOR, PMS_VERSION_MINOR);
	mpd.mpd_connect(config.host, config.port);
	if (!init_curses())
	{
		perror("Fatal: failed to initialise ncurses.\n");
		return 1;
	}
	while(true);
	shutdown_curses();
}

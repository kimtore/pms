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

#ifndef _PMS_PMS_H_
#define _PMS_PMS_H_

#include "input.h"
#include "song.h"

/*
 * This class contains all user interface actions,
 * and could probably be used for plugin interfacing
 * in the future.
 */
class PMS
{
	public:
		/* This function handles input events from main(). */
		int		run_event(Inputevent * ev);

		/* Set options */
		int		set_opt(Inputevent * ev);

		/* Run a text command */
		int		run_cmd(string cmd, unsigned int multiplier = 1, bool batch = false);

		/* Run a search */
		int		run_search(string terms, unsigned int multiplier = 1);

		/* Quit the program. */
		int		quit();

		/* Move cursor in current window */
		int		move_cursor(int offset);

		/* Move cursor N pages in current window */
		int		move_cursor_page(int offset);

		/* Scroll the current window */
		int		scroll_window(int offset);

		/* Set cursor absolute position */
		int		set_cursor_home();
		int		set_cursor_end();
		int		set_cursor_top();
		int		set_cursor_bottom();
		int		set_cursor_currentsong();
		int		set_cursor_random();

		/* Change windows */
		int		cycle_windows(int offset);
		int		activate_songlist();

		/* List management */
		int		add(int count);
		int		remove(int count);

		/* MPD options */
		int		update(string dir = "/");
		int		set_crossfade(string secs);
		int		change_volume(int offset);
		int		set_volume(string volume);
		int		set_password(string password);
		
		/* Playback */
		int		toggle_play();
		int		play();
		int		stop();
		int		change_song(int steps);
		int		seek(int seconds);
};

/* Return the song pointed to by the cursor, otherwise NULL */
Song * cursorsong();

#endif /* _PMS_PMS_H_ */

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

		/* Run a text command */
		int		run_cmd(string cmd);

		/* Quit the program. */
		int		quit();

		/* Move cursor in current window*/
		int		move_cursor(int offset);

		/* Scroll the current window */
		int		scroll_window(int offset);

		/* Set scroll absolute position */
		int		set_cursor_home();
		int		set_cursor_end();
		int		set_cursor_top();
		int		set_cursor_bottom();

		/* Change windows */
		int		cycle_windows(int offset);

		/* MPD options */
		int		set_consume(bool consume);
		int		set_crossfade(string secs);
		int		set_random(bool random);
		int		set_repeat(bool repeat);
		int		set_volume(string volume);
		int		set_single(bool single);
		
		/* Playback */
		int		toggle_play();
};


#endif /* _PMS_PMS_H_ */

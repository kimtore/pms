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

#include "pms.h"
#include "console.h"
#include "curses.h"
#include "config.h"
#include "window.h"
#include "mpd.h"
#include "input.h"

extern Config		config;
extern MPD		mpd;
extern Curses		curses;
extern Windowmanager	wm;
extern Input		input;

int PMS::run_event(input_event * ev)
{
	if (!ev) return false;

	if (ev->result == INPUT_RESULT_RUN)
	switch(ev->action)
	{
		case ACT_MODE_INPUT:
			input.setmode(INPUT_MODE_INPUT);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_MODE_COMMAND:
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_RUN_CMD:
			/* TODO: add a command parser */
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_QUIT:
			return quit();

		case ACT_RESIZE:
			curses.detect_dimensions();
			wm.draw();
			curses.flush();
			return true;

		case ACT_SCROLL_UP:
			return scroll_window(-ev->multiplier);

		case ACT_SCROLL_DOWN:
			return scroll_window(ev->multiplier);

		case ACT_CURSOR_UP:
			return move_cursor(-ev->multiplier);

		case ACT_CURSOR_DOWN:
			return move_cursor(ev->multiplier);

		case ACT_CURSOR_HOME:
			return set_cursor_home();

		case ACT_CURSOR_END:
			return set_cursor_end();

		default:
			return false;
	}

	else if (ev->result == INPUT_RESULT_BUFFERED)
	{
		wm.statusbar->draw();
		curses.flush();
	}


	return false;
}

int PMS::quit()
{
	config.quit = true;
	return true;
}

int PMS::scroll_window(int offset)
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->scroll_window(offset);
	return true;
}

int PMS::move_cursor(int offset)
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->move_cursor(offset);
	return true;
}

int PMS::set_cursor_home()
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->set_cursor(0);
	return true;
}

int PMS::set_cursor_end()
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->set_cursor(-1);
	return true;
}

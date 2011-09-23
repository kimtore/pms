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

#include "input.h"
#include "curses.h"
#include "command.h"
#include "window.h"
#include <cstring>

Keybindings * keybindings;
extern Windowmanager wm;

Input::Input()
{
	mode = INPUT_MODE_COMMAND;
	chbuf = 0;
	multiplier = 1;
	strbuf.clear();
	buffer.clear();
	keybindings = new Keybindings();
}

input_event * Input::next()
{
	int m;

	if ((chbuf = getch()) == -1)
		return NULL;

	memset(&ev, 0, sizeof ev);

	/* This is a global signal/event non-dependant on anything else. */
	if (chbuf == KEY_RESIZE)
	{
		ev.result = INPUT_RESULT_RUN;
		ev.action = ACT_RESIZE;
		return &ev;
	}

	switch(mode)
	{
		/* Command-mode */
		default:
		case INPUT_MODE_COMMAND:
			buffer.push_back(chbuf);
			strbuf.push_back(chbuf);
			m = keybindings->find(wm.context, &buffer, &ev.action);

			if (m == KEYBIND_FIND_EXACT)
				ev.result = INPUT_RESULT_RUN;
			else if (m == KEYBIND_FIND_NOMATCH)
				buffer.clear();

			break;

		/* Text input of some sorts */
		case INPUT_MODE_INPUT:
		case INPUT_MODE_SEARCH:
			switch(chbuf)
			{
				case 10:
				case KEY_ENTER:
					if (mode == INPUT_MODE_INPUT)
						ev.action = ACT_RUN_CMD;

					ev.result = INPUT_RESULT_RUN;
					break;

				case 21:		/* ^U */
					buffer.clear();
					strbuf.clear();
					ev.result = INPUT_RESULT_BUFFERED;
					break;

				case 27:		/* Escape */
					ev.result = INPUT_RESULT_RUN;
					ev.action = ACT_MODE_COMMAND;
					break;

				case 8:			/* ^H -- backspace */
				case 127:		/* ^? -- delete */
				case KEY_BACKSPACE:
					if (buffer.size() > 0)
					{
						buffer.pop_back();
						strbuf.erase(--string::iterator(strbuf.end()));
						ev.result = INPUT_RESULT_BUFFERED;
					}
					else
					{
						ev.result = INPUT_RESULT_RUN;
						ev.action = ACT_MODE_COMMAND;
					}
					break;

				default:
					buffer.push_back(chbuf);
					strbuf.push_back(chbuf);
					ev.result = INPUT_RESULT_BUFFERED;
			}
			break;
	}

	if (ev.result != INPUT_RESULT_NOINPUT)
	{
		if (ev.result != INPUT_RESULT_BUFFERED)
		{
			buffer.clear();
			strbuf.clear();
		}

		ev.multiplier = multiplier;
		return &ev;
	}
	
	return NULL;
}

void Input::setmode(int nmode)
{
	if (nmode == mode)
		return;
	
	strbuf.clear();
	buffer.clear();
	chbuf = 0;
	multiplier = 1;
	mode = nmode;

	if (mode == INPUT_MODE_COMMAND)
		curs_set(0);
	else
		curs_set(1);
}

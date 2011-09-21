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

Keybindings keybindings;

Input::Input()
{
	mode = INPUT_MODE_COMMAND;
	chbuf = 0;
	multiplier = 1;
	buffer = "";
}

int Input::next()
{
	action_t action;
	int m;

	if ((chbuf = getch()) == INPUT_NOINPUT)
		return chbuf;

	switch(mode)
	{
		default:
		case INPUT_MODE_COMMAND:
			buffer += chbuf;
			m = keybindings.find(CONTEXT_ALL, buffer, &action);

			if (m == KEYBIND_FIND_EXACT)
				return INPUT_RUN;
			else if (m == KEYBIND_FIND_NOMATCH)
				buffer.clear();

			return INPUT_NOINPUT;
		case INPUT_MODE_INPUT:
		case INPUT_MODE_SEARCH:
			buffer += chbuf;
			return INPUT_BUFFERED;
	}
}

void Input::setmode(int nmode)
{
	if (nmode == mode)
		return;
	
	buffer = "";
	chbuf = 0;
	multiplier = 1;
	mode = nmode;
}

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
Commandlist commandlist;
extern Windowmanager wm;

Inputevent::Inputevent()
{
	clear();
}

void Inputevent::clear()
{

	context = 0;
	multiplier = 0;
	action = ACT_NOACTION;
	context = 0;
	text.clear();
}

Input::Input()
{
	mode = INPUT_MODE_COMMAND;
	chbuf = 0;
	multiplier = 1;
	strbuf.clear();
	buffer.clear();
	is_tab_completing = false;
	tab_complete_index = 0;
	keybindings = new Keybindings();
}

Inputevent * Input::next()
{
	int m;

	if ((chbuf = getch()) == -1)
		return NULL;

	ev.clear();

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
			else if (m == KEYBIND_FIND_BUFFERED)
				ev.result = INPUT_RESULT_BUFFERED;
			else if (m == KEYBIND_FIND_NOMATCH)
			{
				buffer.clear();
				strbuf.clear();
			}

			break;

		/* Text input of some sorts */
		case INPUT_MODE_INPUT:
		case INPUT_MODE_SEARCH:
			handle_text_input();
			break;
	}

	if (ev.result != INPUT_RESULT_NOINPUT)
	{
		ev.context = wm.context;
		ev.text = strbuf;
		ev.multiplier = multiplier;

		if (ev.result != INPUT_RESULT_BUFFERED)
		{
			buffer.clear();
			strbuf.clear();
		}

		return &ev;
	}
	
	return NULL;
}

void Input::handle_text_input()
{
	string::iterator si;

	if (chbuf != 9)
		is_tab_completing = false;

	switch(chbuf)
	{
		case 10:
		case KEY_ENTER:
			if (mode == INPUT_MODE_INPUT)
				ev.action = ACT_RUN_CMD;

			ev.result = INPUT_RESULT_RUN;
			return;

		case 21:		/* ^U */
			buffer.clear();
			strbuf.clear();
			ev.result = INPUT_RESULT_BUFFERED;
			return;

		case 27:		/* Escape */
			ev.result = INPUT_RESULT_RUN;
			ev.action = ACT_MODE_COMMAND;
			return;

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
			return;

		case 9:			/* TODO: Tab-completion */
			if (is_tab_completing)
			{
				if (++tab_complete_index >= tab_results->size())
					tab_complete_index = 0;
			}
			else
			{
				tab_results = commandlist.grep(wm.context, strbuf);
				if (tab_results->size() > 0)
				{
					is_tab_completing = true;
					tab_complete_index = 0;
				}
			}

			if (is_tab_completing)
			{
				strbuf = tab_results->at(tab_complete_index)->name;
				buffer.clear();
				for (si = strbuf.begin(); si != strbuf.end(); ++si)
					buffer.push_back(*si);
			}
			ev.result = INPUT_RESULT_BUFFERED;
			return;

		default:
			buffer.push_back(chbuf);
			strbuf.push_back(chbuf);
			ev.result = INPUT_RESULT_BUFFERED;
	}
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

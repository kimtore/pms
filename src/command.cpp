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

#include "command.h"
#include <string>
#include <vector>

using namespace std;

Commandlist::Commandlist()
{
	add(CONTEXT_ALL, ACT_SET, "se");
	add(CONTEXT_ALL, ACT_SET, "set");
	add(CONTEXT_ALL, ACT_MODE_INPUT, "cmd");
	add(CONTEXT_SONGLIST, ACT_MODE_SEARCH, "search");
	add(CONTEXT_ALL, ACT_REPEATACTION, "repeat-action");

	add(CONTEXT_ALL, ACT_SOURCE, "source");
	add(CONTEXT_ALL, ACT_QUIT, "quit");
	add(CONTEXT_ALL, ACT_QUIT, "q");
	add(CONTEXT_ALL, ACT_NEXT_WINDOW, "next-window");
	add(CONTEXT_ALL, ACT_PREVIOUS_WINDOW, "prev-window");
	add(CONTEXT_SONGLIST, ACT_ACTIVATE_SONGLIST, "activate-list");
	add(CONTEXT_SONGLIST, ACT_ADD, "add");
	add(CONTEXT_SONGLIST, ACT_REMOVE, "remove");

	add(CONTEXT_LIST, ACT_SCROLL_UP, "scroll-up");
	add(CONTEXT_LIST, ACT_SCROLL_DOWN, "scroll-down");
	add(CONTEXT_LIST, ACT_CURSOR_UP, "cursor-up");
	add(CONTEXT_LIST, ACT_CURSOR_DOWN, "cursor-down");
	add(CONTEXT_LIST, ACT_CURSOR_PGUP, "cursor-pageup");
	add(CONTEXT_LIST, ACT_CURSOR_PGDOWN, "cursor-pagedown");
	add(CONTEXT_LIST, ACT_CURSOR_HOME, "cursor-home");
	add(CONTEXT_LIST, ACT_CURSOR_END, "cursor-end");
	add(CONTEXT_LIST, ACT_CURSOR_TOP, "cursor-top");
	add(CONTEXT_LIST, ACT_CURSOR_BOTTOM, "cursor-bottom");
	add(CONTEXT_SONGLIST, ACT_CURSOR_CURRENTSONG, "cursor-currentsong");
	add(CONTEXT_SONGLIST, ACT_CURSOR_RANDOM, "cursor-random");

	add(CONTEXT_ALL, ACT_UPDATE, "update");
	add(CONTEXT_ALL, ACT_CROSSFADE, "crossfade");
	add(CONTEXT_ALL, ACT_PASSWORD, "password");

	add(CONTEXT_ALL, ACT_TOGGLEPLAY, "toggle-play");
	add(CONTEXT_SONGLIST, ACT_PLAY, "play");
	add(CONTEXT_ALL, ACT_STOP, "stop");
	add(CONTEXT_ALL, ACT_NEXT, "next");
	add(CONTEXT_ALL, ACT_PREVIOUS, "previous");
	add(CONTEXT_ALL, ACT_SEEK_FORWARD, "seek-forward");
	add(CONTEXT_ALL, ACT_SEEK_BACK, "seek-back");
}

Command * Commandlist::add(int context, action_t action, string name)
{
	Command * c = new Command;
	c->context = context;
	c->action = action;
	c->name = name;
	cmds.push_back(c);
	return c;
}

vector<Command *> * Commandlist::grep(int context, string name)
{
	vector<Command *>::iterator i;

	grepcmds.clear();

	for (i = cmds.begin(); i != cmds.end(); i++)
	{
		if (name.size() > (*i)->name.size() || !(context & (*i)->context))
			continue;

		if (name == (*i)->name.substr(0, name.size()))
			grepcmds.push_back(*i);
	}

	return &grepcmds;
}

Command * Commandlist::find(int context, string name)
{
	vector<Command *>::iterator i;

	for (i = cmds.begin(); i != cmds.end(); i++)
	{
		if (context & (*i)->context && name == (*i)->name)
			return *i;
	}

	return NULL;
}

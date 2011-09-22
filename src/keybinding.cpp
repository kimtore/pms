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
#include "command.h"
#include <string>
#include <vector>

using namespace std;

Keybindings::Keybindings()
{
	add(CONTEXT_ALL, ACT_QUIT, "q");
	add(CONTEXT_LIST, ACT_SCROLL_UP, "k");
	add(CONTEXT_LIST, ACT_SCROLL_DOWN, "j");
}

Keybinding * Keybindings::add(int context, action_t action, string sequence)
{
	Keybinding * c;

	if (find_conflict(sequence))
		return NULL;

	c = new Keybinding();
	c->context = context;
	c->action = action;
	c->sequence = sequence;
	bindings.push_back(c);
	return c;
}

Keybinding * Keybindings::find_conflict(string sequence)
{
	vector<Keybinding *>::iterator i;

	for (i = bindings.begin(); i != bindings.end(); i++)
	{
		if (sequence == (*i)->sequence)
		{
			return *i;
		}
		else if (sequence.size() < (*i)->sequence.size())
		{
			if (sequence == (*i)->sequence.substr(0, sequence.size() - 1))
				return *i;
		}
		else if (sequence.size() > (*i)->sequence.size())
		{
			if (sequence.substr(0, (*i)->sequence.size()) == (*i)->sequence)
				return *i;
		}
	}

	return NULL;
}

int Keybindings::find(int context, string sequence, action_t * action)
{
	vector<Keybinding *>::iterator i;
	int found = KEYBIND_FIND_NOMATCH;

	for (i = bindings.begin(); i != bindings.end(); i++)
	{
		if (!((*i)->context & context) || (*i)->sequence.size() < sequence.size())
			continue;

		if ((*i)->sequence == sequence)
		{
			*action = (*i)->action;
			return KEYBIND_FIND_EXACT;
		}

		if ((*i)->sequence.size() > sequence.size() && sequence == (*i)->sequence.substr(0, sequence.size()))
			found = KEYBIND_FIND_BUFFERED;
	}

	return found;
}

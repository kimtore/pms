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
#include "console.h"
#include "curses.h"
#include <string>
#include <algorithm>
#include <vector>

using namespace std;

Keybindings::Keybindings()
{
	add(CONTEXT_ALL, ACT_QUIT, "q");
	add(CONTEXT_ALL, ACT_QUIT, "<C-d<>");
	add(CONTEXT_LIST, ACT_CURSOR_UP, "k");
	add(CONTEXT_LIST, ACT_CURSOR_DOWN, "j");
	add(CONTEXT_LIST, ACT_CURSOR_HOME, "gg");
	add(CONTEXT_LIST, ACT_CURSOR_END, "G");
}

Keybinding * Keybindings::add(int context, action_t action, string sequence)
{
	Keybinding * c;

	sequence = conv_sequence(sequence);
	if (sequence.size() == 0 || find_conflict(sequence))
		return NULL;

	c = new Keybinding();
	c->context = context;
	c->action = action;
	c->sequence = sequence;
	bindings.push_back(c);
	return c;
}

string Keybindings::conv_sequence(string seq)
{
	string r = "";
	string subseq;
	size_t i;
	size_t epos;
	int escape = false;
	int ex = false;
	int ctrl = false;

	for (i = 0; i < seq.size(); i++)
	{
		/* Escape sequence */
		if (seq[i] == '\\')
		{
			if (escape)
				r += seq[i];
			escape = !escape;
			continue;
		}
		/* Start a special key */
		else if (seq[i] == '<')
		{
			if (escape)
			{
				escape = false;
				r += seq[i];
				continue;
			}
			if (ex)
			{
				sterr("Bind: unexpected '%c' near ...%s", seq[i], seq.substr(i - 1).c_str());
				return "";
			}
			ex = true;
		}
		/* End a special key */
		else if (seq[i] == '>')
		{
			if (escape)
			{
				escape = false;
				r += seq[i];
				continue;
			}
			if (!ex)
			{
				sterr("Bind: unexpected '%c' near ...%s", seq[i], seq.substr(i).c_str());
				return "";
			}
		}
		else
		{
			/* Just a normal character */
			if (!ex)
			{
				r += seq[i];
				continue;
			}
			else if (ctrl)
			{
				seq[i] = ::toupper(seq[i]);
				if (!(seq[i] >= 'A' || seq[i] <= 'Z'))
				{
					sterr("Bind: unexpected %c, expected one letter between A-Z near ...%s", seq[i], seq.substr(i - 3).c_str());
					return "";
				}
				r += (seq[i] - 64);
				ctrl = false;
				continue;
			}

			/* Switch on Ctrl-mode */
			if (seq.substr(i, 2) == "C-")
			{
				ctrl = true;
				++i;
				continue;
			}

			epos = seq.find('>', i);
			if (epos == string::npos)
			{
				sterr("Bind: unclosed tag near ...%s", seq[i], seq.substr(i).c_str());
				return "";
			}

			subseq = seq.substr(i, epos - i);
			std::transform(subseq.begin(), subseq.end(), subseq.begin(), ::tolower);
			if (subseq == "cr")
				r += 13;

			else
			{
				sterr("Bind: invalid identifier '%s'", subseq.c_str());
				return "";
			}

			sterr("Bind: unexpected '%c' near ...%s", seq[i], seq.substr(i - 1).c_str());
			return "";
		}
	}

	return r;
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

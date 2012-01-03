/* vi:set ts=8 sts=8 sw=8 noet:
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

void Keybindings::load_defaults()
{
	truncate();

	add(CONTEXT_ALL, ACT_MODE_INPUT, ":");
	add(CONTEXT_SONGLIST, ACT_MODE_LIVESEARCH, "/");
	add(CONTEXT_SONGLIST, ACT_MODE_SEARCH, "<C-f>");
	add(CONTEXT_SONGLIST, ACT_RESET_SEARCH, "<C-g>");
	add(CONTEXT_ALL, ACT_REPEATACTION, ".");
	add(CONTEXT_ALL, ACT_QUIT, "q");
	add(CONTEXT_ALL, ACT_UPDATE, "U");

	add(CONTEXT_ALL, ACT_NEXT_WINDOW, "l");
	add(CONTEXT_ALL, ACT_NEXT_WINDOW, "gt");
	add(CONTEXT_ALL, ACT_PREVIOUS_WINDOW, "h");
	add(CONTEXT_ALL, ACT_PREVIOUS_WINDOW, "gT");
	add(CONTEXT_ALL, ACT_TOGGLE_WINDOW, "<Tab>");
	add(CONTEXT_ALL, ACT_GOTO_WINDOW_POS, "B");

	add(CONTEXT_SONGLIST, ACT_ADD, "a");
	add(CONTEXT_SONGLIST, ACT_ADD_SAME, "b", "album");
	add(CONTEXT_SONGLIST, ACT_REMOVE, "d");
	add(CONTEXT_SONGLIST, ACT_REMOVE, "<Delete>");
	add(CONTEXT_SONGLIST, ACT_REMOVE, "x");
	add(CONTEXT_SONGLIST, ACT_VISUAL, "v");
	add(CONTEXT_SONGLIST, ACT_VISUAL, "V");
	add(CONTEXT_SONGLIST, ACT_YANK, "y");
	add(CONTEXT_SONGLIST, ACT_YANK, "Y");
	add(CONTEXT_SONGLIST, ACT_PUT, "p");

	add(CONTEXT_LIST, ACT_SCROLL_UP, "<C-y>");
	add(CONTEXT_LIST, ACT_SCROLL_DOWN, "<C-e>");
	add(CONTEXT_LIST, ACT_CURSOR_UP, "k");
	add(CONTEXT_LIST, ACT_CURSOR_UP, "<Up>");
	add(CONTEXT_LIST, ACT_CURSOR_PGUP, "<PgUp>");
	add(CONTEXT_LIST, ACT_CURSOR_TOP, "H");
	add(CONTEXT_LIST, ACT_CURSOR_BOTTOM, "L");
	add(CONTEXT_LIST, ACT_CURSOR_DOWN, "j");
	add(CONTEXT_LIST, ACT_CURSOR_DOWN, "<Down>");
	add(CONTEXT_LIST, ACT_CURSOR_PGDOWN, "<PgDn>");
	add(CONTEXT_LIST, ACT_CURSOR_HOME, "gg");
	add(CONTEXT_LIST, ACT_CURSOR_HOME, "<Home>");
	add(CONTEXT_LIST, ACT_CURSOR_END, "G");
	add(CONTEXT_LIST, ACT_CURSOR_END, "<End>");
	add(CONTEXT_LIST, ACT_CURSOR_CURRENTSONG, "C");
	add(CONTEXT_LIST, ACT_CURSOR_CURRENTSONG, "gc");
	add(CONTEXT_LIST, ACT_CURSOR_RANDOM, "R");
	add(CONTEXT_LIST, ACT_CURSOR_RANDOM, "gr");

	add(CONTEXT_ALL, ACT_SET, "c", "consume!");
	add(CONTEXT_ALL, ACT_SET, "z", "random!");
	add(CONTEXT_ALL, ACT_SET, "r", "repeat!");
	add(CONTEXT_ALL, ACT_SET, "s", "single!");
	add(CONTEXT_ALL, ACT_SET, "+", "volume+=2");
	add(CONTEXT_ALL, ACT_SET, "-", "volume-=2");
	add(CONTEXT_ALL, ACT_SET, "m", "mute!");

	add(CONTEXT_ALL, ACT_TOGGLEPLAY, "<Space>");
	add(CONTEXT_SONGLIST, ACT_PLAY, "<Enter>");
	add(CONTEXT_ALL, ACT_STOP, "<Backspace>");
	add(CONTEXT_ALL, ACT_NEXT, "<Right>");
	add(CONTEXT_ALL, ACT_PREVIOUS, "<Left>");
	add(CONTEXT_ALL, ACT_SEEK_FORWARD, "\\>");
	add(CONTEXT_ALL, ACT_SEEK_BACK, "\\<");
}

Keybinding * Keybindings::add(int context, action_t action, string sequence, string params)
{
	Keybinding * c;
	vector<int> * seq;

	seq = conv_sequence(sequence);
	if (!seq)
	{
		return NULL;
	}
	if ((c = find_conflict(seq)))
	{
		sterr("Key binding '%s' conflicts with already configured '%s'.", sequence.c_str(), c->seqstr.c_str());
		return NULL;
	}

	c = new Keybinding();
	c->context = context;
	c->action = action;
	c->sequence = *seq;
	c->seqstr = sequence;
	c->params = params;
	bindings.push_back(c);
	return c;
}

vector<int> * Keybindings::conv_sequence(string seq)
{
	vector<int> * r;
	string subseq;
	size_t i;
	size_t epos;
	int escape = false;
	int ex = false;
	int ctrl = false;

	r = new vector<int>;

	for (i = 0; i < seq.size(); i++)
	{
		/* Escape sequence */
		if (seq[i] == '\\')
		{
			if (escape)
				r->push_back(seq[i]);
			escape = !escape;
			continue;
		}
		/* Start a special key */
		else if (seq[i] == '<')
		{
			if (escape)
			{
				escape = false;
				r->push_back(seq[i]);
				continue;
			}
			if (ex)
			{
				sterr("Bind: unexpected '%c' near ...%s, declaration dropped.", seq[i], seq.substr(i - 1).c_str());
				delete r;
				return NULL;
			}
			ex = true;
		}
		/* End a special key */
		else if (seq[i] == '>')
		{
			if (escape)
			{
				escape = false;
				r->push_back(seq[i]);
				continue;
			}
			if (!ex)
			{
				sterr("Bind: unexpected '%c' near ...%s, declaration dropped.", seq[i], seq.substr(i).c_str());
				delete r;
				return NULL;
			}
			ex = false;
		}
		else
		{
			/* Just a normal character */
			if (!ex)
			{
				r->push_back(seq[i]);
				continue;
			}
			else if (ctrl)
			{
				seq[i] = ::toupper(seq[i]);
				if (!(seq[i] >= 'A' || seq[i] <= 'Z'))
				{
					sterr("Bind: unexpected %c, expected one letter between A-Z near ...%s, declaration dropped.", seq[i], seq.substr(i - 3).c_str());
					delete r;
					return NULL;
				}
				r->push_back(seq[i] - 64);
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
				sterr("Bind: unclosed tag near ...%s, declaration dropped.", seq.substr(i - 1).c_str());
				delete r;
				return NULL;
			}

			subseq = seq.substr(i, epos - i);
			std::transform(subseq.begin(), subseq.end(), subseq.begin(), ::tolower);
			if (subseq == "up")
				r->push_back(KEY_UP);
			else if (subseq == "down")
				r->push_back(KEY_DOWN);
			else if (subseq == "left")
				r->push_back(KEY_LEFT);
			else if (subseq == "right")
				r->push_back(KEY_RIGHT);
			else if (subseq == "pgup")
				r->push_back(KEY_PPAGE);
			else if (subseq == "pgdn")
				r->push_back(KEY_NPAGE);
			else if (subseq == "home")
				r->push_back(KEY_HOME);
			else if (subseq == "end")
				r->push_back(KEY_END);
			else if (subseq == "backspace")
				r->push_back(KEY_BACKSPACE);
			else if (subseq == "delete")
				r->push_back(KEY_DC);
			else if (subseq == "insert")
				r->push_back(KEY_IC);
			else if (subseq == "cr" || subseq == "return" || subseq == "enter")
				r->push_back(KEY_ENTER);
			else if (subseq == "space")
				r->push_back(' ');
			else if (subseq == "tab")
				r->push_back('\t');

			else
			{
				sterr("Bind: invalid identifier '<%s>', declaration dropped.", subseq.c_str());
				delete r;
				return NULL;
			}

			i = epos - 1;
		}
	}

	return r;
}

Keybinding * Keybindings::find_conflict(vector<int> * sequence)
{
	vector<Keybinding *>::iterator i;
	unsigned int s;
	unsigned int limit;

	for (i = bindings.begin(); i != bindings.end(); i++)
	{
		limit = sequence->size() > (*i)->sequence.size() ? (*i)->sequence.size() : sequence->size();
		for (s = 0; s < limit; ++s)
		{
			if (sequence->at(s) != (*i)->sequence.at(s))
				break;
		}
		if (s == limit)
			return *i;
	}

	return NULL;
}

bool Keybindings::remove(string sequence)
{
	vector<Keybinding *>::iterator it;
	Keybinding * c;
	vector<int> * seq;

	if ((seq = conv_sequence(sequence)) == NULL)
		return false;

	c = find_conflict(seq);
	delete seq;

	if (c)
	for (it = bindings.begin(); it != bindings.end(); ++it)
	{
		if (*it == c)
		{
			debug("Removed key mapping `%s'.", sequence.c_str());
			delete c;
			bindings.erase(it);
			return true;
		}
	}

	sterr("Key sequence `%s' is not mapped.", sequence.c_str());
	return false;
}

int Keybindings::find(int context, vector<int> * sequence, action_t * action, string * params)
{
	vector<Keybinding *>::iterator i;
	int found = KEYBIND_FIND_NOMATCH;
	unsigned int s;

	for (i = bindings.begin(); i != bindings.end(); i++)
	{
		if (!((*i)->context & context) || (*i)->sequence.size() < sequence->size())
			continue;

		for (s = 0; s < sequence->size(); ++s)
		{
			if (sequence->at(s) == (*i)->sequence.at(s))
				found = KEYBIND_FIND_BUFFERED;
			else
				break;
		}

		if (s == (*i)->sequence.size())
		{
			*action = (*i)->action;
			*params = (*i)->params;
			return KEYBIND_FIND_EXACT;
		}
	}

	return found;
}

void Keybindings::truncate()
{
	vector<Keybinding *>::const_iterator i;

	for (i = bindings.begin(); i != bindings.end(); ++i)
		delete *i;

	bindings.clear();
}

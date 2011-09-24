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

#ifndef _PMS_INPUT_H_
#define _PMS_INPUT_H_

#include "command.h"
#include <string>
using namespace std;

#define INPUT_RESULT_NOINPUT 0
#define INPUT_RESULT_BUFFERED 1
#define INPUT_RESULT_RUN 2

#define INPUT_MODE_COMMAND 0
#define INPUT_MODE_INPUT 1
#define INPUT_MODE_SEARCH 2

#define KEYBIND_FIND_NOMATCH -1
#define KEYBIND_FIND_EXACT 0
#define KEYBIND_FIND_BUFFERED 1

/* This class is returned by Input::next(), defining what to do next */
class Inputevent
{
	public:
		Inputevent();
		void clear();

		/* Which context are we in? */
		int context;

		/* How many? */
		unsigned int multiplier;

		/* What kind of action to run */
		action_t action;

		/* one of INPUT_RESULT_* */
		int result;

		/* The full character buffer or command/input */
		string text;
};

class Keybinding
{
	public:
		string		seqstr;
		vector<int>	sequence;
		action_t	action;
		int		context;
};

class Keybindings
{
	private:
		vector<Keybinding *>	bindings;
	public:
		Keybindings();

		/* Add and check for duplicate sequences */
		Keybinding *	add(int context, action_t action, string sequence);
		Keybinding *	find_conflict(vector<int> * sequence);

		/* Convert a string sequence to an int sequence */
		vector<int> *	conv_sequence(string seq);

		/* Find an action based on the key sequence */
		int		find(int context, vector<int> * sequence, action_t * action);
};

class Input
{
	private:

		int		chbuf;
		bool		is_tab_completing;
		Inputevent	ev;

		void		handle_text_input();

	public:

		int		mode;
		unsigned long	multiplier;
		vector<int> 	buffer;
		string		strbuf;

		Input();

		/* Read next character from ncurses buffer */
		Inputevent *	next();

		/* Setter and getter for mode */
		void		setmode(int nmode);
		int		getmode() { return mode; }
};

#endif /* _PMS_INPUT_H_ */

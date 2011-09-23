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

/* This struct is returned by Input::next(), defining what to do next */
typedef struct
{
	/* Which context are we in? */
	int context;

	/* How many? */
	unsigned int multiplier;

	/* What kind of action to run */
	action_t action;

	/* one of INPUT_RESULT_* */
	int result;
}

input_event;


class Keybinding
{
	public:
		string		sequence;
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
		Keybinding *	find_conflict(string sequence);

		/* Convert a string sequence to a binary sequence */
		string		conv_sequence(string seq);

		/* Find an action based on the key sequence */
		int		find(int context, string sequence, action_t * action);
};

class Input
{
	private:

		int		chbuf;
		input_event	ev;

	public:

		int		mode;
		unsigned long	multiplier;
		string		buffer;

		Input();

		/* Read next character from ncurses buffer */
		input_event *	next();

		/* Setter and getter for mode */
		void		setmode(int nmode);
		int		getmode() { return mode; }
};

#endif /* _PMS_INPUT_H_ */

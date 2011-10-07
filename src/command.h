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

#ifndef _PMS_COMMAND_H_
#define _PMS_COMMAND_H_

#include <string>
#include <vector>
using namespace std;

enum
{
	CONTEXT_CONSOLE		= (1 << 0),
	CONTEXT_SONGLIST	= (1 << 1),

	CONTEXT_LIST		= (1 << 1) |
				  (1 << 0),

	CONTEXT_ALL		= (1 << 2) - 1
};

typedef enum
{
	/* Configuration and internals */
	ACT_NOACTION = 0,
	ACT_RESIZE,
	ACT_RUN_CMD,
	ACT_SET,

	/* Global commands */
	ACT_QUIT,
	ACT_MODE_INPUT,
	ACT_MODE_COMMAND,
	ACT_NEXT_WINDOW,
	ACT_PREVIOUS_WINDOW,

	/* List actions */
	ACT_SCROLL_UP,
	ACT_SCROLL_DOWN,
	ACT_CURSOR_UP,
	ACT_CURSOR_DOWN,
	ACT_CURSOR_PGUP,
	ACT_CURSOR_PGDOWN,
	ACT_CURSOR_HOME,
	ACT_CURSOR_END,
	ACT_CURSOR_TOP,
	ACT_CURSOR_BOTTOM,
	ACT_CURSOR_CURRENTSONG,

	/* MPD options */
	ACT_CROSSFADE,
	ACT_SETVOL,
	ACT_VOLUME_UP,
	ACT_VOLUME_DOWN,

	/* Playback */
	ACT_TOGGLEPLAY,
	ACT_PLAY,
	ACT_STOP,
	ACT_NEXT,
	ACT_PREVIOUS,
	ACT_SEEK_FORWARD,
	ACT_SEEK_BACK
}

action_t;

typedef enum
{
	COMMAND_PARAM_INT,
	COMMAND_PARAM_STRING
}

commandparam_t;

class Command
{
	private:
		vector<commandparam_t>	params;

	public:
		/* Where the command can be executed. Combine several by bitwise or */
		int			context;

		/* Map to which action? */
		action_t		action;

		/* String representation of command */
		string			name;

		/* Add a parameter to this command */
		void			addparam(commandparam_t param) { params.push_back(param); };

		/* Get number of parameters */
		unsigned int		numparams() { return params.size(); };


};

/*
 * Contains a list of all commands, used to search for a command.
 */
class Commandlist
{
	private:
		vector<Command *>	cmds;
		vector<Command *>	grepcmds;
		
		Command *		add(int context, action_t action, string name);

	public:
		/* Set up all available pms commands */
		Commandlist();

		/* Search for a command, used for tab completion */
		vector<Command *> *	grep(int context, string name);

		/* Exact match for command */
		Command *		find(int context, string name);
};

#endif /* _PMS_COMMAND_H_ */

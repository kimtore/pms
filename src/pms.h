/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2010  Kim Tore Jensen
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
 *
 * pms.h - The PMS main class
 *
 */

#ifndef _PMS_H_
#define _PMS_H_

#define PMS_NAME "Practical Music Search"
#define PMS_COPYRIGHT "(c) 2006-2010 by Kim Tore Jensen <kimtjen@gmail.com>"

#define PMS_EXIT_SUCCESS 0
#define PMS_EXIT_NODISPLAY 1
#define PMS_EXIT_NOCOMMAND 2
#define PMS_EXIT_CANTCONNECT 3
#define PMS_EXIT_BADARGS 4
#define PMS_EXIT_BADPASS 5
#define PMS_EXIT_NOWINDOWS 6
#define PMS_EXIT_CONFIGERR 7
#define PMS_EXIT_LOMEM 8

#include "../config.h"
#include <assert.h>
#include <cstdarg>
#include <string>
#include <sstream>
#include <stdio.h>
#include "i18n.h"
#include "libmpdclient.h"

#include "types.h"
#include "settings.h"
#include "song.h"
#include "message.h"
#include "topbar.h"

#include "config.h"
#include "color.h"
#include "list.h"
#include "action.h"
#include "input.h"
#include "mediator.h"

#ifdef __FreeBSD__
	#include <sys/wait.h>
#endif

using namespace std;

/*
 * Global functions
 */
void					debug(const char *, ...);

/*
 * This is the program itself, everything is run within here.
 */
class Pms
{
private:

	/* Variables */
	int				argc;
	char **				argv;

	bool				_shutdown;
	Connection *			conn;
	pms_win_playlist *		playlist;
	pms_win_playlist *		library;
	vector<Message *>		msglog;

	/* Private functions */
	void				init_default_keymap();
	bool				connect_window_list();
	bool				progress_nextsong();

	/* Options/arguments */
	void				print_version();
	void				print_usage();
	bool				require_arg(char);
	bool				parse_args(int, char **);

public:

					Pms(int, char **);
					~Pms();

	/* Public variables */
	//FIXME: program should be rewritten so that none of these should have to be public
	Options	*			options;
	Control *			comm;
	Display *			disp;
	Input *				input;
	Commandmap *			commands;
	Bindings *			bindings;
	Fieldtypes *			fieldtypes;
	Formatter *			formatter;
	Configurator *			config;

	/* FIXME: this is an attempt on the above */
	Mediator *			mediator;
	Message *			msg;
	Interface *			interface;

	/* Global public functions */
	static string			tostring(long);
	static string			tostring(int);
	static string			tostring(size_t);
	static vector<string> *		splitstr(string, string = " ");
	static string			joinstr(vector<string> *, vector<string>::iterator, vector<string>::iterator, string = " ");
	static string			timeformat(int);
	static string			pluralformat(unsigned int);
	static string			zeropad(int, unsigned int);
	static string			formtext(string);
	static bool			unicode();

	/* Public member functions */
	void				setstatus(statusbar_mode, const char *, ...);
	void				drawstatus();
	int				resetstatus(int);
	void				shutdown() { _shutdown = true; };
	bool				run_shell(string);
	Song *				cursong();
	string				playstring();
	void				putlog(Message *);			// Put an arbitrary message into the message log.
	void				log(int, long, const char *, ...);

	/* Main loop and initialization */
	int				init();
	int				main();
};

#endif

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

#ifndef _PMS_CONFIG_H_
#define _PMS_CONFIG_H_

#include <string>
#include "color.h"
#include "field.h"
using namespace std;

enum
{
	OPT_CHANGE_NONE		= 0,
	OPT_CHANGE_MPD		= 1 << 0,
	OPT_CHANGE_DIMENSIONS	= 1 << 1,
	OPT_CHANGE_COLUMNS	= 1 << 2,
	OPT_CHANGE_DRAWLIST	= 1 << 3,
	OPT_CHANGE_REDRAW	= 1 << 4,
	OPT_CHANGE_TOPBAR	= 1 << 5,
	OPT_CHANGE_PLAYMODE	= 1 << 6,
	OPT_CHANGE_ALL		= (1 << 7) - 1
};

/* Standard option types */
typedef enum
{
	OPTION_TYPE_INVALID,
	OPTION_TYPE_BOOL,
	OPTION_TYPE_UINT,
	OPTION_TYPE_INT,
	OPTION_TYPE_STRING,
	OPTION_TYPE_COLOR,
	OPTION_TYPE_COLORLIST,

	/* More exotic stuff */
	OPTION_TYPE_COLUMNHEADERS,
	OPTION_TYPE_TOPBAR,
	OPTION_TYPE_SEARCHFIELDS,
	OPTION_TYPE_SCROLLMODE
}

option_type_t;

typedef struct
{
	string name;
	option_type_t type;
	void * ptr;
	int mask;
}

option_t;

typedef enum
{
	SCROLL_MODE_NORMAL,
	SCROLL_MODE_CENTERED
}

scroll_mode_t;



class Config
{
	private:
		void			setup_default_connection_info();
		void			set_column_headers(string hdr);
		void			set_search_fields(string fields);
		void			set_scroll_mode(string mode);

		vector<option_t *>	options;
		
		/* Add an option to the options vector. */
		option_t *		add_option(string name, option_type_t type, void * ptr, int mask);

	public:

		Config();

		/* Load all default config files */
		void		source_default_config();

		/* Load a config file */
		bool		source(string filename, bool suppress_errmsg = false);

		/* Parse "option=value" */
		option_t *	readline(string line, bool verbose = true);

		/* Option string getter and setter */
		int		add_opt_str(option_t * opt, string value, int arithmetic);
		int		set_opt_str(option_t * opt, string value);
		string		get_opt_str(option_t * opt);

		/* Print option values to the console. */
		void		print_option(option_t * opt);
		int		print_all_options();
		int		print_all_colors();

		/* Return the option_t struct of the option in question. */
		option_t *	get_opt_ptr(string opt);

		/* Tab-complete search, return a list of option_t */
		unsigned int	grep_opt(string opt, vector<option_t *> * list, string * prefix);


		/* Connection parameters */
		string		host;
		string		port;
		string		password;

		/* Default songlist column headers */
		vector<Field *>	songlist_columns;

		/* Main loop variable */
		bool		quit;

		/* Autoconnect timeout */
		unsigned int	reconnect_delay;

		/* How many seconds left of song before adding next song */
		unsigned int	add_next_interval;

		/* System bell */
		bool		use_bell;
		bool		visual_bell;

		/* Use column headers */
		bool		show_column_headers;
		bool		show_window_title;

		/* Scroll/cursor mode */
		scroll_mode_t	scroll_mode;

		/* Topbar stuff */
		unsigned int	topbar_height;

		/* Auto-advance to next song? */
		bool		autoadvance;

		/* Playback follows window when switched. */
		bool		playback_follows_window;

		/* Default sort string */
		string		default_sort;

		/* Ignore case on searching/sorting? */
		bool		sort_case;
		bool		search_case;

		/* We need to handle MPD's options in config, too. */
		bool		random;
		bool		repeat;
		bool		consume;
		bool		single;
		bool		mute;
		int		volume;

		/* Redraw play string in statusbar after this long */
		bool		status_reset_interval;

		/* What fields to search by default */
		long		search_field_mask;

		/* Advance cursor on add actions */
		bool		advance_cursor;

		/* The entire color collection */
		Colortable	colors;
};

#endif /* _PMS_CONFIG_H_ */

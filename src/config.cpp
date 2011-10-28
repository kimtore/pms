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

#include "config.h"
#include "field.h"
#include "console.h"
#include "song.h"
#include "window.h"
#include "topbar.h"
#include "pms.h"
#include <cstring>
#include <stdlib.h>
#include <algorithm>
#include <string>
#include <iostream>
#include <fstream>

using namespace std;

extern Fieldtypes fieldtypes;
extern Windowmanager wm;
extern PMS pms;
extern Keybindings keybindings;
Topbar topbar;

void Config::load_default_config()
{
	setup_default_connection_info();

	quit = false;
	reconnect_delay = 5;
	use_bell = true;
	visual_bell = false;
	show_column_headers = true;
	show_window_title = true;
	topbar_height = 2;
	add_next_interval = 5;
	autoadvance = true;
	status_reset_interval = 2;
	playback_follows_window = true;
	advance_cursor = true;
	split_search_terms = true;
	random = false;
	repeat = false;
	consume = false;
	single = false;
	mute = false;
	volume = 100;
	sort_case = false;
	search_case = false;
	autoconnect = true;
	default_sort = "track disc album date albumartistsort";
	set_column_headers("artist track title album year length");
	set_search_fields("artist title album");
	set_scroll_mode("normal");
	topbar.set("{PMS $if(connected){$if(song){$volume $state [$modes] $elapsed / $remaining}}$else{disconnected}}"
			"{$if(song){$artist / $title / $album / $year}}"
			"{$if(connected){Q:$queuesize/$queuelength S:$listsize/$listlength}}"
			"{$progressbar}{}{}");

	colors.load_defaults();
	keybindings.load_defaults();
}

Config::Config()
{
	/* Load internal defaults */
	load_default_config();

	/* Set up options array */
	add_option("host", OPTION_TYPE_STRING, (void *)&host, OPT_CHANGE_NONE);
	add_option("port", OPTION_TYPE_STRING, (void *)&port, OPT_CHANGE_NONE);
	add_option("password", OPTION_TYPE_STRING, (void *)&password, OPT_CHANGE_NONE);
	add_option("autoconnect", OPTION_TYPE_BOOL, (void *)&autoconnect, OPT_CHANGE_NONE);

	add_option("reconnectdelay", OPTION_TYPE_UINT, (void *)&reconnect_delay, OPT_CHANGE_NONE);
	add_option("addnextinterval", OPTION_TYPE_UINT, (void *)&add_next_interval, OPT_CHANGE_NONE);

	add_option("advancecursor", OPTION_TYPE_BOOL, (void *)&advance_cursor, OPT_CHANGE_NONE);
	add_option("bell", OPTION_TYPE_BOOL, (void *)&use_bell, OPT_CHANGE_NONE);
	add_option("visualbell", OPTION_TYPE_BOOL, (void *)&visual_bell, OPT_CHANGE_NONE);
	add_option("columnheaders", OPTION_TYPE_BOOL, (void *)&show_column_headers, OPT_CHANGE_DRAWLIST);
	add_option("windowtitle", OPTION_TYPE_BOOL, (void *)&show_window_title, OPT_CHANGE_DRAWLIST);
	add_option("autoadvance", OPTION_TYPE_BOOL, (void *)&autoadvance, OPT_CHANGE_NONE);
	add_option("followwindow", OPTION_TYPE_BOOL, (void *)&playback_follows_window, OPT_CHANGE_NONE);
	add_option("resetstatus", OPTION_TYPE_UINT, (void *)&status_reset_interval, OPT_CHANGE_NONE);

	add_option("random", OPTION_TYPE_BOOL, (void *)&random, OPT_CHANGE_MPD | OPT_CHANGE_TOPBAR | OPT_CHANGE_PLAYMODE);
	add_option("repeat", OPTION_TYPE_BOOL, (void *)&repeat, OPT_CHANGE_MPD);
	add_option("consume", OPTION_TYPE_BOOL, (void *)&consume, OPT_CHANGE_MPD);
	add_option("single", OPTION_TYPE_BOOL, (void *)&single, OPT_CHANGE_MPD);
	add_option("mute", OPTION_TYPE_BOOL, (void *)&mute, OPT_CHANGE_MPD);
	add_option("volume", OPTION_TYPE_VOLUME, (void *)&volume, OPT_CHANGE_MPD);

	add_option("sort", OPTION_TYPE_STRING, (void *)&default_sort, OPT_CHANGE_NONE);
	add_option("casesort", OPTION_TYPE_BOOL, (void *)&sort_case, OPT_CHANGE_NONE);
	add_option("casesearch", OPTION_TYPE_BOOL, (void *)&search_case, OPT_CHANGE_NONE);
	add_option("wordsearch", OPTION_TYPE_BOOL, (void *)&split_search_terms, OPT_CHANGE_NONE);

	add_option("scroll", OPTION_TYPE_SCROLLMODE, (void *)&scroll_mode, OPT_CHANGE_DRAWLIST);
	add_option("searchfields", OPTION_TYPE_SEARCHFIELDS, (void *)&search_field_mask, OPT_CHANGE_NONE);
	add_option("columns", OPTION_TYPE_COLUMNHEADERS, (void *)&songlist_columns, OPT_CHANGE_COLUMNS | OPT_CHANGE_DRAWLIST);
	add_option("topbar", OPTION_TYPE_TOPBAR, (void *)&topbar, OPT_CHANGE_DIMENSIONS | OPT_CHANGE_REDRAW);
	add_option("topbarlines", OPTION_TYPE_UINT, (void *)&topbar_height, OPT_CHANGE_DIMENSIONS | OPT_CHANGE_REDRAW);

	/*
	 * Set up all colors
	 */

	add_option("color", OPTION_TYPE_COLORLIST, NULL, OPT_CHANGE_NONE);
	add_option("color.topbar", OPTION_TYPE_COLOR, (void *)colors.topbar, OPT_CHANGE_NONE);
	add_option("color.statusbar", OPTION_TYPE_COLOR, (void *)colors.statusbar, OPT_CHANGE_NONE);
	add_option("color.windowtitle", OPTION_TYPE_COLOR, (void *)colors.windowtitle, OPT_CHANGE_NONE);
	add_option("color.columnheaders", OPTION_TYPE_COLOR, (void *)colors.columnheader, OPT_CHANGE_NONE);
	add_option("color.console", OPTION_TYPE_COLOR, (void *)colors.console, OPT_CHANGE_NONE);
	add_option("color.error", OPTION_TYPE_COLOR, (void *)colors.error, OPT_CHANGE_NONE);
	add_option("color.readout", OPTION_TYPE_COLOR, (void *)colors.readout, OPT_CHANGE_NONE);
	add_option("color.cursor", OPTION_TYPE_COLOR, (void *)colors.cursor, OPT_CHANGE_NONE);
	add_option("color.playing", OPTION_TYPE_COLOR, (void *)colors.playing, OPT_CHANGE_NONE);
	add_option("color.directory", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_DIRECTORY], OPT_CHANGE_NONE);
	add_option("color.file", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_FILE], OPT_CHANGE_NONE);
	add_option("color.pos", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_POS], OPT_CHANGE_NONE);
	add_option("color.id", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ID], OPT_CHANGE_NONE);
	add_option("color.time", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_TIME], OPT_CHANGE_NONE);
	add_option("color.name", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_NAME], OPT_CHANGE_NONE);
	add_option("color.artist", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ARTIST], OPT_CHANGE_NONE);
	add_option("color.artistsort", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ARTISTSORT], OPT_CHANGE_NONE);
	add_option("color.album", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ALBUM], OPT_CHANGE_NONE);
	add_option("color.title", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_TITLE], OPT_CHANGE_NONE);
	add_option("color.track", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_TRACK], OPT_CHANGE_NONE);
	add_option("color.disc", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_DISC], OPT_CHANGE_NONE);
	add_option("color.date", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_DATE], OPT_CHANGE_NONE);
	add_option("color.genre", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_GENRE], OPT_CHANGE_NONE);
	add_option("color.albumartist", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ALBUMARTIST], OPT_CHANGE_NONE);
	add_option("color.albumartistsort", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ALBUMARTISTSORT], OPT_CHANGE_NONE);
	add_option("color.year", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_YEAR], OPT_CHANGE_NONE);
	add_option("color.trackshort", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_TRACKSHORT], OPT_CHANGE_NONE);
	add_option("color.elapsed", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_ELAPSED], OPT_CHANGE_NONE);
	add_option("color.remaining", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_REMAINING], OPT_CHANGE_NONE);
	add_option("color.volume", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_VOLUME], OPT_CHANGE_NONE);
	add_option("color.progressbar", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_PROGRESSBAR], OPT_CHANGE_NONE);
	add_option("color.modes", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_MODES], OPT_CHANGE_NONE);
	add_option("color.state", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_STATE], OPT_CHANGE_NONE);
	add_option("color.queuesize", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_QUEUESIZE], OPT_CHANGE_NONE);
	add_option("color.queuelength", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_QUEUELENGTH], OPT_CHANGE_NONE);
	add_option("color.listsize", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_LISTSIZE], OPT_CHANGE_NONE);
	add_option("color.listlength", OPTION_TYPE_COLOR, (void *)colors.field[FIELD_LISTLENGTH], OPT_CHANGE_NONE);
}

void Config::source_default_config()
{
	string home;
	string s;
	char * env;
	const string suffix = "/pms/pms.conf";
	size_t start = 0, end = 0;

	debug("Reading configuration files...", NULL);

	if ((env = getenv("HOME")) != NULL)
		home = env;
	else
		home.clear();

	// XDG config dirs (colon-separated priority list, defaults to just /etc/xdg)
	if ((env = getenv("XDG_CONFIG_DIRS")) == NULL)
	{
		source("/usr/local/etc/xdg" + suffix, true);
		source("/etc/xdg" + suffix, true);
	}
	else
	{
		s = env;
		while ((end = s.find(':', start)) != string::npos)
		{
			source(s.substr(start, end - start) + suffix, true);
			start = end + 1;
		}

		if (start < s.size())
			source(s.substr(start) + suffix, true);
	}

	// XDG config home (usually $HOME/.config)
	if ((env = getenv("XDG_CONFIG_HOME")) == NULL)
	{
		if (home.size() > 0)
			source(home + "/.config" + suffix);
	}
	else
	{
		source(env + suffix);
	}
}

bool Config::source(string filename, bool suppress_errmsg)
{
	ifstream fd;
	char line[1024];

	fd.open(filename.c_str(), ifstream::in);
	if (!fd.good())
	{
		if (!suppress_errmsg)
			sterr("Cannot open file `%s'", filename.c_str());
		return false;
	}

	while (fd.good())
	{
		fd.getline(line, sizeof line);
		pms.run_cmd(line, 1, true);
	}

	fd.close();

	return true;
}

option_t * Config::readline(string line, bool verbose)
{
	string optstr;
	string optval = "";
	option_t * opt;
	size_t pos;
	bool invert = false;
	bool show = false;
	bool negative = false;
	bool * bopt;
	int arithmetic = 0;

	/* Locate the identifier */
	if (line.size() == 0)
	{
		print_all_options();
		return NULL;
	}
	else if ((pos = line.find_first_of("=:")) != string::npos)
	{
		optstr = line.substr(0, pos);
		if (line.size() > pos + 1)
			optval = line.substr(pos + 1);
	}
	else
	{
		optstr = line;
	}

	/* Invert or return value? */
	switch(optstr[optstr.size()-1] )
	{
		case '?':
			show = true;
			break;
		case '!':
			invert = true;
			break;
		default:;
	}

	/* Detect += or -= values */
	if (pos != string::npos && line[pos] == '=')
	{
		switch(optstr[optstr.size()-1])
		{
			case '+':
				arithmetic = 1;
				break;
			case '-':
				arithmetic = -1;
				break;
			default:;
		}
	}

	/* Cut away operators */
	if (show || invert || arithmetic != 0)
		optstr = optstr.substr(0, optstr.size() - 1);

	/* Return the option struct if this is a valid option */
	if ((opt = get_opt_ptr(optstr)) == NULL)
	{
		/* Check if this is a negative boolean (no<option>) */
		if (optstr.size() > 2 && optstr.substr(0, 2) == "no" && ((opt = get_opt_ptr(optstr.substr(2))) != NULL) && opt->type == OPTION_TYPE_BOOL)
		{
			negative = !invert;
			optstr = optstr.substr(2);
		}
		/* Check if this is an invertion (inv<option>) */
		else if (optstr.size() > 3 && optstr.substr(0, 3) == "inv" && ((opt = get_opt_ptr(optstr.substr(3))) != NULL) && opt->type == OPTION_TYPE_BOOL)
		{
			invert = true;
			optstr = optstr.substr(3);
		}
		else
		{
			sterr("Unknown option: %s", line.c_str());
			return NULL;
		}
	}

	/* Print the option to statusbar */
	if (show)
	{
		print_option(opt);
		return NULL;
	}

	if (arithmetic != 0)
	{
		/* Add the new string value to the previous values */
		if (add_opt_str(opt, optval, arithmetic))
		{
			if (verbose)
				print_option(opt);
			return opt;
		}
		return NULL;
	}

	/* Invert an option if boolean */
	if (invert)
	{
		if (opt->type != OPTION_TYPE_BOOL)
		{
			debug("%s=%s", optstr.c_str(), get_opt_str(opt).c_str());
			sterr("Trailing characters: %s", line.c_str());
			return NULL;
		}
		bopt = (bool *)opt->ptr;
		*bopt = !(*bopt);
		if (verbose)
			print_option(opt);
		return opt;
	}

	/* Check for (negative) boolean options */
	if (optval.size() == 0 && pos == string::npos)
	{
		/* Show option instead */
		if (opt->type != OPTION_TYPE_BOOL)
		{
			if (verbose)
				print_option(opt);
			return NULL;
		}
		bopt = (bool *)opt->ptr;
		*bopt = !negative;
		if (verbose)
			print_option(opt);
		return opt;
	}

	/* Set the new string value */
	if (set_opt_str(opt, optval))
	{
		if (verbose)
			print_option(opt);
		return opt;
	}

	return NULL;
}

option_t * Config::add_option(string name, option_type_t type, void * ptr, int mask)
{
	option_t * o = new option_t;
	o->name = name;
	o->type = type;
	o->ptr = ptr;
	o->mask = mask;
	options.push_back(o);
	return o;
}

string Config::get_opt_str(option_t * opt)
{
	vector<Field *>::const_iterator field_it;

	Color * c;
	string str = "";
	unsigned int * ui;
	int * i;
	bool * b;

	if (opt == NULL)
		return str;

	switch(opt->type)
	{
		case OPTION_TYPE_STRING:
			str = (*(string *)opt->ptr);
			break;

		case OPTION_TYPE_BOOL:
			b = (bool *)opt->ptr;
			str = !(*b) ? "no" : "";
			str += opt->name;
			break;

		case OPTION_TYPE_UINT:
		case OPTION_TYPE_VOLUME:
			ui = (unsigned int *)opt->ptr;
			str = tostring(*ui);
			break;

		case OPTION_TYPE_INT:
			i = (int *)opt->ptr;
			str = tostring(*i);
			break;

		case OPTION_TYPE_COLOR:
			c = (Color *)opt->ptr;
			str = c->getstrname();
			break;

		/* Exotic data types */

		case OPTION_TYPE_SCROLLMODE:
			if (scroll_mode == SCROLL_MODE_NORMAL)
				str = "normal";
			else if (scroll_mode == SCROLL_MODE_CENTERED)
				str = "centered";
			break;

		case OPTION_TYPE_COLUMNHEADERS:
			for (field_it = songlist_columns.begin(); field_it != songlist_columns.end(); ++field_it)
				str = str + (*field_it)->str + " ";
			str = str.substr(0, str.size() - 1);
			break;

		case OPTION_TYPE_SEARCHFIELDS:
			for (field_it = fieldtypes.fields.begin(); field_it != fieldtypes.fields.end(); ++field_it)
				if (search_field_mask & (1 << (*field_it)->type))
					str = str + (*field_it)->str + " ";
			str = str.substr(0, str.size() - 1);
			break;

		case OPTION_TYPE_TOPBAR:
			str = topbar.cached_format;
			break;

		default:
			str = "<unknown>";
			break;
	}

	return str;
}

int Config::add_opt_str(option_t * opt, string value, int arithmetic)
{
	string s;
	int * i;
	unsigned int * ui;

	if (opt == NULL)
		return false;

	if (arithmetic == 0)
		return set_opt_str(opt, value);

	switch(opt->type)
	{
		case OPTION_TYPE_COLOR:
		case OPTION_TYPE_COLORLIST:
			/* Easter egg for those who like to play with their apps */
			sterr("If you want to play with colors, please buy a palette.", NULL);
			return false;

		case OPTION_TYPE_COLUMNHEADERS:
		case OPTION_TYPE_SEARCHFIELDS:
			if (arithmetic == 1)
				value = " " + value;
			/* break intentionally omitted */

		case OPTION_TYPE_STRING:
		case OPTION_TYPE_TOPBAR:
			s = get_opt_str(opt);
			if (arithmetic == 1)
				s = s + value;
			else if (arithmetic == -1)
				s = str_replace(value, "", s);
			set_opt_str(opt, s);
			return true;

		case OPTION_TYPE_INT:
			i = (int *)opt->ptr;
			*i = *i + (arithmetic * atoi(value.c_str()));
			return true;

		case OPTION_TYPE_UINT:
			ui = (unsigned int *)opt->ptr;
			*ui = *ui + (arithmetic * atoi(value.c_str()));
			return true;

		case OPTION_TYPE_VOLUME:
			ui = (unsigned int *)opt->ptr;
			*ui = *ui + (arithmetic * atoi(value.c_str()));
			if (*ui > 100)
				*ui = 100;
			if (*ui < 0)
				*ui = 0;
			return true;

		default:
			return false;
	}
}

int Config::set_opt_str(option_t * opt, string value)
{
	Color * c;
	string * s;
	int * i;
	unsigned int * ui;

	if (opt == NULL)
		return false;

	switch(opt->type)
	{
		case OPTION_TYPE_STRING:
			s = (string *)opt->ptr;
			*s = value;
			return true;

		case OPTION_TYPE_INT:
			i = (int *)opt->ptr;
			*i = atoi(value.c_str());
			return true;

		case OPTION_TYPE_UINT:
			ui = (unsigned int *)opt->ptr;
			*ui = atoi(value.c_str());
			return true;

		case OPTION_TYPE_VOLUME:
			ui = (unsigned int *)opt->ptr;
			*ui = atoi(value.c_str());
			if (*ui > 100)
				*ui = 100;
			if (*ui < 0)
				*ui = 0;
			return true;

		case OPTION_TYPE_COLUMNHEADERS:
			set_column_headers(value);
			wm.update_column_length();
			return true;

		case OPTION_TYPE_SEARCHFIELDS:
			set_search_fields(value);
			return true;

		case OPTION_TYPE_SCROLLMODE:
			set_scroll_mode(value);
			return true;

		case OPTION_TYPE_COLOR:
			c = (Color *)opt->ptr;
			c->set(value);
			return true;

		case OPTION_TYPE_TOPBAR:
			topbar.set(value);
			if (topbar.lines[0].size() > topbar_height)
			{
				topbar_height = topbar.lines[0].size();
				print_option(get_opt_ptr("topbarlines"));
			}
			return true;

		default:
			return false;
	}

	return false;
}

option_t * Config::get_opt_ptr(string opt)
{
	vector<option_t *>::const_iterator i;

	for (i = options.begin(); i != options.end(); ++i)
		if ((*i)->name == opt)
			return *i;
	
	return NULL;
}

unsigned int Config::grep_opt(string opt, vector<option_t *> * list, string * prefix)
{
	vector<option_t *>::const_iterator i;

	if (!list) return 0;
	list->clear();

	/* Check for "no..." and "inv..." options, which also needs to be tab-completed. */
	if (opt.size() >= 2 && opt.substr(0, 2) == "no")
		*prefix = "no";
	else if (opt.size() >= 3 && opt.substr(0, 3) == "inv")
		*prefix = "inv";
	else
		prefix->clear();

	if (prefix->size() > 0)
	{
		if (opt.size() == prefix->size())
			opt.clear();
		else
			opt = opt.substr(prefix->size());
	}

	for (i = options.begin(); i != options.end(); i++)
	{
		if (opt.size() > (*i)->name.size())
			continue;

		if (opt == (*i)->name.substr(0, opt.size()))
		{
			if (prefix->size() == 0 || (*i)->type == OPTION_TYPE_BOOL
				|| ((*i)->name.size() > prefix->size() && (*i)->name.substr(0, prefix->size()) == *prefix))
				list->push_back(*i);
		}
	}

	return list->size();
}

void Config::print_option(option_t * opt)
{
	if (opt == NULL)
		return;
	else if (opt->type == OPTION_TYPE_BOOL)
		debug("  %s", get_opt_str(opt).c_str());
	else if (opt->type == OPTION_TYPE_COLORLIST)
		print_all_colors();
	else
		debug("  %s=%s", opt->name.c_str(), get_opt_str(opt).c_str());
}

int Config::print_all_options()
{
	vector<option_t *>::const_iterator i;

	debug("--- Options ---", NULL);

	for (i = options.begin(); i != options.end(); ++i)
		if ((*i)->type != OPTION_TYPE_COLOR && (*i)->type != OPTION_TYPE_COLORLIST)
			print_option(*i);

	return true;
}

int Config::print_all_colors()
{
	vector<option_t *>::const_iterator i;

	debug("--- Colors ---", NULL);

	for (i = options.begin(); i != options.end(); ++i)
		if ((*i)->type == OPTION_TYPE_COLOR)
			print_option(*i);

	return true;
}

void Config::set_column_headers(string hdr)
{
	size_t start = 0;
	size_t pos = 0;
	string f;
	Field * field;

	songlist_columns.clear();

	while (start + 1 < hdr.size())
	{
		if (pos == string::npos)
			break;

		if ((pos = hdr.find(' ', start)) != string::npos)
			f = hdr.substr(start, pos - start);
		else
			f = hdr.substr(start);

		if ((field = fieldtypes.find(f)) != NULL && field->type < FIELD_COLUMN_VALUES)
			songlist_columns.push_back(field);
		else
			sterr("Ignoring invalid header field '%s'.", f.c_str());

		start = pos + 1;
	}

	if (songlist_columns.size() == 0)
	{
		f = "title";
		sterr("Warning: at least one column type needs to be specified, falling back to `%s'.", f.c_str());
		songlist_columns.push_back(fieldtypes.find(f));
	}
}

void Config::set_search_fields(string fields)
{
	size_t start = 0;
	size_t pos = 0;
	string f;
	Field * field;

	search_field_mask = 0;

	while (start + 1 < fields.size())
	{
		if (pos == string::npos)
			break;

		if ((pos = fields.find(' ', start)) != string::npos)
			f = fields.substr(start, pos - start);
		else
			f = fields.substr(start);

		if ((field = fieldtypes.find(f)) != NULL && field->type < FIELD_COLUMN_VALUES)
			search_field_mask |= (1 << field->type);
		else
			sterr("Ignoring invalid header field '%s'.", f.c_str());

		start = pos + 1;
	}

	if (search_field_mask == 0)
	{
		search_field_mask = FIELD_FILE;
		sterr("Warning: at least one field needs to be specified, falling back to `file'.", NULL);
	}
}

void Config::set_scroll_mode(string mode)
{
	if (mode == "normal")
		scroll_mode = SCROLL_MODE_NORMAL;
	else if (mode == "centered")
		scroll_mode = SCROLL_MODE_CENTERED;
	else
		sterr("Invalid scroll mode `%s', expected one of `normal', `centered'", mode.c_str());
}

void Config::setup_default_connection_info()
{
	char *	env;
	size_t	i;

	password = "";

	if ((env = getenv("MPD_HOST")) == NULL)
	{
		host = "localhost";
	}
	else
	{
		host = env;
		if ((i = host.rfind('@')) != string::npos)
		{
			password = host.substr(0, i);
			host = host.substr(i + 1);
		}
	}

	if ((env = getenv("MPD_PORT")) == NULL)
		port = "6600";
	else
		port = env;
}

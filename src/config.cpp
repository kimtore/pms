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
#include <stdlib.h>
#include <algorithm>

using namespace std;

extern Fieldtypes fieldtypes;
extern Windowmanager wm;

Config::Config()
{
	setup_default_connection_info();

	quit = false;
	reconnect_delay = 5;
	use_bell = true;
	visual_bell = false;
	show_column_headers = true;
	show_window_title = true;
	set_column_headers("artist track title album year length");

	/* Set up options array */
	add_option("host", OPTION_TYPE_STRING, (void *)&host);
	add_option("port", OPTION_TYPE_STRING, (void *)&port);
	add_option("password", OPTION_TYPE_STRING, (void *)&password);

	add_option("reconnectdelay", OPTION_TYPE_UINT, (void *)&reconnect_delay);
	add_option("bell", OPTION_TYPE_BOOL, (void *)&use_bell);
	add_option("visualbell", OPTION_TYPE_BOOL, (void *)&visual_bell);
	add_option("columnheaders", OPTION_TYPE_BOOL, (void *)&show_column_headers);
	add_option("windowtitle", OPTION_TYPE_BOOL, (void *)&show_window_title);

	add_option("columns", OPTION_TYPE_COLUMNHEADERS, (void *)&songlist_columns);
}

int Config::readline(string line)
{
	string optstr;
	string optval = "";
	option_t * opt;
	size_t pos;
	bool invert = false;
	bool show = false;
	bool negative = false;
	bool * bopt;
	int result;

	/* Locate the identifier */
	if (line.size() == 0)
	{
		return print_all_options();
	}
	else if ((pos = line.find('=')) != string::npos)
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
			optstr = optstr.substr(0, optstr.size() - 1);
			break;
		case '!':
			invert = true;
			optstr = optstr.substr(0, optstr.size() - 1);
			break;
		default:
			break;
	}

	/* Return the option struct if this is a valid option */
	if ((opt = get_opt_ptr(optstr)) == NULL)
	{
		/* Check if this is a negative boolean (no<option>) */
		if (optstr.size() > 2 && optstr.substr(0, 2) == "no" && ((opt = get_opt_ptr(optstr.substr(2))) != NULL) && opt->type == OPTION_TYPE_BOOL)
		{
			negative = !invert;
			optstr = optstr.substr(2);
		}
		else
		{
			sterr("Unknown option: %s", line.c_str());
			return false;
		}
	}

	/* Print the option to statusbar */
	if (show)
	{
		print_option(opt);
		return true;
	}

	/* Check for (negative) boolean options */
	if (optval.size() == 0 && pos == string::npos)
	{
		/* Show option instead */
		if (opt->type != OPTION_TYPE_BOOL)
		{
			print_option(opt);
			return true;
		}
		bopt = (bool *)opt->ptr;
		*bopt = !negative;
		print_option(opt);
		return true;
	}

	/* Invert an option if boolean */
	if (invert)
	{
		if (opt->type != OPTION_TYPE_BOOL)
		{
			stinfo("%s=%s", optstr.c_str(), get_opt_str(opt).c_str());
			sterr("Trailing characters: %s", line.c_str());
			return false;
		}
		bopt = (bool *)opt->ptr;
		*bopt = !(*bopt);
		print_option(opt);
		return true;
	}

	/* Set the new string value */
	result = set_opt_str(opt, optval);
	if (result)
		print_option(opt);

	return result;
}

option_t * Config::add_option(string name, option_type_t type, void * ptr)
{
	option_t * o = new option_t;
	o->name = name;
	o->type = type;
	o->ptr = ptr;
	options.push_back(o);
	return o;
}

string Config::get_opt_str(option_t * opt)
{
	vector<Field *>::const_iterator field_it;

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
			ui = (unsigned int *)opt->ptr;
			str = tostring(*ui);
			break;

		case OPTION_TYPE_INT:
			i = (int *)opt->ptr;
			str = tostring(*i);
			break;

		/* Exotic data types */

		case OPTION_TYPE_COLUMNHEADERS:
			for (field_it = songlist_columns.begin(); field_it != songlist_columns.end(); ++field_it)
				str = str + (*field_it)->str + " ";
			str = str.substr(0, str.size() - 1);
			break;

		default:
			str = "<unknown>";
			break;
	}

	return str;
}

int Config::set_opt_str(option_t * opt, string value)
{
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

		case OPTION_TYPE_COLUMNHEADERS:
			set_column_headers(value);
			wm.update_column_length();
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

unsigned int Config::grep_opt(string opt, vector<option_t *> * list, bool * negate)
{
	vector<option_t *>::const_iterator i;

	if (!list) return 0;
	list->clear();

	*negate = false;
	if (opt.size() >= 2 && opt.substr(0, 2) == "no")
	{
		if (opt.size() == 2)
			opt.clear();
		else
			opt = opt.substr(2);
		*negate = true;
	}

	for (i = options.begin(); i != options.end(); i++)
	{
		if (opt.size() > (*i)->name.size())
			continue;

		if (opt == (*i)->name.substr(0, opt.size()))
		{
			if (!(*negate) || (*i)->type == OPTION_TYPE_BOOL || ((*i)->name.size() > 2 && (*i)->name.substr(0, 2) == "no"))
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
	else
		debug("  %s=%s", opt->name.c_str(), get_opt_str(opt).c_str());
}

int Config::print_all_options()
{
	vector<option_t *>::const_iterator i;

	debug("--- Options ---", NULL);

	for (i = options.begin(); i != options.end(); ++i)
		print_option(*i);

	return true;
}

void Config::set_column_headers(string hdr)
{
	size_t start = 0;
	size_t pos;
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

		if ((field = fieldtypes.find(f)) != NULL)
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

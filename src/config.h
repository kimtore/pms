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
 * config.h
 * 	configuration file parser
 */


#ifndef _CONFIG_H_
#define _CONFIG_H_

#include <sys/types.h>
#include <fcntl.h>
#include <string>
#include "settings.h"
#include "input.h"
#include "display.h"
#include "message.h"


typedef enum
{
	KW_NONE,
	KW_ERR,
	KW_BIND,
	KW_UNBIND,
	KW_SET,
	KW_COLOR

} pms_config_keyword;

/*
 * Holds a variety of information about each field type
 */
class Fieldtypes
{
public:
	vector<string>				name;
	vector<string>				header;
	vector<unsigned int>			minlen;
	vector<Item>				type;
	vector<bool (*) (Song *, Song *)>	sortfunc;

	bool			add(string, string, Item, unsigned int, bool (*) (Song *, Song *));
	int			lookup(string);
};


/* Holds information about which command belonging to which action */
class Commandmap
{
private:
	vector<string>			command;
	vector<string>			description;
	vector<pms_pending_keys>	action;
public:
	Commandmap() {};
	~Commandmap() {};

	bool			add(string, string, pms_pending_keys);
	pms_pending_keys	act(string);
	string			desc(string);
};


class Bindings
{
private:
	vector<string>			strkey;
	vector<int>			key;
	vector<string>			straction;
	vector<pms_pending_keys>	action;
	vector<string>			param;
	Commandmap *			cmap;
public:
	Bindings(Commandmap * c) { cmap = c; };
	~Bindings() {};

	bool			add(string, string);
	bool			remove(string);
	pms_pending_keys	act(int, string *);
	unsigned int		size() { return key.size(); };
	unsigned int		list(vector<string> *, vector<string> *, vector<string> *);
	void			clear();
};


class Configurator
{
private:
	Options *			opt;
	Bindings *			bindings;

	bool				set_color(string, string);


public:
					Configurator(Options *, Bindings *);
					~Configurator() {};

	/* Public static functions */
	static bool			is_whitespace(char);
	static bool			strtobool(string);
	static vector<string> *		splitline(string);
	static string			getparamopt(string);
	static bool			verify_columns(string);
	
	/* Public members */
	bool				source(string);
	bool				readline(string);
	bool				loadconfigs();
};

#endif

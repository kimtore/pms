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
 * settings.h - configuration option class
 *
 */


#ifndef _SETTINGS_H_
#define _SETTINGS_H_


#include <vector>
#include <string>
#include "color.h"
#include "topbar.h"
#include "message.h"

using namespace std;


/*
 * Types of values
 */
typedef enum
{
	SETTING_TYPE_EINVAL,
	SETTING_TYPE_ALIAS,
	SETTING_TYPE_STRING,
	SETTING_TYPE_LONG,
	SETTING_TYPE_BOOLEAN,
	SETTING_TYPE_FIELDLIST,
	SETTING_TYPE_SCROLL,
	SETTING_TYPE_PLAYMODE,
	SETTING_TYPE_REPEATMODE
}
SettingType;


/*
 * One explicit setting
 */
class Setting
{
public:
	string			v_string;
	long			v_long;
	bool			v_bool;

				Setting();

	Setting *		alias;
	SettingType		type;
	string			key;
};


/*
 * Array of all settings
 */
class Options
{
private:
	vector<Setting *>	vals;

	Setting *		lookup(string);
	Setting *		add(string, SettingType);

	void			destroy();

public:
	/* These are special settings that can't be contained in a Setting class */
	vector<Topbarline *>	topbar;			// Topbar draw information
	Colortable *		colors;

	/* And special option setters */
	bool			set_topbar_values(string, string);
	void			clear_topbar();

				Options();
				~Options();

	/* Add an alias to another option */
	bool			alias(string, string);

	/* Set value - perform some types of error checking */
	bool			set(string, string);
	Setting *		set(string, SettingType, string);
	Setting *		set_string(string, string);
	Setting *		set_long(string, long);
	Setting *		set_bool(string, bool);
	bool			toggle(string);

	/* Returns setting type */
	SettingType		get_type(string);

	/* Returns the setting itself */
	string			get_string(string);
	long 			get_long(string);
	bool 			get_bool(string);

	/* Dump everything into a long string */
	bool			dump(string);
	string			dump(Setting *);
	string			dump_all();

	/* Reset everything to defaults */
	void			reset();
};


#endif /* _SETTINGS_H_ */

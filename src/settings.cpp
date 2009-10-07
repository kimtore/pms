/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2009  Kim Tore Jensen
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


#include "settings.h"

using namespace std;


/*
 * Constructors and destructors
 */
Setting::Setting()
{
	alias = NULL;
	type = SETTING_TYPE_EINVAL;
	key = "";
}

Options::Options()
{
	colors = NULL;
	reset();
}

Options::~Options()
{
	vector<Topbarline *>::iterator	i;

	if (colors != NULL)
		delete colors;

	i = topbar.begin();
	while (i++ != topbar.end())
		delete *i;
	topbar.clear();
}


/*
 * Find a setting based on keyword
 */
Setting *	Options::lookup(string key)
{
	unsigned int	i;

	for (i = 0; i < vals.size(); i++)
	{
		if (vals[i]->key == key)
			return vals[i];
	}

	return NULL;
}


/*
 * Initialize a setting or return an existing one with type t.
 * Returns NULL if the setting exists but has a different SettingType,
 * else return the pointer to the new object.
 */
Setting *	Options::add(string key, SettingType t)
{
	Setting *	s;

	s = lookup(key);
	if (s != NULL)
	{
		if (s->type != t)
			return NULL;
		return s;
	}

	s = new Setting();
	if (s == NULL)
		return s;
	
	s->key = key;
	s->type = t;

	vals.push_back(s);

	return s;
}


/*
 * Alias a keyword to point to another keyword.
 * Limitations: needs to be added _after_ the original keyword.
 * Can be nested indefinately.
 */
bool		Options::alias(string key, string dest)
{
	Setting *	s_key, s_dest;

	s_dest = lookup(dest);
	if (s_dest == NULL)
		return false;
	
	s_key = add(key, SETTING_TYPE_ALIAS);
	if (s_key == NULL)
		return false;
	
	s_key->alias = s_dest;

	return true;
}


/*
 * Set a special kind of variable
 */
bool		Options::set(string key, SettingType t, string val)
{
	Setting *	s;

	s = add(key, t);
	if (s == NULL)
		return false;

	switch(t)
	{
		SETTING_TYPE_STRING:
			return set_string(key, val);
		SETTING_TYPE_LONG:
			return set_long(key, atoi(val.c_str()));
		SETTING_TYPE_BOOLEAN:
			return set_bool(key, Configurator::strtobool(val));
		default:
			s->v_string = val;
			s->v_long = 0;
			s->v_bool = 0;
	}

	return true;
}

/*
 * Set a string variable
 */
bool		Options::set_string(string key, string val)
{
	Setting *	s;

	s = add(key, SETTING_TYPE_STRING);
	if (s == NULL)
		return false;

	s->v_string = val;
	s->v_long = 0;
	s->v_bool = 0;

	return true;
}

/*
 * Set a numeric variable
 */
bool		Options::set_long(string key, long val)
{
	Setting *	s;

	s = add(key, SETTING_TYPE_LONG);
	if (s == NULL)
		return false;

	s->v_long = val;
	s->v_string = "";
	s->v_bool = 0;

	return true;
}

/*
 * Set a boolean variable
 */
bool		Options::set_bool(string key, bool val)
{
	Setting *	s;

	s = add(key, SETTING_TYPE_BOOLEAN);
	if (s == NULL)
		return false;

	s->v_bool = val;
	s->v_string = "";
	s->v_long = 0;

	return true;
}


/*
 * Read functions
 */
string		Options::get_string(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return NULL;

	while (s->alias != NULL)
		s = s->alias;
	
	return s->v_string;
}

long		Options::get_long(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return NULL;
	
	while (s->alias != NULL)
		s = s->alias;
	
	return s->v_long;
}

bool		Options::get_bool(string key)
{
	Setting *	s;

	s = lookup(key);
	if (s == NULL)
		return NULL;
	
	while (s->alias != NULL)
		s = s->alias;
	
	return s->v_bool;
}


/*
 * Dump all options to a long string
 */
string		Options::dump_all()
{
	string		output = "";

	unsigned int	i;
	Setting *	s;

	for (i = 0; i < vals.size(); i++)
	{
		s = vals[i];
		output += "set ";
		output += s;
		output += "=";
		switch(s->type)
		{
			default:
			case SETTING_TYPE_STRING:
				output += s->v_string;
				break;

			case SETTING_TYPE_LONG:
				output += Pms::tostring(s->v_long);
				break;

			case SETTING_TYPE_BOOLEAN:
				output += (s->v_bool ? "true" : "false");
				break;
		}
		output += "\n";
	}

	return output;
}


/*
 * Reset to defaults
 */
void		Options::reset()
{
	vector<Setting *>::iterator	i;
	vector<Topbarline*>::iterator	j;

	/* Truncate old settings array */
	i = vals.begin();
	while (i++ != vals.end())
		delete *i;
	vals.clear();

	/* Truncate topbar */
	j = topbar.begin();
	while (j++ != topbar.end())
		delete *j;
	topbar.clear();

	if (colors != NULL)
		delete colors;

	set("scroll", SETTING_TYPE_SCROLL, "normal");
	set("playmode", SETTING_TYPE_PLAYMODE, "linear");
	set("repeatmode", SETTING_TYPE_REPEATMODE, "none");
	set("columns", SETTING_TYPE_FIELDLIST, "artist track title album length");

	set_long("nextinterval", 5);
	set_long("crossfade", 5);
	set_long("mpd_timeout", 30);
	set_long("repeatonedelay", 1);
	set_long("stopdelay", 1);
	set_long("reconnectdelay", 30);
	set_long("directoryminlen", 30);
	set_long("resetstatus", 3);
	set_long("scrolloff", 0);

	set_bool("debug", false);
	set_bool("addtoreturns", false);
	set_bool("ignorecase", true);
	set_bool("regexsearch", false);
	set_bool("followwindow", false);
	set_bool("followcursor", false);
	set_bool("followplayback", false);
	set_bool("nextafteraction", true);
	set_bool("showtopbar", true);
	set_bool("topbarborders", false);
	set_bool("topbarspace", true);
	set_bool("columnspace", true);
	set_bool("mouse", false);

	set_string("directoryformat", "%artist% - %title%");
	set_string("xtermtitle", "PMS: %ifplaying% %artist% - %title% %else% Not playing %endif%");
	set_string("onplaylistfinish", "");
	set_string("libraryroot", "");
	set_string("startuplist", "playlist");
	set_string("librarysort", "default");
	set_string("albumclass", "artist album date"); //FIXME: implement this

	set_string("status_unknown", "??");
	set_string("status_play", "|>");
	set_string("status_pause", "||");
	set_string("status_stop", "[]");

	/*
	 * Set up option aliases
	 */
	alias("ic", "ignorecase");
	alias("so", "scrolloff");

	//TODO: would be nice to have the commented alteratives default if 
	//Unicode is available
	//status_unknown		= "??"; //?
	//status_play		= "|>"; //▶
	//status_pause		= "||"; //‖
	//status_stop		= "[]"; //■
	
	/* Set up default top bar values */
	topbar.clear();
	while(topbar.size() < 3)
		topbar.push_back(new Topbarline());

	topbar[0]->strings[0] = _("%time_elapsed% %playstate% %time%%ifcursong% (%progresspercentage%%%)%endif%");
	topbar[0]->strings[1] = _("%ifcursong%%artist%%endif%");
	topbar[0]->strings[2] = _("Vol: %volume%%%  Mode: %muteshort%%repeatshort%%randomshort%%manualshort%");
	topbar[1]->strings[1] = _("%ifcursong%==> %title% <==%else%No current song%endif%");
	topbar[2]->strings[0] = _("%listsize%");
	topbar[2]->strings[1] = _("%ifcursong%%album% (%year%)%endif%");
	topbar[2]->strings[2] = _("Q: %livequeuesize%");
}

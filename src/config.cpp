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
 * config.cpp - configuration parser
 *
 */


#include <string>
#include "mycurses.h"
#include "config.h"
#include "pms.h"

using namespace std;

extern Pms *		pms;

bool			Fieldtypes::add(string nname, string nheader, Item ntype, unsigned int nminlen, bool (*nsortfunc) (Song *, Song *))
{
	if (nname.size() == 0)
		return false;

	name.push_back(nname);
	header.push_back(nheader);
	type.push_back(ntype);
	minlen.push_back(nminlen);
	sortfunc.push_back(nsortfunc);

	return true;
}

int			Fieldtypes::lookup(string s)
{
	unsigned int	i;

	for (i = 0; i < name.size(); i++)
	{
		if (name[i] == s)
			return (int)i;
	}

	return -1;
}

bool			Commandmap::add(string com, string des, pms_pending_keys act)
{
	if (com.size() == 0 && act != PEND_NONE)
		return false;
	command.push_back(com);
	description.push_back(des);
	action.push_back(act);

	return true;
}

pms_pending_keys	Commandmap::act(string key)
{
	unsigned int		i;

	for (i = 0; i < command.size(); i++)
	{
		if (command[i] == key)
		{
			return action[i];
		}
	}

	return PEND_NONE;
}

string			Commandmap::desc(string com)
{
	unsigned int		i;

	for (i = 0; i < command.size(); i++)
	{
		if (command[i] == com)
		{
			return description[i];
		}
	}

	return "";
}

/*
 * Deletes all key bindings
 */
void			Bindings::clear()
{
	key.clear();
	action.clear();
	param.clear();
	straction.clear();
	strkey.clear();
}

/*
 * Remove a key binding
 */
bool			Bindings::remove(string b)
{
	unsigned int		i;
	bool			m = false;

	if (b == "all")
	{
		clear();
		return true;
	}

	for (i = 0; i < strkey.size(); ++i)
	{
		if (i >= strkey.size()) break;
		if (b == strkey[i])
		{
			strkey.erase(strkey.begin() + i);
			key.erase(key.begin() + i);
			straction.erase(straction.begin() + i);
			action.erase(action.begin() + i);
			param.erase(param.begin() + i);
			--i;
			m = true;
		}
	}

	return m;
}

/*
 * Associate a key with a binding
 */
bool			Bindings::add(string b, string command)
{
	string			par = "";
	pms_pending_keys	k;
	size_t			l;
	int			i = 0;

	/* Find parameter */
	l = command.find(" ");
	if (l != string::npos && l > 0)
	{
		par = command.substr(l + 1);
		command = command.substr(0, l);
	}

	k = cmap->act(command);
	if (k == PEND_NONE)
	{
		if (command[0] == '!' && command.size() > 1)
		{
			par = command.substr(1) + " " + par;
			command = "!";
			k = cmap->act(command);
		}
		else
		{
			pms->msg->clear();
			pms->msg->code = CERR_INVALID_COMMAND;
			pms->msg->str = _("invalid command");
			pms->msg->str += " '" + command + "'";
			return false;
		}
	}

	/* Simple character */
	if (b.size() == 1)
	{
		i = b[0];
	}
	/* String -> character */
	else
	{
		/* Control + character. This is always lowercase. */
		if (b[0] == '^' && (b[1] >= 'A' && b[1] <= 'Z'))
		{
			i = b[1] - 64;
		}
		else
		{
			if (b == "up")
				i = KEY_UP;
			else if (b == "down")
				i = KEY_DOWN;
			else if (b == "left")
				i = KEY_LEFT;
			else if (b == "right")
				i = KEY_RIGHT;
			else if (b == "pageup")
				i = KEY_PPAGE;
			else if (b == "pagedown")
				i = KEY_NPAGE;
			else if (b == "home")
				i = KEY_HOME;
			else if (b == "end")
				i = KEY_END;
			else if (b == "backspace")
				i = KEY_BACKSPACE;
			else if (b == "delete")
				i = KEY_DC;
			else if (b == "insert")
				i = KEY_IC;
			else if (b == "return")
				i = 10;
			else if (b == "kpenter")
				i = 343;
			else if (b == "space")
				i = ' ';
			else if (b == "tab")
				i = '\t';
			else if (b == "mouse1")
				i = BUTTON1_CLICKED;
			else if (b == "mouse2")
				i = BUTTON2_CLICKED;
			else if (b == "mouse3")
				i = BUTTON3_CLICKED;
			else if (b == "mouse4")
				i = BUTTON4_CLICKED;
#if NCURSES_MOUSE_VERSION > 1
			else if (b == "mouse5")
				i = BUTTON5_CLICKED;
#endif
			else if ((b[0] == 'F' || b[0] == 'f') && b.size() > 1)
			{
				i = atoi(string(b.substr(1)).c_str());
				if (i > 0 && i < 64)
					i = KEY_F(i);
				else
				{
					pms->msg->clear();
					pms->msg->code = CERR_INVALID_KEY;
					pms->msg->str = _("function key out of range");
					pms->msg->str += ": " + b;
					return false;
				}
			}
			else
			{
				pms->msg->clear();
				pms->msg->code = CERR_INVALID_KEY;
				pms->msg->str = _("invalid key");
				pms->msg->str += " '" + b + "'";
				return false;
			}
		}
	}

	/* Remove any old bind for this key */
	remove(b);

	strkey.push_back(b);
	key.push_back(i);
	action.push_back(k);
	straction.push_back(command);
	param.push_back(par);
	/*
	if (pms->options)
	{
		pms->log(MSG_DEBUG, 0, "Mapping key %3d to action %d '%s' with parameter '%s'\n", i, k, command.c_str(), par.c_str());
	}
	*/

	return true;
}

pms_pending_keys	Bindings::act(int k, string * parm)
{
	unsigned int		i;

	/* Standardize ambiguous keys */
	if (k == 8 || k == 127)
		k = KEY_BACKSPACE;

	for (i = 0; i < key.size(); i++)
	{
		if (key[i] == k)
		{
			*parm = param[i];
			return action[i];
		}
	}

	return PEND_NONE;
}

/*
 * Creates a list of all key mappings
 */
unsigned int		Bindings::list(vector<string> * k, vector<string> * com, vector<string> * desc)
{
	unsigned int		i;

	if (!k || !com || !desc)
		return 0;

	k->clear();
	com->clear();
	desc->clear();

	for (i = 0; i < key.size(); i++)
	{
		k->push_back(strkey[i]);
		com->push_back(straction[i] + " " + param[i]);
		desc->push_back(cmap->desc(straction[i]));
	}

	return i;
}






/*
 * Tells whether a character is whitespace or not
 */
bool			Configurator::is_whitespace(char n)
{
	return (n == ' ' || n == '\t' || n == '\0' || n == '\n' ? true : false);
}

/*
 * Converts a string into a boolean value
 */
bool			Configurator::strtobool(string s)
{
	//convert to lowercase
	transform(s.begin(), s.end(), s.begin(), ::tolower);

	return s == "yes" || s == "true" || s == "on" || s == "1";
}

/*
 * Verify that a columns string is OK
 */
bool			Configurator::verify_columns(string s)
{
	unsigned int		i;
	vector<string> *	v;

	if (s.size() == 0)
		return false;

	v = Pms::splitstr(s, " ");

	for (i = 0; i < v->size(); i++)
	{
		if (pms->fieldtypes->lookup((*v)[i]) == -1)
		{
			pms->msg->clear();
			pms->msg->code = CERR_INVALID_COLUMN;
			pms->msg->str = _("invalid column type");
			pms->msg->str += " '" + (*v)[i] + "'";
			delete v;
			return false;
		}
	}

	delete v;

	return true;
}





/*
 * Constructor
 */
Configurator::Configurator(Options * o, Bindings * b)
{
	opt = o;
	bindings = b;
}

/*
 * Loads a configuration file
 */
bool			Configurator::source(string fn)
{
	FILE *		fd;
	char		buffer[1024];
	int		line = 0;

	pms->msg->clear();
	pms->msg->code = CERR_NONE;
	fd = fopen(fn.c_str(), "r");

	if (fd == NULL)
	{
		pms->msg->code = CERR_NO_FILE;
		pms->msg->str = fn + _(": could not open file.\n");
		return false;
	}

	pms->log(MSG_CONSOLE, 0, _("Reading configuration file %s\n"), fn.c_str());

	while (fgets(buffer, 1024, fd) != NULL)
	{
		++line;
		if (!readline(buffer))
		{
			pms->log(MSG_CONSOLE, STERR, _("Encountered an error on line %d.\n"), line);
			break;
		}
	}

	if (pms->msg->code != 0)
	{
		pms->msg->str = "line " + Pms::tostring(line) + ": " + pms->msg->str;
	}

	pms->log(MSG_CONSOLE, 0, _("Finished reading configuration file.\n"));

	fclose(fd);

	return (pms->msg->code == 0);
}

/*
 * Splits a line into segments
 */
vector<string> *	Configurator::splitline(string line)
{
	string::iterator	i;
	vector<string> *	v;
	string			buf = "";

	v = new vector<string>;

	i = line.begin();
	while (i != line.end())
	{
		if (Configurator::is_whitespace(*i) || *i == '=' || *i == ':')
		{
			if (buf.size() > 0)
			{
				v->push_back(buf);
				buf.clear();
			}
			if (*i == '=')
				v->push_back("=");
			else if (*i == ':')
				v->push_back(":");
		}
		else
		{
			buf += *i;
		}
		++i;
	}

	if (buf.size() > 0)
		v->push_back(buf);

	return v;
}

/*
 * Gets parameter from a command line string
 */
string			Configurator::getparamopt(string buffer)
{
	size_t				n;
	size_t				epos;
	size_t				cpos;

	//set n to the first position of = or :, return empty string if neither 
	//is found
	epos = buffer.find_first_of("=");
	cpos = buffer.find_first_of(":");
	if (epos == string::npos && cpos == string::npos)
		return "";
	else if (cpos == string::npos || epos < cpos)
		n = epos;
	else
		n = cpos;
	if (n == string::npos || n == buffer.size() - 1)
		return "";

	buffer = buffer.substr(n + 1);
	while (buffer.size() > 0 && Configurator::is_whitespace(buffer[buffer.size()-1]))
		buffer = buffer.substr(0, buffer.size() - 1);

	return buffer;
}

/*
 * Interprets a command line
 */
bool			Configurator::readline(string buffer)
{
	vector<string> *		tok;
	vector<string>::iterator	it;
	bool				state = false;
	string				proc;
	string				val;

	/* No errors by default */
	pms->msg->clear();

	/* Empty lines pass through */
	if (buffer.size() == 0)
		return true;

	/* Split into tokens delimited by whitespace */
	tok = Configurator::splitline(buffer);
	if (tok->size() == 0)
	{
		delete tok;
		return true;
	}
	
	/* Comments start with '#' */
	if ((*tok)[0][0] == '#')
		return true;

	proc = (*tok)[0];

	/* Process first keyword */
	if (proc == "set" || proc == "se")
	{
		proc.clear();
		val.clear();

		if (tok->size() < 2)
			pms->msg->code = CERR_MISSING_IDENTIFIER;
		else
		{
			proc = (*tok)[1];
			if (tok->size() == 2)
			{
				//check for various prefixes/suffixes
				if (pms->options->get_type(proc) == SETTING_TYPE_BOOLEAN)
					return pms->options->set(proc, "true");
				else if (proc.substr(proc.length() - 1, 1) == "?" && pms->options->get_type(proc.substr(0, proc.length() - 1)) != SETTING_TYPE_EINVAL)
				{
					pms->options->dump(proc.substr(0, proc.length() - 1));
					return false;
				}
				else if (proc.substr(0, 2) == "no" && pms->options->get_type(proc.substr(2)) == SETTING_TYPE_BOOLEAN)
					return pms->options->set(proc.substr(2), "false");
				else if (proc.substr(0, 3) == "inv" && pms->options->get_type(proc.substr(3)) == SETTING_TYPE_BOOLEAN)
					return pms->options->toggle(proc.substr(3));
				else if (proc.substr(proc.length() - 1, 1) == "!" && pms->options->get_type(proc.substr(0, proc.length() - 1)) == SETTING_TYPE_BOOLEAN)
					return pms->options->toggle(proc.substr(0, proc.length() - 1));
				else if (pms->options->get_type(proc) == SETTING_TYPE_EINVAL)
					pms->msg->code = CERR_INVALID_IDENTIFIER;
				else
				{
					pms->options->dump(proc);
					return false;
				}
			}
			else if (pms->options->get_type(proc) == SETTING_TYPE_BOOLEAN || tok->at(2) != "=" && tok->at(2) != ":")
				pms->msg->code = CERR_UNEXPECTED_TOKEN;
		}

		if (pms->msg->code == CERR_NONE)
			val = Configurator::getparamopt(buffer);

		delete tok;

		switch(pms->msg->code)
		{
			case CERR_NONE:
				return pms->options->set(proc, val);
			case CERR_INVALID_IDENTIFIER:
				pms->msg->str = _("invalid identifier");
				pms->msg->str += " '" + proc + "'";
				return false;
			case CERR_MISSING_IDENTIFIER:
				pms->msg->str = _("missing name identifier after 'set'");
				return false;
			case CERR_MISSING_VALUE:
				pms->msg->str = _("missing value for configuration option");
				pms->msg->str += " '" + proc + "'";
				return false;
			case CERR_UNEXPECTED_TOKEN:
				pms->msg->str = _("unexpected token after identifier");
				return false;
			default:
				return false;
		}
	}
	else if (proc == "bind" || proc == "map")
	{
		proc.clear();
		val.clear();

		if (tok->size() == 1)
		{
			pms->msg->code = CERR_MISSING_IDENTIFIER;
			pms->msg->str = _("missing key after 'bind'");
			return false;
		}

		if (tok->size() == 2)
		{
			pms->msg->code = CERR_MISSING_VALUE;
			pms->msg->str = _("missing command to bind to key");
			pms->msg->str += " '" + tok->at(1) + "'";
			return false;
		}

		proc = tok->at(1);
		val = Pms::joinstr(tok, tok->begin() + 2, tok->end());
		delete tok;

		return bindings->add(proc, val);
	}
	else if (proc == "unbind" || proc == "unmap" || proc == "unm")
	{
		it = tok->begin() + 1;
		while (it != tok->end())
		{
			if (!bindings->remove(*it))
			{
				pms->msg->code = CERR_INVALID_KEY;
				pms->msg->str = _("Can't remove binding for key");
				pms->msg->str += " '" + *it + "'";
				delete tok;
				return false;
			}
			++it;
		}
		delete tok;
	}
	else if (proc == "color" || proc == "colour")
	{
		proc.clear();
		val.clear();

		if (tok->size() == 1)
		{
			pms->msg->code = CERR_MISSING_IDENTIFIER;
			pms->msg->str = _("missing names after 'color'");
			return false;
		}

		if (tok->size() == 2)
		{
			pms->msg->code = CERR_MISSING_VALUE;
			pms->msg->str = _("missing colors to add to ");
			pms->msg->str += tok->at(1);
			return false;
		}

		proc = tok->at(1);
		val = Pms::joinstr(tok, tok->begin() + 2, tok->end());
		delete tok;

		return set_color(proc, val);
	}
	else
	{
		pms->msg->code = CERR_SYNTAX;
		pms->msg->str = _("syntax error: unexpected");
		pms->msg->str += " '" + proc + "'";
		return false;
	}

	return true;
}

/*
 * Load all configuration files
 */
bool			Configurator::loadconfigs()
{
	string			str;

	char *			homedir;
	char *			xdgconfighome;
	char *			xdgconfigdirs;
	string			next;
	vector<string>		configfiles;

	homedir = getenv("HOME");
	xdgconfighome = getenv("XDG_CONFIG_HOME");
	xdgconfigdirs = getenv("XDG_CONFIG_DIRS");

	/* Make a list of possible configuration files */
	// commandline argument
	if (pms->options->get_string("configfile").length() > 0)
		configfiles.push_back(pms->options->get_string("configfile"));

	// XDG config home (usually $HOME/.config)
	if (xdgconfighome == NULL || strlen(xdgconfighome) == 0)
	{
		if (homedir != NULL && strlen(homedir) > 0)
		{
			str = homedir;
			configfiles.push_back(str + "/.config/pms/rc");
		}
	}
	else
	{
		str = xdgconfighome;
		configfiles.push_back(str + "/pms/rc");
	}
	// XDG config dirs (colon-separated priority list, defaults to just /etc/xdg)
	if (xdgconfigdirs == NULL || strlen(xdgconfigdirs) == 0)
	{
		configfiles.push_back("/usr/local/etc/xdg/pms/rc");
		configfiles.push_back("/etc/xdg/pms/rc");
	}
	else
	{
		next = "";
		str = xdgconfigdirs;
		for (string::const_iterator it = str.begin(); it != str.end(); it++)
		{
			if (*it == ':')
			{
				if (next.length() > 0)
				{
					configfiles.push_back(next + "/pms/rc");
					next = "";
				}
			}
			else
				next += *it;
		}
		if (next.length() > 0)
			configfiles.push_back(next + "/pms/rc");
	}

	/* Load configuration files in reverse order */
	for (int i = configfiles.size() - 1; i >= 0; i--)
	{
		if (!source(configfiles[i]))
		{
			if (pms->msg->code != CERR_NO_FILE)
			{
				pms->log(MSG_CONSOLE, 0, _("\nConfiguration error in file %s:\n%s\n"), configfiles[i].c_str(), pms->msg->str.c_str());
				return false;
			}
			pms->log(MSG_CONSOLE, 0, _("Didn't find configuration file %s\n"), configfiles[i].c_str());
		}
	}

	return true;
}

/*
 * Set a color pair for a field
 */
bool			Configurator::set_color(string name, string pairs)
{
	vector<string> *	pair;
	string 			str;
	color *			dest;
	colortable_fields *	field;
	Colortable *		c;
	int			colors[2];
	int			attr = 0;
	unsigned int		cur = 0;
	size_t			found;

	if (pairs.size() == 0) return false;
	c = opt->colors;

	pms->msg->clear();

	/* Standard colors */
	if (name == "background")
		dest = (c->back);
	else if (name == "foreground")
		dest = (c->standard);
	else if (name == "statusbar")
		dest = (c->status);
	else if (name == "error")
		dest = (c->status_error);
	else if (name == "borders")
		dest = (c->border);
	else if (name == "headers")
		dest = (c->headers);
	else if (name == "title")
		dest = (c->title);

	/* Topbar statuses */
	else if (name == "topbar.time_elapsed")
		dest = (c->topbar.time_elapsed);
	else if (name == "topbar.time_remaining")
		dest = (c->topbar.time_remaining);
	else if (name == "topbar.playstate")
		dest = (c->topbar.playstate);
	else if (name == "topbar.progressbar")
		dest = (c->topbar.progressbar);
	else if (name == "topbar.progresspercentage")
		dest = (c->topbar.progresspercentage);
	else if (name == "topbar.librarysize")
		dest = (c->topbar.librarysize);
	else if (name == "topbar.listsize")
		dest = (c->topbar.listsize);
	else if (name == "topbar.queuesize")
		dest = (c->topbar.queuesize);
	else if (name == "topbar.livequeuesize")
		dest = (c->topbar.livequeuesize);
	else if (name == "topbar.foreground")
		dest = (c->topbar.standard);

	else if (name == "topbar.repeat")
		dest = (c->topbar.repeat);
	else if (name == "topbar.random")
		dest = (c->topbar.random);
	else if (name == "topbar.mute")
		dest = (c->topbar.mute);
	else if (name == "topbar.randomshort")
		dest = (c->topbar.randomshort);
	else if (name == "topbar.repeatshort")
		dest = (c->topbar.repeatshort);
	else if (name == "topbar.muteshort")
		dest = (c->topbar.randomshort);

	/* List colors */
	else if (name == "current")
		dest = (c->current);
	else if (name == "cursor")
		dest = (c->cursor);
	else if (name == "selection")
		dest = (c->selection);
	else if (name == "lastlist")
		dest = (c->lastlist);
	else if (name == "playinglist")
		dest = (c->playinglist);

	/* Fields for topbar and others */
	else if (name.size() > 7)
	{
		found = name.find("fields.");
		if (found != string::npos)
		{
			if (found == 0)
				field = &c->fields;
			else if (name.find("topbar.fields.") == 0)
				field = &c->topbar.fields;
			else
				field = NULL;

			if (field)
			{
				name = name.substr(found + 7);
				if (name == "file")
					dest = field->file;
				else if (name == "artist")
					dest = field->artist;
				else if (name == "artistsort")
					dest = field->artistsort;
				else if (name == "albumartist")
					dest = field->albumartist;
				else if (name == "albumartistsort")
					dest = field->albumartistsort;
				else if (name == "title")
					dest = field->title;
				else if (name == "album")
					dest = field->album;
				else if (name == "genre")
					dest = field->genre;
				else if (name == "track")
					dest = field->track;
				else if (name == "trackshort")
					dest = field->trackshort;
				else if (name == "time")
					dest = field->time;
				else if (name == "date")
					dest = field->date;
				else if (name == "year")
					dest = field->year;
				else if (name == "name")
					dest = field->name;
				else if (name == "composer")
					dest = field->composer;
				else if (name == "performer")
					dest = field->performer;
				else if (name == "disc")
					dest = field->disc;
				else if (name == "comment")
					dest = field->comment;
				else
					pms->msg->code = CERR_INVALID_IDENTIFIER;
			}
		}
		else
		{
			pms->msg->code = CERR_INVALID_IDENTIFIER;
		}
	}

	/* No valid color field */
	else
	{
		pms->msg->code = CERR_INVALID_IDENTIFIER;
	}

	if (pms->msg->code == CERR_INVALID_IDENTIFIER)
	{
		pms->msg->str = _("invalid identifier");
		pms->msg->str += " '" + name + "'";
		return false;
	}

	pair = Pms::splitstr(pairs);
	if (pair->size() > 2)
	{
		delete pair;
		pms->msg->code = CERR_EXCESS_ARGUMENTS;
		pms->msg->str = _("excess arguments: expected 1 or 2, got ");
		pms->msg->str += Pms::tostring(pair->size());
		return false;
	}

	for (cur = 0; cur < pair->size(); cur++)
	{
		str = (*pair)[cur];
		if (str == "black")
			colors[cur] = COLOR_BLACK;
		else if (str == "red")
			colors[cur] = COLOR_RED;
		else if (str == "green")
			colors[cur] = COLOR_GREEN;
		else if (str == "brown")
			colors[cur] = COLOR_YELLOW;
		else if (str == "blue")
			colors[cur] = COLOR_BLUE;
		else if (str == "magenta")
			colors[cur] = COLOR_MAGENTA;
		else if (str == "cyan")
			colors[cur] = COLOR_CYAN;
		else if (str == "brightgray" || str == "brightgrey" ||
				str == "lightgray" || str == "lightgrey" ||
				(str == "gray" || str == "grey") && cur == 1)
			colors[cur] = COLOR_WHITE;

		/* Back color only */
		else if (cur == 1)
		{
			if (str == "trans")
				colors[cur] = -1;
		}

		/* Front color only */
		else if (cur == 0)
		{
			if (str == "white")
			{
				colors[cur] = COLOR_WHITE;
				attr = A_BOLD;
			}
			else if (str == "gray" || str == "grey")
			{
				colors[cur] = COLOR_BLACK;
				attr = A_BOLD;
			}
			else if (str == "brightred" || str == "lightred")
			{
				colors[cur] = COLOR_RED;
				attr = A_BOLD;
			}
			else if (str == "brightgreen" || str == "lightgreen")
			{
				colors[cur] = COLOR_GREEN;
				attr = A_BOLD;
			}
			else if (str == "yellow")
			{
				colors[cur] = COLOR_YELLOW;
				attr = A_BOLD;
			}
			else if (str == "brightblue" || str == "lightblue")
			{
				colors[cur] = COLOR_BLUE;
				attr = A_BOLD;
			}
			else if (str == "brightmagenta" || str == "lightmagenta")
			{
				colors[cur] = COLOR_MAGENTA;
				attr = A_BOLD;
			}
			else if (str == "brightcyan" || str == "lightcyan")
			{
				colors[cur] = COLOR_CYAN;
				attr = A_BOLD;
			}
		}
		else
		{
			delete pair;
			pms->msg->clear();
			pms->msg->code = CERR_INVALID_COLOR;
			pms->msg->str = _("invalid color name");
			pms->msg->str += " '" + str + "'";
			return false;
		}
	}

	dest->set(colors[0], (pair->size() == 2 ? colors[1] : -1), attr);

	delete pair;

	return true;
}

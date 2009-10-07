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
bool			Bindings::add(string b, string command, Error & err)
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
			err.code = CERR_INVALID_COMMAND;
			err.str = _("invalid command");
			err.str += " '" + command + "'";
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
					err.code = CERR_INVALID_KEY;
					err.str = _("function key out of range");
					err.str += ": " + b;
					return false;
				}
			}
			else
			{
				err.code = CERR_INVALID_KEY;
				err.str = _("invalid key");
				err.str += " '" + b + "'";
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
		debug("Mapping key %3d to action %d '%s' with parameter '%s'\n", i, k, command.c_str(), par.c_str());
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
bool			Configurator::verify_columns(string s, Error & err)
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
			err.code = CERR_INVALID_COLUMN;
			err.str = _("invalid column type");
			err.str += " '" + (*v)[i] + "'";
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
bool			Configurator::source(string fn, Error & err)
{
	FILE *		fd;
	char		buffer[1024];
	int		line = 0;

	err.code = CERR_NONE;
	fd = fopen(fn.c_str(), "r");

	if (fd == NULL)
	{
		err.code = CERR_NO_FILE;
		err.str = fn + ": could not open file.";
		return false;
	}

	debug("Reading configuration file %s\n", fn.c_str());

	while (fgets(buffer, 1024, fd) != NULL)
	{
		++line;
		if (!readline(buffer, err))
			break;
	}

	if (err.code != 0)
	{
		err.str = "line " + Pms::tostring(line) + ": " + err.str;
	}

	debug("Finished reading configuration file.\n");

	fclose(fd);

	return (err.code == 0);
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
bool			Configurator::readline(string buffer, Error & err)
{
	vector<string> *		tok;
	vector<string>::iterator	it;
	bool				state = false;
	string				proc;
	string				val;

	/* No errors by default */
	err.clear();

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
			err.code = CERR_MISSING_IDENTIFIER;
		else
		{
			proc = (*tok)[1];
			if (tok->size() == 2)
			{
				//check for various prefixes/suffixes
				if (proc.substr(proc.length() - 1, 1) == "?" && get_opt_type(proc.substr(0, proc.length() - 1)) != OPT_NONE)
				{
					show_option(proc.substr(0, proc.length() - 1), err);
					return false;
				}
				else if (proc.substr(0, 2) == "no" && get_opt_type(proc.substr(2)) == OPT_BOOL)
					return set_option(proc.substr(2), "false", err);
				else if (proc.substr(0, 3) == "inv" && get_opt_type(proc.substr(3)) == OPT_BOOL)
					return toggle_option(proc.substr(3), err);
				else if (proc.substr(proc.length() - 1, 1) == "!" && get_opt_type(proc.substr(0, proc.length() - 1)) == OPT_BOOL)
					return toggle_option(proc.substr(0, proc.length() - 1), err);
				else if (get_opt_type(proc) == OPT_BOOL)
					return set_option(proc, "true", err);
				else if (get_opt_type(proc) == OPT_NONE)
					err.code = CERR_INVALID_IDENTIFIER;
				else
				{
					show_option(proc, err);
					return false;
				}
			}
			else if (get_opt_type(proc) == OPT_BOOL || tok->at(2) != "=" && tok->at(2) != ":")
				err.code = CERR_UNEXPECTED_TOKEN;
		}

		if (err.code == CERR_NONE)
			val = Configurator::getparamopt(buffer);

		delete tok;

		switch(err.code)
		{
			case CERR_NONE:
				break;
			case CERR_INVALID_IDENTIFIER:
				err.str = _("invalid identifier");
				err.str += " '" + proc + "'";
				return false;
			case CERR_MISSING_IDENTIFIER:
				err.str = _("missing name identifier after 'set'");
				return false;
			case CERR_MISSING_VALUE:
				err.str = _("missing value for configuration option");
				err.str += " '" + proc + "'";
				return false;
			case CERR_UNEXPECTED_TOKEN:
				err.str = _("unexpected token after identifier");
				return false;
			default:
				return false;
		}

		return set_option(proc, val, err);
	}
	else if (proc == "bind" || proc == "map")
	{
		proc.clear();
		val.clear();

		if (tok->size() == 1)
		{
			err.code = CERR_MISSING_IDENTIFIER;
			err.str = _("missing key after 'bind'");
			return false;
		}

		if (tok->size() == 2)
		{
			err.code = CERR_MISSING_VALUE;
			err.str = _("missing command to bind to key");
			err.str += " '" + tok->at(1) + "'";
			return false;
		}

		proc = tok->at(1);
		val = Pms::joinstr(tok, tok->begin() + 2, tok->end());
		delete tok;

		return bindings->add(proc, val, err);
	}
	else if (proc == "unbind" || proc == "unmap" || proc == "unm")
	{
		it = tok->begin() + 1;
		while (it != tok->end())
		{
			if (!bindings->remove(*it))
			{
				err.code = CERR_INVALID_KEY;
				err.str = _("Can't remove binding for key");
				err.str += " '" + *it + "'";
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
			err.code = CERR_MISSING_IDENTIFIER;
			err.str = _("missing names after 'color'");
			return false;
		}

		if (tok->size() == 2)
		{
			err.code = CERR_MISSING_VALUE;
			err.str = _("missing colors to add to ");
			err.str += tok->at(1);
			return false;
		}

		proc = tok->at(1);
		val = Pms::joinstr(tok, tok->begin() + 2, tok->end());
		delete tok;

		return set_color(proc, val, err);
	}
	else
	{
		err.code = CERR_SYNTAX;
		err.str = _("syntax error: unexpected");
		err.str += " '" + proc + "'";
		return false;
	}

	return true;
}

/*
 * Set a color pair for a field
 */
bool			Configurator::set_color(string name, string pairs, Error & err)
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
					err.code = CERR_INVALID_IDENTIFIER;
			}
		}
		else
		{
			err.code = CERR_INVALID_IDENTIFIER;
		}
	}

	/* No valid color field */
	else
	{
		err.code = CERR_INVALID_IDENTIFIER;
	}

	if (err.code == CERR_INVALID_IDENTIFIER)
	{
		err.str = _("invalid identifier");
		err.str += " '" + name + "'";
		return false;
	}

	pair = Pms::splitstr(pairs);
	if (pair->size() > 2)
	{
		delete pair;
		err.code = CERR_EXCESS_ARGUMENTS;
		err.str = _("excess arguments: expected 1 or 2, got ");
		err.str += Pms::tostring(pair->size());
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
			err.code = CERR_INVALID_COLOR;
			err.str = _("invalid color name");
			err.str += " '" + str + "'";
			return false;
		}
	}

	dest->set(colors[0], (pair->size() == 2 ? colors[1] : -1), attr);

	delete pair;

	return true;
}

/*
 * Configuration options.
 * Get an option pointer based on a name. Returns option type.
 */
int			Configurator::get_opt_ptr(string name, void *& dest)
{
	dest = NULL;

	/*
	 * Boolean values
	 */

	if (name == "debug")
		dest = &(opt->debug);
#ifdef HAVE_BOOST
	else if (name == "regexsearch")
		dest = &(opt->regexsearch);
#endif
	else if (name == "ignorecase" || name == "ic")
		dest = &(opt->ignorecase);
	else if (name == "followwindow")
		dest = &(opt->followwindow);
	else if (name == "followcursor")
		dest = &(opt->followcursor);
	else if (name == "followplayback")
		dest = &(opt->followplayback);
	else if (name == "nextafteraction")
		dest = &(opt->nextafteraction);
	else if (name == "topbarvisible")
		dest = &(opt->showtopbar);
	else if (name == "topbarborders")
		dest = &(opt->topbarborders);
	else if (name == "topbarspace")
		dest = &(opt->topbarspace);
	else if (name == "columnspace")
		dest = &(opt->columnspace);
	else if (name == "mouse")
		dest = &(opt->mouse);
	else if (name == "addtoreturns")
		dest = &(opt->addtoreturns);

	if (dest != NULL)
		return OPT_BOOL;

	/*
	 * Integer values
	 */

	if (name == "port")
		dest = &(opt->port);
	else if (name == "repeatonedelay")
		dest = &(opt->repeatonedelay);
	else if (name == "stopdelay")
		dest = &(opt->stopdelay);
	else if (name == "reconnectdelay")
		dest = &(opt->reconnectdelay);
	else if (name == "crossfade")
		dest = &(opt->crossfade);
	else if (name == "mpd_timeout")
		dest = &(opt->mpd_timeout);
	else if (name == "resetstatus")
		dest = &(opt->resetstatus);
	else if (name == "scrolloff" || name == "so")
		dest = &(opt->scrolloff);

	if (dest != NULL)
		return OPT_INT;

	/*
	 * String values
	 */

	if (name == "host")
		dest = &(opt->hostname);
	else if (name == "password")
		dest = &(opt->password);
	else if (name == "onplaylistfinish")
		dest = &(opt->onplaylistfinish);
	else if (name == "startuplist")
		dest = &(opt->startuplist);
	else if (name == "columns")
		dest = &(opt->columns);
	else if (name == "sort")
		dest = &(opt->librarysort);
	else if (name == "status_unknown")
		dest = &(opt->status_unknown);
	else if (name == "status_play")
		dest = &(opt->status_play);
	else if (name == "status_pause")
		dest = &(opt->status_pause);
	else if (name == "status_stop")
		dest = &(opt->status_stop);
	else if (name == "libraryroot")
		dest = &(opt->libraryroot);
	else if (name == "xtermtitle")
		dest = &(opt->xtermtitle);

	if (dest != NULL)
		return OPT_STRING;

	if (name == "follow")
		return OPT_SPECIAL;
	else if (name == "scroll")
		return OPT_SPECIAL;
	else if (name == "playmode")
		return OPT_SPECIAL;

	return OPT_NONE;
}

/*
 * Get an option's type
 */
int			Configurator::get_opt_type(string name)
{
	void *		ptr;
	return get_opt_ptr(name, ptr);
}

/*
 * Option getting
 */
string			Configurator::get_option(string name, Error & err)
{
	int		type;
	void *		ptr;

	type = get_opt_ptr(name, ptr);

	switch(type)
	{
		case OPT_BOOL:
			return (*(static_cast<bool *>( ptr )) ? _("on") : _("off"));
		case OPT_INT:
			return Pms::tostring(*(static_cast<int *>( ptr )));
		case OPT_STRING:
			return *(static_cast<string *>( ptr ));
		case OPT_SPECIAL:
			return _("<special>");
		default:
			err.code = CERR_INVALID_IDENTIFIER;
			return "";
	}
}

/*
 * Option setting
 */
bool			Configurator::set_option(string name, string value, Error & err)
{
	int		type;
	void *		ptr;
	bool		inter;

	if (name.size() == 0)
	{
		err.code = CERR_MISSING_IDENTIFIER;
		err.str = _("missing name identifier");
		return false;
	}

	if (	name == "sort" ||
		name == "columns"
	   )
	{
		if (!Configurator::verify_columns(value, err))
		{
			return false;
		}
	}


	type = get_opt_ptr(name, ptr);

	if (type == OPT_BOOL)
	{
		*(static_cast<bool *>(ptr)) = Configurator::strtobool(value);
	}
	else if (type == OPT_STRING)
	{
		*(static_cast<string *>(ptr)) = value;
	}
	else if (type == OPT_INT)
	{
		*(static_cast<int *>(ptr)) = atoi(value.c_str());
	}
	else if (type == OPT_SPECIAL)
	{
		if (name == "follow")
		{
			opt->followwindow = Configurator::strtobool(value);
			opt->followcursor = Configurator::strtobool(value);
		}
		else if (name == "scroll")
		{
			if (value == "centered" || value == "centred")
				opt->scroll_mode = SCROLL_CENTERED;
			else if (value == "relative")
				opt->scroll_mode = SCROLL_RELATIVE;
			else if (value == "normal")
				opt->scroll_mode = SCROLL_NORMAL;
			else
			{
				err.code = CERR_INVALID_VALUE;
				err.str = _("invalid scroll mode, expected 'normal', 'centered' or 'relative'");
				return false;
			}
		}
		else if (name == "playmode")
		{
			if (value == "manual")
				opt->playmode = PLAYMODE_MANUAL;
			else if (value == "linear")
				opt->playmode = PLAYMODE_LINEAR;
			else if (value == "random")
				opt->playmode = PLAYMODE_RANDOM;
			else
			{
				err.code = CERR_INVALID_VALUE;
				err.str = _("invalid play mode, expected 'manual', 'linear' or 'random'");
				return false;
			}
		}
	}
	else if (name == "topbarclear")
	{
		opt->topbar.clear();
		if (pms->disp) pms->disp->resized();
	}
	/*
	 * Topbar strings
	 */
	else if (name.size() > 6 && name.substr(0, 6) == "topbar")
	{
		inter = set_topbar_values(name, value, err);
		if (inter && pms->disp) pms->disp->resized();
		return inter;
	}
	else
	{
		err.code = CERR_INVALID_OPTION;
		err.str = _("invalid option");
		err.str += " '" + name + "'";
		return false;
	}

	/*
	 * Special cases with more handling.
	 * TODO: these should not be here, but somewhere else
	 */

	if (name == "sort")
	{
		if (pms->comm && pms->comm->library())
			pms->comm->library()->sort(opt->librarysort);
	}
	else if (name == "columns")
	{
		if (pms->disp && pms->disp->actwin())
			pms->disp->actwin()->set_column_size();
	}
	else if (name == "mouse")
	{
		if (pms->disp) pms->disp->setmousemask();
	}
	else if (name == "topbarvisible" || name == "topbarborders" || name == "topbarspace" || name == "columnspace")
	{
		if (pms->disp) pms->disp->resized();
	}

//	debug("Setting option '%s' = '%s'\n", name.c_str(), value.c_str());

	return true;
}

/*
 * Toggle a boolean option
 */
bool			Configurator::toggle_option(string name, Error & err)
{
	void *		b;

	get_opt_ptr(name, b);
	return set_option(name, *(static_cast<bool *>(b)) ? "false" : "true", err);
}

/*
 * Show an option's name and value
 */
bool			Configurator::show_option(string name, Error & err)
{
	string		val;

	val = get_option(name, err);
	if (err.code == CERR_NONE)
	{
		if (get_opt_type(name) == OPT_BOOL)
		{
			if (val == "on")
				err.str = name;
			else
				err.str = "no" + name;
		}
		else
			err.str = name + "=" + val;
		return true;
	}
	return false;
}

/*
 * Set which fields to show in the topbar.
 */
bool			Configurator::set_topbar_values(string name, string value, Error & err)
{
	int		column, row;

	column = atoi(name.substr(6).c_str());

	if (column <= 0)
		return false;
	else if (column < 10)
		name = name.substr(7);
	else if (column < 100)
		name = name.substr(8);
	else
	{
		err.code = CERR_INVALID_TOPBAR_INDEX;
		err.str = _("invalid topbar line");
		err.str += " '" + Pms::tostring(column) + "', ";
		err.str += _("expected range is 1-99");
		return false;
	}

	if (name.size() == 0)
	{
		err.code = CERR_INVALID_TOPBAR_POSITION;
		err.str = _("expected placement after topbar index");
		return false;
	}

	if (name == ".left")
		row = 0;
	else if (name == ".center" || name == ".centre")
		row = 1;
	else if (name == ".right")
		row = 2;
	else
	{
		err.code = CERR_INVALID_TOPBAR_POSITION;
		err.str = _("invalid topbar position");
		err.str += " '" + name.substr(1) + "', ";
		err.str += "expected one of: left center right";
		return false;
	}

	/* Arguments are now sanitized, should have row = 0-2, and column = 1-99 */

	while (opt->topbar.size() < column)
		opt->topbar.push_back(new Topbarline());

	opt->topbar[column-1]->strings[row] = value;

	return true;
}


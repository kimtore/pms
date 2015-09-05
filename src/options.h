/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS	<<Practical Music Search>>
 * Copyright (C) 2006-2015  Kim Tore Jensen
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
 *
 */

#ifndef _OPTIONS_H_
#define _OPTIONS_H_

#include <string>
#include <vector>

#include "types.h"
#include "color.h"
#include "topbar.h"
#include "field.h"

using namespace std;


/**
 * The option type defines which basic type an option is.
 */
typedef enum
{
	OPT_NONE,
	OPT_BOOL,
	OPT_LONG,
	OPT_STRING
}
option_t;


/**
 * Bitmasks denoting which option group an option belongs to.
 *
 * Option groups are used as an indicator of which parts of the user interface
 * needs to be re-initialized after an option change.
 */
#define OPT_GROUP_NONE		0
#define OPT_GROUP_CONNECTION	1 << 0
#define OPT_GROUP_DISPLAY	1 << 1
#define OPT_GROUP_CROSSFADE	1 << 2
#define OPT_GROUP_COLUMNS	1 << 3
#define OPT_GROUP_SORT		1 << 4
#define OPT_GROUP_MOUSE		1 << 5


class Options;


/**
 * The Option class represents a single option, regardless of type. It can be
 * used for matching option names to actual variables, as well as validating
 * the contents of them.
 */
class Option
{
public:
	Option(Options * parent_, option_t type_, uint32_t groups_, const char * name_, void * pointer_);

	Options *		parent;
	option_t		type;
	uint32_t		groups;
	char			name[32];
	void *			pointer;

	/**
	 * Return this Option as a bool pointer.
	 */
	bool *			as_bool_ptr();

	/**
	 * Return this Option as a long pointer.
	 */
	long *			as_long_ptr();

	/**
	 * Return this Option as a string pointer.
	 */
	string *		as_string_ptr();

	/**
	 * Set an Option value. This function must transfer the value of the
	 * void pointer into the correct option in the parent Options class.
	 *
	 * Subclasses should override this method, and implement their own
	 * custom parsing and validation.
	 *
	 * Returns true on success, false on failure.
	 */
	virtual bool		set(void * value);
};


/**
 * The TagListOption represents a string option containing only an array of
 * tags a.k.a. Fields.
 */
class TagListOption : public Option
{
public:
	TagListOption(Options * parent_, option_t type_, uint32_t groups_, const char * name_, void * pointer_) :
		Option(parent_, type_, groups_, name_, pointer_) { };

	bool			set(void * value);
};


/**
 * The TopbarOption represents a string option containing a Topbar formatting
 * string. The string is delimited by the literal strings "\n" and "\t", used
 * to split the topbar formatting string into lines and fields, respectively.
 *
 * Any topbar line can contain up to three fields.
 *
 * Example: the following string:
 *
 *	foo\tbar\nfoobar\t\tbaz\n\ttest
 *
 * Will produce the following formatted output:
 *	foo		bar
 *	foobar				baz
 *			test
 *
 */
class TopbarOption : public Option
{
public:
	TopbarOption(Options * parent_, option_t type_, uint32_t groups_, const char * name_, void * pointer_) :
		Option(parent_, type_, groups_, name_, pointer_) { };

	vector<Topbarline *>	topbar;

	/**
	 * Delete all the Topbarline pointers in the topbar vector, and clear
	 * the vector itself.
	 */
	void			clear();

	/**
	 * Parse the string value given, and convert it into a set of
	 * Topbarlines. The results are put into the topbar vector.
	 */
	bool			parse(const string * value);

	bool			set(void * value);
};


/**
 * The ScrollModeOption converts the string literal values "centered", "normal"
 * and "relative" to their enumerated equivalents.
 */
class ScrollModeOption : public Option
{
public:
	ScrollModeOption(Options * parent_, option_t type_, uint32_t groups_, const char * name_, void * pointer_) :
		Option(parent_, type_, groups_, name_, pointer_) { };

	bool			set(void * value);
};


/**
 * The Options class holds all the program options that PMS supports.
 */
class Options
{
private:
	vector<Option *>	option_index;

	uint32_t		changed_flags;

public:
				Options();
				~Options();
	
	void			reset();

	Colortable *		colors;
	vector<Topbarline *>	topbar_lines;
	pms_scroll_mode		scroll_mode;

	/* FIXME: refactor this elsewhere */
	string			configfile;

	long			crossfade;
	long			mpd_timeout;
	long			msg_buffer_size;
	long			nextinterval;
	long			port;
	long			reconnectdelay;
	long			resetstatus;
	long			scrolloff;

	bool			addtoreturns;
	bool			columnborders;
	bool			debug;
	bool			followcursor;
	bool			followplayback;
	bool			followwindow;
	bool			ignorecase;
	bool			mouse;
	bool			nextafteraction;
	bool			regexsearch;
	bool			topbarborders;
	bool			topbarvisible;

	string			columns;
	string			host;
	string			libraryroot;
	string			onplaylistfinish;
	string			password;
	string			scroll;
	string			sort;
	string			startuplist;
	string			status_pause;
	string			status_play;
	string			status_stop;
	string			status_unknown;
	string			topbar;
	string			xtermtitle;

	/**
	 * Look up an option by its name.
	 *
	 * Returns the option type if the option is known, or OPT_NONE if there is no
	 * such option.
	 *
	 * The second parameter is optional. If a pointer is passed, the value of the
	 * pointer will be modified to point to the option, or NULL if the option was
	 * not known.
	 */
	Option *		lookup_option(const char * varname);

	uint32_t		get_changed_flags();
	uint32_t		set_changed_flags(uint32_t flags);
	uint32_t		add_changed_flags(uint32_t flags);
};


#endif /* _OPTIONS_H_ */

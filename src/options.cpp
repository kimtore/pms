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
 * options.cpp - the options class
 *
 */

#include "options.h"
#include "color.h"
#include "i18n.h"
#include "pms.h"

extern Pms * pms;

/**
 * Macros for easy creation of options
 */
#define NEW_OPTION(CLASS, TYPE, GROUP, NAME) option_index.push_back(new CLASS(this, TYPE, GROUP, #NAME, &NAME))
#define NEW_BOOL(NAME) NEW_OPTION(Option, OPT_BOOL, OPT_GROUP_NONE, NAME)
#define NEW_BOOL_GROUPED(NAME, GROUP) NEW_OPTION(Option, OPT_BOOL, GROUP, NAME)
#define NEW_LONG(NAME) NEW_OPTION(Option, OPT_LONG, OPT_GROUP_NONE, NAME)
#define NEW_LONG_GROUPED(NAME, GROUP) NEW_OPTION(Option, OPT_LONG, GROUP, NAME)
#define NEW_STRING(NAME) NEW_OPTION(Option, OPT_STRING, OPT_GROUP_NONE, NAME)
#define NEW_STRING_GROUPED(NAME, GROUP) NEW_OPTION(Option, OPT_STRING, GROUP, NAME)
#define NEW_TAG_LIST(NAME, GROUP) NEW_OPTION(TagListOption, OPT_STRING, OPT_GROUP_NONE, NAME)
#define NEW_TOPBAR(NAME, GROUP) NEW_OPTION(TopbarOption, OPT_STRING, GROUP, NAME)
#define NEW_SCROLL_MODE(NAME, GROUP) NEW_OPTION(ScrollModeOption, OPT_STRING, GROUP, NAME)


Option::Option(Options * parent_, option_t type_, uint32_t groups_, const char * name_, void * pointer_)
{
	parent = parent_;
	type = type_;
	groups = groups_;
	strncpy(name, name_, 31);
	name[31] = '\0';
	pointer = pointer_;
}

bool *
Option::as_bool_ptr()
{
	return static_cast<bool *>(pointer);
}

long *
Option::as_long_ptr()
{
	return static_cast<long *>(pointer);
}

string *
Option::as_string_ptr()
{
	return static_cast<string *>(pointer);
}

bool
Option::set(void * value)
{
	switch(type) {
		case OPT_BOOL:
			*(static_cast<bool *>(pointer)) = *(static_cast<bool *>(value));
			break;
		case OPT_LONG:
			*(static_cast<long *>(pointer)) = *(static_cast<long *>(value));
			break;
		case OPT_STRING:
			*(static_cast<string *>(pointer)) = *(static_cast<string *>(value));
			break;
		default:
			abort();
	};

	parent->add_changed_flags(groups);

	return true;
}

bool
TagListOption::set(void * value)
{
	vector<string> *		v;
	vector<string>::iterator	iter;
	string *			s;

	s = static_cast<string *>(value);

	if (!s->size()) {
		return false;
	}

	v = Pms::splitstr(*s, " ");

	iter = v->begin();

	while (iter != v->end()) {
		if (pms->fieldtypes->lookup(*iter) == -1) {
			delete v;
			return false;
		}
		++iter;
	}

	delete v;

	return Option::set(value);
}

void
TopbarOption::clear()
{
	vector<Topbarline *>::const_iterator	iter;

	iter = topbar.begin();
	while (iter != topbar.end()) {
		delete *iter;
		++iter;
	}

	topbar.clear();
}

bool
TopbarOption::parse(const string * value)
{
	uint32_t		position = 0;
	bool			escaped = false;
	Topbarline *		line = NULL;
	string			buffer = "";
	string::const_iterator	iter;

	clear();
	iter = value->begin();

	while (iter != value->end()) {
		if (!line) {
			line = new Topbarline();
			position = 0;
		}
		if (!escaped) {
			if (*iter == '\\') {
				escaped = true;
			} else {
				buffer += *iter;
			}
		} else {
			if (*iter == 't' || *iter == 'n') {
				line->strings[position] = buffer;
				buffer.clear();
			}
			if (*iter == 'n') {
				topbar.push_back(line);
				line = NULL;
			} else if (*iter == 't') {
				if (++position > 2) {
					break;
				}
			} else {
				break;
			}
			escaped = false;
		}
		++iter;
	}

	if (line) {
		line->strings[position] = buffer;
		topbar.push_back(line);
	}

	if (iter != value->end()) {
		clear();
		return false;
	}

	return true;
}

bool
TopbarOption::set(void * value)
{
	string * s;

	s = static_cast<string *>(value);

	if (!parse(s)) {
		return false;
	}

	parent->topbar_lines = topbar;
	topbar.clear();

	return Option::set(value);
}

bool
ScrollModeOption::set(void * value)
{
	string * s = static_cast<string *>(value);

	if (*s == "normal") {
		parent->scroll_mode = SCROLL_NORMAL;
	} else if (*s == "center" || *s == "centered" || *s == "centre" || *s == "centred") {
		parent->scroll_mode = SCROLL_CENTERED;
	} else if (*s == "relative") {
		parent->scroll_mode = SCROLL_RELATIVE;
	} else {
		return false;
	}

	return Option::set(value);
}

Options::Options()
{
	colors = NULL;
	changed_flags = 0;

	NEW_BOOL(addtoreturns);
	NEW_BOOL_GROUPED(columnborders, OPT_GROUP_DISPLAY);
	NEW_BOOL(debug);
	NEW_BOOL(followcursor);
	NEW_BOOL(followplayback);
	NEW_BOOL(followwindow);
	NEW_BOOL_GROUPED(ignorecase, OPT_GROUP_SORT);
	NEW_BOOL(mouse);
	NEW_BOOL(nextafteraction);
	NEW_BOOL(regexsearch);
	NEW_BOOL_GROUPED(topbarborders, OPT_GROUP_DISPLAY);
	NEW_BOOL_GROUPED(topbarvisible, OPT_GROUP_DISPLAY);

	NEW_LONG(crossfade);
	NEW_LONG_GROUPED(mpd_timeout, OPT_GROUP_CONNECTION);
	NEW_LONG(msg_buffer_size);
	NEW_LONG(nextinterval);
	NEW_LONG_GROUPED(port, OPT_GROUP_CONNECTION);
	NEW_LONG(reconnectdelay);
	NEW_LONG(resetstatus);
	NEW_LONG_GROUPED(scrolloff, OPT_GROUP_DISPLAY);

	NEW_TAG_LIST(columns, OPT_GROUP_COLUMNS);
	NEW_STRING_GROUPED(host, OPT_GROUP_CONNECTION);
	NEW_STRING(libraryroot);
	NEW_STRING(onplaylistfinish);
	NEW_STRING(password);
	NEW_SCROLL_MODE(scroll, OPT_GROUP_DISPLAY);
	NEW_TAG_LIST(sort, OPT_GROUP_SORT);
	NEW_STRING(startuplist);
	NEW_STRING_GROUPED(status_pause, OPT_GROUP_DISPLAY);
	NEW_STRING_GROUPED(status_play, OPT_GROUP_DISPLAY);
	NEW_STRING_GROUPED(status_stop, OPT_GROUP_DISPLAY);
	NEW_STRING_GROUPED(status_unknown, OPT_GROUP_DISPLAY);
	NEW_TOPBAR(topbar, OPT_GROUP_DISPLAY);
	NEW_STRING(xtermtitle);

	reset();
}

Options::~Options()
{
	vector<Option *>::iterator iter;

	iter = option_index.begin();
	while (iter != option_index.end()) {
		delete *iter;
		++iter;
	}
}

void
Options::reset()
{
	TopbarOption * topbar_option;
	ScrollModeOption * scroll_option;

	if (colors != NULL) {
		delete colors;
	}

	colors = new Colortable();

	addtoreturns = false;
	columnborders = false;
	columns = "artist track title album length";
	crossfade = 5;
	debug = false;
	followcursor = false;
	followplayback = false;
	followwindow = false;
	host = "localhost";
	ignorecase = true;
	libraryroot = "";
	mouse = false;
	mpd_timeout = 2;
	msg_buffer_size = 1024;
	nextafteraction = true;
	nextinterval = 5;
	onplaylistfinish = "";
	password = "";
	port = 6600;
	reconnectdelay = 10;
	regexsearch = false;
	resetstatus = 3;
	scroll = "normal";
	scrolloff = 0;
	sort = "track disc album date albumartistsort";
	startuplist = "playlist";
	status_pause = Pms::unicode() ? "‖" : "||";
	status_play = Pms::unicode() ? "▶" : "|>";
	status_stop = Pms::unicode() ? "■" : "[]";
	status_unknown = Pms::unicode() ? "?" : "??";
	topbar = "\\n"
		"%volume%%% Mode: %muteshort%%consumeshort%%repeatshort%%randomshort%%singleshort%%ifcursong% %playstate% %time_elapsed% / %time_remaining%%endif%\\t"
		"%ifcursong%%artist% - %title% on %album% from %date%%else%Not playing anything%endif%\\t"
		"Queue has %livequeuesize%\\n"
		"\\t\\t%listsize%\\n"
		"%progressbar%";
	topbarborders = false;
	topbarvisible = true;
	xtermtitle = "PMS: %playstate%%ifcursong% %artist% – %title%%endif%";
	
	/* Derive scroll mode */
	scroll_option = dynamic_cast<ScrollModeOption *>(lookup_option("scroll"));
	scroll_option->set(&scroll);

	/* Set up default top bar values */
	topbar_option = dynamic_cast<TopbarOption *>(lookup_option("topbar"));
	topbar_option->set(&topbar);
}

Option *
Options::lookup_option(const char * varname)
{
	vector<Option *>::const_iterator iter;

	iter = option_index.begin();

	while (iter != option_index.end()) {
		if (!strcmp((*iter)->name, varname)) {
			return *iter;
		}
		++iter;
	}

	return NULL;
}

uint32_t
Options::get_changed_flags()
{
	return changed_flags;
}

uint32_t
Options::set_changed_flags(uint32_t flags)
{
	changed_flags = flags;
	return get_changed_flags();
}

uint32_t
Options::add_changed_flags(uint32_t flags)
{
	changed_flags |= flags;
	return get_changed_flags();
}

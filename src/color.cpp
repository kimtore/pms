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
 * color.cpp
 * Color class, holds information about color + attribute values for one color pair
 */

#include "mycurses.h"
#include "color.h"
#include "pms.h"

static short			totalcolors = 0;
extern Pms *			pms;


Colortable::Colortable()
{
	isset = false;
	defaults();
}

Colortable::~Colortable()
{
	clear();
}

void			Colortable::clear()
{
	if (!isset) return;

	totalcolors = 0;

	/* standard colors */
	delete back;
	delete standard;
	delete status;
	delete status_error;
	delete position;
	delete border;
	delete headers;
	delete title;

	/* topbar colors */
	delete topbar.standard;
	delete topbar.time_elapsed;
	delete topbar.time_remaining;
	delete topbar.progressbar;
	delete topbar.progresspercentage;
	delete topbar.librarysize;
	delete topbar.listsize;
	delete topbar.queuesize;
	delete topbar.livequeuesize;
	delete topbar.playstate;
	delete topbar.volume;
	delete topbar.bitrate;
	delete topbar.samplerate;
	delete topbar.bits;
	delete topbar.channels;

	delete topbar.repeat;
	delete topbar.random;
	delete topbar.manualprogression;
	delete topbar.mute;
	delete topbar.repeatshort;
	delete topbar.randomshort;
	delete topbar.manualprogressionshort;
	delete topbar.muteshort;

	/* field types */
	delete fields.num;
	delete fields.file;
	delete fields.artist;
	delete fields.albumartist;
	delete fields.artistsort;
	delete fields.albumartistsort;
	delete fields.title;
	delete fields.album;
	delete fields.track;
	delete fields.trackshort;
	delete fields.time;
	delete fields.date;
	delete fields.year;
	delete fields.name;
	delete fields.genre;
	delete fields.composer;
	delete fields.performer;
	delete fields.disc;
	delete fields.comment;

	/* field types for the topbar */
	delete topbar.fields.num;
	delete topbar.fields.file;
	delete topbar.fields.artist;
	delete topbar.fields.albumartist;
	delete topbar.fields.artistsort;
	delete topbar.fields.albumartistsort;
	delete topbar.fields.title;
	delete topbar.fields.album;
	delete topbar.fields.track;
	delete topbar.fields.trackshort;
	delete topbar.fields.time;
	delete topbar.fields.date;
	delete topbar.fields.year;
	delete topbar.fields.name;
	delete topbar.fields.genre;
	delete topbar.fields.composer;
	delete topbar.fields.performer;
	delete topbar.fields.disc;
	delete topbar.fields.comment;

	/* list colors */
	delete current;
	delete cursor;
	delete selection;
	delete lastlist;
	delete playinglist;

	isset = false;
}

void			Colortable::defaults()
{
	if (isset) clear();

	/* standard colors */
	back					= new color(COLOR_BLACK, -1, 0);
	standard				= new color(COLOR_WHITE, -1, 0);
	status					= new color(COLOR_WHITE, -1, 0);
	status_error				= new color(COLOR_WHITE, COLOR_RED, 0);
	position				= new color(COLOR_WHITE, -1, 0);
	border					= new color(COLOR_BLACK, -1, 0);
	headers					= new color(COLOR_BLACK, -1, A_BOLD);
	title					= new color(COLOR_CYAN, -1, A_BOLD);

	/* topbar colors */
	topbar.standard				= new color(COLOR_WHITE, -1, 0);
	topbar.time_elapsed			= new color(COLOR_GREEN, -1, 0);
	topbar.time_remaining			= new color(COLOR_WHITE, -1, 0);
	topbar.progressbar			= new color(COLOR_WHITE, -1, 0);
	topbar.progresspercentage		= new color(COLOR_WHITE, -1, 0);
	topbar.librarysize			= new color(COLOR_WHITE, -1, 0);
	topbar.listsize				= new color(COLOR_WHITE, -1, 0);
	topbar.queuesize			= new color(COLOR_WHITE, -1, 0);
	topbar.livequeuesize			= new color(COLOR_WHITE, -1, 0);
	topbar.playstate			= new color(COLOR_WHITE, -1, 0);
	topbar.volume				= new color(COLOR_YELLOW, -1, 0);
	topbar.bitrate				= new color(COLOR_WHITE, -1, 0);
	topbar.samplerate			= new color(COLOR_WHITE, -1, 0);
	topbar.bits				= new color(COLOR_WHITE, -1, 0);
	topbar.channels				= new color(COLOR_WHITE, -1, 0);

	topbar.repeat				= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.random				= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.manualprogression		= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.mute				= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.repeatshort			= new color(COLOR_CYAN, -1, 0);
	topbar.randomshort			= new color(COLOR_CYAN, -1, 0);
	topbar.manualprogressionshort		= new color(COLOR_CYAN, -1, 0);
	topbar.muteshort			= new color(COLOR_CYAN, -1, 0);

	/* field types */
	fields.num				= new color(COLOR_BLACK, -1, A_BOLD);
	fields.file				= new color(COLOR_WHITE, -1, 0);
	fields.artist				= new color(COLOR_YELLOW, -1, 0);
	fields.albumartist			= new color(COLOR_YELLOW, -1, 0);
	fields.artistsort			= new color(COLOR_YELLOW, -1, 0);
	fields.albumartistsort			= new color(COLOR_YELLOW, -1, 0);
	fields.title				= new color(COLOR_WHITE, -1, A_BOLD);
	fields.album				= new color(COLOR_CYAN, -1, 0);
	fields.track				= new color(COLOR_BLACK, -1, A_BOLD);
	fields.trackshort			= new color(COLOR_BLACK, -1, A_BOLD);
	fields.time				= new color(COLOR_MAGENTA, -1, 0);
	fields.date				= new color(COLOR_YELLOW, -1, 0);
	fields.year				= new color(COLOR_YELLOW, -1, 0);
	fields.name				= new color(COLOR_WHITE, -1, 0);
	fields.genre				= new color(COLOR_WHITE, -1, 0);
	fields.composer				= new color(COLOR_WHITE, -1, 0);
	fields.performer			= new color(COLOR_WHITE, -1, 0);
	fields.disc				= new color(COLOR_BLACK, -1, 0);
	fields.comment				= new color(COLOR_WHITE, -1, 0);

	/* field types for the topbar */
	topbar.fields.num			= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.fields.file			= new color(COLOR_WHITE, -1, 0);
	topbar.fields.artist			= new color(COLOR_YELLOW, -1, A_BOLD);
	topbar.fields.albumartist		= new color(COLOR_YELLOW, -1, A_BOLD);
	topbar.fields.artistsort		= new color(COLOR_YELLOW, -1, A_BOLD);
	topbar.fields.albumartistsort		= new color(COLOR_YELLOW, -1, A_BOLD);
	topbar.fields.title			= new color(COLOR_WHITE, -1, A_BOLD);
	topbar.fields.album			= new color(COLOR_CYAN, -1, 0);
	topbar.fields.track			= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.fields.trackshort		= new color(COLOR_BLACK, -1, A_BOLD);
	topbar.fields.time			= new color(COLOR_MAGENTA, -1, 0);
	topbar.fields.date			= new color(COLOR_YELLOW, -1, 0);
	topbar.fields.year			= new color(COLOR_YELLOW, -1, 0);
	topbar.fields.name			= new color(COLOR_WHITE, -1, 0);
	topbar.fields.genre			= new color(COLOR_WHITE, -1, 0);
	topbar.fields.composer			= new color(COLOR_WHITE, -1, 0);
	topbar.fields.performer			= new color(COLOR_WHITE, -1, 0);
	topbar.fields.disc			= new color(COLOR_BLACK, -1, 0);
	topbar.fields.comment			= new color(COLOR_WHITE, -1, 0);

	/* list colors */
	current					= new color(COLOR_BLACK, COLOR_YELLOW, 0);
	cursor					= new color(COLOR_BLACK, COLOR_WHITE, 0);
	selection				= new color(COLOR_BLACK, COLOR_GREEN, 0);
	lastlist				= new color(COLOR_WHITE, COLOR_BLUE, A_BOLD);
	playinglist				= new color(COLOR_BLACK, COLOR_YELLOW, 0);

	isset = true;
}






color::color()
{
	clean();
}

color::color(int fg, int bg, int at)
{
	clean();
	set(fg, bg, at);
}

color::~color()
{
}

void		color::clean()
{
	front = -1;
	back = -1;
	id = -1;
	attr = 0;
	initialized = false;
	isset = false;
}

bool		color::set(int fg, int bg, int at)
{
	if (fg < COLOR_BLACK || fg > COLOR_WHITE)
		return false;

	front = fg;
	back = bg;
	attr = at;
	isset = true;
	initialized = false;

	return true;
}

int		color::pair()
{
	if (!isset)
		return 0;
	if (!initialized)
	{
		if (back == -1)
			back = pms->options->colors->back->back;
		if (id == -1)
			id = ++totalcolors;
		init_pair(id, front, back);
	}

	initialized = true;

	return (COLOR_PAIR(id) | attr);
}

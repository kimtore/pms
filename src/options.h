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
 * options.h - the options class 
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


enum
{
	OPT_NONE,
	OPT_BOOL,
	OPT_INT,
	OPT_STRING,
	OPT_SPECIAL
};


class Options
{
public:
				Options();
	
	void			reset();

	/* These are all the options */

	string			hostname;
	unsigned int		port;
	string			password;
	string			configfile;

	bool			debug;
	Colortable *		colors;

	int			nextinterval;		// Time before crossfade to add next song from library to queue
	int			crossfade;		// Default crossfade time (for toggling)
	int			repeatonedelay;		// How many seconds remaining in the song before playing it again in repeat-1 mode
	int			stopdelay;		// How many seconds remaining in the song before stopping, when in manual progression mode
	int			mpd_timeout;		// MPD server connection timeout
	int			reconnectdelay;		// How many seconds to wait before retrying connection
	int			resetstatus;		// Time to wait before resetting status text
	bool			regexsearch;		// Whether or not to search with regular expressions
	int			directoryminlen;	// Minimum length of directory view panes
	string			directoryformat;	// How to format songs within directory view
	string			onplaylistfinish;	// Optional shell command to run when playlist finishes
	bool			mouse;			// Whether or not to use mouse functions
	bool			followwindow;		// Playback follows active window
	bool			followcursor;		// Playback follows cursor position
	bool			followplayback;		// Cursor position follows playback
	bool			nextafteraction;	// Move cursor to next item after selecting/unselecting/toggling selection or adding to playlist
	bool			showtopbar;		// Whether or not to display the topbar
	bool			topbarspace;		// Whether or not to leave a blank row between the topbar and the playlist windows
	bool			columnspace;		// Whether or not to add an extra unit of width to fixed-width columns
	bool			topbarborders;		// Draw borders on topbar window?
	bool			addtoreturns;		// Return to original window when using addto?
	bool			ignorecase;		// Perform case-insensitive sorts and matches?
	vector<Topbarline *>	topbar;			// Topbar draw information
	string			startuplist;		// This list is focused when PMS starts
	string			stralbumclass;		// Textual representation of albumclass
	vector<Item>		albumclass;		// What fields needs to be the same to call a group of songs 'album' ?
	string			librarysort;		// How to sort the library
	string			columns;		// Which columns to show in main view
	string			xtermtitle;		// Title to set in graphical terminals
	pms_play_mode		playmode;
	pms_repeat_mode		repeatmode;
	pms_scroll_mode		scroll_mode;
	string			libraryroot;		// Path to the library root, for prepending to filenames when they are being used in shell commands
	int			scrolloff;		// Number of items to keep in view above and below cursor in normal scroll mode

	/* Status indicator strings */
	string			status_unknown;
	string			status_play;
	string			status_stop;
	string			status_pause;
};


#endif /* _OPTIONS_H_ */

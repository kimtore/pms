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
 * options.cpp - the options class
 *
 */

#include "options.h"
#include "color.h"
#include "i18n.h"


Options::Options()
{
	colors = NULL;
	reset();
}

void			Options::reset()
{
	if (colors != NULL)
		delete colors;

	colors = new Colortable();

	scroll_mode		= SCROLL_NORMAL;
	nextinterval		= 5;
	crossfade		= 5;
	mpd_timeout		= 30;
	repeatonedelay		= 1;
	stopdelay		= 1;
	reconnectdelay		= 30;
	debug			= false;
	addtoreturns		= false;
	directoryminlen		= 30;
	directoryformat		= "%artist% - %title%";
	playmode		= PLAYMODE_LINEAR;
	repeatmode		= REPEAT_NONE;
	regexsearch		= false;
	resetstatus		= 3;
	startuplist		= "playlist";
	librarysort		= "default";
	columns			= "artist track title album length";
	onplaylistfinish	= "";
	followwindow		= false;
	followcursor		= false;
	followplayback		= false;
	nextafteraction		= true;
	showtopbar		= true;
	topbarborders		= false;
	topbarspace		= true;
	mouse			= false;
	libraryroot		= "";
	scrolloff		= 0;

	//TODO: would be nice to have the commented alteratives default if 
	//Unicode is available
	status_unknown		= "??"; //?
	status_play		= "|>"; //▶
	status_pause		= "||"; //‖
	status_stop		= "[]"; //■
	
	/* Album classification */
	stralbumclass		= "artist album date";
	albumclass.clear();

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

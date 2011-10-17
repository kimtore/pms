/* vi:set ts=8 sts=8 sw=8:
 *
 * Practical Music Search
 * Copyright (c) 2006-2011  Kim Tore Jensen
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
 */

#include "window.h"
#include "console.h"
#include "curses.h"
#include "input.h"
#include "config.h"
#include "mpd.h"
#include <cstring>
#include <string>
#include <vector>
#include <sys/time.h>

using namespace std;

extern vector<Logline *> logbuffer;
extern Curses curses;
extern Windowmanager wm;
extern Input input;
extern Config config;
extern MPD mpd;

Wstatusbar::Wstatusbar()
{
	memset(&cl, 0, sizeof cl);
	cl_isreset = false;
}

void Wstatusbar::drawline(int rely)
{
	string pstr;
	size_t vscroll = 0;
	size_t width = rect->right - rect->left;
	vector<Logline *>::reverse_iterator i;
	Color * c;

	switch(input.mode)
	{
		case INPUT_MODE_COMMAND:
			gettimeofday(&cl, NULL);
			for (i = logbuffer.rbegin(); i != logbuffer.rend(); i++)
			{
				if ((*i)->level > MSG_LEVEL_INFO)
					continue;

				memcpy(&cl_reset, &((*i)->tm), sizeof cl_reset);

				/* Expired message - draw playstring instead */
				if (cl.tv_sec - cl_reset.tv_sec >= (int)config.status_reset_interval)
				{
					cl_isreset = true;
					curses.wipe(rect, config.colors.statusbar);
					curses.print(rect, config.colors.statusbar, rely, 0, mpd.playstring.c_str());
					break;
				}

				/* Draw last statusbar message */
				cl_isreset = false;
				c = (*i)->level == MSG_LEVEL_ERR ? config.colors.error : config.colors.statusbar;
				curses.wipe(rect, c);
				curses.print(rect, c, rely, 0, (*i)->line.c_str());
				break;
			}
			break;

		case INPUT_MODE_INPUT:
		case INPUT_MODE_SEARCH:
		case INPUT_MODE_LIVESEARCH:
			if (input.strbuf.size() >= width)
				vscroll = input.strbuf.size() - width;
			if (vscroll > input.cursorpos)
				vscroll = input.cursorpos;
			curses.wipe(rect, config.colors.standard);

			if (input.mode == INPUT_MODE_INPUT)
				pstr = ":";
			else if (input.mode == INPUT_MODE_SEARCH)
				pstr = "Search: ";
			else if (input.mode == INPUT_MODE_LIVESEARCH)
				pstr = "/";

			curses.print(rect, config.colors.statusbar, rely, 0, pstr.c_str());
			curses.print(rect, config.colors.statusbar, rely, pstr.size(), input.strbuf.substr(vscroll).c_str());
			curses.setcursor(rect, rely, input.cursorpos - vscroll + pstr.size());
			break;
	}
}

void Wreadout::drawline(int rely)
{
	char		buf[4];
	Wmain *		win;

	win = WMAIN(wm.active);

	if (win->content_size() <= win->height())
		strcpy(buf, "All");
	else if (win->position == 0)
		strcpy(buf, "Top");
	else if (win->position >= win->content_size() - win->height() - 1)
		strcpy(buf, "Bot");
	else
		sprintf(buf, "%2d%%%%", 100 * win->position / (win->content_size() - win->height() - 1));

	curses.print(rect, config.colors.readout, 0, 0, buf);
}

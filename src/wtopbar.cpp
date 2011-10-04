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
#include "config.h"
#include "curses.h"
#include "topbar.h"
#include "field.h"
#include "mpd.h"

extern Config config;
extern Curses curses;
extern Topbar topbar;
extern MPD mpd;

void Wtopbar::drawline(int rely)
{
	vector<Topbarsegment *>::iterator seg;
	vector<Topbarchunk *>::iterator chunk;
	unsigned int pos;
	unsigned int y;
	int x;
	unsigned int strl;

	if (!rect || rely < 0 || rely > (int)height())
		return;

	y = (unsigned int)rely;

	curses.clearline(rect, y, config.colors.topbar);

	/* Cycle left, center, right */
	for (pos = 0; pos < 3; pos++)
	{
		if (topbar.lines[pos].size() <= y)
			break;

		for (seg = topbar.lines[pos].at(y)->segments.begin(); seg != topbar.lines[pos].at(y)->segments.end(); ++seg)
		{
			strl = (*seg)->compile(mpd.currentsong);
			if (pos == 0)
				x = 0;
			else if (pos == 1)
				x = ((rect->right - rect->left) / 2) - (strl / 2);
			else if (pos == 2)
				x = rect->right - rect->left - strl;

			for (chunk = (*seg)->chunks.begin(); chunk != (*seg)->chunks.end(); ++chunk)
			{
				curses.print(rect, (*chunk)->color, rely, x, (*chunk)->str.c_str());
				x += (*chunk)->str.size();
			}
		}
	}
}

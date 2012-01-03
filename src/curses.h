/* vi:set ts=8 sts=8 sw=8 noet:
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

#ifndef _PMS_CURSES_H_
#define _PMS_CURSES_H_

#include <ncurses.h>
#include "color.h"

typedef struct
{
	int	left;
	int	top;
	int	right;
	int	bottom;
}

Rect;

class Curses
{
	public:

		Curses();
		~Curses();

		/*
		 * Prints formatted, color output onto a rectangle.
		 *
		 * %s		= char *
		 * %d		= int
		 * %f		= double
		 * %B %/B	= bold on/off
		 * %R %/R	= reverse on/off
		 * %0-n% %/0-n%	= color on/off
		 *
		 */
		void		print(Rect * rect, Color * c, int y, int x, const char * fmt, ...);

		/* Refresh the screen. */
		void		flush();

		/* Clear a line, relative to rect. */
		void		clearline(Rect * rect, int line, Color * c);

		/* Clear the rectangle. */
		void		wipe(Rect * rect, Color * c);

		/* Set cursor position relative to rect */
		void		setcursor(Rect * rect, int y, int x);

		/* Set left/right/top/bottom layout for all panels */
		void		detect_dimensions();

		/* Trigger the bell */
		void		bell();

		Rect		self;
		Rect		topbar;
		Rect		main;
		Rect		statusbar;
		Rect		readout;

		bool		ready;
		bool		hascolors;
};

#endif /* _PMS_CURSES_H_ */

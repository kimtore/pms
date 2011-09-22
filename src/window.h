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

#ifndef _PMS_WINDOW_H_
#define _PMS_WINDOW_H_

#include "curses.h"
#include <vector>

using namespace std;

#define WWINDOW(x)	dynamic_cast<Window *>(x)
#define WCONSOLE(x)	dynamic_cast<Wconsole *>(x)
#define WMAIN(x)	dynamic_cast<Wmain *>(x)

class Window
{
	protected:
		Rect *		rect;

	public:
		void		set_rect(Rect * r) { rect = r; };

		/* Draw all lines on rect */
		void		draw();

		/* Is this window visible? */
		virtual bool	visible() { return true; };

		/* Draw one line on rect */
		virtual void	drawline(int y) = 0;

};

class Wmain : public Window
{
	protected:

		/* Scroll position */
		unsigned int	position;

		/* Cursor position */
		unsigned int	cursor;

	public:

		Wmain();

		/* Scroll window */
		virtual void	scroll_window(int offset);

		/* Move cursor inside window */
		virtual void	move_cursor(int offset);

		/* List size */
		virtual unsigned int content_size() = 0;

		/* Is this window visible? */
		bool		visible();
};

class Wconsole : public Wmain
{
	public:
		void		drawline(int rely);
		unsigned int	content_size();
};

class Wtopbar : public Window
{
	public:
		void		drawline(int rely) {};
};

class Wstatusbar : public Window
{
	public:
		void		drawline(int rely);
};

class Windowmanager
{
	private:
		vector<Window *>	windows;
	
	public:
		Windowmanager();

		/* What kind of input events are accepted right now */
		int			context;

		/* Redraw all visible windows */
		void			draw();

		Wconsole *		console;
		Window *		active;
		Wtopbar *		topbar;
		Wstatusbar *		statusbar;
};

#endif /* _PMS_WINDOW_H_ */

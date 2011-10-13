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
#include "songlist.h"
#include "command.h"
#include <vector>

using namespace std;

#define WWINDOW(x)	dynamic_cast<Window *>(x)
#define WCONSOLE(x)	dynamic_cast<Wconsole *>(x)
#define WMAIN(x)	dynamic_cast<Wmain *>(x)
#define WSONGLIST(x)	dynamic_cast<Wsonglist *>(x)

class Window
{
	protected:
		Rect *		rect;

	public:
		void		set_rect(Rect * r) { rect = r; };

		/* Window height */
		unsigned int	height();

		/* Draw all lines on rect */
		void		draw();

		/* Clear this window */
		void		clear();

		/* Is this window visible? */
		virtual bool	visible() { return true; };

		/* Draw one line on rect */
		virtual void	drawline(int y) = 0;

};

class Wmain : public Window
{
	protected:

	public:
		/* Which context should commands be accepted in */
		int		context;

		/* Scroll position */
		unsigned int	position;

		/* Cursor position */
		unsigned int	cursor;

		/* Window title */
		string		title;


		Wmain();
		virtual unsigned int	height();

		/* Draw all lines and update readout */
		virtual void	draw();

		/* Scroll window */
		virtual void	scroll_window(int offset);

		/* Move cursor inside window */
		virtual void	move_cursor(int offset);

		/* Set absolute window/cursor position */
		virtual void	set_position(unsigned int absolute);
		virtual void	set_cursor(unsigned int absolute);

		/* List size */
		virtual unsigned int content_size() = 0;

		/* Is this window visible? */
		bool		visible();
};

class Wconsole : public Wmain
{
	public:
		Wconsole() { context = CONTEXT_CONSOLE; };

		void		drawline(int rely);
		unsigned int	content_size();
		void		move_cursor(int offset);
		void		set_cursor(unsigned int absolute);
};

class Wsonglist : public Wmain
{
	private:
		vector<unsigned int>	column_len;

	public:
		Wsonglist() { context = CONTEXT_SONGLIST; };

		void		draw();
		void		drawline(int rely);
		unsigned int	height();
		unsigned int	content_size();

		/* Pointer to connected songlist */
		Songlist *	songlist;

		/* Pointer to song beneath cursor */
		Song *		cursorsong();

		/* Update column lengths */
		void		update_column_length();
};

class Wtopbar : public Window
{
	public:
		void		drawline(int rely);
};

class Wstatusbar : public Window
{
	public:
		Wstatusbar();

		void		drawline(int rely);
		struct timeval	cl;
};

class Wreadout : public Window
{
	public:
		void		drawline(int rely);
};

class Windowmanager
{
	private:
		vector<Wmain *>		windows;

		/* Active window index */
		unsigned int		active_index;
	
	public:
		Windowmanager();

		/* What kind of input events are accepted right now */
		int			context;

		/* Redraw all visible windows */
		void			draw();

		/* Flush ncurses buffer */
		void			flush();

		/* Cycle window list */
		void			cycle(int offset);

		/* Activate a window */
		bool			activate(Wmain * nactive);

		/* Activate the last used window */
		bool			toggle();

		/* Activate a window with given title, case insensitive */
		bool			go(string title);

		/* Update column lengths in all windows */
		void			update_column_length();

		Wconsole *		console;
		Wsonglist *		playlist;
		Wsonglist *		library;

		Wmain *			last_active;
		Wmain *			active;
		Wtopbar *		topbar;
		Wstatusbar *		statusbar;
		Wreadout *		readout;
};

#endif /* _PMS_WINDOW_H_ */

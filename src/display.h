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
 * display.h - ncurses, display and window management
 *
 */

#ifndef _PMS_DISPLAY_H_
#define _PMS_DISPLAY_H_

#include <cmath>
#include <cstdlib>
#include <string>
#include <vector>
#include "mycurses.h"
#include "types.h"
#include "field.h"
#include "string.h"
#include "command.h"
#include "color.h"

using namespace std;


class Display;

enum pms_win_role
{
	WIN_ROLE_STATIC = 0,
	WIN_ROLE_TOPBAR,
	WIN_ROLE_PLAYLIST,
	WIN_ROLE_STATUSBAR,
	WIN_ROLE_POSITIONREADOUT,
	WIN_ROLE_WINDOWLIST,
	WIN_ROLE_BINDLIST,
	WIN_ROLE_DIRECTORYLIST
};

class pms_column
{
private:
	unsigned long		median;
	unsigned int		items;
public:
				pms_column(string, Item, unsigned int);
				~pms_column() {};
	string			title;
	Item			type;
	unsigned int		minlen;
	int			abslen;
	void			addmedian(unsigned int);
	unsigned int		len();
};


/*
 * Window class and derived classes
 */
class pms_window
{
protected:
	WINDOW			*handle;

	int			border[4];
	string			title;

	void			drawborders();
	void			drawtitle();
public:
	int			x;
	int			y;
	int			width;
	int			height;
	bool			stretch;

	bool			wantdraw;
	int			cursor; //songlist offset of the song the cursor is on
	int			scrolloffset; //songlist offset of top song visible (only used for normal scrolling mode)

				pms_window();
				~pms_window();

	WINDOW			*h() { return handle; };

	void			setborders(bool, bool, bool, bool); // top, right, bottom, left
	bool			resize(int, int, int, int); // x, y, width, height
	void			clear(bool, color *);
	void			settitle(string);

	int			left() { return x + border[3]; };
	int			top() { return y + border[0]; };
	int			right() { return x + width - 1 - border[1] - border[3]; };
	int			bottom() { return y + height - 1 - border[0] - border[2]; };
	int			bwidth() { return width - border[1] - border[3]; };
	int			bheight() { return height - border[0] - border[2]; };
	int			centered(string s) { return (bwidth() / 2) - (s.size() / 2); };
	int			hasborder(int i) { if (i >= 0 && i <= 3) return border[i]; else return 0; };


	virtual Songlist *	plist() { return NULL; };
	virtual void		setplist(Songlist *) {};
	virtual void		set_column_size() {};
	virtual void		draw() {};
	virtual int		posof_jump(string, int, bool = false) { return -1; };
	virtual bool		jumpto(string, int, bool = false) { return false; };
	virtual string		fulltitle() { return title; };
	virtual unsigned int	size() { return 0; };
	virtual int		type() = 0;
	virtual pms_window *	current() { return NULL; };
	virtual pms_window *	lastwin() { return NULL; };
	virtual void		switchlastwin() {};
	virtual bool		gotocurrent() { return false; };

	/*
	 * Scroll
	 */

	unsigned int		cursordrawstart();
	virtual int		scursor() { return cursor; };

	virtual void		movecursor(int);
	virtual void		setcursor(int);
	virtual void		scrollwin(int);
};

class pms_scroller
{
protected:
};




/*
 * Different types of windows
 */


class pms_win_playlist : public pms_window, public pms_scroller
{
	vector<pms_column *>	column;
public:
				pms_win_playlist();

	Songlist		*list;

	int			posof_jump(string, int, bool = false);
	bool			jumpto(string, int, bool = false);
	void			set_column_size();

	/* Virtual override */
	unsigned int		listsize() { return (list ? list->size() : 0); };
	virtual int		scursor() { return (list ? list->cursor() : 0); };
	Songlist *		plist() { return list; };
	void			setplist(Songlist *);
	unsigned int		size() { return (list ? list->size() : 0); };
	void			draw();
	void			movecursor(int);
	void			setcursor(int);
	string			fulltitle();
	bool			gotocurrent();
	int			type() { return WIN_ROLE_PLAYLIST; };
};

class pms_win_topbar : public pms_window
{
	Control *			comm;
public:
					pms_win_topbar(Control *);
	void				movecursor(int) {};
	void				setcursor(int) {};
	void				scrollwin(int) {};
	void				draw();
	int				type() { return WIN_ROLE_TOPBAR; };
	int				height();
};

class pms_win_statusbar : public pms_window
{
private:
	string				text;
public:
	void				set(statusbar_mode, string) {};
	void				movecursor(int) {};
	void				setcursor(int) {};
	void				scrollwin(int) {};
	void				draw() {};
	int				type() { return WIN_ROLE_STATUSBAR; };
};

class pms_win_positionreadout : public pms_window
{
private:
	string				text;
public:
	void				set(string) {};
	void				draw();
	int				type() { return WIN_ROLE_POSITIONREADOUT; };
};

class pms_win_windowlist : public pms_window
{
private:
	Display *			mydisp;
	vector<pms_column *>		column;
	vector<pms_win_playlist *> *	wlist;
	pms_win_playlist *		originwin;
	pms_window *			clastwin;
	pms_window *			selected;
public:
					pms_win_windowlist(Display *, vector<pms_win_playlist *> *);

	int				type() { return WIN_ROLE_WINDOWLIST; };

	pms_window *			current();
	pms_window *			lastwin();
	void				switchlastwin();

	unsigned int			size() { return (wlist ? wlist->size() : 0); };
	void				draw();
};

class pms_win_bindings : public pms_window
{
private:
	vector<pms_column *>		column;
	vector<string>			key, command, desc;
public:
					pms_win_bindings();

	int				type() { return WIN_ROLE_BINDLIST; };

	unsigned int			size() { return key.size(); };
	void				draw();
};



/*
 * Display class: manages ncurses and windows
 */
class Display
{
private:
	vector<pms_window *>		windows;
	vector<pms_win_playlist *>	playlists;
	pms_window *			curwin;		// Pointer to active window

	Control *			comm;
	mmask_t				oldmmask;
	mmask_t				mmask;

public:
	pms_win_topbar *		topbar;
	pms_win_statusbar *		statusbar;
	pms_win_positionreadout *	positionreadout;
	pms_window *			lastwin;

					Display(Control *);
					~Display();

	pms_window *			actwin() { return curwin; };
	pms_window *			playingwin();
	
	mmask_t				setmousemask();
	bool				init();
	void				uninit();
	void				resized();
	void				refresh();
	void				movecursor(int);
	void				setcursor(int);
	void				scrollwin(int);
	Song *				cursorsong();
	void				draw();
	void				forcedraw();
	void				set_xterm_title();

	pms_window *			findwlist(Songlist *);

	pms_window *			nextwindow();
	pms_window *			prevwindow();
	bool				activate(pms_window *);

	pms_win_bindings *		create_bindlist();
	pms_win_windowlist *		create_windowlist();
	pms_win_playlist *		create_playlist();
	bool				delete_window(pms_window *);
};
 
void	colprint(pms_window * w, int y, int x, color * c, const char *fmt, ...);
mmask_t	setmousemask();


#endif /* _PMS_DISPLAY_H_ */

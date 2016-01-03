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
#include "column.h"

using namespace std;


class Display;


/**
 * 2D coordinates
 */
class Point
{
public:
	uint16_t	x;
	uint16_t	y;

	Point();
	Point(uint16_t x_, uint16_t y_);

	Point &		operator=(const Point & src);
};

/**
 * Bounding box
 */
class BBox
{
private:
	Point		tl;
	Point		br;

public:
	WINDOW *	window;

	BBox();

	uint16_t	top();
	uint16_t	bottom();
	uint16_t	left();
	uint16_t	right();
	uint16_t	width();
	uint16_t	height();

	bool		clear(color * c);
	bool		refresh();
	bool		resize(const Point & tl_, const Point & br_);
};


/*
 * Display class: manages ncurses and windows
 */
class Display
{
private:
	vector<List *>		lists;

	mmask_t			oldmmask;
	mmask_t			mmask;

public:
	BBox			topbar;
	BBox			titlebar;
	BBox			main_window;
	BBox			statusbar;
	BBox			position_readout;

	List *			active_list;
	List *			last_list;

				Display(Control *);
				~Display();

	mmask_t			setmousemask();

	bool			add_list(List * list);
	bool			activate_list(List * list);
	bool			delete_list(List *);

	/**
	 * @return List* The next list of all lists. Wraps around if the active
	 * list is the last in series.
	 */
	List *			next_list();

	/**
	 * @return List* The previous list of all lists. Wraps around if the
	 * active list is the first in series.
	 */
	List *			previous_list();

	/**
	 * Given a list title, return a matching List object.
	 *
	 * Returns a pointer to a List, or NULL if not found.
	 */
	List *			find(const char * title);

	bool			draw();
	bool			draw_topbar();
	bool			draw_titlebar();
	bool			draw_main_window();
	bool			draw_position_readout();

	bool				init();
	void				uninit();
	void				resized();
	void				refresh();
	void				scrollwin(int);
	Song *				cursorsong();
	void				forcedraw();
	void				set_xterm_title();
};
 
void	colprint(BBox * bbox, int y, int x, color * c, const char *fmt, ...);
mmask_t	setmousemask();


#endif /* _PMS_DISPLAY_H_ */

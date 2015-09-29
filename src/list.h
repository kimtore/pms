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
 */


#ifndef _PMS_LIST_H_
#define _PMS_LIST_H_


#include <vector>
#include <stdint.h>

#include "listitem.h"

using namespace std;


/* Forward declaration of bounding box */
class BBox;


/**
 * Superclass for SongList, WindowList, OutputList, FileList, BindingList.
 */
class List
{
private:
	void					init();

protected:
	vector<ListItem *>			items;
	vector<ListItem *>::iterator		seliter;
	vector<ListItem *>::reverse_iterator	rseliter;
	int32_t					top_position_;

public:
					List();
					List(BBox * bbox_);
	virtual				~List();

	BBox *				bbox;
	uint32_t			cursor_position;

	bool				add(ListItem * item);
	bool				remove(uint32_t position);
	void				clear();

	uint32_t			size();
	int32_t				top_position();
	int32_t				bottom_position();
	int32_t				min_top_position();
	int32_t				max_top_position();
	bool				scroll_window(int32_t delta);
	bool				set_viewport_position(int32_t delta);
	bool				move_cursor(int32_t delta);
	bool				set_cursor(int32_t position);
	ListItem *			cursor_item();

	/**
	 * Make sure that the cursor is in the correct place according to
	 * scroll mode and viewport.
	 *
	 * Returns true if cursor position changed, false otherwise.
	 */
	bool				adjust_cursor_to_viewport();

	/**
	 * Make sure that the viewport is showing the correct items, according
	 * to the cursor position and scroll mode.
	 *
	 * Returns true if viewport position changed, false otherwise.
	 */
	bool				adjust_viewport_to_cursor();

	virtual const char *		title() = 0;

	/* FIXME: needed? */
	bool				move(uint32_t, uint32_t);

	vector<ListItem *>::iterator	begin();
	vector<ListItem *>::iterator	end();
	ListItem *			first();
	ListItem *			last();

	/* Selection iterator emulation */
	/* FIXME: rewrite to std::iterator */
	ListItem *			lastget;
	ListItem *			get_next_selected();
	ListItem *			get_prev_selected();
	ListItem *			popnextselected();
	void				resetgets();

	/**
	 * Dynamically configure the width of the columns.
	 */
	virtual void			set_column_size() = 0;

	/**
	 * Assign a bounding box.
	 */
	void				set_bounding_box(BBox * bbox_);

	/**
	 * Draw the contents of the list into bounding box.
	 */
	virtual bool			draw() = 0;
};


#endif /* _PMS_LIST_H_ */

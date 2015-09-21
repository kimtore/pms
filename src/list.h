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

using namespace std;


/* Forward declaration of bounding box */
class BBox;


/**
 * FIXME
 */
class ListItem
{
public:
	bool			selected;

	virtual			~ListItem();
};


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
	uint32_t				top_position_;

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
	uint32_t			top_position();
	uint32_t			bottom_position();
	uint32_t			min_top_position();
	uint32_t			max_top_position();
	bool				scroll_window(int32_t delta);
	bool				set_cursor(uint32_t position);
	bool				move_cursor(int32_t delta);
	virtual const char *		title() = 0;

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

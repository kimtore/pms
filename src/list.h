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
	void				init();

protected:
	vector<ListItem *>		items;
	vector<ListItem *>		selection;
	int32_t				top_position_;
	bool				selection_cache_valid_;

	/**
	 * Build a cache of selected list items, and put them in the selection
	 * vector.
	 */
	void				build_selection_cache();

	/**
	 * Returns true if the selection cache is up to date, false otherwise.
	 */
	bool				selection_cache_valid();

	/**
	 * Set the selection cache status.
	 */
	void				set_selection_cache_valid(bool state);

public:
					List();
					List(BBox * bbox_);
	virtual				~List();

	BBox *				bbox;
	int32_t				cursor_position;

	bool				add(ListItem * item);
	void				clear();

	/**
	 * Return the item at position N, or NULL if out of bounds.
	 */
	ListItem *			item(uint32_t index);

	/**
	 * Return the size of the list.
	 */
	uint32_t			size();

	/**
	 * Remove an item from the list.
	 *
	 * Returns true on success, false on failure.
	 */
	bool				remove_local(uint32_t position);

	/**
	 * Remove an item asynchronously. Used by subclasses for MPD
	 * communication when a song removal is requested.
	 */
	virtual bool			remove_async(ListItem * i) = 0;

	/**
	 * Crop the list to the list selection.
	 *
	 * Returns true on success, false on failure.
	 */
	bool				crop_to_selection();

	/**
	 * Remove the selected items from the list.
	 *
	 * Returns true on success, false on failure.
	 */
	bool				remove_selection();

	/**
	 * Return the absolute position of the list item visible in the top of
	 * the viewport.
	 */
	int32_t				top_position();

	/**
	 * Return the absolute position of the list item visible in the bottom
	 * of the viewport.
	 */
	int32_t				bottom_position();

	/**
	 * Return the minimum possible value for top_position(), according to
	 * list and viewport size.
	 */
	int32_t				min_top_position();

	/**
	 * Return the maximum possible value for top_position(), according to
	 * list and viewport size.
	 */
	int32_t				max_top_position();

	/**
	 * Scroll the viewport by the number of lines specified in delta.
	 *
	 * Returns false if the scroll went out of bounds, true otherwise.
	 */
	bool				scroll_window(int32_t delta);

	/**
	 * Set the top position of the viewport to the specified position.
	 *
	 * Returns false if the viewport went out of bounds, true otherwise.
	 */
	bool				set_viewport_position(int32_t position);

	/**
	 * Move the cursor by the number of lines specified in delta.
	 *
	 * Returns false if the cursor went out of bounds, true otherwise.
	 */
	bool				move_cursor(int32_t delta);

	/**
	 * Move the cursor to the specified position.
	 *
	 * Returns false if the cursor went out of bounds, true otherwise.
	 */
	bool				set_cursor(int32_t position);

	/**
	 * Return the list item beneath the cursor, or NULL if the list is empty.
	 */
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

	/**
	 * Return the title of this list as a C string.
	 */
	virtual const char *		title() = 0;

	/* FIXME: needed? */
	bool				move(uint32_t, uint32_t);

	vector<ListItem *>::iterator	begin();
	vector<ListItem *>::iterator	end();
	ListItem *			first();
	ListItem *			last();

	/**
	 * Start iterating over the list selection.
	 *
	 * Returns an iterator to the start of the list selection.
	 */
	vector<ListItem *>::iterator	selection_begin();

	/**
	 * Return an iterator to the end of the list selection.
	 */
	vector<ListItem *>::iterator	selection_end();

	/**
	 * Start iterating in reverse over the list selection.
	 *
	 * Returns an iterator to the reverse start of the list selection.
	 */
	vector<ListItem *>::reverse_iterator	selection_rbegin();

	/**
	 * Return an iterator to the reverse end of the list selection.
	 */
	vector<ListItem *>::reverse_iterator	selection_rend();

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

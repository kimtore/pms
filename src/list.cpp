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

#include <stdlib.h>
#include <assert.h>

#include "display.h"
#include "list.h"

List::List(BBox * bbox_)
{
	init();
	set_bounding_box(bbox_);
}

List::List()
{
	init();
	set_bounding_box(NULL);
}

List::~List()
{
}

void
List::init()
{
	seliter = items.begin();
	rseliter = items.rbegin();
	top_position_ = 0;
}

bool
List::add(ListItem * item)
{
}

/*
 * Remove item in position N from the list.
 *
 * Returns true on success, false on failure.
 */
bool
List::remove(uint32_t position)
{
	vector<ListItem *>::iterator iter;

	assert(position >= 0);
	assert(position < size());

	iter = items.begin() + position;
	assert(iter != items.end());

	delete *iter;
	items.erase(iter);

	// FIXME
	seliter = items.begin();
	rseliter = items.rbegin();

	return true;
}

/*
 * Move a list item inside the list to position dest
 */
/*
bool
List::move(uint32_t from, uint32_t to)
{
	ListItem * tmp;
	vector<ListItem *>::iterator from_iter;
	vector<ListItem *>::iterator to_iter;

	assert(from < size());
	assert(to < size());
	assert(from != to);

	from_iter = items.begin() + from;
	to_iter = items.begin() + to;

	assert(from_iter != items.end());
	assert(to_iter != items.end());

	tmp = *from_iter;

	/* Set direction FIXME
	if (dst == songpos)
		return false;
	else if (dst > songpos)
		direction = 1;
	else
		direction = -1;

	/* Swap every element on its way
	while (songpos != dst)
	{
		if (!this->swap(songpos, (songpos + direction)))
			return false;

		songpos += direction;
	}

	/* Clear queue length
	{
		qlen = 0;
		qpos = 0;
		qnum = 0;
		qsize = 0;
	}

	return true;
}
	*/

inline
uint32_t
List::size()
{
	return items.size();
}

inline
uint32_t
List::top_position()
{
	return top_position_;
}

uint32_t
List::bottom_position()
{
	uint32_t p;

	assert(bbox);

	if (!size()) {
		return 0;
	}

	p = top_position() + bbox->height();

	if (p >= size()) {
		p = size() - 1;
	}

	return p;
}

inline
uint32_t
List::min_top_position()
{
	return 0;
}

uint32_t
List::max_top_position()
{
	assert(bbox);

	if (bbox->height() > size()) {
		return 0;
	}

	return size() - bbox->height();
}

bool
List::move_cursor(int32_t delta)
{
	return set_cursor(cursor_position + delta);
}

bool
List::scroll_window(int32_t delta)
{
	top_position_ += delta;

	if (top_position_ < min_top_position()) {
		top_position_ = min_top_position();
	} else if (top_position_ < max_top_position()) {
		top_position_ = max_top_position();
	}

	return true;
}

bool
List::set_cursor(uint32_t position)
{
	cursor_position = position;
	if (cursor_position < 0) {
		cursor_position = 0;
	} else if (cursor_position >= size()) {
		cursor_position = size() - 1;
	}
}

void
List::set_bounding_box(BBox * bbox_)
{
	bbox = bbox_;
}

void
List::clear()
{
	vector<ListItem *>::iterator iter;

	iter = items.begin();

	while (iter != items.end()) {
		delete *iter;
	}

	items.clear();

	init();
}

ListItem *
List::cursor_item()
{
	if (!size()) {
		return NULL;
	}

	assert(cursor_position < size());

	return items[cursor_position];
}

inline
vector<ListItem *>::iterator
List::begin()
{
	return items.begin();
}

inline
vector<ListItem *>::iterator
List::end()
{
	return items.end();
}

ListItem *
List::first()
{
	if (!size()) {
		return NULL;
	}

	return items.front();
}

ListItem *
List::last()
{
	if (!size()) {
		return NULL;
	}

	return items.back();
}

/*
 * Returs a consecutive list of selected songs each call
 */
ListItem *
List::get_next_selected()
{
	if (lastget == NULL) {
		seliter = items.begin();
	}

	while (seliter != items.end()) {
		if ((*seliter)->selected) {
			lastget = *seliter;
			++seliter;
			return lastget;
		}
		++seliter;
	}

	/* No selection, return cursor */
	/* FIXME: override in Songlist
	if (lastget == NULL) {
		if (lastget == cursorsong()) {
			lastget = NULL;
		} else {
			lastget = cursorsong();
		}

		return lastget;
	}
	*/

	lastget = NULL;
	return NULL;
}

/*
 * Returs a consecutive list of selected songs, and unselects them
 */
ListItem *
List::get_prev_selected()
{
	if (lastget == NULL) {
		rseliter = items.rbegin();
	}

	while (rseliter != items.rend()) {
		if ((*rseliter)->selected) {
			lastget = *rseliter;
			++rseliter;
			return lastget;
		}
		++rseliter;
	}

	/* No selection, return cursor */
	/* FIXME: override in Songlist
	if (lastget == NULL) {
		if (lastget == cursorsong()) {
			lastget = NULL;
		} else {
			lastget = cursorsong();
		}

		return lastget;
	}
	*/

	lastget = NULL;
	return NULL;
}

/*
 * Returs a consecutive list of selected songs, and unselects them
 */
ListItem *
List::popnextselected()
{
	ListItem * item;

	item = get_next_selected();
	if (item) {
		item->selected = false;
	}
	return item;
}

/*
 * Reset iterators
 */
void
List::resetgets()
{
	lastget = NULL;
	seliter = items.begin();
	rseliter = items.rbegin();
}


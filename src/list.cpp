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
#include "pms.h"

extern Pms * pms;

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
	clear();
}

void
List::init()
{
	top_position_ = 0;
	cursor_position = 0;
}

void
List::remove_local(uint32_t position)
{
	vector<ListItem *>::iterator iter;

	assert(position >= 0);
	assert(position < size());

	iter = items.begin() + position;
	assert(iter != items.end());
	assert(*iter);

	(*iter)->set_selected(false);
	delete *iter;
	items.erase(iter);

	if (cursor_position >= size()) {
		set_cursor(size() - 1);
	}

	set_selection_cache_valid(false);
}

bool
List::crop_to_selection()
{
	vector<ListItem *>::reverse_iterator iter;
	ListItem * item;
	Song * song;

	iter = selection_rbegin();
	while (iter != selection_rend()) {
		if (!item->selected()) {
			if (!remove(item)) {
				return false;
			}
		} else {
			item->set_selected(false);
		}
		++iter;
	}

	return true;
}

bool
List::remove_selection()
{
	vector<ListItem *>::reverse_iterator iter;

	iter = selection_rbegin();
	while (iter != selection_rend()) {
		if (!remove(*iter)) {
			return false;
		}
		++iter;
	}

	return true;
}


/*
 * Move a list item inside the list to position dest
 */
bool
List::move(uint32_t from, uint32_t to)
{
/*
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

	// Set direction FIXME
	if (dst == songpos)
		return false;
	else if (dst > songpos)
		direction = 1;
	else
		direction = -1;

	// Swap every element on its way
	while (songpos != dst)
	{
		if (!this->swap(songpos, (songpos + direction)))
			return false;

		songpos += direction;
	}

	// Clear queue length
	{
		qlen = 0;
		qpos = 0;
		qnum = 0;
		qsize = 0;
	}

	return true;
	*/
}

ListItem *
List::item(uint32_t index)
{
	if (index >= size()) {
		return NULL;
	}

	return items[index];
}

inline
uint32_t
List::size()
{
	return items.size();
}

inline
int32_t
List::top_position()
{
	return top_position_;
}

int32_t
List::bottom_position()
{
	int32_t p;

	assert(bbox);

	if (!size()) {
		return 0;
	}

	p = top_position() + bbox->height() - 2;

	if (p >= size()) {
		p = size() - 1;
	}

	return p;
}

inline
int32_t
List::min_top_position()
{
	return 0;
}

int32_t
List::max_top_position()
{
	assert(bbox);

	if (bbox->height() > size()) {
		return 0;
	}

	return size() - bbox->height() + 1;
}

bool
List::scroll_window(int32_t delta)
{
	return set_viewport_position(top_position_ + delta);
}

bool
List::set_viewport_position(int32_t position)
{
	bool rc = true;

	top_position_ = position;
	if (top_position_ < min_top_position()) {
		top_position_ = min_top_position();
		rc = false;
	} else if (top_position_ > max_top_position()) {
		top_position_ = max_top_position();
		rc = false;
	}

	adjust_cursor_to_viewport();

	return rc;
}

bool
List::move_cursor(int32_t delta)
{
	return set_cursor(cursor_position + delta);
}

bool
List::set_cursor(int32_t position)
{
	bool rc = true;

	cursor_position = position;
	if (cursor_position < 0) {
		cursor_position = 0;
		rc = false;
	} else if (cursor_position >= size()) {
		cursor_position = size() - 1;
		rc = false;
	}

	adjust_viewport_to_cursor();

	return rc;
}

bool
List::adjust_cursor_to_viewport()
{
	int32_t new_position;

	switch (pms->options->scroll_mode) {
		case SCROLL_NORMAL:
			if (cursor_position < top_position()) {
				cursor_position = top_position();
			} else if (cursor_position > bottom_position()) {
				cursor_position = bottom_position();
			} else {
				return false;
			}
			return true;

		case SCROLL_CENTERED:
			new_position = top_position() + (bbox->height() / 2);
			if (new_position >= size()) {
				if (cursor_position >= size()) {
					new_position = size() - 1;
				} else {
					new_position = cursor_position;
				}
			}
			if (new_position == cursor_position) {
				return false;
			}
			cursor_position = new_position;
			return true;

		default:
			abort();
	}
}

bool
List::adjust_viewport_to_cursor()
{
	int32_t new_position;

	switch (pms->options->scroll_mode) {
		case SCROLL_NORMAL:
			if (cursor_position < top_position()) {
				top_position_ = cursor_position;
			} else if (cursor_position > bottom_position()) {
				top_position_ = cursor_position - bbox->height() + 2;
			} else {
				return false;
			}
			return true;

		case SCROLL_CENTERED:
			new_position = cursor_position - (bbox->height() / 2);
			if (new_position < min_top_position()) {
				new_position = 0;
			} if (new_position > max_top_position()) {
				new_position = max_top_position();
			}
			if (new_position == top_position()) {
				return false;
			}
			top_position_ = new_position;
			return true;

		default:
			abort();
	}
}

const char *
List::title()
{
	return title_.c_str();
}

void
List::set_title(string new_title)
{
	title_ = new_title;
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
		++iter;
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

vector<ListItem *>::iterator
List::begin()
{
	return items.begin();
}

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

void
List::build_selection_cache()
{
	vector<ListItem *>::const_iterator iter;
	ListItem * item;

	selection.clear();

	iter = items.begin();
	while (iter != items.end()) {
		if ((*iter)->selected()) {
			selection.push_back(*iter);
		}
		++iter;
	}

	if (!selection.size()) {
		if ((item = cursor_item()) != NULL) { 
			selection.push_back(item);
		}
	}

	set_selection_cache_valid(true);
}

inline
bool
List::selection_cache_valid()
{
	return selection_cache_valid_;
}

void
List::set_selection_cache_valid(bool state)
{
	selection_cache_valid_ = state;
}

vector<ListItem *>::iterator
List::selection_begin()
{
	if (!selection_cache_valid()) {
		build_selection_cache();
	}

	return selection.begin();
}

vector<ListItem *>::iterator
List::selection_end()
{
	return selection.end();
}

vector<ListItem *>::reverse_iterator
List::selection_rbegin()
{
	if (!selection_cache_valid()) {
		build_selection_cache();
	}

	return selection.rbegin();
}

vector<ListItem *>::reverse_iterator
List::selection_rend()
{
	return selection.rend();
}

ListItem *
List::match(string pattern, unsigned int from, unsigned int to, long flags)
{
	ListItem * it;
	int i;

	if (!size()) {
		return NULL;
	}

	assert(from < size());
	assert(to < size());

	if (flags & MATCH_REVERSE) {
		i = from;
		from = to;
		to = i;
	}

	i = from;

	while (true)
	{
		if (i < 0) {
			i = size() - 1;
		} else if (i >= size()) {
			i = 0;
		}

		it = item(i);

		assert(it);

		if (it->match(pattern, flags)) {
			return it;
		}

		if (i == to) {
			break;
		}

		i += (flags & MATCH_REVERSE ? -1 : 1);
	}

	return NULL;
}

ListItem *
List::match_until_cursor(string pattern, long flags)
{
	int32_t from;
	int32_t to;

	if (!size()) {
		return NULL;
	}

	from = cursor_position;
	if (from == 0) {
		to = size() - 1;
	} else {
		to = from - 1;
	}

	return match(pattern, from, to, flags);
}

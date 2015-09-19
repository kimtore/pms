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
#include "list.h"

List::List(BBox * bbox_)
{
	bbox = bbox_;
}

List::List()
{
	bbox = NULL;
}

List::~List()
{
}

inline
uint32_t
List::size()
{
	return items.size();
}

bool
List::move_cursor(int32_t delta)
{
	return set_cursor(cursor_position + delta);
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

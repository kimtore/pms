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
};


/**
 * Superclass for SongList, WindowList, OutputList, FileList, BindingList.
 */
class List
{
private:
	vector<ListItem *>	items;

public:
				List();
				List(BBox * bbox_);
	virtual			~List();

	BBox *			bbox;
	uint32_t		cursor_position;
	uint32_t		top_position;

	uint32_t		size();
	bool			set_cursor(uint32_t position);
	bool			move_cursor(int32_t delta);

	void			set_column_size();
};


#endif /* _PMS_LIST_H_ */

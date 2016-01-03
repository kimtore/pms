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

#ifndef _PMS_LISTITEM_H_
#define _PMS_LISTITEM_H_

#include <string>

using namespace std;

class List;

class ListItem
{
private:
	bool			selected_;

protected:
	/**
	 * Pointer to a List owning this ListItem.
	 */
	List *			list;

public:
				ListItem(List * owner);
	virtual			~ListItem();

	/**
	 * Set whether this ListItem is selected in the list context.
	 */
	void			set_selected(bool state);

	/**
	 * @return Selection state.
	 */
	bool			selected();

	/**
	 * Match this ListItem against search criteria. This method should be
	 * overridden in subclasses.
	 *
	 * Returns true if the match succeeds, false otherwise.
	 */
	virtual bool		match(string term, long flags);
};

#endif /* _PMS_LISTITEM_H_ */

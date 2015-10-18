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

#ifndef _PMS_LIBRARY_H_
#define _PMS_LIBRARY_H_

#include "songlist.h"

#define LIBRARY(x) dynamic_cast<Library *>(x)

class Library : public Songlist
{
public:
	/**
	 * Called when the user requests to remove a song from the MPD library.
	 * The library is read-only, so an error is returned.
	 *
	 * Always returns false.
	 */
	bool			remove_async(ListItem * i);
};

#endif /* _PMS_LIBRARY_H_ */

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

#ifndef _PMS_QUEUE_H_
#define _PMS_QUEUE_H_

#include "songlist.h"

#define QUEUE(x) dynamic_cast<Queue *>(x)

class Queue : public Songlist
{
public:
				Queue();

	/**
	 * Remove a song from the MPD queue.
	 *
	 * Returns true on success, false on failure.
	 */
	bool			remove(ListItem * i);
};

#endif /* _PMS_QUEUE_H_ */

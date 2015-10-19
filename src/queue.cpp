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

#include "queue.h"
#include "pms.h"

#include <mpd/client.h>

extern Pms * pms;

#define EXIT_IDLE if (!pms->comm->exit_idle()) { return false; }

bool
Queue::remove(ListItem * i)
{
	ListItemSong * item_song;
	
	item_song = LISTITEMSONG(i);
	assert(item_song->song->id != MPD_SONG_NO_ID);

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Removing song from queue: id=%d pos=%d uri=%s\n", item_song->song->id, item_song->song->pos, item_song->song->file.c_str());

	return mpd_run_delete_id(pms->conn->h(), item_song->song->id);
}

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

#ifndef _PMS_PLAYLIST_H_
#define _PMS_PLAYLIST_H_


#include "songlist.h"
#include <mpd/client.h>

using namespace std;


/**
 * Class representing a playlist in MPD
 */
class Playlist : public Songlist
{
private:
	time_t			_last_modified;
	bool			_synchronized;
	bool			_exists_in_mpd;

public:
	Playlist();

	void			assign_metadata_from_mpd(mpd_playlist * playlist);
	time_t			get_last_modified();
	bool			is_synchronized();
	void			set_synchronized(bool is_sync);
	bool			exists_in_mpd();
	void			set_exists_in_mpd(bool exists);

	/**
	 * Remove a song within this playlist from the MPD version.
	 *
	 * Returns true on success, false on failure.
	 */
	bool			remove_async(ListItem * i);
};

#endif /* _PMS_PLAYLIST_H_ */

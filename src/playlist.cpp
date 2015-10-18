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

#include <mpd/client.h>

#include "playlist.h"
#include "pms.h"

extern Pms * pms;

Playlist::Playlist()
{
	_last_modified = 0;
	_synchronized = false;
	_exists_in_mpd = true;
}

/**
 * Return playlist modification time.
 */
time_t
Playlist::get_last_modified()
{
	return _last_modified;
}

/**
 * Set filename and last_modified tags from MPD playlist.
 */
void
Playlist::assign_metadata_from_mpd(mpd_playlist * playlist)
{
	time_t last_mod = _last_modified;

	filename = mpd_playlist_get_path(playlist);
	_last_modified = mpd_playlist_get_last_modified(playlist);
	_exists_in_mpd = true;

	if (_last_modified > last_mod) {
		_synchronized = false;
	}
}

/**
 * Get synchronization state.
 *
 * Returns true if playlist is up-to-date with MPD, false if not.
 */
bool
Playlist::is_synchronized()
{
	return _synchronized;
}

/**
 * Set synchronization state.
 */
void
Playlist::set_synchronized(bool is_sync)
{
	_synchronized = is_sync;
}

/**
 * Return true if playlist exists in MPD, false if not.
 */
bool
Playlist::exists_in_mpd()
{
	return _exists_in_mpd;
}

/**
 * Set existing state.
 */
void
Playlist::set_exists_in_mpd(bool exists)
{
	_exists_in_mpd = exists;
}

bool
Playlist::remove_async(Song * song)
{
	assert(song);
	return mpd_run_playlist_delete(pms->conn->h(), filename.c_str(), song->pos);
}

/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS  <<Practical Music Search>>
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
 *
 * command.cpp
 * 	mediates all commands between UI and mpd
 */

#include <unistd.h>
#include <math.h>
#include <mpd/client.h>

#include "command.h"
#include "pms.h"

extern Pms *			pms;

#define EXIT_IDLE		if (is_idle() && (!noidle() || !wait_until_noidle())) { return false; }

#define NOIDLE_POLL_TIMEOUT	20  /* time to wait for zmq_poll() to finish after calling noidle() */


/*
 * Status class
 */
Mpd_status::Mpd_status()
{
	muted			= false;
	volume			= 0;
	repeat			= false;
	single			= false;
	random			= false;
	playlist_length		= 0;
	playlist		= -1;
	state			= MPD_STATE_UNKNOWN;
	crossfade		= 0;
	song			= MPD_SONG_NO_NUM;
	songid			= MPD_SONG_NO_ID;
	time_elapsed		= 0;
	time_elapsed_hires.tv_sec = 0;
	time_elapsed_hires.tv_nsec = 0;
	time_total		= 0;
	db_updating		= false;
	error			= 0;
	errstr.clear();

	bitrate			= 0;
	samplerate		= 0;
	bits			= 0;
	channels		= 0;

	artists_count		= 0;
	albums_count		= 0;
	songs_count		= 0;

	uptime			= 0;
	db_update_time		= 0;
	playtime		= 0;
	db_playtime		= 0;

	last_playlist		= playlist;
	last_db_update_time	= db_update_time;
	last_db_updating	= db_updating;
	update_job_id		= -1;
}

void
Mpd_status::assign_status(struct mpd_status * status)
{
	const struct mpd_audio_format	*format;
	uint32_t ms;

	volume			= mpd_status_get_volume(status);
	repeat			= mpd_status_get_repeat(status);
	single			= mpd_status_get_single(status);
	random			= mpd_status_get_random(status);
	playlist_length		= mpd_status_get_queue_length(status);
	playlist		= mpd_status_get_queue_version(status);
	state			= mpd_status_get_state(status);
	crossfade		= mpd_status_get_crossfade(status);
	song			= mpd_status_get_song_pos(status);
	songid			= mpd_status_get_song_id(status);
	time_total		= mpd_status_get_total_time(status);
	db_updating		= mpd_status_get_update_id(status);

	/* Time elapsed */
	ms = mpd_status_get_elapsed_ms(status);
	set_time_elapsed_ms(ms);

	/* Audio format */
	bitrate			= mpd_status_get_kbit_rate(status);
	format			= mpd_status_get_audio_format(status);

	if (!format) {
		return;
	}

	samplerate		= format->sample_rate;
	bits			= format->bits;
	channels		= format->channels;
}

void
Mpd_status::assign_stats(struct mpd_stats * stats)
{
	artists_count		= mpd_stats_get_number_of_artists(stats);
	albums_count		= mpd_stats_get_number_of_albums(stats);
	songs_count		= mpd_stats_get_number_of_songs(stats);

	uptime			= mpd_stats_get_uptime(stats);
	db_update_time		= mpd_stats_get_db_update_time(stats);
	playtime		= mpd_stats_get_play_time(stats);
	db_playtime		= mpd_stats_get_db_play_time(stats);
}

void
Mpd_status::set_time_elapsed_ms(uint32_t ms)
{
	time_elapsed_hires.tv_sec = ms / 1000;
	time_elapsed_hires.tv_nsec = (ms * 10e5) - (time_elapsed_hires.tv_sec * 10e8);
	time_elapsed = round(time_elapsed_hires.tv_sec);
	// pms->log(MSG_DEBUG, 0, "Time elapsed %dms converted to %lus %luns\n", ms, time_elapsed_hires.tv_sec, time_elapsed_hires.tv_nsec);
}

void
Mpd_status::increase_time_elapsed(struct timespec ts)
{
	time_t seconds;

	// pms->log(MSG_DEBUG, 0, "Increasing time elapsed by %lus %9luns\n", ts.tv_sec, ts.tv_nsec);
	time_elapsed_hires.tv_sec += ts.tv_sec;
	time_elapsed_hires.tv_nsec += ts.tv_nsec;

	seconds = time_elapsed_hires.tv_nsec / 10e8;
	if (seconds > 0) {
		time_elapsed_hires.tv_sec += seconds;
		time_elapsed_hires.tv_nsec -= (seconds * 10e8);
	}

	time_elapsed = round(time_elapsed_hires.tv_sec);
	// pms->log(MSG_DEBUG, 0, "Time elapsed set to %lus %9luns\n", time_elapsed_hires.tv_sec, time_elapsed_hires.tv_nsec);
}

bool
Mpd_status::alive() const
{
	/* FIXME: what is this? */
	assert(0);
}



/*
 * Command class manages commands sent to and from mpd
 */
Control::Control(Connection * n_conn)
{
	conn = n_conn;
	st = new Mpd_status();
	rootdir = new Directory(NULL, "");
	_song = NULL;
	st->last_playlist = -1;
	last_song = MPD_SONG_NO_NUM;
	oldsong = MPD_SONG_NO_NUM;
	_playlist = new Songlist;
	_library = new Songlist;
	_playlist->role = LIST_ROLE_MAIN;
	_library->role = LIST_ROLE_LIBRARY;
	_active = NULL;
	_is_idle = false;
	command_mode = 0;
	mutevolume = 0;
	crossfadetime = pms->options->get_long("crossfade");

	/* Set all bits in mpd_idle event */
	set_mpd_idle_events((enum mpd_idle) 0xffffffff);
	finished_idle_events = 0;
}

Control::~Control()
{
	delete _library;
	delete _playlist;
	delete st;
}

/*
 * Finishes a command and debugs any errors
 */
bool		Control::finish()
{
	// FIXME: this function must die
	assert(0);

	/*
	mpd_finishCommand(conn->h());
	st->error = conn->h()->error;
	st->errstr = conn->h()->errorStr;

	if (st->error != 0)
	{
		pms->log(MSG_CONSOLE, STERR, "MPD returned error %d: %s\n", st->error, st->errstr.c_str());

		// Connection closed
		if (st->error == MPD_ERROR_CONNCLOSED)
		{
			conn->disconnect();
		}

		clearerror();

		return false;
	}

	return true;
	*/
}

/*
 * Clears any error
 */
void
Control::clearerror()
{
	mpd_connection_clear_error(conn->h());
}

/*
 * Have a usable connection?
 */
bool		Control::alive()
{
	return (conn != NULL && conn->connected());
}

/*
 * Reports any error from the last command
 */
const char *	Control::err()
{
	static char * buffer = static_cast<char *>(malloc(1024));

	if (st->errstr.size() == 0)
	{
		if (pms->msg->code == 0)
			sprintf(buffer, _("Error: %s"), pms->msg->str.c_str());
		else
			sprintf(buffer, _("Error %d: %s"), pms->msg->code, pms->msg->str.c_str());

		return buffer;
	}

	return st->errstr.c_str();
}

/**
 * Return the success or error status of the last MPD command sent.
 */
bool
Control::get_error_bool()
{
	return (mpd_connection_get_error(conn->h()) == MPD_ERROR_SUCCESS);
}

/**
 * Set pending updates based on which IDLE events were returned from the server.
 */
void
Control::set_mpd_idle_events(enum mpd_idle idle_reply)
{
	uint32_t event = 1;
	const char *idle_name;
	char buffer[2048];
	char *ptr = buffer;

	idle_events |= idle_reply;

	/* Code below only prints debug statement. TODO: return if not debugging? */
	do {
		idle_name = mpd_idle_name((enum mpd_idle) event);
		if (!idle_name) {
			break;
		}
		if (!(idle_reply & event)) {
			continue;
		}
		ptr += sprintf(ptr, "%s ", idle_name);

	} while(event = event << 1);

	*ptr = '\0';

	pms->log(MSG_DEBUG, 0, "Set pending MPD IDLE events: %s\n", buffer);
}

/**
 * Run all pending updates.
 */
bool
Control::run_pending_updates()
{
	/* MPD has new current song */
	if (idle_events & MPD_IDLE_PLAYER) {
		if (!get_current_playing()) {
			return false;
		}
		/* MPD_IDLE_PLAYER will be subtracted below */
	}

	/* MPD has new status information */
	if (idle_events & MPD_IDLE_PLAYER || idle_events & MPD_IDLE_MIXER || idle_events & MPD_IDLE_OPTIONS) {
		if (!get_status()) {
			return false;
		}
		set_update_done(MPD_IDLE_PLAYER);
		set_update_done(MPD_IDLE_MIXER);
		set_update_done(MPD_IDLE_OPTIONS);
	}

	/* MPD has new playlist */
	if (idle_events & MPD_IDLE_QUEUE) {
		if (!update_playlist()) {
			return false;
		}
		set_update_done(MPD_IDLE_QUEUE);
	}

	/* MPD has new song database */
	if (idle_events & MPD_IDLE_DATABASE) {
		if (!update_library()) {
			return false;
		}
		set_update_done(MPD_IDLE_DATABASE);
	}

	return true;
}

/**
 * Mark an MPD IDLE update as retrieved.
 */
void
Control::set_update_done(enum mpd_idle flags)
{
	idle_events &= ~flags;
	finished_idle_events |= flags;
}

/**
 * Check whether an MPD IDLE update is retrieved.
 */
bool
Control::has_finished_update(enum mpd_idle flags)
{
	return (finished_idle_events & flags);
}

/**
 * Remove a finished MPD IDLE update.
 */
void
Control::clear_finished_update(enum mpd_idle flags)
{
	finished_idle_events &= ~flags;
}

/*
 * Return authorisation level in mpd server
 */
int		Control::authlevel()
{
	int	a = AUTH_NONE;

	if (commands.update)
		a += AUTH_ADMIN;
	if (commands.play)
		a += AUTH_CONTROL;
	if (commands.add)
		a += AUTH_ADD;
	if (commands.status)
		a += AUTH_READ;

	return a;
}

/*
 * Retrieve available commands
 */
bool
Control::get_available_commands()
{
	mpd_pair * pair;

	EXIT_IDLE;

	if (!mpd_send_allowed_commands(conn->h())) {
		return false;
	}

	memset(&commands, 0, sizeof(commands));

	while ((pair = mpd_recv_command_pair(conn->h())) != NULL)
	{
		// FIXME: any other response is not expected
		assert(!strcmp(pair->name, "command"));

		if (!strcmp(pair->value, "add")) {
			commands.add = true;
		} else if (!strcmp(pair->value, "addid")) {
			commands.addid = true;
		} else if (!strcmp(pair->value, "clear")) {
			commands.clear = true;
		} else if (!strcmp(pair->value, "clearerror")) {
			commands.clearerror = true;
		} else if (!strcmp(pair->value, "close")) {
			commands.close = true;
		} else if (!strcmp(pair->value, "commands")) {
			commands.commands = true;
		} else if (!strcmp(pair->value, "count")) {
			commands.count = true;
		} else if (!strcmp(pair->value, "crossfade")) {
			commands.crossfade = true;
		} else if (!strcmp(pair->value, "currentsong")) {
			commands.currentsong = true;
		} else if (!strcmp(pair->value, "delete")) {
			commands.delete_ = true;
		} else if (!strcmp(pair->value, "deleteid")) {
			commands.deleteid = true;
		} else if (!strcmp(pair->value, "disableoutput")) {
			commands.disableoutput = true;
		} else if (!strcmp(pair->value, "enableoutput")) {
			commands.enableoutput = true;
		} else if (!strcmp(pair->value, "find")) {
			commands.find = true;
		} else if (!strcmp(pair->value, "idle")) {
			commands.idle = true;
		} else if (!strcmp(pair->value, "kill")) {
			commands.kill = true;
		} else if (!strcmp(pair->value, "list")) {
			commands.list = true;
		} else if (!strcmp(pair->value, "listall")) {
			commands.listall = true;
		} else if (!strcmp(pair->value, "listallinfo")) {
			commands.listallinfo = true;
		} else if (!strcmp(pair->value, "listplaylist")) {
			commands.listplaylist = true;
		} else if (!strcmp(pair->value, "listplaylistinfo")) {
			commands.listplaylistinfo = true;
		} else if (!strcmp(pair->value, "listplaylists")) {
			commands.listplaylists = true;
		} else if (!strcmp(pair->value, "load")) {
			commands.load = true;
		} else if (!strcmp(pair->value, "lsinfo")) {
			commands.lsinfo = true;
		} else if (!strcmp(pair->value, "move")) {
			commands.move = true;
		} else if (!strcmp(pair->value, "moveid")) {
			commands.moveid = true;
		} else if (!strcmp(pair->value, "next")) {
			commands.next = true;
		} else if (!strcmp(pair->value, "notcommands")) {
			commands.notcommands = true;
		} else if (!strcmp(pair->value, "outputs")) {
			commands.outputs = true;
		} else if (!strcmp(pair->value, "password")) {
			commands.password = true;
		} else if (!strcmp(pair->value, "pause")) {
			commands.pause = true;
		} else if (!strcmp(pair->value, "ping")) {
			commands.ping = true;
		} else if (!strcmp(pair->value, "play")) {
			commands.play = true;
		} else if (!strcmp(pair->value, "playid")) {
			commands.playid = true;
		} else if (!strcmp(pair->value, "playlist")) {
			commands.playlist = true;
		} else if (!strcmp(pair->value, "playlistadd")) {
			commands.playlistadd = true;
		} else if (!strcmp(pair->value, "playlistclear")) {
			commands.playlistclear = true;
		} else if (!strcmp(pair->value, "playlistdelete")) {
			commands.playlistdelete = true;
		} else if (!strcmp(pair->value, "playlistfind")) {
			commands.playlistfind = true;
		} else if (!strcmp(pair->value, "playlistid")) {
			commands.playlistid = true;
		} else if (!strcmp(pair->value, "playlistinfo")) {
			commands.playlistinfo = true;
		} else if (!strcmp(pair->value, "playlistmove")) {
			commands.playlistmove = true;
		} else if (!strcmp(pair->value, "playlistsearch")) {
			commands.playlistsearch = true;
		} else if (!strcmp(pair->value, "plchanges")) {
			commands.plchanges = true;
		} else if (!strcmp(pair->value, "plchangesposid")) {
			commands.plchangesposid = true;
		} else if (!strcmp(pair->value, "previous")) {
			commands.previous = true;
		} else if (!strcmp(pair->value, "random")) {
			commands.random = true;
		} else if (!strcmp(pair->value, "rename")) {
			commands.rename = true;
		} else if (!strcmp(pair->value, "repeat")) {
			commands.repeat = true;
		} else if (!strcmp(pair->value, "single")) {
			commands.single = true;
		} else if (!strcmp(pair->value, "rm")) {
			commands.rm = true;
		} else if (!strcmp(pair->value, "save")) {
			commands.save = true;
		} else if (!strcmp(pair->value, "filter")) {
			commands.filter = true;
		} else if (!strcmp(pair->value, "seek")) {
			commands.seek = true;
		} else if (!strcmp(pair->value, "seekid")) {
			commands.seekid = true;
		} else if (!strcmp(pair->value, "setvol")) {
			commands.setvol = true;
		} else if (!strcmp(pair->value, "shuffle")) {
			commands.shuffle = true;
		} else if (!strcmp(pair->value, "stats")) {
			commands.stats = true;
		} else if (!strcmp(pair->value, "status")) {
			commands.status = true;
		} else if (!strcmp(pair->value, "stop")) {
			commands.stop = true;
		} else if (!strcmp(pair->value, "swap")) {
			commands.swap = true;
		} else if (!strcmp(pair->value, "swapid")) {
			commands.swapid = true;
		} else if (!strcmp(pair->value, "tagtypes")) {
			commands.tagtypes = true;
		} else if (!strcmp(pair->value, "update")) {
			commands.update = true;
		} else if (!strcmp(pair->value, "urlhandlers")) {
			commands.urlhandlers = true;
		} else if (!strcmp(pair->value, "volume")) {
			commands.volume = true;
		}

		mpd_return_pair(conn->h(), pair);
	}

	return get_error_bool();
}

/*
 * Play, pause, toggle, stop, next, prev
 */
bool
Control::play()
{
	EXIT_IDLE;

	return mpd_run_play(conn->h());
}

bool
Control::playid(song_t songid)
{
	EXIT_IDLE;

	return mpd_run_play_id(conn->h(), songid);
}

bool
Control::playpos(song_t songpos)
{
	EXIT_IDLE;

	return mpd_run_play_pos(conn->h(), songpos);
}

bool
Control::pause(bool tryplay)
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Toggling pause, tryplay=%d\n", tryplay);

	switch(st->state)
	{
		case MPD_STATE_PLAY:
			return mpd_run_pause(conn->h(), true);
		case MPD_STATE_PAUSE:
			return mpd_run_pause(conn->h(), false);
		case MPD_STATE_STOP:
		case MPD_STATE_UNKNOWN:
		default:
			return (tryplay && play());
	}
}

bool
Control::stop()
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Stopping playback.\n");

	return mpd_run_stop(conn->h());
}

/*
 * Shuffles the playlist.
 */
bool
Control::shuffle()
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Shuffling playlist.\n");

	return mpd_run_shuffle(conn->h());
}

/*
 * Sets repeat mode
 */
bool
Control::repeat(bool on)
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Set repeat to %d\n", on);

	return mpd_run_repeat(conn->h(), on);
}

/*
 * Sets single mode
 */
bool
Control::single(bool on)
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Set single to %d\n", on);

	return mpd_run_single(conn->h(), on);
}

/*
 * Set the volume to an integer between 0 and 100.
 *
 * Return true on success, false on failure.
 */
bool
Control::setvolume(int vol)
{
	if (vol < 0) {
		vol = 0;
	} else if (vol > 100) {
		vol = 100;
	}

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Setting volume to %d%%\n", vol);

	if (mpd_run_set_volume(conn->h(), vol)) {
		mutevolume = 0;
	}

	return get_error_bool();
}

/*
 * Changes volume
 */
bool
Control::volume(int offset)
{
	return setvolume(st->volume + offset);
}

/*
 * Mute/unmute volume
 */
bool
Control::mute()
{
	if (muted()) {
		return setvolume(mutevolume);
	}

	mutevolume = st->volume;
	return setvolume(0);
}

/*
 * Is muted?
 */
bool
Control::muted()
{
	return (st->volume == -1 || mutevolume != 0);
}

/*
 * Toggles MPDs built-in random mode
 */
bool
Control::random(int set)
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Set random to %d\n", set);

	if (set == -1) {
		set = (st->random == false ? true : false);
	}

	return mpd_run_random(conn->h(), set);
}

/*
 * Appends a playlist to another playlist
 */
song_t		Control::add(Songlist * source, Songlist * dest)
{
	song_t			first = MPD_SONG_NO_ID;
	song_t			result;
	unsigned int		i;

	assert(source != NULL);
	assert(dest != NULL);

	for (i = 0; i < source->size(); i++)
	{
		result = add(dest, source->song(i));
		if (result == MPD_SONG_NO_ID) {
			return result;
		}
		if (first == MPD_SONG_NO_ID) {
			first = result;
		}
	}

	return first;
}

/*
 * Add a song to a playlist
 * FIXME: return value
 */
song_t
Control::add(Songlist * list, Song * song)
{
	song_t		i = MPD_SONG_NO_ID;
	Song *		nsong;

	assert(list != NULL);
	assert(song != NULL);

	EXIT_IDLE;

	if (list == _playlist) {
		return mpd_run_add_id(conn->h(), song->file.c_str());
	} else if (list != _library) {
		if (list->filename.size() == 0) {
			return i;
		}
		return mpd_run_playlist_add(conn->h(), list->filename.c_str(), song->file.c_str());
	}

	return i;

	/* FIXME
	if (command_mode != 0) return i;
	if (finish())
	{
		nsong = new Song(song);
		if (list == _playlist)
		{
			nsong->id = i;
			nsong->pos = playlist()->size();
			increment();
		}
		else
		{
			nsong->id = MPD_SONG_NO_ID;
			nsong->pos = MPD_SONG_NO_NUM;
			i = list->size();
		}
		list->add(nsong);
	}

	return i;
	*/
}

/*
 * Remove a song from the playlist
 */
int
Control::remove(Songlist * list, Song * song)
{
	int		pos = MATCH_FAILED;

	assert(song != NULL);
	assert(list != NULL);

	if (list == _library) {
		// FIXME: error message
		return false;
	}

	if (list == _playlist && song->id == MPD_SONG_NO_ID) {
		// All songs must have ID's
		// FIXME: version requirement
		assert(song->id != MPD_SONG_NO_ID);
	}

	if (list != _playlist) {
		if (list->filename.size() == 0) {
			// FIXME: what does this check?
			return false;
		}
		pos = list->locatesong(song);
		if (pos == MATCH_FAILED) {
			// FIXME: error message
			return false;
		}
		pms->log(MSG_DEBUG, 0, "Removing song %d from list.\n", pos);
	}

	if (list == _playlist) {
		return mpd_run_delete_id(conn->h(), song->id);
	} else {
		return mpd_run_playlist_delete(conn->h(), (char *)list->filename.c_str(), pos);
	}

	// FIXME: remove from list?
	/*
	if (command_mode != 0) return true;
	if (finish())
	{
		list->remove(pos == MATCH_FAILED ? song->pos : pos);
		if (list == _playlist)
			increment();
		return true;
	}

	return false;
	*/
}

/*
 * Crops the playlist
 * FIXME: de-duplicate
 * FIXME: split into two functions
 */
bool
Control::crop(Songlist * list, int mode)
{
	unsigned int		i;
	int			pos;
	unsigned int		upos;
	Song *			song;

	assert(list != NULL);

	if (list == _library) {
		pms->msg->assign(STOK, _("The library is read-only."));
		return false;
	}

	/* Crop to currently playing song */
	if (mode == CROP_PLAYING)
	{
		song = pms->cursong();
		if (!song)
		{
			pms->msg->assign(STOK, _("No song is playing: can't crop to playing song."));
			return false;
		}

		pos = list->match(song->file, 0, list->end(), MATCH_FILE | MATCH_EXACT);
		if (pos == MATCH_FAILED)
		{
			pms->msg->assign(STOK, _("The currently playing song is not in this list."));
			return false;
		}
		upos = static_cast<unsigned int>(pos);

		//list_start();
		for (i = list->end(); i < list->size(); i--)
		{
			if (upos != i) {
				if (list == _playlist) {
					mpd_run_delete_id(conn->h(), list->song(i)->id);
				} else {
					mpd_run_playlist_delete(conn->h(), list->filename.c_str(), static_cast<int>(i));
				}
				list->remove(i);
				increment();
			}
		}
		return list_end();
	}
	/* Crop to selection */
	else if (mode == CROP_SELECTION)
	{
		list->resetgets();
		if (list->getnextselected() == list->cursorsong())
		{
			if (list->getnextselected() == NULL)
			{
				list->selectsong(list->cursorsong(), true);
			}
		}

		//list_start();
		for (i = list->end(); i < list->size(); i--)
		{
			if (list->song(i)->selected == false)
			{
				if (list == _playlist) {
					mpd_run_delete_id(conn->h(), list->song(i)->id);
				} else {
					mpd_run_playlist_delete(conn->h(), list->filename.c_str(), static_cast<int>(i));
				}
				list->remove(i);
				increment();
			}
			else
			{
				list->selectsong(list->song(i), false);
			}
		}
		return list_end();
	}

	return false;
}

/*
 * Clears the playlist
 */
int
Control::clear(Songlist * list)
{
	bool r;

	assert(list != NULL);

	/* FIXME: error message */
	if (list == _library) {
		return false;
	}

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Clearing playlist\n");

	if (list == _playlist) {
		if ((r = mpd_run_clear(conn->h()))) {
			st->last_playlist = -1;
		}
		return r;
	}

	return mpd_run_playlist_clear(conn->h(), list->filename.c_str());
}

/*
 * Seeks in the stream
 */
bool
Control::seek(int offset)
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Seeking by %d seconds\n", offset);

	/* FIXME: perhaps this check should be performed at an earlier stage? */
	if (!song()) {
		return false;
	}

	offset = st->time_elapsed + offset;

	return mpd_run_seek_id(conn->h(), song()->id, offset);
}

/*
 * Toggles crossfading
 * FIXME: return value changed from crossfadetime to boolean
 */
int
Control::crossfade()
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Toggling crossfade\n");

	if (st->crossfade == 0) {
		return mpd_run_crossfade(conn->h(), crossfadetime);
	}

	crossfadetime = st->crossfade;
	return mpd_run_crossfade(conn->h(), 0);
}

/*
 * Set crossfade time in seconds
 * FIXME: return value changed from crossfadetime to boolean
 */
int
Control::crossfade(int interval)
{
	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Set crossfade to %d seconds\n", interval);

	if (interval < 0) {
		return false;
	}

	crossfadetime = interval;
	return mpd_run_crossfade(conn->h(), crossfadetime);
}

/*
 * Move selected songs
 */
unsigned int
Control::move(Songlist * list, int offset)
{
	Song *		song;
	int		oldpos;
	int		newpos;
	const char *	filename;
	unsigned int	moved = 0;

	/* Library is read only */
	/* FIXME: error message */
	if (list == _library || !list)
		return 0;

	filename = list->filename.c_str();
	
	if (offset < 0) {
		song = list->getnextselected();
	} else {
		song = list->getprevselected();
	}

	EXIT_IDLE;

	//list_start();

	while (song != NULL)
	{
		if (song->pos == MPD_SONG_NO_NUM)
		{
			oldpos = list->match(song->file, 0, list->end(), MATCH_FILE | MATCH_EXACT);
			if (oldpos == MATCH_FAILED)
				break;
		}
		else
		{
			oldpos = song->pos;
		}

		newpos = oldpos + offset;

		if (!list->move(oldpos, newpos)) {
			break;
		}

		++moved;

		/* FIXME: return values? */
		if (list != _playlist) {
			mpd_send_playlist_move(conn->h(), filename, oldpos, newpos);
		} else {
			mpd_run_move(conn->h(), song->pos, oldpos);
		}

		if (offset < 0)
			song = list->getnextselected();
		else
			song = list->getprevselected();

	}

	list->resetgets();

	if (!list_end() || moved == 0)
	{
		return 0;
	}

	if (list == _playlist)
	{
		st->last_playlist += moved;
	}

	return moved;
}


/*
 * Removes all songs from list1 not found in list2
 */
int		Control::prune(Songlist * list1, Songlist * list2)
{
	unsigned int		i;
	int			pruned = 0;

	if (!list1 || !list2) return pruned;

	for (i = 0; i < list1->size(); i++)
	{
		if (list2->match(list1->song(i)->file, 0, list2->size() - 1, MATCH_FILE) == MATCH_FAILED)
		{
			pms->log(MSG_DEBUG, 0, "Pruning '%s' from list.\n", list1->song(i)->file.c_str());
			list1->remove(i);
			++pruned;
		}
	}

	return pruned;
}


/*
 * Starts mpd command list/queue mode
 * FIXME: not implemented
 */
bool
Control::list_start()
{
	assert(0);

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Entering command list mode.\n");

	if (mpd_command_list_begin(conn->h(), true)) {
		command_mode = 1;
	}

	return get_error_bool();
}

/*
 * Ends mpd command list/queue mode
 * FIXME: not implemented
 */
bool
Control::list_end()
{
	assert(0);

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Leaving command list mode.\n");

	if (mpd_command_list_end(conn->h())) {
		command_mode = 0;
	}

	return get_error_bool();
}

/*
 * Retrieves status about the state of MPD.
 */
bool
Control::get_status()
{
	mpd_status *	status;
	mpd_stats *	stats;

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Retrieving MPD status from server.\n");

	if ((status = mpd_run_status(conn->h())) == NULL) {
		/* FIXME: error handling? */
		pms->log(MSG_DEBUG, 0, "mpd_run_status returned NULL pointer.\n");
		delete _song;
		_song = NULL;
		st->song = MPD_SONG_NO_NUM;
		st->songid = MPD_SONG_NO_ID;
		last_song = MPD_SONG_NO_ID;
		return false;
	}

	st->assign_status(status);
	mpd_status_free(status);

	if ((stats = mpd_run_stats(conn->h())) == NULL) {
		/* FIXME ? */
		pms->log(MSG_DEBUG, 0, "mpd_run_stats returned NULL pointer.\n");
		return false;
	}

	st->assign_stats(stats);
	mpd_stats_free(stats);

	/* Override local settings if MPD mode changed */
	if (st->random) {
		pms->options->set_long("playmode", PLAYMODE_RANDOM);
	}

	if (st->repeat) {
		if (st->single) {
			pms->options->set_long("repeat", REPEAT_ONE);
		} else {
			pms->options->set_long("repeat", REPEAT_LIST);
		}
	}

	return true;
}

Directory::Directory(Directory * par, string n)
{
	parent_ = par;
	name_ = n;
	cursor = 0;
}

/*
 * Return full path from top-level to here
 */
string				Directory::path()
{
	if (parent_ == NULL)
		return "";
	else if (parent_->name().size() == 0)
		return name_;
	else
		return (parent_->path() + '/' + name_);
}

/*
 * Adds a directory entry to the tree
 */
Directory *			Directory::add(string s)
{
	size_t				i;
	string				t;
	vector<Directory *>::iterator	it;
	Directory *			d;

	if (s.size() == 0)
		return NULL;

	i = s.find_first_of('/');

	/* Within this directory */
	if (i == string::npos)
	{
		d = new Directory(this, s);
		children.push_back(d);
		return d;
	}

	t = s.substr(0, i);		// top-level
	s = s.substr(i + 1);		// all sub-level

	/* Search for top-level string in subdirectories */
	it = children.begin();
	while (it != children.end())
	{
		if ((*it)->name() == t)
		{
			return (*it)->add(s);
		}
		++it;
	}

	/* Not found, this should _not_ happen */
	pms->log(MSG_DEBUG, 0, "BUG: directory not found in hierarchy: '%s', '%s'\n", t.c_str(), s.c_str());

	return NULL;
}

/*
void		Directory::debug_tree()
{
	vector<Directory *>::iterator	it;
	vector<Song *>::iterator	is;

	pms->log(MSG_DEBUG, 0, "Printing contents of %s\n", path().c_str());

	is = songs.begin();
	while (is != songs.end())
	{
		pms->log(MSG_DEBUG, 0, "> %s\n", (*is)->file.c_str());
		++is;
	}

	it = children.begin();
	while (it != children.end())
	{
		(*it)->debug_tree();
		++it;
	}
}
*/

/*
 * Retrieves the entire song library from MPD
 */
bool
Control::update_library()
{
	uint32_t			total = 0;
	Song *				song;
	struct mpd_entity *		ent;
	const struct mpd_directory *	ent_directory;
	const struct mpd_song *		ent_song;
	const struct mpd_playlist *	ent_playlist;
	Directory *			dir = rootdir;

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Updating library from DB time %d to %d\n", st->last_db_update_time, st->db_update_time);
	st->last_db_update_time = st->db_update_time;

	if (!mpd_send_list_all_meta(conn->h(), "")) {
		return false;
	}

	_library->clear();

	while ((ent = mpd_recv_entity(conn->h())) != NULL)
	{
		switch(mpd_entity_get_type(ent))
		{
			case MPD_ENTITY_TYPE_SONG:
				ent_song = mpd_entity_get_song(ent);
				song = new Song(ent_song);
				song->id = MPD_SONG_NO_ID;
				song->pos = MPD_SONG_NO_NUM;
				_library->add(song);
				dir->songs.push_back(song);
				break;
			case MPD_ENTITY_TYPE_PLAYLIST:
				/* Issue #8: https://github.com/ambientsound/pms/issues/8 */
				ent_playlist = mpd_entity_get_playlist(ent);
				pms->log(MSG_DEBUG, 0, "NOT IMPLEMENTED in update_library(): got playlist entity in update_library(): %s\n", mpd_playlist_get_path(ent_playlist));
				break;
			case MPD_ENTITY_TYPE_DIRECTORY:
				ent_directory = mpd_entity_get_directory(ent);
				dir = rootdir->add(mpd_directory_get_path(ent_directory));
				assert(dir != NULL);
				/*
				if (dir == NULL)
				{
					dir = rootdir;
				}
				*/
				break;
			case MPD_ENTITY_TYPE_UNKNOWN:
				pms->log(MSG_DEBUG, 0, "BUG in update_library(): entity type not implemented by libmpdclient\n");
				break;
			default:
				pms->log(MSG_DEBUG, 0, "BUG in update_library(): entity type not implemented by PMS\n");
				break;
		}

		mpd_entity_free(ent);

		++total;
	}

	pms->log(MSG_DEBUG, 0, "Processed a total of %d entities during library update\n", total);

	return get_error_bool();
}

/*
 * Synchronizes playlists with MPD server, overwriting local versions
 */
unsigned int
Control::update_playlists()
{
	struct mpd_playlist *		playlist;
	Songlist *			list;
	vector<Songlist *>		newlist;
	vector<Songlist *>::iterator	i;

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Refreshing playlists.\n");

	if (!mpd_send_list_playlists(conn->h())) {
		/* FIXME */
		return -1;
	}

	/* FIXME: store in a temporary list instead */
	while ((playlist = mpd_recv_playlist(conn->h())) != NULL)
	{
		pms->log(MSG_DEBUG, 0, "Got playlist entity: %s\n", mpd_playlist_get_path(playlist));
		list = findplaylist(mpd_playlist_get_path(playlist));
		if (!list) {
			list = new Songlist();
			list->filename = mpd_playlist_get_path(playlist);
			newlist.push_back(list);
		}
		mpd_playlist_free(playlist);
	}

	/* FIXME: check for errors */

	retrieve_lists(newlist);
	{
		i = newlist.begin();
		while (i != newlist.end())
		{
			playlists.push_back(*i);
			++i;
		}

		pms->log(MSG_DEBUG, 0, "Server returned %d new playlists, sums to total of of %d custom playlists.\n", newlist.size(), playlists.size());
	}

	return playlists.size();
}

/*
 * Get all contents from server playlists playlists
 */
bool
Control::retrieve_lists(vector<Songlist *> &lists)
{
	vector<Songlist *>::iterator	i;
	Song *				song;
	mpd_entity *			ent;
	const mpd_song *		ent_song;

	EXIT_IDLE;

	i = lists.begin();

	while (i != lists.end())
	{
		if (!mpd_send_list_playlist_meta(conn->h(), (*i)->filename.c_str())) {
			return false;
		}

		(*i)->clear();

		while ((ent = mpd_recv_entity(conn->h())) != NULL)
		{
			switch(mpd_entity_get_type(ent))
			{
				case MPD_ENTITY_TYPE_SONG:
					ent_song = mpd_entity_get_song(ent);
					song = new Song(ent_song);
					(*i)->add(song);
					break;
				case MPD_ENTITY_TYPE_UNKNOWN:
					pms->log(MSG_DEBUG, 0, "BUG in retrieve_lists(): entity type not implemented by libmpdclient\n");
					break;
				default:
					pms->log(MSG_DEBUG, 0, "BUG in retrieve_lists(): entity type not implemented by PMS\n");
					break;
			}
			mpd_entity_free(ent);
		}

		if (!get_error_bool()) {
			return false;
		}

		++i;
	}

	return get_error_bool();
}

/*
 * Returns a playlist with the specified filename
 */
Songlist *	Control::findplaylist(string fn)
{
	vector<Songlist *>::iterator	i;

	i = playlists.begin();
	while (i != playlists.end())
	{
		if ((*i)->filename == fn)
		{
			return *i;
		}
		++i;
	}

	return NULL;
}

/*
 * Creates or locates a new playlist
 * FIXME: dubious return value
 */
Songlist *
Control::newplaylist(string fn)
{
	Songlist * list;

	list = findplaylist(fn);
	if (list != NULL) {
		return list;
	}

	list = new Songlist();
	assert(list != NULL);

	if (mpd_run_save(conn->h(), fn.c_str())) {
		list = new Songlist();
		assert(list != NULL);
		pms->log(MSG_DEBUG, 0, "newplaylist(): created playlist '%s'\n", fn.c_str());
		list->filename = fn;
		playlists.push_back(list);
	}

	return list;
}

/*
 * Deletes a playlist
 */
bool
Control::deleteplaylist(string fn)
{
	vector<Songlist *>::iterator	i;
	Songlist *			lst;

	/* FIXME: implement PlaylistList for this functionality */
	i = playlists.begin();
	do {
		if ((*i)->filename != fn) {
			continue;
		}

		if (mpd_run_rm(conn->h(), (*i)->filename.c_str())) {
			lst = *i;
			delete *i;
			i = playlists.erase(i);

			if (lst != _active) {
				return true;
			}

			/* Change active list */
			if (i == playlists.end())
			{
				if (playlists.size() == 0)
					_active = *i;
				else
					--i;
			}

			_active = *i;
			return true;
		}

		break;

	} while (++i != playlists.end());

	return false;
}

/*
 * Returns the active playlist
 */
Songlist *
Control::activelist()
{
	return _active;
}

/*
 * Sets the active playlist
 */
bool		Control::activatelist(Songlist * list)
{
	vector<Songlist *>::iterator	i;
	bool				changed = false;

	if (list == _playlist || list == _library)
	{
		_active = list;
		changed = true;
	}
	else
	{
		i = playlists.begin();
		while (i != playlists.end())
		{
			if (*i == list)
			{
				_active = list;
				changed = true;
				break;
			}
			++i;
		}
	}

	/* Have MPD manage random inside playlist */
	/* FIXME: custom function */
	/* FIXME: not our responsibility! */
	/* FIXME: wrong return codes */
	bool set_repeat;
	bool set_single;
	bool set_random;

	if (changed)
	{
		set_repeat = ((pms->options->get_long("repeat") == REPEAT_LIST || pms->options->get_long("repeat") == REPEAT_ONE) && activelist() == playlist());
		set_single = (pms->options->get_long("repeat") == REPEAT_ONE && activelist() == playlist());
		set_random = (pms->options->get_long("playmode") == PLAYMODE_RANDOM && activelist() == playlist());

		if (set_repeat != st->repeat && !repeat(set_repeat)) {
			return false;
		}

		if (set_single != st->single && !single(set_single)) {
			return false;
		}

		if (set_random != st->random && !random(set_random)) {
			return false;
		}
	}

	return changed;
}

/*
 * Retrieves current playlist from MPD
 * TODO: implement more entity types
 */
bool
Control::update_playlist()
{
	bool			rc;
	Song *			song;
	struct mpd_entity *	ent;
	const struct mpd_song *	ent_song;

	EXIT_IDLE;

	pms->log(MSG_DEBUG, 0, "Updating playlist from version %d to %d\n", st->last_playlist, st->playlist);

	if (st->last_playlist == -1) {
		_playlist->clear();
	}

	if (!mpd_send_queue_changes_meta(conn->h(), st->last_playlist)) {
		return false;
	}

	while ((ent = mpd_recv_entity(conn->h())) != NULL)
	{
		switch(mpd_entity_get_type(ent))
		{
			case MPD_ENTITY_TYPE_SONG:
				ent_song = mpd_entity_get_song(ent);
				song = new Song(ent_song);
				_playlist->add(song);
				break;
			case MPD_ENTITY_TYPE_UNKNOWN:
				pms->log(MSG_DEBUG, 0, "BUG in update_playlist(): entity type not implemented by libmpdclient\n");
				break;
			default:
				pms->log(MSG_DEBUG, 0, "BUG in update_playlist(): entity type not implemented by PMS\n");
				break;
		}
		mpd_entity_free(ent);
	}

	if ((rc = get_error_bool()) == true) {
		_playlist->truncate(st->playlist_length);
		st->last_playlist = st->playlist;
	}

	return rc;
}

/*
 * Tells whether the currently playing song has changed since last call
 */
bool		Control::song_changed()
{
	if (last_song == oldsong)
		return false;

	oldsong = last_song;
	return true;
}

/*
 * Stores the currently playing song in _song
 * FIXME: dubious return value
 */
int
Control::get_current_playing()
{
	Song *			current_song;
	struct mpd_song *	song;

	EXIT_IDLE;

	if ((song = mpd_run_current_song(conn->h())) == NULL) {
		return MPD_SONG_NO_ID;
	}

	/* FIXME: wtf is this?
	ent = mpd_getNextInfoEntity(conn->h());
	if (ent == NULL || ent->type != MPD_INFO_ENTITY_TYPE_SONG) {
		last_song = MPD_SONG_NO_NUM;
		_song = NULL;
		return MPD_SONG_NO_ID;
	}
	*/

	if (_song != NULL) {
		delete _song;
	}

	_song = new Song(song);

	/* FIXME: sketchy */
	/* better implement set_current_song or something */
	if (_song->id != last_song) {
		oldsong = last_song;
		last_song = _song->id;
	}

	mpd_song_free(song);

	return _song->id;
}

/*
 * Rescans entire library
 * FIXME: runs "update", there is also a "rescan" that can be implemented
 * FIXME: dubious return value
 */
bool
Control::rescandb(string dest)
{
	/* we can handle an MPD error if this is not supported */
	/*
	if (st->db_updating) {
		// FIXME: error message
		return false;
	}
	*/

	unsigned int job_id;

	EXIT_IDLE;

	job_id = mpd_run_update(conn->h(), dest.c_str());
	if (job_id == 0) {
		/* FIXME: handle errors */
		return false;
	}

	// FIXME?
	st->update_job_id = job_id;
	return job_id;
	//st->update_job_id = mpd_getUpdateId(conn->h());

	//return finish();
}

/*
 * Sends a password to the mpd server
 * FIXME: should retrieve updated privileges list?
 */
bool
Control::sendpassword(string pw)
{
	EXIT_IDLE;

	return mpd_run_password(conn->h(), pw.c_str());
}

/*
 * Notifies command system that an update from server is unneccessary as PMS already has done it.
 * FIXME: this command is probably dangerous and a cause of bugs due to PMS drifting out of synch.
 * FIXME: remove this function and all dependencies on it!
 */
bool
Control::increment()
{
	if (st->last_playlist == -1)
	{
		return false;
	}
	++(st->last_playlist);
	return true;
}

/**
 * Set client in IDLE mode
 */
bool
Control::idle()
{
	if (is_idle()) {
		return true;
	}

	pms->log(MSG_DEBUG, 0, "Entering IDLE mode.\n");
	set_is_idle(mpd_send_idle(conn->h()));

	return is_idle();
}

/**
 * Take client out of IDLE mode
 */
bool
Control::noidle()
{
	if (!is_idle()) {
		return true;
	}

	pms->log(MSG_DEBUG, 0, "Leaving IDLE mode.\n");
	set_is_idle(mpd_send_noidle(conn->h()));

	return is_idle();
}

/**
 * Block until MPD is ready to receieve requests.
 */
bool
Control::wait_until_noidle()
{
	pms->zeromq->poll_events(NOIDLE_POLL_TIMEOUT);
	return pms->run_has_idle_events();
}

bool
Control::is_idle()
{
	return _is_idle;
}

bool
Control::set_is_idle(bool i)
{
	_is_idle = i;
}

/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2010  Kim Tore Jensen
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

#include "command.h"
#include "pms.h"

extern Pms *			pms;


/*
 * Status class
 */
Mpd_status::Mpd_status()
{
	status			= NULL;
	stats			= NULL;

	muted			= false;
	volume			= 0;
	repeat			= false;
	single			= false;
	random			= false;
	playlist_length		= 0;
	playlist		= -1;
	storedplaylist		= -1;
	state			= MPD_STATUS_STATE_UNKNOWN;
	crossfade		= 0;
	song			= MPD_SONG_NO_NUM;
	songid			= MPD_SONG_NO_ID;
	time_elapsed		= 0;
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
	last_state		= state;
	last_db_update_time	= db_update_time;
	last_db_updating	= db_updating;
	update_job_id		= -1;
}

Mpd_status::~Mpd_status()
{
	if (status != NULL)
		mpd_freeStatus(status);

	if (stats != NULL)
		mpd_freeStats(stats);
}

void		Mpd_status::assign_status(mpd_Status * st)
{
	if (status != NULL)
		mpd_freeStatus(status);

	status = st;
	if (status == NULL)
		return;

	volume			= status->volume;
	repeat			= (status->repeat == 1 ? true : false);
	single			= (status->single == 1 ? true : false);
	random			= (status->random == 1 ? true : false);
	playlist_length		= status->playlistLength;
	playlist		= status->playlist;
	storedplaylist		= status->storedplaylist;
	state			= status->state;
	crossfade		= status->crossfade;
	song			= status->song;
	songid			= status->songid;
	time_elapsed		= status->elapsedTime;
	time_total		= status->totalTime;
	db_updating		= status->updatingDb;
	errstr			= (status->error ? status->error : "");

	/* Audio decoded properties */
	bitrate			= status->bitRate;
	samplerate		= status->sampleRate;
	bits			= status->bits;
	channels		= status->channels;
}

void		Mpd_status::assign_stats(mpd_Stats * st)
{
	if (stats != NULL)
		mpd_freeStats(stats);

	stats = st;
	if (stats == NULL)
		return;

	artists_count		= stats->numberOfArtists;
	albums_count		= stats->numberOfAlbums;
	songs_count		= stats->numberOfSongs;

	uptime			= stats->uptime;
	db_update_time		= stats->dbUpdateTime;
	playtime		= stats->playTime;
	db_playtime		= stats->dbPlayTime;
}

bool		Mpd_status::alive() const
{
	return (status != NULL);
}



/*
 * Command class manages commands sent to and from mpd
 */
Control::Control(Connection * n_conn)
{
	conn = n_conn;
	st = new Mpd_status();
	rootdir = new Directory(NULL, "");
	_stats = NULL;
	_song = NULL;
	st->last_playlist = -1;
	last_song = MPD_SONG_NO_NUM;
	oldsong = MPD_SONG_NO_NUM;
	_has_new_playlist = false;
	_has_new_library = false;
	_playlist = new Songlist;
	_library = new Songlist;
	_playlist->role = LIST_ROLE_MAIN;
	_library->role = LIST_ROLE_LIBRARY;
	_active = NULL;
	command_mode = 0;
	mutevolume = 0;
	crossfadetime = pms->options->get_long("crossfade");

	usetime = 0;
	time(&(mytime[0]));
	mytime[1] = 0; // Update immedately
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
	mpd_finishCommand(conn->h());
	st->error = conn->h()->error;
	st->errstr = conn->h()->errorStr;

	if (st->error != 0)
	{
		pms->log(MSG_CONSOLE, STERR, "MPD returned error %d: %s\n", st->error, st->errstr.c_str());

		/* Connection closed */
		if (st->error == MPD_ERROR_CONNCLOSED)
		{
			conn->disconnect();
		}

		clearerror();

		return false;
	}

	return true;
}

/*
 * Clears any error
 */
void		Control::clearerror()
{
	if (conn->h())
		mpd_clearError(conn->h());
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
void		Control::get_available_commands()
{
	char *		c;
	string		s;

	memset(&commands, 0, sizeof(commands));

	mpd_sendCommandsCommand(conn->h());
	while ((c = mpd_getNextCommand(conn->h())) != NULL)
	{
		s = c;
		free(c);

		if (s == "add")
			commands.add = true;
		else if (s == "addid")
			commands.addid = true;
		else if (s == "clear")
			commands.clear = true;
		else if (s == "clearerror")
			commands.clearerror = true;
		else if (s == "close")
			commands.close = true;
		else if (s == "commands")
			commands.commands = true;
		else if (s == "count")
			commands.count = true;
		else if (s == "crossfade")
			commands.crossfade = true;
		else if (s == "currentsong")
			commands.currentsong = true;
		else if (s == "delete")
			commands.delete_ = true;
		else if (s == "deleteid")
			commands.deleteid = true;
		else if (s == "disableoutput")
			commands.disableoutput = true;
		else if (s == "enableoutput")
			commands.enableoutput = true;
		else if (s == "find")
			commands.find = true;
		else if (s == "idle")
			commands.idle = true;
		else if (s == "kill")
			commands.kill = true;
		else if (s == "list")
			commands.list = true;
		else if (s == "listall")
			commands.listall = true;
		else if (s == "listallinfo")
			commands.listallinfo = true;
		else if (s == "listplaylist")
			commands.listplaylist = true;
		else if (s == "listplaylistinfo")
			commands.listplaylistinfo = true;
		else if (s == "listplaylists")
			commands.listplaylists = true;
		else if (s == "load")
			commands.load = true;
		else if (s == "lsinfo")
			commands.lsinfo = true;
		else if (s == "move")
			commands.move = true;
		else if (s == "moveid")
			commands.moveid = true;
		else if (s == "next")
			commands.next = true;
		else if (s == "notcommands")
			commands.notcommands = true;
		else if (s == "outputs")
			commands.outputs = true;
		else if (s == "password")
			commands.password = true;
		else if (s == "pause")
			commands.pause = true;
		else if (s == "ping")
			commands.ping = true;
		else if (s == "play")
			commands.play = true;
		else if (s == "playid")
			commands.playid = true;
		else if (s == "playlist")
			commands.playlist = true;
		else if (s == "playlistadd")
			commands.playlistadd = true;
		else if (s == "playlistclear")
			commands.playlistclear = true;
		else if (s == "playlistdelete")
			commands.playlistdelete = true;
		else if (s == "playlistfind")
			commands.playlistfind = true;
		else if (s == "playlistid")
			commands.playlistid = true;
		else if (s == "playlistinfo")
			commands.playlistinfo = true;
		else if (s == "playlistmove")
			commands.playlistmove = true;
		else if (s == "playlistsearch")
			commands.playlistsearch = true;
		else if (s == "plchanges")
			commands.plchanges = true;
		else if (s == "plchangesposid")
			commands.plchangesposid = true;
		else if (s == "previous")
			commands.previous = true;
		else if (s == "random")
			commands.random = true;
		else if (s == "rename")
			commands.rename = true;
		else if (s == "repeat")
			commands.repeat = true;
		else if (s == "single")
			commands.single = true;
		else if (s == "rm")
			commands.rm = true;
		else if (s == "save")
			commands.save = true;
		else if (s == "filter")
			commands.filter = true;
		else if (s == "seek")
			commands.seek = true;
		else if (s == "seekid")
			commands.seekid = true;
		else if (s == "setvol")
			commands.setvol = true;
		else if (s == "shuffle")
			commands.shuffle = true;
		else if (s == "stats")
			commands.stats = true;
		else if (s == "status")
			commands.status = true;
		else if (s == "stop")
			commands.stop = true;
		else if (s == "swap")
			commands.swap = true;
		else if (s == "swapid")
			commands.swapid = true;
		else if (s == "tagtypes")
			commands.tagtypes = true;
		else if (s == "update")
			commands.update = true;
		else if (s == "urlhandlers")
			commands.urlhandlers = true;
		else if (s == "volume")
			commands.volume = true;
	}
	finish();
}

/*
 * Play, pause, toggle, stop, next, prev
 */
bool		Control::play()
{
	bool		r;
	if (!alive())	return false;
	mpd_sendPlayCommand(conn->h());
	r = finish();
	get_status();
	return r;
}

bool		Control::playid(song_t songid)
{
	bool		r;
	if (!alive())	return false;
	mpd_sendPlayIdCommand(conn->h(), songid);
	r = finish();
	get_status();
	return r;
}

bool		Control::playpos(song_t songpos)
{
	bool		r;
	if (!alive())	return false;
	mpd_sendPlayPosCommand(conn->h(), songpos);
	r = finish();
	get_status();
	return r;
}

bool		Control::pause(bool tryplay)
{
	bool		r;
	if (!alive())	return false;

	switch(st->state)
	{
		case MPD_STATUS_STATE_PLAY:
			mpd_sendPauseCommand(conn->h(), 1);
			break;
		case MPD_STATUS_STATE_PAUSE:
			mpd_sendPauseCommand(conn->h(), 0);
			break;
		case MPD_STATUS_STATE_STOP:
		case MPD_STATUS_STATE_UNKNOWN:
		default:
			return (tryplay && play());
	}

	r = finish();
	get_status();
	return r;
}

bool		Control::stop()
{
	bool		r;
	if (!alive())	return false;
	mpd_sendStopCommand(conn->h());
	r = finish();
	get_status();
	return r;
}

/*
 * Shuffles the playlist.
 */
bool		Control::shuffle()
{
	bool		r;
	if (!alive())	return false;
	mpd_sendShuffleCommand(conn->h());
	r = finish();
	get_status();
	return r;
}

/*
 * Sets repeat mode
 */
bool		Control::repeat(bool on)
{
	bool		r;
	if (!alive())	return false;
	mpd_sendRepeatCommand(conn->h(), on);
	r = finish();
	get_status();
	return r;
}

/*
 * Sets single mode
 */
bool		Control::single(bool on)
{
	bool		r;
	if (!alive())	return false;
	mpd_sendSingleCommand(conn->h(), on);
	r = finish();
	get_status();
	return r;
}

/*
 * Set an absolute volume
 */
bool		Control::setvolume(int vol)
{
	if (!alive())	return false;

	if (vol < 0)
		vol = 0;
	else if (vol > 100)
		vol = 100;

	mpd_sendSetvolCommand(conn->h(), vol);
	if (finish())
	{
		st->volume = vol;
		return true;
	}
	return false;
}

/*
 * Changes volume
 */
bool		Control::volume(int offset)
{
	bool		r;
	if (!alive())	return false;

	if (st->volume == MPD_STATUS_NO_VOLUME)
		return false;

	if (st->volume + offset > 100)
		offset = 100 - st->volume;
	else if (st->volume + offset < 0)
		offset = -st->volume;

	mpd_sendSetvolCommand(conn->h(), st->volume + offset);
	r = finish();
	if (r)
	{
		mutevolume = 0;
		st->volume += offset;
		if (st->volume < 0)
			st->volume = 0;
		else if (st->volume > 100)
			st->volume = 100;
	}

	return r;
}

/*
 * Mute/unmute volume
 */
bool		Control::mute()
{
	bool		success;
	if (!alive())	return false;

	if (st->volume == MPD_STATUS_NO_VOLUME)
		return false;

	if (muted())
	{
		mpd_sendSetvolCommand(conn->h(), mutevolume);
		success = finish();
		if (success)
			st->volume = mutevolume;
		mutevolume = 0;
	}
	else
	{
		mutevolume = st->volume;
		mpd_sendSetvolCommand(conn->h(), 0);
		success = finish();
		if (success)
			st->volume = 0;
	}

	return success;
}

/*
 * Is muted?
 */
bool		Control::muted()
{
	return (mutevolume > 0 && st->volume == 0);
}

/*
 * Toggles MPDs built-in random mode
 */
bool		Control::random(int set)
{
	bool		r;
	if (!alive())	return false;

	if (set == -1)
		set = (st->random == false ? 1 : 0);
	else
		if (set > 1) set = 1;

	mpd_sendRandomCommand(conn->h(), set);
	r = finish();
	get_status();
	return r;
}

/*
 * Appends a playlist to another playlist
 */
song_t		Control::add(Songlist * source, Songlist * dest)
{
	song_t			first = MPD_SONG_NO_ID;
	song_t			result;
	unsigned int		i;

	if (source == NULL || dest == NULL)
		return MPD_SONG_NO_ID;

	list_start();
	for (i = 0; i < source->size(); i++)
	{
		result = add(dest, source->song(i));
		if (first == MPD_SONG_NO_ID && result != MPD_SONG_NO_ID)
			first = result;
	}
	if (!list_end())
		return MPD_SONG_NO_ID;

	return first;
}

/*
 * Add a song to a playlist
 */
song_t		Control::add(Songlist * list, Song * song)
{
	song_t		i = MPD_SONG_NO_ID;
	Song *		nsong;

	if (!alive() || list == NULL || song == NULL)
		return i;

	if (list == _playlist)
	{
		i = mpd_sendAddIdCommand(conn->h(), song->file.c_str());
	}
	else if (list != _library)
	{
		if (list->filename.size() == 0)
			return i;
		mpd_sendPlaylistAddCommand(conn->h(), (char *)list->filename.c_str(), (char *)song->file.c_str());
	}
	else
	{
		return i;
	}

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
}

/*
 * Remove a song from the playlist
 */
int		Control::remove(Songlist * list, Song * song)
{
	int		pos = MATCH_FAILED;

	if (!alive() || song == NULL || list == NULL)
		return false;
	if (list == _library)
		return false;
	if (list == _playlist && song->id == MPD_SONG_NO_ID)
		return false;
	if (list != _playlist)
	{
		if (list->filename.size() == 0)
			return false;
		pos = list->locatesong(song);
		if (pos == MATCH_FAILED)
			return false;
		pms->log(MSG_DEBUG, 0, "Removing song %d from list.\n", pos);
	}

	if (list == _playlist)
		mpd_sendDeleteIdCommand(conn->h(), song->id);
	else
		mpd_sendPlaylistDeleteCommand(conn->h(), (char *)list->filename.c_str(), pos);

	if (command_mode != 0) return true;
	if (finish())
	{
		list->remove(pos == MATCH_FAILED ? song->pos : pos);
		if (list == _playlist)
			increment();
		return true;
	}

	return false;
}

/*
 * Crops the playlist
 */
bool		Control::crop(Songlist * list, int mode)
{
	unsigned int		i;
	int			pos;
	unsigned int		upos;
	Song *			song;

	if (!alive())		return false;
	if (!list)		return false;
	if (list == _library)
	{
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

		list_start();
		for (i = list->end(); i < list->size(); i--)
		{
			if (upos != i)
			{
				if (list == _playlist)
					mpd_sendDeleteIdCommand(conn->h(), list->song(i)->id);
				else
					mpd_sendPlaylistDeleteCommand(conn->h(), (char *)list->filename.c_str(), static_cast<int>(i));
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

		list_start();
		for (i = list->end(); i < list->size(); i--)
		{
			if (list->song(i)->selected == false)
			{
				if (list == _playlist)
					mpd_sendDeleteIdCommand(conn->h(), list->song(i)->id);
				else
					mpd_sendPlaylistDeleteCommand(conn->h(), (char *)list->filename.c_str(), static_cast<int>(i));
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
int		Control::clear(Songlist * list)
{
	if (!alive())		return false;
	if (!list)		return false;
	if (list == _library)	return false;

	if (list == _playlist)
	{
		mpd_sendClearCommand(conn->h());
		if (finish())
		{
			st->last_playlist = -1;
			return true;
		}
		else	return false;
	}

	mpd_sendPlaylistClearCommand(conn->h(), (char *)(list->filename.c_str()));
	return finish();
}

/*
 * Seeks in the stream
 */
bool		Control::seek(int offset)
{
	if (!alive() || !song()) return false;

	offset = st->time_elapsed + offset;

	if (song()->id == MPD_SONG_NO_ID)
	{
		if (song()->pos == MPD_SONG_NO_NUM)
			return false;

		mpd_sendSeekCommand(conn->h(), song()->pos, offset);
	}
	else
	{
		mpd_sendSeekIdCommand(conn->h(), song()->id, offset);
	}

	return finish();
}

/*
 * Toggles or sets crossfading
 */
int		Control::crossfade()
{
	if (!alive()) return -1;

	if (st->crossfade == 0)
	{
		mpd_sendCrossfadeCommand(conn->h(), crossfadetime);
	}
	else
	{
		crossfadetime = st->crossfade;
		mpd_sendCrossfadeCommand(conn->h(), 0);
	}

	if (finish())
	{
		return (st->crossfade == 0 ? crossfadetime : 0);
	}

	return -1;
}

/*
 * Set crossfade time in seconds
 */
int		Control::crossfade(int interval)
{
	if (!alive()) return false;

	if (interval < 0)
		return false;

	crossfadetime = interval;
	mpd_sendCrossfadeCommand(conn->h(), crossfadetime);

	if (finish())
	{
		st->crossfade = crossfadetime;
		return st->crossfade;
	}
	return -1;
}

/*
 * Move selected songs
 */
unsigned int	Control::move(Songlist * list, int offset)
{
	Song *		song;
	int		oldpos;
	int		newpos;
	char *		filename;
	unsigned int	moved = 0;

	/* Library is read only */
	if (list == _library || !list)
		return 0;

	filename = const_cast<char *>(list->filename.c_str());
	
	if (offset < 0)
		song = list->getnextselected();
	else
		song = list->getprevselected();

	list_start();

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

		if (!list->move(oldpos, newpos))
			break;

		++moved;

		if (list != _playlist)
			mpd_sendPlaylistMoveCommand(conn->h(), filename, oldpos, newpos);
		else
			mpd_sendMoveCommand(conn->h(), song->pos, oldpos);

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
 */
bool		Control::list_start()
{
	if (!alive())	return false;

	mpd_sendCommandListBegin(conn->h());
	if (finish())
	{
		command_mode = 1;
		return true;
	}
	return false;
}

/*
 * Ends mpd command list/queue mode
 */
bool		Control::list_end()
{
	if (!alive())	return false;

	mpd_sendCommandListEnd(conn->h());
	if (finish())
	{
		command_mode = 0;
		return true;
	}
	return false;
}

/*
 * Retrieves status about the state of MPD.
 */
bool		Control::get_status()
{
	mpd_Status *	sta;
	mpd_Stats *	stat;

	if (!alive())	return false;

	mpd_sendStatusCommand(conn->h());
	sta = mpd_getStatus(conn->h());
	finish();
	st->assign_status(sta);

	if (!st->alive())
	{
		pms->log(MSG_DEBUG, 0, "get_status returned NULL pointer.\n");
		delete _song;
		_song = NULL;
		st->song = MPD_SONG_NO_NUM;
		st->songid = MPD_SONG_NO_ID;
		last_song = MPD_SONG_NO_ID;
		return false;
	}

	mpd_sendStatsCommand(conn->h());
	stat = mpd_getStats(conn->h());
	finish();
	st->assign_stats(stat);

	/* Override local settings if MPD mode changed */
	if (st->random)
		pms->options->set_long("playmode", PLAYMODE_RANDOM);
	if (st->repeat)
	{
		if (st->single)
			pms->options->set_long("repeat", REPEAT_ONE);
		else
			pms->options->set_long("repeat", REPEAT_LIST);
	}

	if (st->db_update_time != st->last_db_update_time)
	{
		pms->log(MSG_DEBUG, 0, "DB time was updated from %d to %d\n", st->db_update_time, st->last_db_update_time);
		pms->log(MSG_DEBUG, 0, "Server playlist version is now %d, local is %d\n", st->playlist, st->last_playlist);
		st->last_db_update_time = st->db_update_time;
		st->playlist = -1;
		st->update_job_id = -1;
		update_library();
	}

	return true;
}

/*
 * Query MPD server for updated information
 */
int		Control::update(bool force)
{
	/* Need >= 1 second to update. */
	time(&(mytime[usetime]));
	if (!force && difftime(mytime[0], mytime[1]) == 0)
	{
		return 1;
	}
	usetime = (usetime + 1) % 2;

	/* Get vital signs */
	if (!get_status())
	{
		return -1;
	}
	get_current_playing();

	/* New playlist? */
	if (st->playlist != st->last_playlist || st->last_playlist == -1)
	{
		pms->log(MSG_DEBUG, 0, "Playlist needs to be updated from version %d to %d\n", st->last_playlist, st->playlist);
		update_playlist();
		get_status();
		st->last_playlist = st->playlist;
	}

	return 0;
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
 * Retrieves the entire library from MPD
 */
void Control::update_library()
{
	Song *			song;
	mpd_InfoEntity *	ent;
	Directory *		dir = rootdir;

	if (!alive())		return;

	pms->log(MSG_DEBUG, 0, "Retrieving library from mpd...\n");
	_library->clear();

	mpd_sendListallInfoCommand(conn->h(), "/");
	while ((ent = mpd_getNextInfoEntity(conn->h())) != NULL)
	{
		switch(ent->type)
		{
			case MPD_INFO_ENTITY_TYPE_SONG:
				song = new Song(ent->info.song);
				_library->add(song);
				dir->songs.push_back(song);
				break;
			case MPD_INFO_ENTITY_TYPE_PLAYLISTFILE:
				/* Should not receive this here. */
				pms->log(MSG_DEBUG, 0, "BUG: Got playlist entity in update_library(): %s\n", ent->info.playlistFile->path);
				break;
			case MPD_INFO_ENTITY_TYPE_DIRECTORY:
				dir = rootdir->add(ent->info.directory->path);
				/* Should not be NULL, ever */
				if (dir == NULL)
				{
					dir = rootdir;
				}
				break;
			default:;
		}
		mpd_freeInfoEntity(ent);
	}
	finish();
	update_playlists();

	_has_new_library = true;
}

/*
 * Synchronizes playlists with MPD server, overwriting local versions
 */
unsigned int	Control::update_playlists()
{
	mpd_InfoEntity *		ent;
	Songlist *			list;
	vector<Songlist *>		newlist;
	vector<Songlist *>::iterator	i;

	if (!alive()) return 0;

	pms->log(MSG_DEBUG, 0, "Refreshing playlists.\n");
	mpd_sendLsInfoCommand(conn->h(), "/");
	while ((ent = mpd_getNextInfoEntity(conn->h())) != NULL)
	{
		if (ent->type == MPD_INFO_ENTITY_TYPE_PLAYLISTFILE)
		{
			pms->log(MSG_DEBUG, 0, "Got playlist entity: %s\n", ent->info.playlistFile->path);
			list = findplaylist(ent->info.playlistFile->path);
			if (!list)
			{
				list = new Songlist();
				list->filename = ent->info.playlistFile->path;
				newlist.push_back(list);
			}
		}
		mpd_freeInfoEntity(ent);
	}
	finish();

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
void		Control::retrieve_lists(vector<Songlist *> &lists)
{
	vector<Songlist *>::iterator	i;
	Song *				song;
	mpd_InfoEntity *		ent;

	i = lists.begin();

	while (i != lists.end())
	{
		(*i)->clear();
		mpd_sendListPlaylistInfoCommand(conn->h(), (char *)(*i)->filename.c_str());
		while ((ent = mpd_getNextInfoEntity(conn->h())) != NULL)
		{
			if (ent->type == MPD_INFO_ENTITY_TYPE_SONG)
			{
				song = new Song(ent->info.song);
				(*i)->add(song);
			}
			mpd_freeInfoEntity(ent);
		}
		++i;
	}	
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
 */
Songlist *	Control::newplaylist(string fn)
{
	Songlist *	list;

	list = findplaylist(fn);
	if (list != NULL)
		return list;

	list = new Songlist();
	if (!list)
		return NULL;

	mpd_sendSaveCommand(conn->h(), fn.c_str());
	if (!finish())
	{
		delete list;
		return NULL;
	}
	pms->log(MSG_DEBUG, 0, "newplaylist(): created playlist '%s'\n", fn.c_str());
	list->filename = fn;
	playlists.push_back(list);
	return list;
}

/*
 * Deletes a playlist
 */
bool		Control::deleteplaylist(string fn)
{
	vector<Songlist *>::iterator	i;
	Songlist *			lst;

	i = playlists.begin();
	while (i != playlists.end())
	{
		if ((*i)->filename == fn)
		{
			mpd_sendRmCommand(conn->h(), (*i)->filename.c_str());
			if (finish())
			{
				lst = *i;
				delete *i;
				i = playlists.erase(i);

				if (lst != _active)
					return true;

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
			else	return false;
		}
		++i;
	}

	return false;
}

/*
 * Returns the active playlist
 */
Songlist *	Control::activelist()
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
	if (changed)
	{
		repeat((pms->options->get_long("repeat") == REPEAT_LIST || pms->options->get_long("repeat") == REPEAT_ONE) && activelist() == playlist());
		single(pms->options->get_long("repeat") == REPEAT_ONE && activelist() == playlist());
		random(pms->options->get_long("playmode") == PLAYMODE_RANDOM && activelist() == playlist());
	}

	return changed;
}

/*
 * Retrieves current playlist from MPD
 */
void		Control::update_playlist()
{
	Song			*song;
	mpd_InfoEntity		*ent;

	if (!alive())		return;

	pms->log(MSG_DEBUG, 0, "Quering playlist changes.\n");

	if (st->last_playlist == -1)
	{
		_playlist->clear();
	}

	mpd_sendPlChangesCommand(conn->h(), st->last_playlist);
	while ((ent = mpd_getNextInfoEntity(conn->h())) != NULL)
	{
		song = new Song(ent->info.song);
		_playlist->add(song);
		mpd_freeInfoEntity(ent);
	}
	finish();

	_playlist->truncate(st->playlist_length);

	_has_new_playlist = true;
}

/*
 * Info for display class whether playlist has changed and needs a redraw
 */
bool Control::has_new_library()
{
	if (_has_new_library)
	{
		_has_new_library = false;
		return true;
	}
	return false;
}
bool Control::has_new_playlist()
{
	if (_has_new_playlist)
	{
		_has_new_playlist = false;
		return true;
	}
	return false;
}

/*
 * Tells whether the currently playing song has changed since last call
 */
bool		Control::song_changed()
{
	if (!alive())	return false;
	if (last_song == oldsong)
		return false;

	oldsong = last_song;
	return true;
}

/*
 * Tells whether the play state changed since last call
 */
bool		Control::state_changed()
{
	if (!alive() || st->last_state == st->state)
		return false;

	st->last_state = st->state;
	return true;
}


/*
 * Stores the currently playing song in _song
 */
int Control::get_current_playing()
{
	mpd_InfoEntity		*ent;

	if (!alive())
	{
		return MPD_SONG_NO_ID;
	}
	mpd_sendCurrentSongCommand(conn->h());

	ent = mpd_getNextInfoEntity(conn->h());
	if (ent == NULL || ent->type != MPD_INFO_ENTITY_TYPE_SONG)
	{
		_has_new_playlist = true;
		last_song = MPD_SONG_NO_NUM;
		_song = NULL;
		return MPD_SONG_NO_ID;
	}

	if (_song != NULL)
		delete _song;

	_song = new Song(ent->info.song);

	if (_song->id != last_song)
	{
		_has_new_playlist = true;
		oldsong = last_song;
		last_song = _song->id;
	}

	mpd_freeInfoEntity(ent);
	finish();

	return 0;
}

/*
 * Rescans entire library
 */
bool		Control::rescandb(string dest)
{
	if (!alive())		return false;
	if (st->db_updating)	return false;

	mpd_sendUpdateCommand(conn->h(), dest.c_str());
	st->update_job_id = mpd_getUpdateId(conn->h());

	return finish();
}

/*
 * Sends a password to the mpd server
 */
bool		Control::sendpassword(string pw)
{
	if (!alive())		return false;
	if (pw.size() == 0)	return false;

	mpd_sendPasswordCommand(conn->h(), pw.c_str());
	return finish();
}

/*
 * Notifies command system that an update from server is unneccessary as PMS already has done it.
 */
bool		Control::increment()
{
	if (st->last_playlist == -1)
	{
		return false;
	}
	++(st->last_playlist);
	return true;
}


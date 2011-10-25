/* vi:set ts=8 sts=8 sw=8:
 *
 * Practical Music Search
 * Copyright (c) 2006-2011  Kim Tore Jensen
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

#include "mpd.h"
#include "console.h"
#include "window.h"
#include "config.h"
#include "field.h"
#include <sys/select.h>
#include <sys/time.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <ctime>
#include <errno.h>
#include <unistd.h>
#include <netdb.h>
#include <cstring>
#include <string>
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <math.h>

using namespace std;

extern Config config;
extern Windowmanager wm;
extern Fieldtypes fieldtypes;
extern Curses curses;
extern MPD mpd;

void update_library_statusbar()
{
	unsigned int percent;
	percent = round(((float)mpd.library.size() / mpd.stats.songs) * 100.00);
	curses.wipe(&curses.statusbar, config.colors.statusbar);
	curses.print(&curses.statusbar, config.colors.statusbar, 0, 0, "Retrieving library: %d%%", percent);
	curses.flush();
}

MPD::MPD()
{
	errno = 0;
	error = "";
	host = "";
	port = "";
	sock = 0;
	autoadvance_playlist = -1;
	currentsong = NULL;
	connected = false;
	is_idle = false;
	playlist.playlist = true;
	playlist.title = "playlist";
	library.readonly = true;
	library.title = "library";
	active_songlist = &playlist;
	memset(&last_update, 0, sizeof last_update);
	memset(&last_clock, 0, sizeof last_clock);
	memset(&status, 0, sizeof status);
	memset(&stats, 0, sizeof stats);
	FD_ZERO(&fdset);
	FD_SET(STDIN_FILENO, &fdset);
	update_playstring();
}

MPD::~MPD()
{
	if (connected) close(sock);
}

bool MPD::set_idle(bool nidle)
{
	if (nidle == is_idle)
		return false;
	
	if (nidle)
	{
		mpd_raw_send("idle");
		is_idle = true;
		return true;
	}

	mpd_raw_send("noidle");
	is_idle = (mpd_getline(NULL) == MPD_GETLINE_OK);

	return true;
}

bool MPD::trigerr(int nerrno, const char * format, ...)
{
	va_list		ap;
	char		buf[1024];

	va_start(ap, format);
	vsprintf(buf, format, ap);
	va_end(ap);

	error = buf;
	errno = nerrno;

	sterr("%s", buf);
	wm.statusbar->draw();

	return false;
}

bool MPD::mpd_connect(string nhost, string nport)
{
	int			status;
	char			buf[32];
	struct addrinfo		hints;
	struct addrinfo *	res;

	host = nhost;
	port = nport;

	if (connected)
		mpd_disconnect();

	stinfo("Connecting to server '%s' port '%s'...", host.c_str(), port.c_str());
	wm.statusbar->draw();

	memset(&hints, 0, sizeof hints);
	hints.ai_family = AF_UNSPEC;
	hints.ai_socktype = SOCK_STREAM;

	if ((status = getaddrinfo(host.c_str(), port.c_str(), &hints, &res)) != 0)
	{
		trigerr(MPD_ERR_CONNECTION, "getaddrinfo error: %s", gai_strerror(status));
		return false;
	}

	sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol);
	if (sock == -1)
	{
		trigerr(MPD_ERR_CONNECTION, "could not create socket!");
		freeaddrinfo(res);
		return false;
	}

	if (connect(sock, res->ai_addr, res->ai_addrlen) == -1)
	{
		trigerr(MPD_ERR_CONNECTION, "could not connect to %s:%s", host.c_str(), port.c_str());
		close(sock);
		freeaddrinfo(res);
		return false;
	}

	freeaddrinfo(res);
	connected = true;

	FD_ZERO(&fdset);
	FD_SET(sock, &fdset);
	FD_SET(STDIN_FILENO, &fdset);

	stinfo("Connected to server '%s' on port '%s'.", host.c_str(), port.c_str());
	recv(sock, &buf, 32, 0);

	if (!set_protocol_version(buf))
	{
		trigerr(MPD_ERR_VERSION, "This MPD server is way too old!");
		debug("This version of PMS requires `idle' support, which this MPD server does not support.", NULL);
		debug("Either upgrade your MPD server to version 0.15.0 or newer, or use an older version of PMS (0.42).", NULL);
		mpd_disconnect();
		debug("Autoconnect has been disabled.", NULL);
		config.autoconnect = false;
	}
	is_idle = false;
	currentsong = NULL;

	return connected;
}

void MPD::mpd_disconnect()
{
	close(sock);
	sock = 0;
	connected = false;
	is_idle = false;
	currentsong = NULL;
	memset(&status, 0, sizeof status);
	memset(&stats, 0, sizeof stats);
	FD_ZERO(&fdset);
	FD_SET(STDIN_FILENO, &fdset);
	trigerr(MPD_ERR_CONNECTION, "Connection closed.");
}

bool MPD::is_connected()
{
	return connected;
}

bool MPD::set_password(string password)
{
	if (!connected)
		return false;

	if (password.size() == 0)
		return true;
	
	mpd_send("password \"%s\"", password.c_str());
	if (mpd_getline(NULL) == MPD_GETLINE_OK)
	{
		stinfo("Password '%s' accepted by server.", password.c_str());
		return true;
	}

	return true;
}

bool MPD::set_protocol_version(string data)
{
	unsigned int i = 7;
	unsigned int last = 7;
	int pos = 0;

	if (data.substr(0, 7) != "OK MPD ")
		return false;

	while (i <= data.size() && pos < 3)
	{
		if (data[i] == '.' || data[i] == '\n')
		{
			protocol_version[pos] = atoi(data.substr(last, i - last).c_str());
			++pos;
			last = i + 1;
		}
		++i;
	}
	debug("MPD server protocol version %d.%d.%d", protocol_version[0], protocol_version[1], protocol_version[2]);

	if (protocol_version[1] < 15)
		return false;

	return true;
}

int MPD::mpd_send(const char * data, ...)
{
	va_list		ap;
	char		buf[1024];

	va_start(ap, data);
	vsprintf(buf, data, ap);
	va_end(ap);

	if (!connected)
		return -1;

	set_idle(false);
	return mpd_raw_send(buf);
}

int MPD::mpd_raw_send(string data)
{
	unsigned int sent = 0;
	int s;

	if (!connected)
		return -1;

	data += '\n';
	while (sent < data.size())
	{
		if ((s = send(sock, data.substr(sent).c_str(), data.size() - sent, 0)) == -1)
		{
			mpd_disconnect();
			return s;
		}
		sent += s;
	}

	// Raw traffic dump
	//debug("-> %s", data.c_str());

	return sent;
}

int MPD::mpd_getline(string * nextline)
{
	int received = 0;
	size_t pos;
	string line = "";

	if (!connected)
		return MPD_GETLINE_ERR;

	while((pos = buffer.find('\n')) == string::npos)
	{
		received = recv(sock, getbuf, 1024, 0);
		if (received == 0)
		{
			mpd_disconnect();
			return MPD_GETLINE_ERR;
		}
		else if (received == -1)
		{
			return MPD_GETLINE_ERR;
		}

		getbuf[received] = '\0';
		buffer += getbuf;
	}

	if (pos == string::npos)
		return MPD_GETLINE_ERR;

	line = buffer.substr(0, pos);
	buffer = buffer.substr(pos + 1);

	if (line.size() == 0)
		return MPD_GETLINE_ERR;

	is_idle = false;

	// Raw traffic dump
	//debug("<- %s", line.c_str());

	if (line == "OK")
		return MPD_GETLINE_OK;

	if (line.size() >= 3 && line.substr(0, 3) == "ACK")
	{
		trigerr(MPD_ERR_ACK, line.c_str());
		return MPD_GETLINE_ACK;
	}

	if (nextline != NULL)
		*nextline = line;

	return MPD_GETLINE_MORE;
}

int MPD::split_pair(string * line, string * param, string * value)
{
	size_t pos;

	if ((pos = line->find(':')) != string::npos)
	{
		if (pos + 2 >= line->size())
			return false;
		*param = line->substr(0, pos);
		*value = line->substr(pos + 2);
		return true;
	}

	return false;
}

int MPD::recv_songs_to_list(Songlist * slist, void (*func) ())
{
	Song * song = NULL;
	Field * field;
	string buf;
	string param;
	string value;
	int status;
	unsigned int count = 0;

	while((status = mpd_getline(&buf)) == MPD_GETLINE_MORE)
	{
		if (!split_pair(&buf, &param, &value))
			continue;

		field = fieldtypes.find_mpd(param);
		if (field == NULL)
		{
			//debug("Unhandled song metadata field '%s' in response from MPD", param.c_str());
			continue;
		}

		if (field->type == FIELD_FILE)
		{
			if (song != NULL)
			{
				song->init();
				slist->add(song);
				//debug(song->f[FIELD_FILE].c_str(), NULL);
				if (func != NULL && ++count % 500 == 0)
					func();
			}

			song = new Song;
		}

		if (song != NULL)
			song->f[field->type] = value;
	}

	if (song != NULL)
	{
		song->init();
		slist->add(song);
	}

	return status;
}

int MPD::get_playlist()
{
	int s;

	/* Ignore duplicate update messages */
	if (playlist.version == status.playlist)
		return false;

	playlist.truncate(status.playlistlength);
	if (playlist.version == -1)
		mpd_send("playlistinfo");
	else
		mpd_send("plchanges %d", playlist.version);

	s = recv_songs_to_list(&playlist, NULL);

	playlist.version = status.playlist;
	update_currentsong();
	wm.playlist->update_column_length();

	debug("Playlist has been updated to version %d", playlist.version);

	return s;
}

int MPD::get_library()
{
	int s;

	get_stats();

	if ((unsigned long long)library.version == stats.db_update)
	{
		debug("Request for library update, but local copy is the same as server.", NULL);
		return MPD_GETLINE_OK;
	}

	library.clear();
	library.version = -1;
	library.truncate(stats.songs);

	mpd_send("listallinfo");
	if ((s = recv_songs_to_list(&library, update_library_statusbar)) == MPD_GETLINE_OK)
	{
		library.version = stats.db_update;
		stinfo("Successfully received library, total %d songs.", library.size());
	}
	else
	{
		sterr("Library update terminated!", NULL);
		debug("Library update was terminated by MPD. We got %d of a total of %d songs.", library.size(), stats.songs);
		debug("This sometimes happens due to the large volume of data transferred by PMS.", NULL);
		debug("If this happens often, you might need to increase MPD's `max_output_buffer_size' setting.", NULL);
	}

	/* Manual draw and sort, which may take a while */
	stinfo("Sorting library by `%s'...", config.default_sort.c_str());
	wm.statusbar->draw();
	library.sort(config.default_sort);
	debug("Library is sorted.", NULL);

	wm.library->update_column_length();
	wm.library->qdraw();

	return s;
}

int MPD::get_stats()
{
	string buf;
	string param;
	string value;
	int s;

	mpd_send("stats");

	while((s = mpd_getline(&buf)) == MPD_GETLINE_MORE)
	{
		if (!split_pair(&buf, &param, &value))
			continue;

		if (param == "artists")
			stats.songs = atol(value.c_str());
		else if (param == "albums")
			stats.albums = atoll(value.c_str());
		else if (param == "songs")
			stats.songs = atoll(value.c_str());
		else if (param == "uptime")
			stats.uptime = atoll(value.c_str());
		else if (param == "playtime")
			stats.playtime = atoll(value.c_str());
		else if (param == "db_playtime")
			stats.db_playtime = atoll(value.c_str());
		else if (param == "db_update")
			stats.db_update = atoll(value.c_str());
	}

	return s;
}

int MPD::get_status()
{
	string buf;
	string param;
	string value;
	int s;
	size_t pos;

	mpd_send("status");

	while((s = mpd_getline(&buf)) == MPD_GETLINE_MORE)
	{
		if (!split_pair(&buf, &param, &value))
			continue;

		if (param == "volume")
			status.volume = atoi(value.c_str());
		else if (param == "repeat")
			status.repeat = atoi(value.c_str());
		else if (param == "random")
			status.random = atoi(value.c_str());
		else if (param == "single")
			status.single = atoi(value.c_str());
		else if (param == "consume")
			status.consume = atoi(value.c_str());
		else if (param == "playlist")
			status.playlist = atoi(value.c_str());
		else if (param == "playlistlength")
			status.playlistlength = atoi(value.c_str());
		else if (param == "xfade")
			status.xfade = atoi(value.c_str());
		else if (param == "mixrampdb")
			status.mixrampdb = atof(value.c_str());
		else if (param == "mixrampdelay")
			status.mixrampdelay = atoi(value.c_str());
		else if (param == "song")
			status.song = atol(value.c_str());
		else if (param == "songid")
			status.songid = atol(value.c_str());
		else if (param == "elapsed")
			status.elapsed = atof(value.c_str());
		else if (param == "bitrate")
			status.bitrate = atoi(value.c_str());
		else if (param == "nextsong")
			status.nextsong = atol(value.c_str());
		else if (param == "nextsongid")
			status.nextsongid = atol(value.c_str());

		else if (param == "state")
		{
			if (value == "play")
				status.state = MPD_STATE_PLAY;
			else if (value == "stop")
				status.state = MPD_STATE_STOP;
			else if (value == "pause")
				status.state = MPD_STATE_PAUSE;
			else
				status.state = MPD_STATE_UNKNOWN;
		}

		else if (param == "time")
		{
			if ((pos = value.find(':')) != string::npos)
			{
				status.elapsed = atof(value.substr(0, pos).c_str());
				status.length = atoi(value.substr(pos + 1).c_str());
			}
		}

		else if (param == "audio")
		{
			if ((pos = value.find(':')) != string::npos)
			{
				status.samplerate = atoi(value.substr(0, pos).c_str());
				status.bits = atoi(value.substr(pos + 1).c_str());
				if ((pos = value.find(':', pos + 1)) != string::npos)
				{
					status.channels = atoi(value.substr(pos + 1).c_str());
				}
			}
		}
	}

	gettimeofday(&last_update, NULL);
	memcpy(&last_clock, &last_update, sizeof last_clock);
	update_currentsong();

	return s;
}

void MPD::run_clock()
{
	Song * song;
	struct timeval tm;
	gettimeofday(&tm, NULL);

	if (status.state == MPD_STATE_PLAY)
	{
		status.elapsed += (tm.tv_sec - last_clock.tv_sec);
		status.elapsed += (tm.tv_usec - last_clock.tv_usec) / 1000000.000000;
		wm.topbar->qdraw();
	}

	memcpy(&last_clock, &tm, sizeof last_clock);

	if (config.autoadvance && autoadvance_playlist < playlist.version && status.length - status.elapsed < (int)config.add_next_interval)
	{
		if ((song = next_auto_song_in_line()) != NULL)
		{
			if ((addid(song->f[FIELD_FILE])) != -1)
				wm.statusbar->qdraw();
			autoadvance_playlist = playlist.version;
		}
	}
}

int MPD::poll()
{
	string line;
	string param;
	string value;
	int updates;
	fd_set set;
	struct timeval timeout;
	int s;

	set_idle(true);

	memset(&timeout, 0, sizeof timeout);
	memcpy(&set, &fdset, sizeof set);
	timeout.tv_sec = 1;
	if ((s = select(sock+1, &set, NULL, NULL, &timeout)) == -1)
	{
		if (errno == EINTR)
			return false;

		mpd_disconnect();
		return true;
	}

	/* Update elapsed time */
	run_clock();

	if (FD_ISSET(STDIN_FILENO, &set))
		return false;
	if (!FD_ISSET(sock, &set))
		return false;

	updates = MPD_UPDATE_NONE;

	while((s = mpd_getline(&line)) == MPD_GETLINE_MORE)
	{
		if (!split_pair(&line, &param, &value))
			continue;

		if (param != "changed")
			continue;

		if (value == "database")
		{
			// the song database has been modified after update. 
			updates |= MPD_UPDATE_LIBRARY;
		}
		else if (value == "update")
		{
			// a database update has started or finished. If the database was modified during the update, the database event is also emitted. 
			updates |= MPD_UPDATE_DB;
		}
		else if (value == "stored_playlist")
		{
			// a stored playlist has been modified, renamed, created or deleted 
		}
		else if (value == "playlist")
		{
			// the current playlist has been modified 
			updates |= MPD_UPDATE_STATUS;
			updates |= MPD_UPDATE_PLAYLIST;
		}
		else if (value == "player")
		{
			// the player has been started, stopped or seeked
			updates |= MPD_UPDATE_STATUS;
		}
		else if (value == "mixer")
		{
			// the volume has been changed 
			updates |= MPD_UPDATE_STATUS;
		}
		else if (value == "output")
		{
			// an audio output has been enabled or disabled 
		}
		else if (value == "options")
		{
			// options like repeat, random, crossfade, replay gain
			updates |= MPD_UPDATE_STATUS;
		}
		else if (value == "sticker")
		{
			// the sticker database has been modified.
		}
		else if (value == "subscription")
		{
			// a client has subscribed or unsubscribed to a channel
		}
		else if (value == "message")
		{
			// a message was received on a channel this client is subscribed to; this event is only emitted when the queue is empty
		}
	}

	if (updates & MPD_UPDATE_STATUS)
		get_status();
	if (updates & MPD_UPDATE_PLAYLIST)
		get_playlist();
	if (updates & MPD_UPDATE_LIBRARY)
		get_library();

	set_idle(true);
	update_playstring();

	return true;
}

Song * MPD::update_currentsong()
{
	wm.topbar->qdraw();
	wm.active->qdraw();
	return (currentsong = (int)playlist.songs.size() > status.song ? playlist.songs[status.song] : NULL);
}

int MPD::apply_opts()
{
	bool r;

	/* FIXME: maybe these should be sanitized before being allowed to set them? */
	if (config.volume > 100)
		config.volume = 100;
	else if (config.volume < 0)
		config.volume = 0;

	if (mpd.status.volume != (config.mute ? 0 : config.volume))
		set_volume(config.mute ? 0 : config.volume);
	if (mpd.status.single != config.single)
		set_single(config.single);
	if (mpd.status.repeat != config.repeat)
		set_repeat(config.repeat);
	if (mpd.status.consume != config.consume)
		set_consume(config.consume);

	/* Shuffle/random is a special case, since PMS has it's own random functions. */
	if (config.random)
	{
		if ((r = active_songlist == &playlist) != mpd.status.random)
			set_random(r);
	}
	else if (mpd.status.random)
		set_random(config.random);

	return true;
}

int MPD::read_opts()
{
	config.volume = mpd.status.volume;
	config.single = mpd.status.single;
	config.repeat = mpd.status.repeat;
	config.consume = mpd.status.consume;
	if (active_songlist == &playlist)
		config.random = mpd.status.random;

	return true;
}

int MPD::activate_songlist(Songlist * list)
{
	if (!list)
		return false;
	
	active_songlist = list;
	apply_opts();
	update_playstring();
	wm.statusbar->qdraw();

	return true;
}

int MPD::remove(Songlist * list, int start, int count)
{
	if (!list || start < 0 || start >= (int)list->size() || count == 0)
		return false;

	/* Reverse deletion */
	if (count < 0)
	{
		start += count;
		count = -count;
		if (start < 0)
			return false;
	}

	if (start + count > (int)list->size())
		count = list->size() - start;

	if (list == &playlist)
	{
		mpd_send("delete %d:%d", start, start + count);
		return (mpd_getline(NULL) == MPD_GETLINE_OK);
	}

	return false;
}

Song * MPD::next_song_in_line(int steps)
{
	int i = -1;

	/* Need songs in active list */
	if (active_songlist->size() == 0)
		return NULL;

	/* Hold on, are there more songs in the playlist? */
	if (currentsong && active_songlist != &playlist)
	{
		i = playlist.spos(currentsong->pos);

		/* Yes there are. */
		if (i + 1 != (int)playlist.size())
		{
			/* Enough room to step up to that song? */
			if (steps + i < (int)playlist.size())
				return playlist.at(steps + i);

			/* Steps exceeds playlist queue size, jump to active list. */
			else
			{
				steps -= (playlist.size() - i);
				if ((i = active_songlist->sfind(currentsong->fhash)) == -1)
					i = 0;
			}
		}
		else
		{
			i = -1;
		}
	}

	if (i == -1)
	{
		if (!config.random && currentsong)
		{
			if (active_songlist == &playlist)
				i = playlist.spos(currentsong->pos);
			else
				i = active_songlist->sfind(currentsong->fhash);
		}
		else if (config.random)
		{
			return active_songlist->at(active_songlist->randpos());
		}
	}

	if (!currentsong)
		return active_songlist->at(0);

	steps += i;

	/* Reached end of list, wrap around. */
	if (steps >= (int)active_songlist->size())
	{
		/* ... or continue, if we are not at the end of the playlist */
		if (active_songlist != &playlist && currentsong->pos + 1 != (int)playlist.size())
		{
			steps %= playlist.size();
		}
		else

		/* ... unless we are not repeating ourselves. */
		{
			if (!status.repeat)
				return NULL;
			steps %= active_songlist->size();
		}
	}

	/* Reached beginning of list, wrap around. */
	while (steps < 0)
	{
		/* ... unless we are not repeating ourselves. */
		if (!status.repeat)
			return NULL;
		steps += active_songlist->size();
	}

	return active_songlist->at(steps);
}

Song * MPD::next_auto_song_in_line()
{
	/* Let MPD manage the playlist itself */
	if (!active_songlist || active_songlist == &playlist || status.random)
		return NULL;

	/* No song in line if single-mode */
	if (status.single)
		return NULL;

	/* Don't auto-progress if playlist still has queued songs */
	if (status.song + 1 < (int)playlist.size())
		return NULL;

	return next_song_in_line();
}

void MPD::update_playstring()
{
	Songlist * list;
	string str;
	bool islast;

	if (!is_connected())
	{
		playstring = "Not connected to MPD server.";
		return;
	}
	
	if (status.state == MPD_STATE_STOP || !currentsong)
	{
		playstring = "Stopped.";
		return;
	}

	if (status.state == MPD_STATE_PAUSE)
	{
		playstring = "Paused...";
		return;
	}

	list = config.autoadvance && !status.random ? active_songlist : &playlist;

	if (status.consume)
		playstring = "Consuming ";
	else
		playstring = "Playing ";

	islast = (status.song + 1 == (int)playlist.size());

	if (status.single || (islast && list == &playlist && !status.random && !status.repeat))
	{
		if (!status.consume && status.repeat)
		{
			playstring += "this song forever.";
			return;
		}
		playstring += "this song, then stopping.";
		return;
	}

	if (status.random)
	{
		playstring += "random songs from playlist";
		if (status.consume)
			playstring += " until empty";
		else if (status.repeat)
			playstring += " forever";
		else
			playstring += ", stopping when all songs have been played";

		playstring += ".";
		return;
	}

	if (!islast && list != &playlist)
	{
		if (!status.consume)
			playstring += "through playlist, then ";
		else
			playstring += "playlist, then playing ";
	}

	if (config.random)
		playstring += "random ";
		
	playstring += "songs from " + list->title;

	if (status.consume && list == &playlist)
		playstring += " until empty";
	else if (status.repeat)
		playstring += " repeatedly";

	playstring += ".";
}

int MPD::update(string dir)
{
	mpd_send("update \"%s\"", dir.c_str());
	while (mpd_getline(NULL) == MPD_GETLINE_MORE);
	return true;
}

int MPD::set_consume(bool nconsume)
{
	mpd_send("consume %d", nconsume);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::set_crossfade(unsigned int nseconds)
{
	mpd_send("crossfade %d", nseconds);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::set_mixrampdb(int ndecibels)
{
	return false;
}

int MPD::set_mixrampdelay(int nseconds)
{
	return false;
}

int MPD::set_random(bool nrandom)
{
	mpd_send("random %d", nrandom);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::set_repeat(bool nrepeat)
{
	mpd_send("repeat %d", nrepeat);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::set_volume(unsigned int nvol)
{
	mpd_send("setvol %d", nvol);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::set_single(bool nsingle)
{
	mpd_send("single %d", nsingle);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::set_replay_gain_mode(replay_gain_mode nrgm)
{
	return false;
}

int MPD::pause(bool npause)
{
	mpd_send("pause %d", npause);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::addid(string uri)
{
	string buf;
	string param;
	string value;
	int status;

	mpd_send("addid \"%s\"", uri.c_str());

	while ((status = mpd_getline(&buf)) == MPD_GETLINE_MORE);

	if (!split_pair(&buf, &param, &value))
		return -1;

	return atoi(value.c_str());
}

int MPD::playid(int id)
{
	mpd_send("playid %d", id);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::stop()
{
	mpd_send("stop");
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::next()
{
	mpd_send("next");
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::previous()
{
	mpd_send("previous");
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

int MPD::seek(int seconds)
{
	Song * song;
	int pos;

	if (status.state != MPD_STATE_PLAY || !currentsong)
	{
		sterr("Cannot seek when not playing anything.", NULL);
		return false;
	}

	if (currentsong->time == -1)
	{
		sterr("Cannot seek in this song. It might be a stream.", NULL);
		return false;
	}

	song = currentsong;
	pos = round(status.elapsed) + seconds;
	while (pos < 0 && song->pos > 0)
	{
		song = playlist.songs[song->pos - 1];
		pos = song->time - pos;
	}
	while (pos > song->time && song->pos + 1 < (int)playlist.size())
	{
		pos = pos - song->time;
		song = playlist.songs[song->pos + 1];
	}

	mpd_send("seek %d %d", song->pos, pos);
	return (mpd_getline(NULL) == MPD_GETLINE_OK);
}

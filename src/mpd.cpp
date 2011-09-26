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
#include <sys/select.h>
#include <sys/time.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>
#include <netdb.h>
#include <cstring>
#include <string>
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>

using namespace std;

extern Config config;
extern Windowmanager wm;

MPD::MPD()
{
	errno = 0;
	error = "";
	waiting = false;
	host = "";
	port = "";
	buffer = "";
	sock = 0;
	connected = false;
	memset(&last_update, 0, sizeof last_update);
	memset(&last_clock, 0, sizeof last_clock);
	memset(&state, 0, sizeof state);
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
	is_idle = false;
	mpd_getline(NULL);

	return true;
}

bool MPD::trigerr(int nerrno, const char * format, ...)
{
	va_list		ap;
	char		buffer[1024];

	va_start(ap, format);
	vsprintf(buffer, format, ap);
	va_end(ap);

	error = buffer;
	errno = nerrno;

	sterr("MPD: %s", buffer);

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

	stinfo("Connected to server '%s' on port '%s'.", host.c_str(), port.c_str());
	recv(sock, &buf, 32, 0);
	set_protocol_version(buf);

	return connected;
}

void MPD::mpd_disconnect()
{
	close(sock);
	sock = 0;
	connected = false;
	trigerr(MPD_ERR_CONNECTION, "connection closed");
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
	debug("MPD server speaking protocol version %d.%d.%d", protocol_version[0], protocol_version[1], protocol_version[2]);

	return true;
}

int MPD::mpd_send(const char * data, ...)
{
	va_list		ap;
	char		buffer[1024];

	va_start(ap, data);
	vsprintf(buffer, data, ap);
	va_end(ap);

	if (!connected)
		return -1;

	set_idle(false);
	return mpd_raw_send(buffer);
}

int MPD::mpd_raw_send(string data)
{
	unsigned int sent;
	int s;

	if (!connected)
		return -1;

	data += '\n';
	if ((s = send(sock, data.c_str(), data.size(), 0)) == -1)
		return s;

	sent = s;

	while (sent < data.size())
	{
		if ((s = send(sock, data.substr(sent).c_str(), data.size() - sent, 0)) == -1)
			return -1;

		sent += s;
	}

	waiting = true;

	// Raw traffic dump
	debug("-> %s", data.c_str());

	return sent;
}

int MPD::mpd_getline(string * nextline)
{
	char buf[1024];
	int received = 0;
	size_t pos;
	string line = "";

	if (!connected)
		return MPD_GETLINE_ERR;

	while(buffer.size() == 0 || buffer[buffer.size()-1] != '\n')
	{
		received = recv(sock, &buf, 1023, 0);
		if (received == 0)
		{
			mpd_disconnect();
			return MPD_GETLINE_ERR;
		}
		else if (received == -1)
		{
			continue;
		}
		buf[received] = '\0';
		buffer += buf;
	}

	if ((pos = buffer.find('\n')) != string::npos)
	{
		line = buffer.substr(0, pos);
		if (buffer.size() == pos + 1)
			buffer = "";
		else
			buffer = buffer.substr(pos + 1);
	}

	if (line.size() == 0)
		return MPD_GETLINE_ERR;

	// Raw traffic dump
	debug("<- %s", line.c_str());

	if (line == "OK")
		return MPD_GETLINE_OK;

	if (line.size() >= 3 && line.substr(0, 3) == "ACK")
	{
		trigerr(MPD_ERR_ACK, buf);
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
		*param = line->substr(0, pos);
		*value = line->substr(pos + 2);
		return true;
	}

	return false;
}

int MPD::get_playlist()
{
	Song * song = NULL;
	string buf;
	string param;
	string value;
	int status;

	/* Ignore duplicate update messages */
	if (playlist.version == state.playlist)
		return false;

	playlist.truncate(state.playlistlength);
	if (playlist.version == -1)
		mpd_send("playlistinfo");
	else
		mpd_send("plchanges %d", playlist.version);

	while((status = mpd_getline(&buf)) == MPD_GETLINE_MORE)
	{
		if (!split_pair(&buf, &param, &value))
			continue;

		if (param == "file")
		{
			if (song != NULL)
				playlist.add(song);

			song = new Song;
			song->file = value;
		}
		else if (param == "Pos")
			song->pos = atol(value.c_str());
		else if (param == "Id")
			song->id = atol(value.c_str());
		else if (param == "Time")
			song->length = atoi(value.c_str());
		else if (param == "Name")
			song->name = value;
		else if (param == "Artist")
			song->artist = value;
		else if (param == "Title")
			song->title = value;
		else if (param == "Album")
			song->album = value;
		else if (param == "Track")
			song->track = value;
		else if (param == "Date")
			song->date = value;
		else if (param == "Genre")
			song->genre = value;
	}
	playlist.add(song);
	playlist.version = state.playlist;
	wm.playlist->draw();

	debug("Playlist has been updated to version %d", playlist.version);

	return status;
}

int MPD::get_status()
{
	string buf;
	string param;
	string value;
	int status;
	size_t pos;

	mpd_send("status");

	while((status = mpd_getline(&buf)) == MPD_GETLINE_MORE)
	{
		if (!split_pair(&buf, &param, &value))
			continue;

		if (param == "volume")
			state.volume = atoi(value.c_str());
		else if (param == "repeat")
			state.repeat = atoi(value.c_str());
		else if (param == "random")
			state.random = atoi(value.c_str());
		else if (param == "single")
			state.single = atoi(value.c_str());
		else if (param == "consume")
			state.consume = atoi(value.c_str());
		else if (param == "playlist")
			state.playlist = atoi(value.c_str());
		else if (param == "playlistlength")
			state.playlistlength = atoi(value.c_str());
		else if (param == "xfade")
			state.xfade = atoi(value.c_str());
		else if (param == "mixrampdb")
			state.mixrampdb = atof(value.c_str());
		else if (param == "mixrampdelay")
			state.mixrampdelay = atoi(value.c_str());
		else if (param == "song")
			state.song = atol(value.c_str());
		else if (param == "songid")
			state.songid = atol(value.c_str());
		else if (param == "elapsed")
			state.elapsed = atof(value.c_str());
		else if (param == "bitrate")
			state.bitrate = atoi(value.c_str());
		else if (param == "nextsong")
			state.nextsong = atol(value.c_str());
		else if (param == "nextsongid")
			state.nextsongid = atol(value.c_str());

		else if (param == "state")
		{
			if (value == "play")
				state.state = MPD_STATE_PLAY;
			else if (value == "stop")
				state.state = MPD_STATE_STOP;
			else if (value == "pause")
				state.state = MPD_STATE_PAUSE;
			else
				state.state = MPD_STATE_UNKNOWN;
		}

		else if (param == "time")
		{
			if ((pos = value.find(':')) != string::npos)
			{
				state.elapsed = atof(value.substr(0, pos).c_str());
				state.length = atoi(value.substr(pos + 1).c_str());
			}
		}

		else if (param == "audio")
		{
			if ((pos = value.find(':')) != string::npos)
			{
				state.samplerate = atoi(value.substr(0, pos).c_str());
				state.bits = atoi(value.substr(pos + 1).c_str());
				if ((pos = value.find(':', pos + 1)) != string::npos)
				{
					state.channels = atoi(value.substr(pos + 1).c_str());
				}
			}
		}
	}

	gettimeofday(&last_update, NULL);
	memcpy(&last_clock, &last_update, sizeof last_clock);

	return status;
}

void MPD::run_clock()
{
	struct timeval tm;
	gettimeofday(&tm, NULL);

	state.elapsed += (tm.tv_sec - last_clock.tv_sec);
	state.elapsed += (tm.tv_usec - last_clock.tv_usec) / 1000000.000000;

	debug("Elapsed: %f", state.elapsed);

	memcpy(&last_clock, &tm, sizeof last_clock);
}

int MPD::poll()
{
	string line;
	string param;
	string value;
	int updates;
	struct timeval timeout;
	fd_set set;
	int s;

	set_idle(true);

	FD_ZERO(&set);
	FD_SET(sock, &set);
	FD_SET(STDIN_FILENO, &set);

	memset(&timeout, 0, sizeof timeout);
	timeout.tv_usec = 1000000;
	if ((s = select(sock+1, &set, NULL, NULL, &timeout)) == -1)
	{
		mpd_disconnect();
		return false;
	}
	else if (s == 0)
	{
		// no data ready to recv(), but let's update our clock
		run_clock();
		return false;
	}

	if (!FD_ISSET(sock, &set))
		return true;

	is_idle = false;
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
		}
		else if (value == "update")
		{
			// a database update has started or finished. If the database was modified during the update, the database event is also emitted. 
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

	set_idle(true);

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

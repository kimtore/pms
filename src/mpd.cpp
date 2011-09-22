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
	memset(&state, 0, sizeof state);
}

bool MPD::set_idle(bool nidle)
{
	if (nidle == is_idle)
		return false;
	
	if (nidle)
	{
		mpd_send("idle");
		return true;
	}

	mpd_send("noidle");
	mpd_getline(NULL);
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
	
	mpd_send("password \"" + password + "\"");
	if (mpd_getline(NULL) == MPD_GETLINE_OK)
	{
		stinfo("Password '%s' accepted by server.", password.c_str());
		return true;
	}

	return true;
}

bool MPD::set_protocol_version(string data)
{
	int i = 7;
	int last = 7;
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

int MPD::mpd_send(string data)
{
	int sent;

	if (!connected)
		return -1;

	data += '\n';
	sent = send(sock, data.c_str(), data.size(), 0);
	while (sent < data.size())
		sent += send(sock, data.substr(sent).c_str(), data.size() - sent, 0);

	waiting = true;

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

int MPD::get_status()
{
	string buf;
	string param;
	string value;
	int status;
	size_t pos;

	set_idle(false);
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
			state.elapsed = atoi(value.c_str());
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
				state.elapsed = atoi(value.substr(0, pos).c_str());
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

	set_idle(true);
	return true;
}

int MPD::poll()
{
	string line;
	struct timeval timeout;
	fd_set set;
	int s;

	FD_ZERO(&set);
	FD_SET(sock, &set);

	memset(&timeout, 0, sizeof timeout);
	if ((s = select(sock+1, &set, NULL, NULL, &timeout)) == -1)
	{
		mpd_disconnect();
		return false;
	}
	else if (s == 0)
	{
		return false;
	}

	is_idle = false;

	while((s = mpd_getline(&line)) == MPD_GETLINE_MORE)
	{
		debug("IDLE returned %s", line.c_str());
	}

	set_idle(true);

	return true;
}

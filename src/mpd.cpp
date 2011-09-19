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
#include "debug.h"
#include "config.h"
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <cstring>
#include <string>
#include <stdlib.h>

using namespace std;

extern Config config;

bool MPD::mpd_connect(string nhost, string nport)
{
	int			status;
	char			buf[32];
	struct addrinfo		hints;
	struct addrinfo *	res;

	this->host = nhost;
	this->port = nport;

	this->mpd_disconnect();

	memset(&hints, 0, sizeof hints);
	hints.ai_family = AF_UNSPEC;
	hints.ai_socktype = SOCK_STREAM;

	if ((status = getaddrinfo(this->host.c_str(), this->port.c_str(), &hints, &res)) != 0)
	{
		debug("getaddrinfo error: %s", gai_strerror(status));
		freeaddrinfo(res);
		return false;
	}

	this->sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol);
	if (this->sock == -1)
	{
		freeaddrinfo(res);
		return false;
	}

	if (connect(this->sock, res->ai_addr, res->ai_addrlen) == -1)
	{
		close(this->sock);
		freeaddrinfo(res);
		return false;
	}

	freeaddrinfo(res);
	this->connected = true;

	debug("Successful connection to %s:%s.", this->host.c_str(), this->port.c_str());
	recv(this->sock, &buf, 32, 0);
	this->set_protocol_version(buf);

	send_and_recv("status\n");

	return this->connected;
}

void MPD::mpd_disconnect()
{
	close(this->sock);
	this->sock = 0;
	this->connected = false;
}

bool MPD::is_connected()
{
	return this->connected;
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
			this->protocol_version[pos] = atoi(data.substr(last, i - last).c_str());
			++pos;
			last = i + 1;
		}
		++i;
	}
	debug("MPD server speaking protocol version %d.%d.%d", protocol_version[0], protocol_version[1], protocol_version[2]);

	return true;
}

bool MPD::send_and_recv(string data)
{
	int sent;
	char buf[1024];

	if (!this->connected)
		return false;

	this->buffer = "";

	sent = send(this->sock, data.c_str(), data.size(), 0);
	while (sent < data.size())
		sent += send(this->sock, data.substr(sent).c_str(), data.size() - sent, 0);

	recv(this->sock, &buf, 1024, 0);
	debug("data %d bytes: %s:", strlen(buf), buf);

	return true;
}

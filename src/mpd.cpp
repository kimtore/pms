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

using namespace std;

extern Config config;

bool MPD::mpd_connect(string nhost, string nport)
{
	int			status;
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
	debug("Successful connection to %s:%s.", this->host.c_str(), this->port.c_str());

	return true;
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

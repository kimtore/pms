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
 * conn.cpp
 * 	connection handle to MPD
 */

#include "conn.h"
#include "pms.h"

extern Pms *	pms;


Connection::Connection(string n_hostname, long n_port, long n_timeout)
{
	this->host = n_hostname;
	this->port = static_cast<unsigned int>(n_port);
	this->timeout = static_cast<int>(n_timeout);

	this->handle = NULL;
}

Connection::~Connection()
{
	this->disconnect();
}

string Connection::errorstr()
{
	string		mystr = "";

	if (handle != NULL)
	{
		mystr += handle->errorStr;
	}

	return mystr;
}

/*
 * Returns connection state
 */
bool		Connection::connected()
{
	if (handle == NULL)
		return false;

	switch(handle->error)
	{
		case 0:
		default:
			return true;

		case MPD_ERROR_CONNCLOSED:
		case MPD_ERROR_TIMEOUT:
		case MPD_ERROR_CONNPORT:
		case MPD_ERROR_NOTMPD:
		case MPD_ERROR_NORESPONSE:
			return false;
	}
}

int Connection::connect()
{
	pms->log(MSG_DEBUG, 0, "Connecting to %s:%d, handle=%p...\n", host.c_str(), port, handle);
	if (handle != NULL)
	{
		if (handle->error == 0)
			return 0;

		mpd_clearError(handle);
		disconnect();
	}
	if (handle == NULL)
	{
		handle = mpd_newConnection(host.c_str(), port, timeout);
		pms->log(MSG_DEBUG, 0, "New connection handle is %p, error code %d\n", handle, handle->error);
		error = handle->error;

		return error;
	}

	return handle->error;
}

int Connection::disconnect()
{
	pms->log(MSG_DEBUG, 0, "Closing connection.\n");
	if (handle != NULL)
	{
		mpd_closeConnection(handle);
		handle = NULL;
	}

	return 0;
}

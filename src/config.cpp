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

#include "config.h"
#include "field.h"
#include "console.h"
#include <stdlib.h>

using namespace std;

extern Fieldtypes fieldtypes;

Config::Config()
{
	setup_default_connection_info();

	quit = false;
	reconnect_delay = 5;
	use_bell = true;
	visual_bell = false;
	set_column_headers("artist track title album year length");
}

void Config::set_column_headers(string hdr)
{
	size_t start = 0;
	size_t pos;
	string f;
	Field * field;

	songlist_columns.clear();

	while (start + 1 < hdr.size())
	{
		if (pos == string::npos)
			break;

		if ((pos = hdr.find(' ', start)) != string::npos)
			f = hdr.substr(start, pos - start);
		else
			f = hdr.substr(start);

		if ((field = fieldtypes.find(f)) == NULL)
		{
			sterr("Ignoring invalid header field '%s'.", f.c_str());
			continue;
		}
		songlist_columns.push_back(field);

		start = pos + 1;
	}
}

void Config::setup_default_connection_info()
{
	char *	env;
	size_t	i;

	password = "";

	if ((env = getenv("MPD_HOST")) == NULL)
	{
		host = "localhost";
	}
	else
	{
		host = env;
		if ((i = host.rfind('@')) != string::npos)
		{
			password = host.substr(0, i);
			host = host.substr(i + 1);
		}
	}

	if ((env = getenv("MPD_PORT")) == NULL)
		port = "6600";
	else
		port = env;
}

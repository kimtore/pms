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
 */

#ifndef _PMS_CONN_H_
#define _PMS_CONN_H_

#include <mpd/client.h>
#include <string>

using namespace std;


class Connection
{
private:
	mpd_connection	*handle;
	string		host;
	unsigned int	port;
	int		timeout;
	int		error;

public:
			Connection(string, long, long);
			~Connection();

	mpd_connection	*h() { return handle; };

	bool		connected();
	int		connect();
	int		disconnect();
	bool		clear_error();
};
 
#endif /* _PMS_CONN_H_ */

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

#ifndef _PMS_MPD_H_
#define _PMS_MPD_H_

#include "songlist.h"
#include "song.h"
#include <string>

/* MPD error codes */
#define MPD_ERR_NONE 0
#define MPD_ERR_CONNECTION 1
#define MPD_ERR_NOTMPD 2
#define MPD_ERR_ACK 3
#define MPD_ERR_BADPASS 4

/* mpd_getline statuses */
#define MPD_GETLINE_ERR -1
#define MPD_GETLINE_ACK -1
#define MPD_GETLINE_OK 0
#define MPD_GETLINE_MORE 1

/* MPD player states */
#define MPD_STATE_UNKNOWN -1
#define MPD_STATE_PLAY 0
#define MPD_STATE_STOP 1
#define MPD_STATE_PAUSE 2

/* IDLE command updates */
enum
{
	MPD_UPDATE_NONE = 0,
	MPD_UPDATE_STATUS = 1 << 0,
	MPD_UPDATE_PLAYLIST = 1 << 1
};

typedef struct
{
	int		volume;
	bool		repeat;
	bool		random;
	bool		single;
	bool		consume;
	long		playlist;
	int		playlistlength;
	int		xfade;
	double		mixrampdb;
	int		mixrampdelay;
	int		state;
	song_t		song;
	song_t		songid;
	song_t		nextsong;
	song_t		nextsongid;
	int		length;
	int		elapsed;
	int		bitrate;
	long		samplerate;
	int		bits;
	int		channels;
}

mpd_state;

using namespace std;

class MPD
{
	private:
		string		host;
		string		port;
		string		password;
		string		buffer;

		string		error;
		int		errno;

		/* Connection variables */
		int		sock;
		bool		connected;
		bool		waiting;
		int		protocol_version[3];
		bool		is_idle;

		/* MPD state */
		mpd_state	state;

		/* Set/unset idle status */
		bool		set_idle(bool nidle);

		/* Trigger an error. Always returns false. */
		bool		trigerr(int nerrno, const char * format, ...);

		/* Parse the initial connection string from MPD */
		bool		set_protocol_version(string data);

		/* Send a command to MPD, turning IDLE off if needed */
		int		mpd_send(const char * data, ...);

		/* Send a command to MPD */
		int		mpd_raw_send(string data);

		/* Get data from MPD and fetch next line. See MPD_GETLINE_* for return codes */
		int		mpd_getline(string * nextline);

		/* Split a "parameter: value" pair */
		int		split_pair(string * line, string * param, string * value);

	public:
		MPD();

		/* Server-side lists */
		Songlist	playlist;
		Songlist	library;

		/* Initialise a connection to an MPD server */
		bool		mpd_connect(string host, string port);

		/* Shut it down, houston. */
		void		mpd_disconnect();

		/* Returns true if there is an active connection. */
		bool		is_connected();

		/* Change password */
		bool		set_password(string password);

		/* Fetch MPD playlist updates since last time */
		int		get_playlist();

		/* Retrieve MPD status */
		int		get_status();

		/* Polls the socket to see if there is any IDLE data to collect. */
		int		poll();

};

#endif /* _PMS_MPD_H_ */

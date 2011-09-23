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

#ifndef _PMS_SONG_H_
#define _PMS_SONG_H_

#include <string>
using namespace std;

typedef long song_t;

class Song
{
	public:
		Song();
		
		/* Common function to initialize special fields that MPD don't return */

		void		init();
		string		dirname();

		/* Custom parameters only used by PMS */
		
		string		trackshort;

		string		file;
		string		artist;
		string		albumartist;
		string		albumartistsort;
		string		artistsort;
		string		title;
		string		album;
		string		track;
		string		name;
		string		date;
		string		year;

		string		genre;
		string		composer;
		string		performer;
		string		disc;
		string		comment;

		int		length;
		song_t		pos;
		song_t		id;
};

/* Case insensitive string comparison */
bool cistrcmp(string &a, string &b);

#endif /* _PMS_SONG_H_ */

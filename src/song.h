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

#include "fields.h"
#include <string>
#include <locale>
using namespace std;

typedef long song_t;

class Song
{
	public:
		Song();

		/* Song fields, see field.h */
		string		f[FIELD_COLUMN_VALUES];

		/* Fields that are not strings */
		int		pos;
		int		id;
		int		time;

		/* For quick lookup through filename */
		long			fhash;

		/* Common function to initialize special fields that MPD don't return */
		void		init();
		string		dirname();
};

/* Case insensitive string comparison */
bool cistrcmp(string &a, string &b);

/* Format time into "[n D ][H:]MM:SS" */
string time_format(int seconds);

/* Pad an int with zeroes */
string zeropad(int i, unsigned int target);

/* Search and replace string */
string str_replace(string search, string replace, string subject);

/* Convert number to string */
string tostring(int number);
string tostring(unsigned int number);

/* Correctly escape a string so that it can be printed */
void escape_printf(string &src);

/* Hash a string */
long songhash(string const &str);

#endif /* _PMS_SONG_H_ */

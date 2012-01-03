/* vi:set ts=8 sts=8 sw=8 noet:
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

#ifndef _PMS_FIELD_H_
#define _PMS_FIELD_H_

#include <vector>
#include <string>
#include "fields.h"
#include "song.h"
using namespace std;

/*
 * Song metadata field
 */
class Field
{
	public:
		Field(field_t nfield, string name, string mpd_name, string tit, unsigned int minl, unsigned int maxl);

		/* Format a field to a specific song */
		string		format(Song * song);

		/* Which kind of field is this? */
		field_t		type;

		/* MPD case string representation, e.g. «artist», «album» */
		string		cstr;

		/* Lowercase string representation */
		string		str;

		/* Title for column headers */
		string		title;

		/* Minimum and maximum length in column view */
		unsigned int	minlen;
		unsigned int	maxlen;
};

class Fieldtypes
{
	public:
		Fieldtypes();
		~Fieldtypes();

		/* All supported field types */
		vector<Field *>	fields;

		/* Locate a field type by MPD string */
		Field *		find_mpd(string &value);

		/* Locate a field type by name */
		Field *		find(string &value);
};


#endif /* _PMS_FIELD_H_ */

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

#ifndef _PMS_TOPBAR_H_
#define _PMS_TOPBAR_H_

#include <string>
#include <vector>
#include "field.h"
#include "song.h"
#include "color.h"
using namespace std;

enum
{
	CONDITION_NONE		= 0,
	CONDITION_PLAYING	= 1 << 0,
	CONDITION_SONG		= 1 << 1,
	CONDITION_CONNECTED	= 1 << 2

};

typedef struct
{
	unsigned int		t; /* conditions that has to evaluate true */
	unsigned int		f; /* conditions that has to evaluate false */
	unsigned int		ctl; /* which condition was used to open this? */
}
condition_t;


/*
 * Compiled topbar segment with color info
 */
class Topbarchunk
{
	public:
		Topbarchunk(string s, Color * c);

		string			str;
		Color *			color;
};


/*
 * One topbar segment. Contains compiled information about fields to print.
 */
class Topbarsegment
{
	public:
		Topbarsegment();
		~Topbarsegment();

		/* Compiled segment */
		vector<Topbarchunk *>	chunks;

		/* Compile segment into string vector */
		unsigned int		compile(Song * song);

		string			format;
		string			src;
		condition_t		condition;
		vector<Field *>		fields;
};

/*
 * This class represents one topbar line.
 */
class Topbarline
{
	public:
		~Topbarline();

		vector<Topbarsegment *>	segments;
};

/*
 * The entire topbar
 */
class Topbar
{
	public:
		/* Clear out the entire topbar */
		void			clear();

		/* Set topbar string */
		int			set(string format);

		/* Original format string (cached from option set) */
		string			cached_format;

		/* Split formatlines into left, center, right */
		vector<Topbarline *>	lines[3];
};

#endif /* _PMS_TOPBAR_H_ */

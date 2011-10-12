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

#ifndef _PMS_SEARCH_H_
#define _PMS_SEARCH_H_

#include "song.h"
#include <string>
#include <vector>
using namespace std;

typedef enum
{
	SEARCH_MODE_NONE,
	SEARCH_MODE_JUMP,
	SEARCH_MODE_FILTER
}

search_mode_t;

/* A single song reference */
class Searchresult
{
	public:
		unsigned int	pos;
		Song *		song;
};

/* A set of search results */
class Searchresultset
{
	public:
		vector<Searchresult *>	results;
		string			terms;
		long			mask;
			
		Searchresultset();
		~Searchresultset();
		void			clear();
		Searchresult *		add(unsigned int pos, Song * song);
		Searchresult *		operator[] (unsigned int spos);
		size_t			size() { return results.size(); };
};

/* Performs a case-insensitive match. */
inline bool cistrmatch(const string & haystack, const string & needle);

#endif /* _PMS_SEARCH_H_ */

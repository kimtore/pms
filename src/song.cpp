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

#include "song.h"
#include <string>
#include <vector>
#include <stdlib.h>

using namespace std;

Song::Song()
{
	f[FIELD_TIME]	= "-1";
	f[FIELD_POS]	= "-1";
	f[FIELD_ID]	= "-1";
}

void Song::init()
{
	size_t			s, e;
	string			src;
	string			tmp;
	vector<string *>	original;
	vector<string *>	rewritten;

	time = atoi(f[FIELD_TIME].c_str());
	pos = atoi(f[FIELD_POS].c_str());
	id = atoi(f[FIELD_ID].c_str());

	/* year from date */
	if (f[FIELD_DATE].size() >= 4)
		f[FIELD_YEAR] = f[FIELD_DATE].substr(0, 4);

	/* trackshort from track */
	if ((s = f[FIELD_TRACK].find_first_not_of('0')) != string::npos)
	{
		if ((e = f[FIELD_TRACK].find('/', s)) != string::npos)
			f[FIELD_TRACKSHORT] = f[FIELD_TRACK].substr(s, e - s);
		else
			f[FIELD_TRACKSHORT] = f[FIELD_TRACK].substr(s);
	}
	else
	{
		f[FIELD_TRACKSHORT] = f[FIELD_TRACK];
	}

	/* generate sort names if there are none available */
	if (f[FIELD_ARTISTSORT].size() == 0)
	{
		original.push_back(&f[FIELD_ARTIST]);
		rewritten.push_back(&f[FIELD_ARTISTSORT]);
	}
	if (f[FIELD_ALBUMARTISTSORT].size() == 0)
	{
		original.push_back(&f[FIELD_ALBUMARTIST]);
		rewritten.push_back(&f[FIELD_ALBUMARTISTSORT]);
	}

	tmp = "the ";
	e = tmp.size();
	for (s = 0; s < original.size(); s++)
	{
		/* Too small */
		if (original[s]->size() <= e)
		{
			*(rewritten[s]) = *(original[s]);
			continue;
		}

		src = original[s]->substr(0, e);
		/* Artist name consists of "the ...", place "The" at the end */
		if (cistrcmp(src, tmp) == true)
			*(rewritten[s]) = original[s]->substr(e) + ", " + original[s]->substr(0, e - 1);
		else
			*(rewritten[s]) = *(original[s]);
	}
}

string Song::dirname()
{
	string		ret = "";
	size_t		p;

	if (f[FIELD_FILE].size() == 0)
		return ret;

	p = f[FIELD_FILE].find_last_of("/\\");
	if (p == string::npos)
		return ret;

	return f[FIELD_FILE].substr(0, p);
}

bool cistrcmp(string &a, string &b)
{
	string::const_iterator ai, bi;

	ai = a.begin();
	bi = b.begin();

	while (ai != a.end() && bi != b.end())
	{
		if (::tolower(*ai) != ::tolower(*bi))
			return false;
		++ai;
		++bi;
	}

	return true;
}

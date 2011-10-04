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
#include <sstream>
#include <locale>

using namespace std;

Song::Song()
{
	f[FIELD_TIME]	= "-1";
	f[FIELD_POS]	= "-1";
	f[FIELD_ID]	= "-1";
	time = -1;
	pos = -1;
	id = -1;
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
	f[FIELD_TIME] = time_format(time);

	/* hash filename */
	fhash = songhash(f[FIELD_FILE]);

	/* replace title on no-title songs */
	if (f[FIELD_TITLE].size() == 0)
	{
		if (f[FIELD_NAME].size() > 0)
			f[FIELD_TITLE] = f[FIELD_NAME];
		else
			f[FIELD_TITLE] = f[FIELD_FILE];
	}

	/* show <Unknown ...> */
	if (f[FIELD_ARTIST].size() == 0)
		f[FIELD_ARTIST] = "<Unknown artist>";
	if (f[FIELD_ALBUM].size() == 0)
		f[FIELD_ALBUM] = "<Unknown album>";
	if (f[FIELD_YEAR].size() == 0)
		f[FIELD_YEAR] = "----";

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

	/* replace % with %% for fast printing */
	for (s = 0; s < FIELD_COLUMN_VALUES; ++s)
		escape_printf(f[s]);

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

string time_format(int seconds)
{
	static const int	day	= (60 * 60 * 24);
	static const int	hour	= (60 * 60);
	static const int	minute	= 60;

	int		i;
	string		s = "";

	/* No time */
	if (seconds < 0)
	{
		s = "--:--";
		return s;
	}

	/* days */
	if (seconds >= day)
	{
		i = seconds / day;
		s = tostring(i) + "d ";
		seconds %= day;
	}

	/* hours */
	if (seconds >= hour)
	{
		i = seconds / hour;
		s += zeropad(i, 1) + ":";
		seconds %= hour;
	}

	/* minutes */
	i = seconds / minute;
	s = s + zeropad(i, 2) + ":";
	seconds %= minute;

	/* seconds */
	s += zeropad(seconds, 2);

	return s;
}

string zeropad(int i, unsigned int target)
{
	string s;
	s = tostring(i);
	while(s.size() < target)
		s = '0' + s;
	return s;
}

string tostring(int number)
{
	ostringstream s;
	s << number;
	return s.str();
}

string tostring(unsigned int number)
{
	ostringstream s;
	s << number;
	return s.str();
}

void escape_printf(string &src)
{
	size_t pos = -2;

	while ((pos = src.find('%', pos + 2)) != string::npos)
		src.replace(pos, 1, "%%");
}

long songhash(string const &str)
{
	static locale loc;
	static const collate<char>& coll = use_facet<collate<char> >(loc);
	return coll.hash(str.data(), str.data() + str.size());
}

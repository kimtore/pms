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

using namespace std;

Song::Song()
{
	file			= "";
	artist			= "";
	albumartist		= "";
	artistsort		= "";
	albumartistsort		= "";
	title			= "";
	album			= "";
	track			= "";
	trackshort		= "";
	name			= "";
	date			= "";
	year			= "";

	genre			= "";
	composer		= "";
	performer		= "";
	disc			= "";
	comment			= "";

	length			= -1;
	pos			= -1;
	id			= -1;
}

void Song::init()
{
	unsigned int		i;
	string			src;
	string			tmp;
	vector<string *>	original;
	vector<string *>	rewritten;

	/* year from date */
	if (date.size() >= 4)
		year = date.substr(0, 4);

	/* trackshort from track */
	trackshort = track;
	while (trackshort[0] == '0')
		trackshort = trackshort.substr(1);
	if ((i = trackshort.find('/')) != string::npos)
		trackshort = trackshort.substr(0, i);

	/* sort names if none available */
	if (artistsort.size() == 0)
	{
		original.push_back(&artist);
		rewritten.push_back(&artistsort);
	}
	if (albumartistsort.size() == 0)
	{
		original.push_back(&albumartist);
		rewritten.push_back(&albumartistsort);
	}

	tmp = "the ";

	for (i = 0; i < original.size(); i++)
	{
		/* Too small */
		if (original[i]->size() > 4)
		{
			src = original[i]->substr(0, 4);
		}
		else
		{
			*(rewritten[i]) = *(original[i]);
			continue;
		}
	
		/* Artist name consists of "the ..." */
		if (cistrcmp(src, tmp) == true)
		{
			*(rewritten[i]) = original[i]->substr(4) + ", " + original[i]->substr(0, 3);
		}
		/* Revert to default */
		else
		{
			*(rewritten[i]) = *(original[i]);
		}
	}
}

string Song::dirname()
{
	string		ret = "";
	size_t		p;

	if (file.size() == 0)
		return ret;

	p = file.find_last_of("/\\");
	if (p == string::npos)
		return ret;

	return file.substr(0, p);
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

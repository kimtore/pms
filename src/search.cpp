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

#include "search.h"
#include "songlist.h"
#include "field.h"
#include "console.h"
#include <stdlib.h>
using namespace std;

extern Fieldtypes fieldtypes;

Searchresultset::Searchresultset()
{
	mask = 0;
}

Searchresultset::~Searchresultset()
{
	clear();
}

void Searchresultset::clear()
{
	vector<Searchresult *>::const_iterator i;
	for (i = results.begin(); i != results.end(); ++i)
		delete *i;
	results.clear();
}

Searchresult * Searchresultset::operator[] (unsigned int spos)
{
	if (spos >= results.size())
		return NULL;
	
	return results[spos];
}

Searchresult * Searchresultset::add(unsigned int pos, Song * song)
{
	Searchresult * s;
	s = new Searchresult;
	s->pos = pos;
	s->song = song;
	results.push_back(s);
	return s;
}

/*
 * Search functions
 */

size_t Songlist::find(long hash, size_t pos)
{
	size_t it;
	for (it = 0; it < songs.size(); ++it)
		if (songs[it]->fhash == hash)
			return it;

	return string::npos;
}

Searchresult * Songlist::search(search_mode_t mode)
{
	searchmode = mode;
	return NULL;
}

Searchresult * Songlist::search(search_mode_t mode, long mask, string terms)
{
	size_t i, count;
	vector<Field *>::const_iterator fit;

	count = songs.size();

	/* Check if we need to clear the current result set */
	//if ((searchresult.mask & mask) != searchresult.mask ||
		//terms.find(searchresult.terms) == string::npos)
	//{
		searchresult.clear();
	//}
	searchmode = mode;

	debug("Searching for `%s' in %d songs...", terms.c_str(), count);

	for (i = 0; i < count; ++i)
	{
		for (fit = fieldtypes.fields.begin(); fit != fieldtypes.fields.end(); ++fit)
		{
			if (!(mask & (1 << (*fit)->type)))
				continue;
			if (!cistrmatch(songs[i]->f[(*fit)->type], terms))
				continue;

			debug("Found match in field %d of %s", (*fit)->type, songs[i]->f[(*fit)->type].c_str()); 
			searchresult.add(i, songs[i]);
			break;
		}
	}

	debug("Search finished.", NULL);

	return NULL;
}

/*
 * Performs a case-insensitive match.
 */
inline bool cistrmatch(const string & haystack, const string & needle)
{
	bool matched = false;

	string::const_iterator	it_haystack;
	string::const_iterator	it_needle;

	for (it_haystack = haystack.begin(), it_needle = needle.begin(); it_haystack != haystack.end() && it_needle != needle.end(); it_haystack++)
	{
		/* exit if there aren't enough characters left to match the string */
		if (haystack.end() - it_haystack < needle.end() - it_needle)
			return false;

		/* check next character in needle with character in haystack */
		if (::toupper(*it_needle) == ::toupper(*it_haystack))
		{
			/* matched a letter -- look for next letter */
			matched = true;
			it_needle++;
		}
		else
		{
			/* didn't match a letter -- start from first letter of needle */
			matched = false;
			it_needle = needle.begin();
		}
	}

	if (it_needle != needle.end())
	{
		/* end of the haystack before getting to the end of the needle */
		return false;
	}

	return matched;
}

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
#include "config.h"
#include <stdlib.h>
using namespace std;

extern Fieldtypes fieldtypes;
extern Config config;

Searchresultset::Searchresultset()
{
	mask = 0;
}

/*
 * Search functions
 */

size_t Songlist::find(long hash, size_t pos)
{
	size_t it;
	for (it = pos; it < songs.size(); ++it)
		if (songs[it]->fhash == hash)
			return it;

	return string::npos;
}

size_t Songlist::sfind(long hash, size_t pos)
{
	size_t it;

	if (!searchresult)
		return find(hash, pos);

	for (it = pos; it < searchresult->songs.size(); ++it)
		if (searchresult->songs[it]->fhash == hash)
			return it;

	return string::npos;
}

size_t Songlist::spos(song_t pos)
{
	size_t it;

	if (!searchresult)
		return pos;

	for (it = 0; it < searchresult->songs.size(); ++it)
		if (searchresult->songs[it]->pos == pos)
			return it;

	return string::npos;
}

Song * Songlist::search(search_mode_t mode)
{
	searchmode = mode;
	if (searchmode == SEARCH_MODE_NONE && searchresult)
	{
		delete searchresult;
		searchresult = NULL;
		return NULL;
	}
	return NULL;
}

Song * Songlist::search(search_mode_t mode, long mask, string terms)
{
	Searchresultset * results;
	vector<Field *>::const_iterator fit;
	vector<Song *>::const_iterator sit;
	vector<Song *> * source;

	results = new Searchresultset;

	/* Check if we can use the current result set */
	if (searchresult)
		source = &(searchresult->songs);
	else
		source = &songs;

	for (sit = source->begin(); sit != source->end(); ++sit)
	{
		for (fit = fieldtypes.fields.begin(); fit != fieldtypes.fields.end(); ++fit)
		{
			if (!(mask & (1 << (*fit)->type)))
				continue;

			if (!strmatch((*sit)->f[(*fit)->type], terms, !config.search_case))
				continue;

			results->songs.push_back(*sit);
			break;
		}
	}

	if (searchresult)
		delete searchresult;
	searchresult = results;

	searchmode = mode;

	if (searchresult->size() > 0)
		return searchresult->songs[0];

	return NULL;
}

/*
 * Performs a case-insensitive match.
 */
inline bool strmatch(const string & haystack, const string & needle, bool ignorecase)
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
		if ((!ignorecase && *it_needle == *it_haystack) ||
			(ignorecase && ::toupper(*it_needle) == ::toupper(*it_haystack)))
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

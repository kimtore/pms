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

Searchresults::Searchresults()
{
	mask = 0;
	songlen = 0;
}

Searchresults * Searchresults::operator= (const Searchresults & source)
{
	songs = source.songs;
	terms = source.terms;
	mask = source.mask;
	mode = source.mode;
	return this;
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
	return search(mode, 0, "");
}

Searchresults * Songlist::search(vector<Song *> * source, long mask, string terms)
{
	Searchresults * results[2];
	unsigned int resultptr = 0;
	vector<Field *>::const_iterator fit;
	vector<Song *>::const_iterator sit;
	vector<string>::const_iterator tit;
	vector<string> * t;

	results[0] = new Searchresults;
	results[1] = new Searchresults;

	poscache = -1;
	lengthcache = 0;

	/* Split search string into words and search for them separately */
	if (config.split_search_terms)
	{
		t = str_split(terms, " ");
	}
	else
	{
		t = new vector<string>;
		t->push_back(terms);
	}

	/* Iterate through search words */
	for (tit = t->begin(); tit != t->end(); ++tit)
	{
		/* Iterate through our song source */
		for (sit = source->begin(); sit != source->end(); ++sit)
		{
			/* Check all search fields one by one */
			for (fit = fieldtypes.fields.begin(); fit != fieldtypes.fields.end(); ++fit)
			{
				if (!(mask & (1 << (*fit)->type)))
					continue;

				if (!strmatch((*sit)->f[(*fit)->type], *tit, !config.search_case))
					continue;

				results[resultptr]->songs.push_back(*sit);
				if ((*sit)->time != -1)
					results[resultptr]->songlen += (*sit)->time;
				break;
			}
		}

		/* Use this search result as next source */
		source = &(results[resultptr]->songs);
		++resultptr %= 2;
		results[resultptr]->songs.clear();
		results[resultptr]->songlen = 0;
	}

	delete t;
	delete results[resultptr];
	++resultptr %= 2;
	results[resultptr]->mask = mask;
	results[resultptr]->terms = terms;

	return results[resultptr];
}

Song * Songlist::search(search_mode_t mode, long mask, string terms)
{
	Searchresults * results = NULL;
	vector<Song *> * source;
	vector<Searchresults *>::iterator sit;
	long hashes[2] = { -1, -1 };

	/* Check if we can use the current result set */
	if (searchresult)
		source = &(searchresult->songs);
	else
		source = &songs;

	if (visual_start != -1 && at(visual_start))
		hashes[0] = source->at(visual_start)->fhash;
	if (visual_stop != -1 && at(visual_stop))
		hashes[1] = source->at(visual_stop)->fhash;

	switch(mode)
	{
		case SEARCH_MODE_NONE:
		default:
			/* Clear all search results */
			liveclear();
			if (searchresult)
				delete searchresult;

			searchresult = NULL;
			searchmode = SEARCH_MODE_NONE;

			return NULL;

		case SEARCH_MODE_FILTER:
			results = search(source, mask, terms);
			break;

		case SEARCH_MODE_LIVE:
			/* No search terms and no cached results, fall back to standard song list */
			if (terms.empty())
			{
				if (!livesource)
					return search(SEARCH_MODE_NONE, mask, terms);

				results = livesource;
				break;
			}

			/* Use current search result set as source for all live searching */
			if (!livesource && liveresults.empty() && searchresult)
				livesource = searchresult;

			/* Re-use any cached live search? */
			for (sit = liveresults.begin(); sit != liveresults.end(); ++sit)
			{
				if ((*sit)->terms == terms)
				{
					results = search(&(*sit)->songs, mask, terms);
					break;
				}
			}

			/* Can't re-use, search through current set */
			if (!results)
				results = search(source, mask, terms);

			/* Cache current search */
			if (!terms.empty() && sit == liveresults.end())
				liveresults.push_back(results);
			break;
	}

	if (searchresult && mode != SEARCH_MODE_LIVE)
		delete searchresult;

	searchresult = results;
	searchmode = mode;

	/* Visual selection out of range? */
	if (visual_start != -1)
	{
		visual_start = sfind(hashes[0]);
		visual_stop = sfind(hashes[1]);
		if (visual_start != -1 || visual_stop == -1)
			visual_stop = searchresult->size() - 1;
	}

	if (searchresult && searchresult->size() > 0)
		return searchresult->songs[0];

	return NULL;
}

void Songlist::liveclear()
{
	vector<Searchresults *>::iterator sit;

	for (sit = liveresults.begin(); sit != liveresults.end(); ++sit)
		if (*sit != searchresult)
			delete *sit;
	liveresults.clear();

	if (livesource && livesource != searchresult)
		delete livesource;

	livesource = NULL;
}

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

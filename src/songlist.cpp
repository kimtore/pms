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

#include "search.h"
#include "songlist.h"
#include <stdlib.h>

using namespace std;

Songlist::Songlist()
{
	readonly = false;
	playlist = false;
	version = -1;
	songlen = 0;
	visual_start = -1;
	visual_stop = -1;
	searchmode = SEARCH_MODE_NONE;
}

Songlist::~Songlist()
{
	search(SEARCH_MODE_NONE);
	clear();
}

Song * Songlist::at(unsigned int spos)
{
	if (searchmode == SEARCH_MODE_NONE)
	{
		if (spos >= songs.size())
			return NULL;
		return songs[spos];
	}
	else
	{
		if (spos >= searchresult->size())
			return NULL;
		return searchresult->songs[spos];
	}
}

void Songlist::add(Song * song)
{
	if (!song)
		return;

	/* Replace a song within this list */
	if (song->pos != -1 && song->pos < (long)songs.size())
	{
		if (songs[song->pos]->time != -1)
		{
			songlen -= songs[song->pos]->time;
			if ((int)poscache <= song->pos)
				lengthcache -= songs[song->pos]->time;
		}
		delete songs[song->pos];
		songs[song->pos] = song;
	}
	else
	{
		songs.push_back(song);
	}

	if (song->time != -1)
	{
		songlen += song->time;
		if ((int)poscache <= song->pos)
			lengthcache += song->time;
	}

	/* Reset search. TODO: inject into search results */
	search(SEARCH_MODE_NONE);
}

void Songlist::clear()
{
	vector<Song *>::iterator i;
	for (i = songs.begin(); i != songs.end(); ++i)
		delete *i;
	songs.clear();
	songlen = 0;
	search(SEARCH_MODE_NONE);
}

size_t Songlist::randpos()
{
	size_t r = 0;
	while (r < size())
		r += rand();
	r %= size() - 1;
	return r;
}

void Songlist::truncate(unsigned long length)
{
	while (songs.size() > length)
	{
		delete songs[songs.size()-1];
		songs.pop_back();
	}

	songs.reserve(length);
}

size_t Songlist::size()
{
	if (searchmode == SEARCH_MODE_NONE)
		return songs.size();
	else
		return searchresult->size();
}

unsigned long Songlist::length()
{
	if (searchmode == SEARCH_MODE_NONE)
		return songlen;
	else
		return searchresult->songlen;
}

unsigned long Songlist::length(size_t pos)
{
	vector<Song *>::const_iterator it;
	vector<Song *> * source;
	size_t cspos = pos;

	if (playlist)
		pos = spos(pos);

	if (poscache == pos)
		return lengthcache;

	if (searchmode == SEARCH_MODE_NONE)
		source = &songs;
	else
		source = &(searchresult->songs);

	poscache = pos;
	lengthcache = 0;

	/* Song position was not found in search results, so we need to count it all. */
	if (pos == string::npos)
	{
		pos = 0;

		/* Add current song time even though it's not part of the search results. */
		if (playlist && songs[cspos]->time != -1)
			lengthcache += songs[cspos]->time;
	}

	if (pos >= source->size())
		return 0;

	/* Calculate the sum of all time */
	it = source->begin() + pos;
	while (it < source->end())
	{
		if ((*it)->time != -1)
			lengthcache += (*it)->time;
		++it;
	}

	return lengthcache;
}

selection_t Songlist::get_selection()
{
	size_t start, stop;
	selection.clear();

	if (visual_start != -1)
	{
		visual_pos(&start, &stop);
		for ( ; start <= stop; ++start)
			selection.push_back(at(start));
	}

	return &selection;
}

void Songlist::visual_pos(size_t * start, size_t * stop)
{
	if (start != NULL)
		*start = (visual_start > visual_stop ? visual_stop : visual_start);
	if (stop != NULL)
		*stop = (visual_start > visual_stop ? visual_start : visual_stop);
}

void Songlist::clear_visual()
{
	visual_start = -1;
	visual_stop = -1;
}

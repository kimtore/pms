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
#include <stdlib.h>

using namespace std;

Songlist::Songlist()
{
	readonly = false;
	playlist = false;
	version = -1;
	searchmode = SEARCH_MODE_NONE;
}

Songlist::~Songlist()
{
	searchresult.clear();
	clear();
}

Song * Songlist::operator[] (unsigned int spos)
{
	//if (searchmode == SEARCH_MODE_NONE)
	//{
		if (spos >= songs.size())
			return NULL;
		return songs[spos];
	//}
	//else
	//{
		//if (spos >= searchresult.size())
			//return NULL;
		//return searchresult[spos]->song;
	//}
}

void Songlist::add(Song * song)
{
	if (!song)
		return;

	/* Replace a song within this list */
	if (song->pos != -1 && song->pos < (long)songs.size())
	{
		delete songs[song->pos];
		songs[song->pos] = song;
	}
	else
	{
		songs.push_back(song);
	}

	/* Reset search. TODO: inject into search results */
	search(SEARCH_MODE_NONE);
	searchresult.clear();
}

void Songlist::clear()
{
	vector<Song *>::iterator i;
	for (i = songs.begin(); i != songs.end(); ++i)
		delete *i;
	songs.clear();
	searchresult.clear();
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
	//if (searchmode == SEARCH_MODE_NONE)
		return songs.size();
	//else
		//return searchresult.size();
}

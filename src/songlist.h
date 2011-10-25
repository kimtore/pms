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

#ifndef _PMS_SONGLIST_H_
#define _PMS_SONGLIST_H_

#include "song.h"
#include "search.h"
#include <string>
#include <vector>

using namespace std;

class Songlist
{
	private:
		unsigned int		poscache;
		unsigned long		lengthcache;

	public:
		Songlist();
		~Songlist();

		Song *			at(unsigned int spos);
		Song *			operator[] (unsigned int spos) { return at(spos); };

		vector<Song *>		songs;
		unsigned long		songlen;
		string			title;
		Searchresults *		searchresult;
		Searchresults *		livesource;
		vector<Searchresults *> liveresults;
		search_mode_t		searchmode;

		/* Can we make local modifications? */
		bool			readonly;

		/* Is this the main playlist? */
		bool			playlist;

		/* Playlist version at MPD side */
		long long		version;

		/* Add or replace a song */
		void			add(Song * song);

		/* Remove all songs from the list */
		void			clear();

		/* Get a random song position within boundaries */
		size_t			randpos();

		/* Truncate the list and resize the vector */
		void			truncate(unsigned long length);

		/* Size */
		size_t			size();

		/* Length of all visible songs */
		unsigned long		length();
		unsigned long		length(size_t pos);

		/* Sort the list */
		void			sort(string sortstr);

		/*
		 * Search functions
		 */
		
		/* Find by hash value */
		size_t			find(long hash, size_t pos = 0);

		/* Same as find(), but looks only through search results */
		size_t			sfind(long hash, size_t pos = 0);

		/* Find relative song->pos in search mode, for use in playlist */
		size_t			spos(song_t pos);

		/* Search for songs using song fields. */
		Song *			search(search_mode_t mode);
		Song *			search(search_mode_t mode, long mask, string terms);

		/* Actual search worker, searches through source and returns a result set */
		Searchresults *		search(vector<Song *> * source, long mask, string terms);

		/* Clear live search cache */
		void			liveclear();

};

#endif /* _PMS_SONGLIST_H_ */

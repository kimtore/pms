/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS	<<Practical Music Search>>
 * Copyright (C) 2006-2015  Kim Tore Jensen
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
 *	songlist.h
 *		Playlist class, holds info about a lot of songs
 */

#ifndef _PMS_SONGLIST_H_
#define _PMS_SONGLIST_H_

#include <algorithm>
#include <string>
#include <vector>
#include <mpd/client.h>

#include "list.h"
#include "song.h"
#include "types.h"
#include "field.h"
#include "filter.h"
#include "column.h"

using namespace std;

#define MATCH_FAILED -1

/**
 * Return `x` cast to a Songlist *, or NULL if the list is not a Songlist.
 */
#define SONGLIST(x) dynamic_cast<Songlist *>(x)
#define LISTITEMSONG(x) dynamic_cast<ListItemSong *>(x)

enum
{
	MATCH_ID		= 1 << 0,
	MATCH_POS		= 1 << 1,
	MATCH_FILE		= 1 << 2,
	MATCH_ARTIST		= 1 << 3,
	MATCH_ARTISTSORT	= 1 << 4,
	MATCH_ALBUMARTIST	= 1 << 5,
	MATCH_ALBUMARTISTSORT	= 1 << 6,
	MATCH_TITLE		= 1 << 7,
	MATCH_ALBUM		= 1 << 8,
	MATCH_TRACKSHORT	= 1 << 9,
	MATCH_DATE		= 1 << 10,
	MATCH_ORIGINALDATE	= 1 << 11,
	MATCH_TIME		= 1 << 12,
	MATCH_GENRE		= 1 << 13,
	MATCH_COMPOSER		= 1 << 14,
	MATCH_PERFORMER		= 1 << 15,
	MATCH_DISCSHORT		= 1 << 16,
	MATCH_COMMENT		= 1 << 17,
	MATCH_YEAR		= 1 << 18,
	MATCH_ORIGINALYEAR	= 1 << 19,

	MATCH_ALL		= (1 << 20) - 1,

	MATCH_NOT		= 1 << 20,
	MATCH_EXACT		= 1 << 21,
	MATCH_REVERSE		= 1 << 22,
	MATCH_ALMOST		= 1 << 23,
	MATCH_LT		= 1 << 24,
	MATCH_LTE		= 1 << 25,
	MATCH_GT		= 1 << 26,
	MATCH_GTE		= 1 << 27
};

struct Selection
{
	song_t				size;
	song_t				length;
};

class ListItemSong : public ListItem
{
public:
	Song *		song;

			ListItemSong(List * l, Song * s);
			~ListItemSong();

	bool		match(string term, long flags);
};

class Songlist : public List
{
private:
	song_t					position;
	song_t					qlen;
	song_t					qpos;
	song_t					qnum;
	song_t					qsize;

	vector<pms_column *>			columns;

protected:
	/*
	 * Appends a songlist to the list.
	 *
	 * Returns the zero-indexed position of the first song added.
	 */
	song_t			add_local(Songlist *);

	/*
	 * Remove song in position N from the list.
	 */
	void			remove_local(uint32_t position);

public:
				Songlist();
				~Songlist();

	string			filename;

	bool			draw();

	/**
	 * Return the first occurrence of a song.
	 */
	ListItemSong *		find(Song *);

	/**
	 * Return the song at the specified position.
	 *
	 * Will raise an assertion error when the position is invalid.
	 */
	Song *			song(uint32_t position);

	/**
	 * Crop the list to a specific song.
	 *
	 * Returns true on success, false on failure.
	 */
	bool			crop_to_song(Song * song);

	unsigned int		length;
	Selection		selection_params;
	void			set(Songlist *);
	void			truncate_local(unsigned int);

	bool			sort(string);

	/**
	 * After a sort procedure, the song->pos are inaccurate. This function
	 * numbers them sequentially.
	 */
	void			renumber_pos();

	unsigned int		cursor();
	Song *			cursorsong();
	void			set_column_size();

	bool			swap(uint32_t, uint32_t);

	/* Pick songs based on playmode */
	Song *			next_song_in_direction(Song * s, uint8_t direction, song_t * id);
	Song *			nextsong(song_t * = NULL);
	Song *			prevsong(song_t * = NULL);
	Song *			randsong(song_t * = NULL);

	/* Next-of and prev-of functions */
	song_t			nextof(string);
	song_t			prevof(string);
	song_t			findentry(Item, bool);

	/* Filter functions */
	/*
	Filter *		filter_add(string param, long fields);
	void			filter_remove(Filter *);
	void			filter_clear();
	void			filter_scan();
	bool			filter_match(Song *);
	Filter *		lastfilter();
	unsigned int		filtercount() { return filters.size(); };
	*/

	/*
	 * Adds or replaces a song to the list, depending on the value of
	 * song->pos. The latter value is asserted to less than or equal to the
	 * list size.
	 *
	 * FIXME: this function should be protected!
	 *
	 * Returns the zero-indexed position of the added song.
	 */
	song_t			add_local(Song * s);

	/**
	 * Remove a song asynchronously, i.e. send a message to MPD and request
	 * to remove it. The base class will only call remove_local(), and thus
	 * always return true.
	 */
	bool			remove(ListItem * i);

	/**
	 * Add a song length to the list's cached length.
	 */
	void			add_song_length(int32_t t);

	/**
	 * Subtract a song length from the list's cached length.
	 */
	void			subtract_song_length(int32_t t);

	/*
	unsigned int		realsize() { return songs.size(); };
	unsigned int		size() { return filtersongs.size(); };
	unsigned int		end() { return size() - 1; };
	*/
	unsigned int		qlength();
	unsigned int		qnumber() { return qnum; };
};


bool		lcstrcmp(const string &, const string &);
bool		icstrsort(const string &, const string &);

/* Sorts */
bool		sort_compare_file(ListItem * a_, ListItem * b_);
bool		sort_compare_artist(ListItem * a_, ListItem * b_);
bool		sort_compare_albumartist(ListItem * a_, ListItem * b_);
bool		sort_compare_albumartistsort(ListItem * a_, ListItem * b_);
bool		sort_compare_artistsort(ListItem * a_, ListItem * b_);
bool		sort_compare_title(ListItem * a_, ListItem * b_);
bool		sort_compare_album(ListItem * a_, ListItem * b_);
bool		sort_compare_track(ListItem * a_, ListItem * b_);
bool		sort_compare_length(ListItem * a_, ListItem * b_);
bool		sort_compare_name(ListItem * a_, ListItem * b_);
bool		sort_compare_date(ListItem * a_, ListItem * b_);
bool		sort_compare_originaldate(ListItem * a_, ListItem * b_);
bool		sort_compare_year(ListItem * a_, ListItem * b_);
bool		sort_compare_originalyear(ListItem * a_, ListItem * b_);
bool		sort_compare_genre(ListItem * a_, ListItem * b_);
bool		sort_compare_composer(ListItem * a_, ListItem * b_);
bool		sort_compare_performer(ListItem * a_, ListItem * b_);
bool		sort_compare_disc(ListItem * a_, ListItem * b_);
bool		sort_compare_comment(ListItem * a_, ListItem * b_);


#endif

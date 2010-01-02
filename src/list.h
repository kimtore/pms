/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2010  Kim Tore Jensen
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
 * 	list.h
 * 		Playlist class, holds info about a lot of songs
 */

#ifndef _PMS_LIST_H_
#define _PMS_LIST_H_


#include <algorithm>
#include <string>
#include <vector>
#include "libmpdclient.h"
#include "song.h"
#include "types.h"
#include "field.h"

using namespace std;

#define MATCH_FAILED -1

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
	MATCH_TIME		= 1 << 11,
	MATCH_GENRE		= 1 << 12,
	MATCH_COMPOSER		= 1 << 13,
	MATCH_PERFORMER		= 1 << 14,
	MATCH_DISC		= 1 << 15,
	MATCH_COMMENT		= 1 << 16,
	MATCH_YEAR		= 1 << 17,

	MATCH_ALL		= (1 << 18) - 1,

	MATCH_NOT		= 1 << 18,
	MATCH_EXACT		= 1 << 19,
	MATCH_REVERSE		= 1 << 20
};

struct Selection
{
	song_t				size;
	song_t				length;
};

class Songlist
{
private:
	song_t					position;
	song_t					qlen;
	song_t					qpos;
	song_t					qnum;
	song_t					qsize;

	Song *					lastget;
	vector<Song *>::iterator		seliter;
	vector<Song *>::reverse_iterator	rseliter;

public:
				Songlist();
				~Songlist();

	bool			ignorecase;
	bool			wrap;
	string			filename;
	
	vector<Song *>		songs;
	unsigned int		length;
	void			clear();
	Selection		selection;
	void			set(Songlist *);
	void			truncate(unsigned int);

	bool			sort(string);

	vector<song_t> *	matchall(string, long);
	song_t			match(string, unsigned int, unsigned int, long);
#ifdef HAVE_LIBBOOST_REGEX
	bool			regexmatch(string *, string *);
#endif
	bool			exactmatch(string *, string *);
	bool			inmatch(string *, string *);
	bool			perform_match(string *, string *, int);

	void			movecursor(song_t);
	int			setcursor(song_t);
	bool			gotocurrent();
	unsigned int		cursor();
	Song *			cursorsong();
	int			locatesong(Song *);

	bool			selectsong(Song *, bool);
	Song *			getnextselected();
	Song *			getprevselected();
	Song *			popnextselected();
	void			resetgets();
	bool			swap(int, int);

	/* Pick songs based on playmode */
	Song *			nextsong(song_t * = NULL);
	Song *			prevsong(song_t * = NULL);
	Song *			randsong(song_t * = NULL);

	/* Next-of and prev-of functions */
	song_t			nextof(string);
	song_t			prevof(string);
	song_t			findentry(Item, bool);

	song_t			add(Song *);
	song_t			add(Songlist *);
	int			remove(Song *);
	int			remove(int);
	bool			move(unsigned int, unsigned int);
	unsigned int		size() { return songs.size(); };
	unsigned int		end() { return songs.size() - 1; };
	unsigned int		qlength();
	unsigned int		qnumber() { return qnum; };
};

bool		lcstrcmp(string &, string &);
bool		icstrsort(string &, string &);

/* Sorts */
bool		sort_compare_file(Song *, Song *);
bool		sort_compare_artist(Song *, Song *);
bool		sort_compare_albumartist(Song *, Song *);
bool		sort_compare_albumartistsort(Song *, Song *);
bool		sort_compare_artistsort(Song *, Song *);
bool		sort_compare_title(Song *, Song *);
bool		sort_compare_album(Song *, Song *);
bool		sort_compare_track(Song *, Song *);
bool		sort_compare_length(Song *, Song *);
bool		sort_compare_name(Song *, Song *);
bool		sort_compare_date(Song *, Song *);
bool		sort_compare_year(Song *, Song *);
bool		sort_compare_genre(Song *, Song *);
bool		sort_compare_composer(Song *, Song *);
bool		sort_compare_performer(Song *, Song *);
bool		sort_compare_disc(Song *, Song *);
bool		sort_compare_comment(Song *, Song *);


#endif

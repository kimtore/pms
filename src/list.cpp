/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2009  Kim Tore Jensen
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
 * 	list.cpp
 * 		Playlist class, holds info about a lot of songs
 */


#include "../config.h"
#ifdef HAVE_LIBBOOST_REGEX
	#include <boost/regex.hpp>
#endif
#include "conn.h"
#include "list.h"
#include "song.h"
#include "config.h"
#include "libmpdclient.h"
#include "pms.h"

extern Pms *			pms;


/*
 * Playlist class
 */
Songlist::Songlist()
{
	lastget = NULL;
	seliter = songs.begin();
	rseliter = songs.rbegin();
	position = 0;
	wrap = false;
	length = 0;
	qlen = 0;
	qpos = 0;
	qnum = 0;
	qsize = 0;
	filename = "";
	selection.size = 0;
	selection.length = 0;
}

Songlist::~Songlist()
{
	this->clear();
	position = 0;
}

/*
 * Returns the next song in line, starting from current song
 */
Song *			Songlist::nextsong(song_t * id)
{
	song_t		i = MATCH_FAILED;
	Song *		s;

	s = pms->cursong();

	/* No current song returns first song in line */
	if (!s)
	{
		if (size() == 0)
			return NULL;

		return songs[0];
	}

	/* Find the current song in this list */
	if (s->pos != MPD_SONG_NO_NUM)
		i = match(Pms::tostring(pms->cursong()->pos), 0, end(), MATCH_POS);

	if (i == MATCH_FAILED)
	{
		i = match(pms->cursong()->file, 0, end(), MATCH_EXACT | MATCH_FILE);
		if (i == MATCH_FAILED && size() == 0)
			return NULL;
	}

	/* Wrap around */
	if (++i >= static_cast<song_t>(size()))
	{
		if (pms->options->repeatmode == REPEAT_LIST)
			i = 0;
		else
			return NULL;
	}

	if (id != NULL)
		*id = i;

	return songs[i];
}

/*
 * Returns the previous song
 */
Song *			Songlist::prevsong(song_t * id)
{
	song_t		i = MATCH_FAILED;
	Song *		s;

	s = pms->cursong();

	/* No current song returns last song in line */
	if (!s)
	{
		if (size() == 0)
			return NULL;

		return songs[end()];
	}

	/* Find the current song in this list */
	if (s->pos != MPD_SONG_NO_NUM)
		i = match(Pms::tostring(pms->cursong()->pos), 0, end(), MATCH_POS);

	if (i == MATCH_FAILED)
	{
		i = match(pms->cursong()->file, 0, end(), MATCH_EXACT | MATCH_FILE);
		if (i == MATCH_FAILED && size() == 0)
			return NULL;
	}

	/* Wrap around */
	if (--i < 0)
	{
		if (pms->options->repeatmode == REPEAT_LIST)
			i = end();
		else
			return NULL;
	}

	if (id != NULL)
		*id = i;

	return songs[i];
}

/*
 * Return a random song
 */
Song *			Songlist::randsong(song_t * id)
{
	song_t		i = 0;
	unsigned long	processed = 0;

	if (size() == 0)
		return NULL;

	while (processed < size())
	{
		i += rand();
		processed += RAND_MAX;
	}

	i %= size();

	if (songs[i] == pms->cursong())
	{
		--i;
		if (i < 0) i = end();
	}

	if (id != NULL)
		*id = i;

	return songs[i];
}


/*
 * Next-of returns next unique field
 */
song_t		Songlist::nextof(string s)
{
	Item		i;

	if (s.size() == 0)
		return MPD_SONG_NO_NUM;

	i = pms->formatter->field_to_item(s);

	return findentry(i, false);
}

/*
 * Prev-of returns previous and last unique field
 */
song_t		Songlist::prevof(string s)
{
	Item		i;

	if (s.size() == 0)
		return MPD_SONG_NO_NUM;

	i = pms->formatter->field_to_item(s);

	return findentry(i, true);
}

/*
 * Finds next or previous entry of any type.
 */
song_t		Songlist::findentry(Item field, bool reverse)
{	
	Song *		song;
	song_t		i = MATCH_FAILED;
	long		mode = 0;
//	string		where;
	string		cmp[2];
	bool		tmp;

	if (field == LITERALPERCENT || field == EINVALID) return i;

	/* Set up our variables */
	mode = pms->formatter->item_to_match(field);
	if (reverse) mode = mode | MATCH_REVERSE;
//	where = (reverse ? _("previous") : _("next"));

	/* Sanity checks on environment */
	song = cursorsong();
	if (!song) return i;
	i = cursor();

	/* Return our search string */
	cmp[0] = pms->formatter->format(song, field, true);

	/* Perform a match */
	i = match(cmp[0], i, i - 1, mode | MATCH_NOT | MATCH_EXACT);
	if (i == MATCH_FAILED)
	{
		debug("gotonextentry() fails with mode = %d\n", mode);
		return i;
	}

	song = songs[i];

	/* Reverse match must match first entry, not last */
	if (reverse)
	{
		cmp[0] = pms->formatter->format(song, field, true);
		i = match(cmp[0], i, i - 1, mode | MATCH_NOT | MATCH_EXACT);
		if (i != MATCH_FAILED)
		{
			if (++i == size())
				i = 0;
		}
	}

	return i;
}

/*
 * Copies a list from another list
 */
void		Songlist::set(Songlist * list)
{
	unsigned int	i;
	Song *		s;

	if (list == NULL)	return;

	this->clear();

	for (i = 0; i < list->size(); i++)
	{
		s = new Song(list->songs[i]);
		s->id = MPD_SONG_NO_ID;
		s->pos = MPD_SONG_NO_NUM;
		add(s);
	}
}

/*
 * Sets the maximum list size
 */
void		Songlist::truncate(unsigned int maxsize)
{
	unsigned int	i;

//	debug("Truncating %p, maxsize=%d: ...\n", this, maxsize);
	
	if (maxsize == 0)
	{
		this->clear();
		return;
	}

	for (i = end(); i >= maxsize; i--)
	{
//		debug("...removing %d: %s\n", i, songs[i]->file.c_str());
		remove(static_cast<int>(i));
	}
}

/*
 * Appends an entire list. Returns the id of the first added song.
 */
song_t		Songlist::add(Songlist * list)
{
	song_t			first = MPD_SONG_NO_ID;
	song_t			result;
	unsigned int		i;

	if (!list) return first;

	for (i = 0; i < list->size(); i++)
	{
		result = add(new Song(list->songs[i]));
		if (first == MPD_SONG_NO_ID && result != MPD_SONG_NO_ID)
			first = result;
	}

	return first;
}

/*
 * Adds a song to the list, either at end or in the middle
 */
song_t		Songlist::add(Song * song)
{
	vector<Song *>::iterator	i;

	if (song == NULL)
		return MPD_SONG_NO_ID;
	
	if (song->pos == MPD_SONG_NO_NUM || song->pos == static_cast<song_t>(songs.size()))
	{
		songs.push_back(song);
	}
	else
	{
		i = songs.begin() + song->pos;
		if (songs[song->pos]->pos == song->pos)
		{
			i = songs.erase(songs.begin() + song->pos);
		}
		songs.insert(i, song);
	}

	if (song->time != MPD_SONG_NO_TIME)
	{
		length += song->time;
	}

	seliter = songs.begin();
	rseliter = songs.rbegin();

	return static_cast<song_t>(songs.size() - 1);
}

/*
 * Removes a song from the list
 */
int		Songlist::remove(Song *song)
{
	if (!song)	return false;

	selectsong(song, false);

	if (song->pos == MPD_SONG_NO_NUM)
	{
		return remove(match(song->file, 0, songs.size() - 1, MATCH_FILE));
	}
	else	return remove(song->pos);
}

/*
 * Remove song by index
 */
int		Songlist::remove(int songpos)
{
	vector<Song *>::iterator	it;

	if (songpos < 0 || (unsigned int)songpos >= songs.size() || songpos == MATCH_FAILED)
	{
		return false;
	}

	if (songs[songpos]->time != MPD_SONG_NO_TIME)
	{
		length -= songs[songpos]->time;
	}

	delete songs[songpos];

	it = songs.begin() + songpos;
	it = songs.erase(it);

	while (it != songs.end())
	{
		--(*it)->pos;
		++it;
	}

	seliter = songs.begin();
	rseliter = songs.rbegin();

	return true;
}

/*
 * Swap two song positions
 */
bool			Songlist::swap(int a, int b)
{
	unsigned int	i, j;
	int		tpos;
	Song *		tmp;

	if (a < 0 || a >= songs.size() || b < 0 || b >= songs.size())
		return false;

	i = static_cast<unsigned int>(a);
	j = static_cast<unsigned int>(b);

	tpos = songs[i]->pos;

	tmp = songs[i];
	songs[i] = songs[j];
	songs[j] = tmp;

	songs[j]->pos = songs[i]->pos;
	songs[i]->pos = tpos;

	return true;
}

/*
 * Move a song inside the list to position dest
 */
bool			Songlist::move(unsigned int from, unsigned int dest)
{
	int		songpos, direction, dst;

	if (dest >= songs.size() || from >= songs.size())
		return false;

	songpos = static_cast<int>(from);
	dst = static_cast<int>(dest);

	/* Set direction */
	if (dst == songpos)
		return false;
	else if (dst > songpos)
		direction = 1;
	else
		direction = -1;

	/* Swap every element on its way */
	while (songpos != dst)
	{
		if (!this->swap(songpos, (songpos + direction)))
			return false;

		songpos += direction;
	}

	/* Clear queue length */
	{
		qlen = 0;
		qpos = 0;
		qnum = 0;
		qsize = 0;
	}

	return true;
}

/*
 * Truncate list
 */
void Songlist::clear()
{
	unsigned int		i;

	for (i = 0; i < songs.size(); i++)
	{
		delete songs[i];
	}

	songs.clear();

	length = 0;
	qlen = 0;
	qpos = 0;
	qnum = 0;
	qsize = 0;
}

/*
 * Set absolute cursor position
 */
int		Songlist::setcursor(song_t pos)
{
	if (pos < 0)
	{
		beep();
		pos = 0;
	}
	else if (pos >= songs.size())
	{
		beep();
		pos = songs.size() - 1;
	}

	position = pos;

	return position;
}

/*
 * Goto current playing song
 */
bool		Songlist::gotocurrent()
{
	song_t		i = MATCH_FAILED;

	if (!pms->cursong()) return false;

	if (pms->cursong()->pos != MPD_SONG_NO_NUM)
		i = match(Pms::tostring(pms->cursong()->pos), 0, end(), MATCH_POS | MATCH_EXACT);
	if (i == MATCH_FAILED)
		i = match(pms->cursong()->file, 0, end(), MATCH_FILE | MATCH_EXACT);
	if (i == MATCH_FAILED) return false;

	setcursor(i);
	return true;
}

/*
 * Returns position of song
 */
int		Songlist::locatesong(Song * song)
{
	unsigned int		i;

	for (i = 0; i < songs.size(); i++)
	{
		if (songs[i] == song)
			return (int)i;
	}

	return MATCH_FAILED;
}

/*
 * Set selection state of a song
 */
bool		Songlist::selectsong(Song * song, bool state)
{
	if (!song) return false;

	if (song->selected != state)
	{
		if (state == true)
		{
			if (song->time != MPD_SONG_NO_TIME)
				selection.length += song->time;
			selection.size++;
		}
		else if (state == false)
		{
			if (song->time != MPD_SONG_NO_TIME)
				selection.length -= song->time;
			selection.size--;
		}
		song->selected = state;
	}

	return true;
}

/*
 * Returs a consecutive list of selected songs, and unselects them
 */
Song *		Songlist::getnextselected()
{
	if (lastget == NULL)
	{
		seliter = songs.begin();
	}

	while (seliter != songs.end())
	{
		if (!(*seliter)) break; // out of bounds
		if ((*seliter)->selected)
		{
			lastget = *seliter;
			++seliter;
			return lastget;
		}
		++seliter;
	}

	/* No selection, return cursor */
	if (lastget == NULL)
	{
		if (lastget == cursorsong())
			lastget = NULL;
		else
			lastget = cursorsong();

		return lastget;
	}

	lastget = NULL;
	return NULL;
}

/*
 * Returs a consecutive list of selected songs, and unselects them
 */
Song *		Songlist::getprevselected()
{
	if (lastget == NULL)
	{
		rseliter = songs.rbegin();
	}

	while (rseliter != songs.rend())
	{
		if (!(*rseliter)) break; // out of bounds
		if ((*rseliter)->selected)
		{
			lastget = *rseliter;
			++rseliter;
			return lastget;
		}
		++rseliter;
	}

	/* No selection, return cursor */
	if (lastget == NULL)
	{
		if (lastget == cursorsong())
			lastget = NULL;
		else
			lastget = cursorsong();

		return lastget;
	}

	lastget = NULL;
	return NULL;
}

/*
 * Returs a consecutive list of selected songs, and unselects them
 */
Song *		Songlist::popnextselected()
{
	Song *		song;

	song = getnextselected();
	if (song)
	{
		selectsong(song, false);
	}
	return song;
}

/*
 * Reset iterators
 */
void		Songlist::resetgets()
{
	lastget = NULL;
	seliter = songs.begin();
	rseliter = songs.rbegin();
}

/*
 * Return song struct at cursor position, or NULL
 */
Song *		Songlist::cursorsong()
{
	if (songs.size() == 0) return NULL;
	return (songs[cursor()]);
}

/*
 * Return cursor position
 */
unsigned int		Songlist::cursor()
{
	if (position < 0)			position = 0;
	else if (position >= songs.size())	position = songs.size() - 1;

	return position;
}

/*
 * Return length of songs after playing position.
 */
unsigned int		Songlist::qlength()
{
	unsigned int		i, songpos;
	
	/* Find current playing song */
	if (!pms->cursong() || pms->cursong()->id == MPD_SONG_NO_ID || pms->cursong()->pos == MPD_SONG_NO_NUM)
	{
		qnum = songs.size();
		qpos = 0;
		qlen = length;
		return qlen;
	}

	if ((int)qpos == pms->cursong()->id && qsize == songs.size())
		return qlen;

	qpos = pms->cursong()->id;
	songpos = pms->cursong()->pos;

	/* Calculate from start */
	qlen = 0;
	qnum = 0;
	qsize = songs.size();
	for (i = songpos + 1; i < songs.size(); i++)
	{
		if (songs[i]->time != MPD_SONG_NO_TIME)
			qlen += songs[i]->time;
		++qnum;
	}
	return qlen;
}

/*
 * Set relative cursor position
 */
void		Songlist::movecursor(song_t offset)
{
	if (wrap == true)
	{
		if (songs.size() == 0)
		{
			beep();
			position = 0;
			return;
		}
		position += offset;
		while(position < 0)
		{
			position += songs.size();
		}
		position %= songs.size();
	}
	else
	{
		offset = position + offset;

		if (offset < 0)
		{
			beep();
			position = 0;
		}
		else if ((unsigned int)offset >= songs.size())
		{
			beep();
			position = songs.size() - 1;
		}
		else
			position = offset;
	}
}

/*
 * Find all matches, returns list of song IDs
 */
vector<song_t> *	Songlist::matchall(string src, long mode)
{
	vector<song_t> *	ret	= new vector<song_t>;
	song_t			m	= 0;

	while(m != songs.size() - 1)	
	{
		m = match(src, m, songs.size() - 1, mode);
		if (m != MATCH_FAILED)
			ret->push_back(m);
		else
			break;
	}

	if (ret->size() == 0)
	{
		delete ret;
		return NULL;
	}

	return ret;
}

/*
 * Find next match
 */
song_t			Songlist::match(string src, unsigned int from, unsigned int to, long mode)
{
	bool				matched;
	int				i;
	unsigned int			j;
	Song *				song;
	vector<string>			sources;

	if (songs.size() == 0)
		return MATCH_FAILED;

	if (from < 0)		from = 0;
	if (to >= songs.size())	to = songs.size() - 1;

	if (mode & MATCH_REVERSE)
	{
		i = from;
		from = to;
		to = i;
	}

	i = from;

	while (true)
	{
		if (i < 0)
			i = songs.size() - 1;
		else if (i >= songs.size())
			i = 0;

		if (songs[i] == NULL)
		{
			i += (mode & MATCH_REVERSE ? -1 : 1);
			continue;
		}
		song = songs[i];
		sources.clear();

		/* try the sources in order of likeliness. ID etc last since if we're 
		 * searching for them we likely won't be searching any of the other 
		 * fields. */
		if (mode & MATCH_TITLE)			sources.push_back(song->title);
		if (mode & MATCH_ARTIST)		sources.push_back(song->artist);
		if (mode & MATCH_ALBUMARTIST)		sources.push_back(song->albumartist);
		if (mode & MATCH_COMPOSER)		sources.push_back(song->composer);
		if (mode & MATCH_PERFORMER)		sources.push_back(song->performer);
		if (mode & MATCH_ALBUM)			sources.push_back(song->album);
		if (mode & MATCH_GENRE)			sources.push_back(song->genre);
		if (mode & MATCH_DATE)			sources.push_back(song->date);
		if (mode & MATCH_COMMENT)		sources.push_back(song->comment);
		if (mode & MATCH_TRACKSHORT)		sources.push_back(song->trackshort);
		if (mode & MATCH_DISC)			sources.push_back(song->disc);
		if (mode & MATCH_FILE)			sources.push_back(song->file);
		if (mode & MATCH_ARTISTSORT)		sources.push_back(song->artistsort);
		if (mode & MATCH_ALBUMARTISTSORT)	sources.push_back(song->albumartistsort);
		if (mode & MATCH_YEAR)			sources.push_back(song->year);
		if (mode & MATCH_ID)			sources.push_back(Pms::tostring(song->id));
		if (mode & MATCH_POS)			sources.push_back(Pms::tostring(song->pos));

		for (j = 0; j < sources.size(); j++)
		{
			if (mode & MATCH_EXACT)
				matched = exactmatch(&(sources[j]), &src);
#ifdef HAVE_LIBBOOST_REGEX
			else if (pms->options->regexsearch)
				matched = regexmatch(&(sources[j]), &src);
#endif
			else
				matched = inmatch(&(sources[j]), &src);

			if (matched)
				if (!(mode & MATCH_NOT))
					return i;
				else
					continue;
			else
				if (mode & MATCH_NOT)
					return i;
		}

		if (i == to)
			break;

		i += (mode & MATCH_REVERSE ? -1 : 1);
	}

	return MATCH_FAILED;
}

/*
 * Perform an exact match
 */
bool		Songlist::exactmatch(string * source, string * pattern)
{
	return perform_match(source, pattern, 1);
}

/*
 * Perform an in-string match
 */
bool		Songlist::inmatch(string * source, string * pattern)
{
	return perform_match(source, pattern, 0);
}

/*
 * Performs a case-insensitive match.
 * type:
 *  0 = match inside string also
 *  1 = match entire string only
 */
bool		Songlist::perform_match(string * haystack, string * needle, int type)
{
	bool			matched = (type == 1);

	string::iterator	it_haystack;
	string::iterator	it_needle;

	for (it_haystack = haystack->begin(), it_needle = needle->begin(); it_haystack != haystack->end() && it_needle != needle->end(); it_haystack++)
	{
		/* exit if there aren't enough characters left to match the string */
		if (haystack->end() - it_haystack < needle->end() - it_needle)
			return false;

		/* check next character in needle with character in haystack */
		if (::toupper(*it_needle) == ::toupper(*it_haystack))
		{
			/* matched a letter -- look for next letter */
			matched = true;
			it_needle++;
		}
		else if (type == 1)
		{
			/* didn't match a letter but need exact match */
			return false;
		}
		else
		{
			/* didn't match a letter -- start from first letter of needle */
			matched = false;
			it_needle = needle->begin();
		}
	}

	if (it_needle != needle->end())
	{
		/* end of the haystack before getting to the end of the needle */
		return false;
	}
	if (type == 1 && it_needle == needle->end() && it_haystack != haystack->end())
	{
		/* need exact and got to the end of the needle but not the end of the 
		 * haystack */
		return false;
	}

	return matched;
}

/*
 * Performs a case-insensitive regular expression match
 */
#ifdef HAVE_LIBBOOST_REGEX
bool		Songlist::regexmatch(string * source, string * pattern)
{
	bool			matched;
	boost::regex		reg;

	try
	{
		reg.assign(*pattern, boost::regex_constants::icase);
		matched = boost::regex_search(source->begin(), source->end(), reg);
	}
	catch (boost::regex_error & err)
	{
		return false;
	}
	return matched;
}
#endif


/*
 * Sort list by sort string.
 * sorts is a space-separated list of sort arguments.
 */
bool		Songlist::sort(string sorts)
{
	vector<Song *>::iterator	start;
	vector<Song *>::iterator	stop;
	vector<Song *>			temp;
	vector<string> *		v;
	unsigned int			i;
	int				ft;
	bool (*func) (Song *, Song *);

	if (sorts.size() == 0)
		return false;
	else if (sorts == "default")
		sorts = "track disc album date albumartistsort";

	temp = songs;

	start = temp.begin();
	stop = temp.end();

	v = Pms::splitstr(sorts, " ");

	for (i = 0; i < v->size(); i++)
	{
		ft = pms->fieldtypes->lookup((*v)[i]);
		if (ft == -1)
			continue;

		func = pms->fieldtypes->sortfunc[(unsigned int)ft];
		if (func == NULL) continue;

		if (i == 0)
			std::sort(start, stop, func);
		else
			std::stable_sort(start, stop, func);
	}

	if (i == v->size())
	{
		songs = temp;
		delete v;
		return true;
	}

	delete v;
	return false;
}

/*
 * Performs a case insensitive string comparison.
 */
bool	lcstrcmp(string & a, string & b)
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

/*
 * Performs a sort comparison based on the 'ignorecase' option.
 */
bool	icstrsort(string & a, string & b)
{
	string::const_iterator ai, bi;

	if (!pms->options->ignorecase)
		return a < b;

	ai = a.begin();
	bi = b.begin();

	while (ai != a.end() && bi != b.end())
	{
		if (::tolower(*ai) < ::tolower(*bi))
			return true;
		else if (::tolower(*ai) > ::tolower(*bi))
			return false;
		++ai;
		++bi;
	}

	return false;
}

/*
 * Sort functions
 */
bool	sort_compare_file(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->file, b->file);
}

bool	sort_compare_artist(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->artist, b->artist);
}

bool	sort_compare_albumartist(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->albumartist, b->albumartist);
}

bool	sort_compare_albumartistsort(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->albumartistsort, b->albumartistsort);
}

bool	sort_compare_artistsort(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->artistsort, b->artistsort);
}

bool	sort_compare_title(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->title, b->title);
}

bool	sort_compare_album(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->album, b->album);
}

bool	sort_compare_track(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return atoi(a->track.c_str()) < atoi(b->track.c_str());
}

bool	sort_compare_length(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return (a->time < b->time);
}

bool	sort_compare_name(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->name, b->name);
}

bool	sort_compare_date(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return a->date < b->date;
}

bool	sort_compare_year(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return a->year < b->year;
}

bool	sort_compare_genre(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->genre, b->genre);
}

bool	sort_compare_composer(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->composer, b->composer);
}

bool	sort_compare_performer(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->performer, b->performer);
}

bool	sort_compare_disc(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return atoi(a->disc.c_str()) < atoi(b->disc.c_str());
}

bool	sort_compare_comment(Song *a, Song *b)
{
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->comment, b->comment);
}


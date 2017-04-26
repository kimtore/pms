/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS  <<Practical Music Search>>
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
 * 	list.cpp
 * 		Playlist class, holds info about a lot of songs
 */


#include "../config.h"
#ifdef HAVE_REGEX
	#include <regex>
#endif
#include "conn.h"
#include "songlist.h"
#include "song.h"
#include "config.h"
#include "queue.h"
#include "pms.h"

extern Pms *			pms;


ListItemSong::ListItemSong(List * l, Song * s) :
ListItem(l)
{
	assert(l);
	assert(s);

	list = l;
	song = s;
}

ListItemSong::~ListItemSong()
{
	delete song;
}

bool
ListItemSong::match(string term, long flags)
{
	return song->match(term, flags);
}

/*
 * Playlist class
 */
Songlist::Songlist()
{
	position = 0;
	length = 0;
	qlen = 0;
	qpos = 0;
	qnum = 0;
	qsize = 0;
	filename = "";
	selection_params.size = 0;
	selection_params.length = 0;
}

Songlist::~Songlist()
{
}

/*
 * Return a pointer to the Nth song in the list.
 */
Song *
Songlist::song(uint32_t position)
{
	assert(position >= 0);
	assert(position < size());

	return LISTITEMSONG(items[position])->song;
}

/*
 * Returns the next song in line, starting from current song
 *
 * FIXME: should probably not be a part of the Songlist class
 */
Song *
Songlist::next_song_in_direction(Song * s, uint8_t direction, song_t * id)
{
	ListItem *	it;
	song_t		i = MATCH_FAILED;

	assert(direction == 1 || direction == -1);

	/* No current song returns first song in line */
	if (!s) {
		if (!size()) {
			return NULL;
		}
		return song(0);
	}

	/* Find the current song in this list */
	if (s->pos != MPD_SONG_NO_NUM && QUEUE(this)) {
		it = match(Pms::tostring(s->pos), 0, size() - 1, MATCH_POS);
	}

	/* Fallback to file path */
	if (!it) {
		it = match(s->file, 0, size() - 1, MATCH_EXACT | MATCH_FILE);
		if (!it && !size()) {
			return NULL;
		}
	}

	/* Wrap around */
	/* FIXME: not our responsibility */
	i = LISTITEMSONG(it)->song->pos + direction;
	if (i < 0 || i >= size()) {
		if (!pms->comm->status()->repeat) {
			return NULL;
		} else if (i < 0) {
			i = size() - 1;
		} else {
			i = 0;
		}
	}

	/* Assign song id to parameter */
	if (id != NULL) {
		*id = i;
	}

	return song(i);
}

Song *
Songlist::nextsong(song_t * id)
{
	return next_song_in_direction(pms->cursong(), 1, id);
}

Song *
Songlist::prevsong(song_t * id)
{
	return next_song_in_direction(pms->cursong(), -1, id);
}

/*
 * Return a random song
 */
Song *			Songlist::randsong(song_t * id)
{
	Song *		s;
	song_t		i = 0;
	unsigned long	processed = 0;

	if (!size()) {
		return NULL;
	}

	while (processed < size()) {
		i += rand();
		processed += RAND_MAX;
	}

	i %= size();

	s = song(i);
	if (s == pms->cursong()) {
		return next_song_in_direction(s, -1, id);
	}

	if (id != NULL) {
		*id = i;
	}

	return s;
}


/*
 * Next-of returns next unique field
 */
song_t
Songlist::nextof(string s)
{
	Item		i;

	if (!s.size()) {
		return MPD_SONG_NO_NUM;
	}

	i = pms->formatter->field_to_item(s);

	return findentry(i, false);
}

/*
 * Prev-of returns previous and last unique field
 */
song_t
Songlist::prevof(string s)
{
	Item		i;

	if (!s.size()) {
		return MPD_SONG_NO_NUM;
	}

	i = pms->formatter->field_to_item(s);

	return findentry(i, true);
}

/*
 * Finds next or previous entry of any type.
 */
song_t		Songlist::findentry(Item field, bool reverse)
{
	ListItem *	it;
	Song *		s;
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
	s = cursorsong();
	assert(s);
	i = cursor_position;

	/* Return our search string */
	cmp[0] = pms->formatter->format(s, field, true);

	/* Perform a match */
	it = match_wrap_around(cmp[0], (unsigned int) i, mode | MATCH_NOT | MATCH_EXACT);
	if (!it) {
		pms->log(MSG_DEBUG, 0, "gotonextentry() fails with mode = %d\n", mode);
		return MPD_SONG_NO_NUM;
	}

	s = LISTITEMSONG(it)->song;

	/* Reverse match must match first entry, not last */
	if (reverse)
	{
		cmp[0] = pms->formatter->format(s, field, true);
		it = match_wrap_around(cmp[0], i, mode | MATCH_NOT | MATCH_EXACT);
		if (it && it == last()) {
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
		s = new Song(list->song(i));
		s->id = MPD_SONG_NO_ID;
		s->pos = MPD_SONG_NO_NUM;
		add_local(s);
	}
}

/*
 * Sets the maximum list size
 */
void		Songlist::truncate_local(unsigned int maxsize)
{
	unsigned int	i;

	if (maxsize == 0)
	{
		this->clear();
		return;
	}

	for (i = size() - 1; i >= maxsize; i--)
	{
		remove_local(static_cast<int>(i));
	}
}

song_t		Songlist::add_local(Songlist * list)
{
	song_t			first = MPD_SONG_NO_ID;
	song_t			result;
	unsigned int		i;

	if (!list) return first;

	for (i = 0; i < list->size(); i++)
	{
		result = add_local(new Song(list->song(i)));
		if (first == MPD_SONG_NO_ID && result != MPD_SONG_NO_ID)
			first = result;
	}

	return first;
}

song_t
Songlist::add_local(Song * s)
{
	Song * existing_song;

	assert(s != NULL);
	assert(s->pos <= size());

	//pms->log(MSG_DEBUG, 0, "Add to queue: id=%d pos=%d uri=%s\n", s->id, s->pos, s->file.c_str());

	/* Append song to end of list */
	if (s->pos == MPD_SONG_NO_NUM || s->pos == size()) {
		items.push_back(new ListItemSong(this, s));
		s->pos = size() - 1;

	/* Insert song into arbitrary position */
	} else {
		existing_song = song(s->pos);
		assert(existing_song);
		assert(existing_song->pos == s->pos);

		subtract_song_length(existing_song->time);
		delete item(s->pos);
		items[s->pos] = new ListItemSong(this, s);
	}

	add_song_length(s->time);

	set_selection_cache_valid(false);

	return s->pos;
}

void
Songlist::add_song_length(int32_t t)
{
	if (t != MPD_SONG_NO_TIME) {
		length += t;
	}
}

void
Songlist::subtract_song_length(int32_t t)
{
	if (t != MPD_SONG_NO_TIME) {
		length -= t;
	}
}

ListItemSong *
Songlist::find(Song * s)
{
	ListItem * it = NULL;

	assert(s);

	if (s->id != MPD_SONG_NO_NUM && QUEUE(this)) {
		it = match(Pms::tostring(pms->cursong()->id), 0, size() - 1, MATCH_ID | MATCH_EXACT);
	} else if (s->pos != MPD_SONG_NO_NUM && QUEUE(this)) {
		it = match(Pms::tostring(pms->cursong()->pos), 0, size() - 1, MATCH_POS | MATCH_EXACT);
	}

	if (!it) {
		it = match(s->file, 0, size() - 1, MATCH_FILE | MATCH_EXACT);
	}

	return LISTITEMSONG(it);
}

void
Songlist::remove_local(uint32_t position)
{
	vector<ListItem *>::iterator iter;
	ListItemSong * list_item;
	Song * s;
	song_t song_length;

	s = song(position);
	assert(s);

	song_length = s->time;

	List::remove_local(position);

	subtract_song_length(song_length);

	iter = items.begin() + position;

	/* Decrease song position of all following song instances */
	while (iter != items.end()) {
		list_item = LISTITEMSONG(*iter);
		assert(list_item);
		assert(list_item->song);
		--list_item->song->pos;
		++iter;
	}
}

bool
Songlist::remove(ListItem * i)
{
	ListItemSong * list_item = LISTITEMSONG(i);

	assert(list_item);
	assert(list_item->song);
	assert(list_item->song->pos != MPD_SONG_NO_NUM);
	assert(list_item->song->pos < size());

	remove_local(list_item->song->pos);

	return true;
}

/*
 * Set selection state of a song
 *
 * FIXME
bool		Songlist::selectsong(Song * song, bool state)
{
	assert(false);
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
 */

/*
 * Return song at cursor position, or NULL if the songlist is empty
 */
Song *
Songlist::cursorsong()
{
	if (!size()) {
		return NULL;
	}

	return song(cursor_position);
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
		qnum = size();
		qpos = 0;
		qlen = length;
		return qlen;
	}

	if ((int)qpos == pms->cursong()->id && qsize == size()) {
		return qlen;
	}

	qpos = pms->cursong()->id;
	songpos = pms->cursong()->pos;

	/* Calculate from start */
	qlen = 0;
	qnum = 0;
	qsize = size();
	for (i = songpos + 1; i < size(); i++)
	{
		if (song(i)->time != MPD_SONG_NO_TIME)
			qlen += song(i)->time;
		++qnum;
	}
	return qlen;
}

/*
 * Calculates table grid size and positions.
 */
void
Songlist::set_column_size()
{
	int			index;
	unsigned int		i;
	unsigned int		ui, j;
	unsigned int		winlen;
	Song			*s;
	string			tmp;
	vector<string> *	v;
	bool			allfixed;

	/* If there are any old columns, remove them */
	/* FIXME: why delete and re-add them? */
	for (ui = 0; ui < columns.size(); ui++)
	{
		delete columns[ui];
	}
	columns.clear();

	v = Pms::splitstr(pms->options->columns, " ");

	for (i = 0; i < v->size(); i++)
	{
		index = pms->fieldtypes->lookup((*v)[i]);
		if (index == -1)
			continue;
		j = (unsigned int)index;
		columns.push_back(new pms_column(	pms->fieldtypes->header[j],
							pms->fieldtypes->type[j],
							pms->fieldtypes->minlen[j]));
	}

	delete v;

	if (columns.size() == 0)
	{
		return;
	}

	/* Maximum length of fields */
	assert(bbox);
	winlen = bbox->width();

	/* Find minimum length needed to display all content */
	for (i = 0; i < size(); i++)
	{
		s = song(i);

		for (j = 0; j < columns.size(); j++)
		{
			ui = 0;

			switch(columns[j]->type)
			{
			case FIELD_NUM:
				ui = Pms::tostring(s->pos).size();
				break;
			case FIELD_FILE:
				ui = s->file.size();
				break;
			case FIELD_ARTIST:
				ui = s->artist.size();
				break;
			case FIELD_ALBUMARTIST:
				ui = s->albumartist.size();
				break;
			case FIELD_ALBUMARTISTSORT:
				ui = s->albumartistsort.size();
				break;
			case FIELD_ARTISTSORT:
				ui = s->artistsort.size();
				break;
			case FIELD_TITLE:
				if (s->title.size())
					ui = s->title.size();
				else if (s->name.size())
					ui = s->name.size();
				else if (s->file.size())
					ui = s->file.size();
				break;
			case FIELD_ALBUM:
				ui = s->album.size();
				break;
			case FIELD_TRACK:
				ui = s->track.size();
				break;
			case FIELD_TRACKSHORT:
				ui = s->trackshort.size();
				break;
			case FIELD_TIME:
				ui = Pms::timeformat(s->time).size();
				break;
			case FIELD_DATE:
				ui = s->date.size();
				break;
			case FIELD_YEAR:
				ui = s->year.size();
				break;
			case FIELD_NAME:
				ui = s->name.size();
				break;
			case FIELD_GENRE:
				ui = s->genre.size();
				break;
			case FIELD_COMPOSER:
				ui = s->composer.size();
				break;
			case FIELD_PERFORMER:
				ui = s->performer.size();
				break;
			case FIELD_DISC:
				ui = s->disc.size();
				break;
			case FIELD_DISCSHORT:
				ui = s->discshort.size();
				break;
			case FIELD_COMMENT:
				ui = s->comment.size();
				break;
			default:
				continue;
			}

			columns[j]->addmedian(ui);
		}
	}

	/* Calculate total length of existing fields */
	j = 0;
	for (ui = 0; ui < columns.size(); ui++)
	{
		j += columns[ui]->len();
	}

	/* Do we have only fixed width fields? */
	allfixed = true;
	for (ui = 0; ui < columns.size(); ui++)
	{
		if (columns[ui]->minlen == 0)
		{
			allfixed = false;
			break;
		}
	}

	/* Resize fields until they fit into the window */
	while (j != winlen)
	{
		for (ui = 0; ui < columns.size(); ui++)
		{
			if (j > winlen && columns[ui]->len() > columns[ui]->minlen)
			{
				--columns[ui]->abslen;
				--j;
			}
			else if (allfixed || j < winlen && columns[ui]->minlen == 0)
			{
				++columns[ui]->abslen;
				++j;
			}
			if (j == winlen)
				break;
		}
	}
}

bool
Songlist::draw()
{
	unsigned int		pair;
	unsigned int		counter = 0;
	unsigned int		i, j, winlen;
	unsigned int		min;
	unsigned int		max;
	int			ii;
	ListItem *		list_item;
	Song *			s;
	string			t;
	color *			hilight;
	color *			c;

	/* Clear window first */
	bbox->clear(NULL);

	/* Zero songs: zero draw */
	if (!size()) {
		return true;
	}

	/* Define range of songs to draw */
	min = top_position();
	max = bottom_position();

	/* Traverse song list and draw lines */
	for (i = min; i <= max; i++)
	{
		++counter;
		hilight = NULL;

		list_item = item(i);
		s = song(i);
		assert(s);

		if (i == cursor_position)
		{
			hilight = pms->options->colors->cursor;
		}
		else if (list_item->selected())
		{
			hilight = pms->options->colors->selection;
		}
		else if (pms->cursong()) {
                        if ((QUEUE(this) && pms->cursong()->id == s->id) || (!QUEUE(this) && s->file == pms->cursong()->file)) {
				hilight = pms->options->colors->current;
			}
		}

		winlen = 0;
		for (j = 0; j < columns.size(); j++)
		{
			pair = 0;

                        /* Draw highlight line */
			if (hilight) wattron(bbox->window, hilight->pair());
			mvwhline(bbox->window, counter, winlen, ' ', columns[j]->len() + 1);
			if (hilight) wattroff(bbox->window, hilight->pair());

			c = pms->formatter->getcolor(columns[j]->type, &(pms->options->colors->fields));
			if (c)
			{
				t = pms->formatter->format(s, columns[j]->type);
				colprint(bbox, counter, (j == 0 ? winlen : winlen + 1),
					(hilight ? hilight : c),
					"%s", t.c_str());

			}

			winlen += columns[j]->len();
		}

		hilight = pms->options->colors->standard;
	}

	/* Draw captions and column borders */
	j = 0;
	for (i = 0; i < columns.size(); i++)
	{
		colprint(bbox, 0, (i == 0 ? j : j + 1),
			pms->options->colors->headers,
			"%s", columns[i]->title.c_str());
		if (i > 0 && pms->options->columnborders)
		{
			wattron(bbox->window, pms->options->colors->border->pair());
			mvwvline(bbox->window, 0, j, ACS_VLINE, bbox->height());
			wattroff(bbox->window, pms->options->colors->border->pair());
		}
		j += columns[i]->len();
	}

	return true;
}

/*
 * Sort list by sort string.
 * sorts is a space-separated list of sort arguments.
 */
bool		Songlist::sort(string sorts)
{
	vector<ListItem *>::iterator	start;
	vector<ListItem *>::iterator	stop;
	vector<ListItem *>		temp;
	vector<string> *		v;
	unsigned int			i;
	int				ft;
	bool (*func) (ListItem *, ListItem *);

	if (sorts.size() == 0)
		return false;

	v = Pms::splitstr(sorts, " ");

	/* Sort the real song list */
	start = items.begin();
	stop = items.end();

	for (i = 0; i < v->size(); i++)
	{
		ft = pms->fieldtypes->lookup((*v)[i]);
		if (ft == -1)
			continue;

		func = pms->fieldtypes->sortfunc[(unsigned int)ft];
		if (func == NULL) continue;

		if (i == 0) {
			std::sort(start, stop, func);
		} else {
			std::stable_sort(start, stop, func);
		}
	}

	renumber_pos();

	delete v;
	return true;
}

void
Songlist::renumber_pos()
{
	ListItemSong * list_item;
	uint32_t i;

	for (i = 0; i < size(); i++) {
		list_item = LISTITEMSONG(items[i]);
		list_item->song->pos = i;
	}
}

/*
 * Performs a case insensitive string comparison.
 */
bool	lcstrcmp(const string & a, const string & b)
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
bool	icstrsort(const string & a, const string & b)
{
	string::const_iterator a_iter;
	string::const_iterator b_iter;
	unsigned char a_lower;
	unsigned char b_lower;

	if (!pms->options->ignorecase) {
		return a < b;
	}

	a_iter = a.begin();
	b_iter = b.begin();

	while (a_iter != a.end() && b_iter != b.end()) {
		a_lower = tolower(*a_iter);
		b_lower = tolower(*b_iter);

		if (a_lower < b_lower) {
			return true;
		} else if (a_lower > b_lower) {
			return false;
		}

		++a_iter;
		++b_iter;
	}

	if (a_iter == a.end() && b_iter != b.end()) {
		return true;
	}

	return false;
}

/*
 * Sort functions
 */
bool	sort_compare_file(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->file, b->file);
}

bool	sort_compare_artist(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->artist, b->artist);
}

bool	sort_compare_albumartist(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->albumartist, b->albumartist);
}

bool	sort_compare_albumartistsort(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->albumartistsort, b->albumartistsort);
}

bool	sort_compare_artistsort(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->artistsort, b->artistsort);
}

bool	sort_compare_title(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->title, b->title);
}

bool	sort_compare_album(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->album, b->album);
}

bool	sort_compare_track(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return atoi(a->track.c_str()) < atoi(b->track.c_str());
}

bool	sort_compare_length(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return (a->time < b->time);
}

bool	sort_compare_name(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->name, b->name);
}

bool	sort_compare_date(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return a->date < b->date;
}

bool	sort_compare_year(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return a->year < b->year;
}

bool	sort_compare_genre(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->genre, b->genre);
}

bool	sort_compare_composer(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->composer, b->composer);
}

bool	sort_compare_performer(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->performer, b->performer);
}

bool	sort_compare_disc(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return atoi(a->disc.c_str()) < atoi(b->disc.c_str());
}

bool	sort_compare_comment(ListItem * a_, ListItem * b_)
{
	Song * a = LISTITEMSONG(a_)->song;
	Song * b = LISTITEMSONG(b_)->song;
	if (a == NULL && b == NULL)			return true;
	else if (a == NULL && b != NULL)		return true;
	else if (a != NULL && b == NULL)		return false;
	else 						return icstrsort(a->comment, b->comment);
}

bool
Songlist::crop_to_song(Song * song)
{
	vector<ListItem *>::reverse_iterator iter;
	ListItem * it;

	it = match(song->file, 0, size() - 1, MATCH_FILE | MATCH_EXACT);

	if (!it) {
		return false;
	}

	iter = items.rbegin();
	while (iter != items.rend()) {
		if (*iter == it) {
			continue;
		}
		assert(LISTITEMSONG(*iter));
		if (!remove(LISTITEMSONG(*iter))) {
			/* FIXME: error reporting */
			return false;
		}
		++iter;
	}

	return true;
}

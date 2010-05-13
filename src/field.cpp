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
 *
 * field.cpp - format a song using field variables
 *
 */

#include "field.h"
#include "pms.h"

using namespace std;


extern Pms *		pms;



/*
 * Evaluates a format string and returns readable text
 */
string			Formatter::format(Song * song, string fmt, unsigned int & printlen, colortable_fields * f, bool clean)
{
	Item				val;
	int				i = 0, next = 0, last = 0;
	string				strval;
	string				condformat;
	string				tmp;
	unsigned int			tmpint;

	condformat = evalconditionals(fmt);
	printlen = 0;

	while (true)
	{
		val = nextitem(condformat, &i, &next);
		if (val == EINVALID)
		{
			//already added the last item. add rest of string
			tmp = condformat.substr(next);
			strval += tmp;
			printlen += tmp.size();
			break;
		}

		//add all text between last item and this one
		strval += condformat.substr(last, i - last);
		printlen += i - last;

		tmp = format(song, val, tmpint, f, clean);
		if (!tmp.empty())
		{
			strval += tmp;
			printlen += tmpint;
		}
		last = next;
		i = next;
	}

	return strval;
}

/*
 * Returns a single field for use in e.g. searches, no printlen is returned
 */
string			Formatter::format(Song * song, Item i, bool clean)
{
	unsigned int		in;

	return format(song, i, in, NULL, clean);
}

/*
 * Formats a string with multiple items into a vector
 * e.g  'artist title track' => FIELD_ARTIST, FIELD_ALBUM, FIELD_TRACK
 */
vector<Item> *		Formatter::multiformat_item(string s)
{
	vector<string> *	split;
	vector<Item> *		v = new vector<Item>;
	unsigned int		i;

	split = Pms::splitstr(s);
	if (split == NULL)
	{
		delete v;
		return NULL;
	}

	for (i = 0; i < split->size(); i++)
	{
		v->push_back(field_to_item(split->at(i)));
	}

	delete split;
	return v;
}

/*
 * Interprets a string value into a field type
 */
Item			Formatter::field_to_item(string f)
{
	/* All field types */

	if (f == "")
		return LITERALPERCENT;
	else if (f == "num")
		return FIELD_NUM;
	else if (f == "file")
		return FIELD_FILE;
	else if (f == "artist")
		return FIELD_ARTIST;
	else if (f == "artistsort")
		return FIELD_ARTISTSORT;
	else if (f == "albumartist")
		return FIELD_ALBUMARTIST;
	else if (f == "albumartistsort")
		return FIELD_ALBUMARTISTSORT;
	else if (f == "title")
		return FIELD_TITLE;
	else if (f == "album")
		return FIELD_ALBUM;
	else if (f == "date")
		return FIELD_DATE;
	else if (f == "year")
		return FIELD_YEAR;
	else if (f == "track")
		return FIELD_TRACK;
	else if (f == "trackshort")
		return FIELD_TRACKSHORT;
	else if (f == "time")
		return FIELD_TIME;
	else if (f == "name")
		return FIELD_NAME;
	else if (f == "genre")
		return FIELD_GENRE;
	else if (f == "composer")
		return FIELD_COMPOSER;
	else if (f == "performer")
		return FIELD_PERFORMER;
	else if (f == "disc")
		return FIELD_COMMENT;

	/* Conditionals */
	else if (f == "ifcursong")
		return COND_IFCURSONG;
	else if (f == "ifplaying")
		return COND_IFPLAYING;
	else if (f == "ifpaused")
		return COND_IFPAUSED;
	else if (f == "ifstopped")
		return COND_IFSTOPPED;
	else if (f == "else")
		return COND_ELSE;
	else if (f == "endif")
		return COND_ENDIF;

	/* The rest */

	else if (f == "repeat")
		return REPEAT;
	else if (f == "random")
		return RANDOM;
	else if (f == "manual")
		return MANUALPROGRESSION;
	else if (f == "mute")
		return MUTE;
	else if (f == "repeatshort")
		return REPEATSHORT;
	else if (f == "randomshort")
		return RANDOMSHORT;
	else if (f == "manualshort")
		return MANUALPROGRESSIONSHORT;
	else if (f == "muteshort")
		return MUTESHORT;

	else if (f == "librarysize")
		return LIBRARYSIZE;
	else if (f == "listsize")
		return LISTSIZE;
	else if (f == "queuesize")
		return QUEUESIZE;
	else if (f == "livequeuesize")
		return LIVEQUEUESIZE;

	else if (f == "bitrate")
		return BITRATE;
	else if (f == "samplerate")
		return SAMPLERATE;
	else if (f == "bits")
		return BITS;
	else if (f == "channels")
		return CHANNELS;

	else if (f == "time_elapsed")
		return TIME_ELAPSED;
	else if (f == "time_remaining")
		return TIME_REMAINING;
	else if (f == "progressbar")
		return PROGRESSBAR;
	else if (f == "progresspercentage")
		return PROGRESSPERCENTAGE;
	else if (f == "playstate")
		return PLAYSTATE;
	else if (f == "volume")
		return VOLUME;
	else
		return EINVALID;
}

/*
 * Return the next Item starting at the given string index of the given format
 * The index of the leading % ends up in i
 * The index of the character after the item is stored in next
 */
Item			Formatter::nextitem(string fmt, int *i, int *next)
{
	int				l;

	if (*i >= fmt.size())
		return EINVALID;

	//keep going until we get to a %
	while (fmt[*i] != '%')
	{
		if (++(*i) >= fmt.size())
			return EINVALID;
	}

	//find number of chars before next %
	for (l = 0; fmt[*i + l + 1] != '%'; l++);

	//set next to index past keyword and the following %
	*next = *i + l + 2;

	return field_to_item(fmt.substr(*i + 1, l));
}

/*
 * Recursive function to return a format string which has had conditionals 
 * applied
 */
string			Formatter::evalconditionals(string fmt)
{
	int				if_start = 0, if_next;
	int				else_start = -1, else_next = -1;
	int				endif_start = 0, endif_next;
	int				depth = 0;
	Item				item1, item2;
	bool				satisfied;
	string				before, after;

	while (true)
	{
		item1 = nextitem(fmt, &if_start, &if_next);
		if (item1 == EINVALID)
			break;

		switch(item1)
		{
			case COND_IFCURSONG:
			case COND_IFPLAYING:
			case COND_IFPAUSED:
			case COND_IFSTOPPED:
				//find matching endif
				endif_start = if_next;
				while (true)
				{
					item2 = nextitem(fmt, &endif_start, &endif_next);
					if (item2 == EINVALID)
						break;

					switch(item2)
					{
						case COND_IFCURSONG:
						case COND_IFPLAYING:
						case COND_IFPAUSED:
						case COND_IFSTOPPED:
							depth++;
							break;

						case COND_ELSE:
							if (depth == 0)
							{
								else_start = endif_start;
								else_next = endif_next;
							}
							break;

						case COND_ENDIF:
							if (depth > 0)
							{
								depth--;
								break;
							}
							switch(item1)
							{
								case COND_IFCURSONG:
									satisfied = pms->cursong() ? true : false;
									break;

								case COND_IFPLAYING:
									satisfied = pms->comm->status()->state == MPD_STATUS_STATE_PLAY;
									break;

								case COND_IFPAUSED:
									satisfied = pms->comm->status()->state == MPD_STATUS_STATE_PAUSE;
									break;

								case COND_IFSTOPPED:
									satisfied = pms->comm->status()->state == MPD_STATUS_STATE_STOP;
									break;

								default:
									//shouldn't be here
									pms->log(MSG_DEBUG, 0, "error: didn't know how to evaluate condition. assuming true\n");
									satisfied = true;
									break;
							}
							before = fmt.substr(0, if_start);
							after = evalconditionals(fmt.substr(endif_next));
							if (satisfied)
								return before + evalconditionals(fmt.substr(if_next, (else_start == -1 ? endif_start : else_start) - if_next)) + after;
							return before + (else_start == -1 ? "" : evalconditionals(fmt.substr(else_next, endif_start - else_next))) + after;

						default:
							//ignore other items
							break;
					}
					endif_start = endif_next;
				}
				pms->log(MSG_DEBUG, 0, "error: no matching endif\n");
				return "";

			case COND_ELSE:
				pms->log(MSG_DEBUG, 0, "error: found else before if\n");
				return "";

			case COND_ENDIF:
				pms->log(MSG_DEBUG, 0, "error: found endif before if\n");
				return "";

			default:
				break;
		}
		if_start = if_next;
	}

	//no conditionals
	return fmt;
}

/*
 * Evaluate a topbar keyword
 */
string			Formatter::format(Song * song, Item keyword, unsigned int & printlen, colortable_fields * f, bool clean)
{
	char			s[30];
	int			tmpint;
	string			retstr;
	string			tmp;
	unsigned int		pint, progress;
	Songlist *		list;
	color *			c;
	long			playmode;
	long			repeatmode;

	c = getcolor(keyword, f);
	playmode = pms->options->get_long("playmode");
	repeatmode = pms->options->get_long("repeat");

	retstr = "";
	switch(keyword)
	{
		/* Field types */

		case FIELD_NUM:
			if (!song) return retstr;
			retstr = Pms::tostring(song->pos);
			break;

		case FIELD_FILE:
			if (!song) return retstr;
			retstr = song->file;
			break;

		case FIELD_ARTIST:
			if (!song) return retstr;
			retstr = song->artist;
			if (clean) break;

			if (!retstr.size())
				retstr = song->albumartist;
			if (!retstr.size())
				retstr = "<Unknown artist>";
			break;

		case FIELD_ALBUMARTIST:
			if (!song) return retstr;
			retstr = song->albumartist;
			if (clean) break;

			if (!retstr.size())
				retstr = song->artist;
			if (!retstr.size())
				retstr = "<Unknown artist>";
			break;

		case FIELD_ARTISTSORT:
			if (!song) return retstr;
			retstr = song->artistsort;
			if (clean) break;

			if (!retstr.size())
				retstr = song->artist;
			if (!retstr.size())
				retstr = "<Unknown artist>";
			break;

		case FIELD_ALBUMARTISTSORT:
			if (!song) return retstr;
			retstr = song->albumartistsort;
			if (clean) break;

			if (!retstr.size())
				retstr = song->albumartist;
			if (!retstr.size())
				retstr = "<Unknown artist>";
			break;

		case FIELD_TITLE:
			if (!song) return retstr;
			retstr = song->title;
			if (clean) break;

			if (!retstr.size())
				retstr = song->name;
			if (!retstr.size())
				retstr = song->file;
			break;

		case FIELD_ALBUM:
			if (!song) return retstr;
			retstr = song->album;
			if (clean) break;

			if (!retstr.size())
				retstr = "<Unknown album>";
			break;

		case FIELD_DATE:
			if (!song) return retstr;
			retstr = song->date;
			if (clean) break;

			if (!retstr.size())
				retstr = "----";
			break;

		case FIELD_YEAR:
			if (!song) return retstr;
			retstr = song->year;
			if (clean) break;

			if (!retstr.size())
				retstr = "----";
			break;

		case FIELD_TRACK:
			if (!song) return retstr;
			retstr = song->track;
			break;

		case FIELD_TRACKSHORT:
			if (!song) return retstr;
			retstr = song->trackshort;
			break;

		case FIELD_TIME:
			if (!song) return retstr;
			if (clean)
				retstr = Pms::tostring(song->time);
			else
				retstr = Pms::timeformat(song->time);
			break;

		case FIELD_NAME:
			if (!song) return retstr;
			retstr = song->name;
			break;

		case FIELD_GENRE:
			if (!song) return retstr;
			retstr = song->genre;
			break;

		case FIELD_COMPOSER:
			if (!song) return retstr;
			retstr = song->composer;
			break;

		case FIELD_PERFORMER:
			if (!song) return retstr;
			retstr = song->performer;
			break;

		case FIELD_DISC:
			if (!song) return retstr;
			retstr = song->disc;
			break;


		/* Times */

		case TIME_ELAPSED:
			if (!pms->comm || !pms->comm->status()) return "";

			retstr = Pms::timeformat(pms->comm->status()->time_elapsed);

			break;

		case TIME_REMAINING:
			if (!pms->comm || !pms->comm->status()) return "";

			retstr = Pms::timeformat(pms->comm->status()->time_total - pms->comm->status()->time_elapsed);

			break;

		case PROGRESSPERCENTAGE:
			if (!pms->disp || !pms->comm || !pms->comm->status() || pms->comm->status()->time_total == 0)
				retstr = "0";
			else
			{
				tmpint = pms->comm->status()->time_elapsed * 100 / pms->comm->status()->time_total;
				retstr = Pms::tostring(tmpint);
			}
			break;

		/* Widgets */

		case PROGRESSBAR:
			if (!pms->disp || !pms->comm || !pms->comm->status() || pms->comm->status()->time_total == 0) return "";

			retstr.clear();
			progress = pms->comm->status()->time_elapsed * pms->disp->topbar->bwidth() / pms->comm->status()->time_total;
			for (pint = 0; pint < progress; pint++)
			{
				retstr += "=";
			}
			retstr += ">";
			break;

		/* Status items */

		case REPEAT:
			switch(repeatmode)
			{
				case REPEAT_NONE:
					retstr = "no";
					break;
				case REPEAT_ONE:
					retstr = "one";
					break;
				case REPEAT_LIST:
					retstr = "yes";
					break;
				default:
					retstr = "unknown";
					break;
			}
			break;

		case RANDOM:
			switch(playmode)
			{
				case PLAYMODE_LINEAR:
					retstr = "no";
					break;
				case PLAYMODE_RANDOM:
					retstr = "yes";
					break;
				default:
					retstr = "unknown";
					break;
			}
			break;

		case MANUALPROGRESSION:
			if (playmode == PLAYMODE_MANUAL)
				retstr = "yes";
			else
				retstr = "no";
			break;

		case MUTE:
			if (!pms->comm) return "";
			if (pms->comm->muted())
				retstr = "yes";
			else
				retstr = "no";
			break;

		case REPEATSHORT:
			switch(repeatmode)
			{
				default:
				case REPEAT_NONE:
					retstr = "-";
					break;
				case REPEAT_ONE:
					retstr = "r";
					break;
				case REPEAT_LIST:
					retstr = "R";
					break;
			}
			break;

		case RANDOMSHORT:
			switch(playmode)
			{
				default:
				case PLAYMODE_LINEAR:
					retstr = "-";
					break;
				case PLAYMODE_RANDOM:
					retstr = "S";
					break;
			}
			break;

		case MANUALPROGRESSIONSHORT:
			if (playmode == PLAYMODE_MANUAL)
				retstr = "1";
			else
				retstr = "-";
			break;

		case MUTESHORT:
			if (!pms->comm) return "-";
			if (pms->comm->muted())
				retstr = "M";
			else
				retstr = "-";
			break;

		case LIBRARYSIZE:
			if (!pms->comm || !pms->comm->library()) return "";

			list = pms->comm->library();
			sprintf(s, "%d %s (%s)", list->size(),
					Pms::pluralformat(list->size()).c_str(),
					Pms::timeformat(list->length).c_str());

			retstr = s;
			break;

		case LISTSIZE:
			if (!pms->disp || !pms->disp->actwin()) return "";

			list = pms->disp->actwin()->plist();
			if (list)
			{
				if (list->selection.size > 0)
				{
					if (list->filtercount() == 0)
					{
						sprintf(s, "%ld/%d %s (%s)", static_cast<unsigned long>(list->selection.size),
								list->size(),
								Pms::pluralformat(list->size()).c_str(),
								Pms::timeformat(list->selection.length).c_str());
					}
					else
					{
						sprintf(s, "%ld/%d/%d %s (%s)", static_cast<unsigned long>(list->selection.size),
								list->size(),
								list->realsize(),
								Pms::pluralformat(list->realsize()).c_str(),
								Pms::timeformat(list->selection.length).c_str());
					}
				}
				else
				{
					if (list->filtercount() == 0)
					{
						sprintf(s, "%d %s (%s)", list->size(),
								Pms::pluralformat(list->size()).c_str(),
								Pms::timeformat(list->length).c_str());
					}
					else
					{
						sprintf(s, "%d/%d %s (%s)",
								list->size(),
								list->realsize(),
								Pms::pluralformat(list->realsize()).c_str(),
								Pms::timeformat(list->selection.length).c_str());
					}
				}

			}
			else
			{
				sprintf(s, "Total of %d items", pms->disp->actwin()->size());
			}

			retstr = s;
			break;

		case QUEUESIZE:
			if (!pms->comm || !pms->comm->playlist()) return "";

			tmpint = pms->comm->playlist()->qnumber();
			retstr = Pms::tostring(tmpint);
			retstr += " " + Pms::pluralformat(tmpint) + " (" + Pms::timeformat(pms->comm->playlist()->qlength()) + ")";
			break;

		case LIVEQUEUESIZE:
			if (!pms->comm || !pms->comm->status() || !pms->comm->playlist()) return "";

			if (pms->comm->playlist()->size() == 0)
			{
				retstr = "0 songs (0:00)";
				break;
			}

			switch (pms->comm->status()->state)
			{
				case MPD_STATUS_STATE_PLAY:
				case MPD_STATUS_STATE_PAUSE:
					tmp = Pms::timeformat(pms->comm->playlist()->qlength() + pms->comm->status()->time_total - pms->comm->status()->time_elapsed);
					tmpint = pms->comm->playlist()->qnumber() + 1;
					break;
				default:
					if (pms->cursong())
					{
						tmpint = pms->comm->playlist()->qnumber() + 1;
						tmp = Pms::timeformat(pms->comm->playlist()->qlength() + pms->cursong()->time);
						break;
					}
					tmpint = pms->comm->playlist()->qnumber();
					tmp = Pms::timeformat(pms->comm->playlist()->qlength());
					break;
			}
			retstr = Pms::tostring(tmpint);
			retstr += " " + Pms::pluralformat(tmpint) + " (" + tmp + ")";
			break;

		case PLAYSTATE:
			if (!pms->comm || !pms->comm->status()) return "";

			switch (pms->comm->status()->state)
			{
				default:
				case MPD_STATUS_STATE_UNKNOWN:
					retstr = pms->options->get_string("status_unknown");
					break;
				case MPD_STATUS_STATE_STOP:
					retstr = pms->options->get_string("status_stop");
					break;
				case MPD_STATUS_STATE_PLAY:
					retstr = pms->options->get_string("status_play");
					break;
				case MPD_STATUS_STATE_PAUSE:
					retstr = pms->options->get_string("status_pause");
					break;
			}
			break;

		case VOLUME:
			if (!pms->comm || !pms->comm->status()) return "";

			retstr = Pms::tostring(pms->comm->status()->volume);
			break;

		case BITRATE:
			if (!pms->comm) return "";

			retstr = Pms::tostring(pms->comm->status()->bitrate);
			break;

		case SAMPLERATE:
			if (!pms->comm) return "";

			retstr = Pms::tostring(static_cast<long>(pms->comm->status()->samplerate));
			break;

		case BITS:
			if (!pms->comm) return "";

			retstr = Pms::tostring(pms->comm->status()->bits);
			break;

		case CHANNELS:
			if (!pms->comm) return "";

			retstr = Pms::tostring(pms->comm->status()->channels);
			break;

		case LITERALPERCENT:
			retstr = "%";
			break;

		case EINVALID:
		default:
			return "";
	}

	// real length of returned string (not including colour codes)
	printlen = retstr.size();

	// escape any percent signs by doubling them
	retstr = Pms::formtext(retstr);

	/* Format string with colors */
	if (c != NULL)
	{
		return "%" + Pms::tostring(c->pair()) + "%" + retstr + "%/" + Pms::tostring(c->pair()) + "%";
	}

	return retstr;
}

color *			Formatter::getcolor(Item i, colortable_fields * f)
{
	color *		c = NULL;

	if (!f)
		return NULL;

	switch(i)
	{
		case FIELD_NUM:
			c = f->num;
			break;

		case FIELD_FILE:
			c = f->file;
			break;

		case FIELD_ARTIST:
			c = f->artist;
			break;

		case FIELD_ALBUMARTIST:
			c = f->albumartist;
			break;

		case FIELD_ARTISTSORT:
			c = f->artistsort;
			break;

		case FIELD_ALBUMARTISTSORT:
			c = f->albumartistsort;
			break;

		case FIELD_TITLE:
			c = f->title;
			break;

		case FIELD_ALBUM:
			c = f->album;
			break;

		case FIELD_DATE:
			c = f->date;
			break;

		case FIELD_YEAR:
			c = f->year;
			break;

		case FIELD_TRACK:
			c = f->track;
			break;

		case FIELD_TRACKSHORT:
			c = f->trackshort;
			break;

		case FIELD_TIME:
			c = f->time;
			break;

		case FIELD_NAME:
			c = f->name;
			break;

		case FIELD_GENRE:
			c = f->genre;
			break;

		case FIELD_COMPOSER:
			c = f->composer;
			break;

		case FIELD_PERFORMER:
			c = f->performer;
			break;

		case FIELD_DISC:
			c = f->disc;
			break;

		case TIME_ELAPSED:
			c = pms->options->colors->topbar.time_elapsed;
			break;

		case TIME_REMAINING:
			c = pms->options->colors->topbar.time_remaining;
			break;

		case PROGRESSPERCENTAGE:
			c = pms->options->colors->topbar.progresspercentage;
			break;

		case PROGRESSBAR:
			c = pms->options->colors->topbar.progressbar;
			break;

		case REPEAT:
			c = pms->options->colors->topbar.repeat;
			break;

		case RANDOM:
			c = pms->options->colors->topbar.random;
			break;

		case MANUALPROGRESSION:
			c = pms->options->colors->topbar.manualprogression;
			break;

		case MUTE:
			c = pms->options->colors->topbar.mute;
			break;

		case REPEATSHORT:
			c = pms->options->colors->topbar.repeatshort;
			break;

		case RANDOMSHORT:
			c = pms->options->colors->topbar.randomshort;
			break;

		case MANUALPROGRESSIONSHORT:
			c = pms->options->colors->topbar.manualprogressionshort;
			break;

		case MUTESHORT:
			c = pms->options->colors->topbar.muteshort;
			break;

		case LIBRARYSIZE:
			c = pms->options->colors->topbar.librarysize;
			break;

		case LISTSIZE:
			c = pms->options->colors->topbar.listsize;
			break;

		case QUEUESIZE:
			c = pms->options->colors->topbar.queuesize;
			break;

		case LIVEQUEUESIZE:
			c = pms->options->colors->topbar.livequeuesize;
			break;

		case PLAYSTATE:
			c = pms->options->colors->topbar.playstate;
			break;

		case VOLUME:
			c = pms->options->colors->topbar.volume;
			break;

		case BITRATE:
			c = pms->options->colors->topbar.bitrate;
			break;

		case SAMPLERATE:
			c = pms->options->colors->topbar.samplerate;
			break;

		case BITS:
			c = pms->options->colors->topbar.bits;
			break;

		case CHANNELS:
			c = pms->options->colors->topbar.channels;
			break;

		case LITERALPERCENT:
			c = pms->options->colors->topbar.standard;
			break;

		default:
			break;
	}

	return c;
}

/*
 * Converts an Item to a match field
 */
long			Formatter::item_to_match(Item i)
{
	long		l;

	switch(i)
	{
		case FIELD_ARTIST:
			l = MATCH_ARTIST;
			break;

		case FIELD_ARTISTSORT:
			l = MATCH_ARTISTSORT;
			break;

		case FIELD_ALBUMARTIST:
			l = MATCH_ALBUMARTIST;
			break;

		case FIELD_ALBUMARTISTSORT:
			l = MATCH_ALBUMARTISTSORT;
			break;

		case FIELD_TITLE:
			l = MATCH_TITLE;
			break;

		case FIELD_ALBUM:
			l = MATCH_ALBUM;
			break;

		case FIELD_TRACK:
		case FIELD_TRACKSHORT:
			/* only match the short one so we don't get every track 
			 * of an album marked up as xx/14 when searching for 
			 * track 14s */
			l = MATCH_TRACKSHORT;
			break;

		case FIELD_TIME:
			l = MATCH_TIME;
			break;

		case FIELD_DATE:
			l = MATCH_DATE;
			break;

		case FIELD_YEAR:
			l = MATCH_YEAR;
			break;

		case FIELD_GENRE:
			l = MATCH_GENRE;
			break;

		case FIELD_COMPOSER:
			l = MATCH_COMPOSER;
			break;

		case FIELD_PERFORMER:
			l = MATCH_PERFORMER;
			break;

		case FIELD_DISC:
			l = MATCH_DISC;
			break;

		default:
			return MATCH_FAILED;
	}

	return l;
}


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
 * field.h - format a song using field variables
 *
 */

#ifndef _PMS_FIELD_H_
#define _PMS_FIELD_H_

#include <string>
#include <vector>
#include "song.h"
#include "color.h"

/*
 * Insertable items and items in playlist
 */
typedef enum
{
	EINVALID = -1,

	/* Field types are present both in library view and topbar view */
	FIELD_NUM,
	FIELD_FILE,
	FIELD_ARTIST,
	FIELD_ARTISTSORT,
	FIELD_ALBUMARTIST,
	FIELD_ALBUMARTISTSORT,
	FIELD_TITLE,
	FIELD_ALBUM,
	FIELD_TRACK,
	FIELD_TRACKSHORT,
	FIELD_TIME,
	FIELD_DATE,
	FIELD_YEAR,
	FIELD_NAME,
	FIELD_GENRE,
	FIELD_COMPOSER,
	FIELD_PERFORMER,
	FIELD_DISC,
	FIELD_COMMENT,

	/* Conditionals */
	COND_IFCURSONG,
	COND_IFPLAYING,
	COND_IFPAUSED,
	COND_IFSTOPPED,
	COND_ELSE,
	COND_ENDIF,

	/* These types are only available to the topbar */
	REPEAT,
	RANDOM,
	MANUALPROGRESSION,
	MUTE,
	REPEATSHORT,
	RANDOMSHORT,
	MANUALPROGRESSIONSHORT,
	MUTESHORT,
	TIME_ELAPSED,
	TIME_REMAINING,
	PLAYSTATE,
	PROGRESSBAR,
	PROGRESSPERCENTAGE,
	VOLUME,
	LIBRARYSIZE,
	LISTSIZE,
	QUEUESIZE,
	LIVEQUEUESIZE,

	/* Audio properties */
	BITRATE,
	SAMPLERATE,
	BITS,
	CHANNELS,

	/* Misc */
	LITERALPERCENT

}
Item;


/*
 * Formatter class formats a song into names, i.e:
 *
 * format(song, "%artist% - %album%");
 * 	returns
 * "U2 - Beautiful Day"
 *
 */
class Formatter
{
private:
	string			fm;

	Item			nextitem(string, int *, int *);
	string			evalconditionals(string);

public:
	string			format(Song *, string, unsigned int &, colortable_fields *, bool = false);
	string			format(Song *, Item, unsigned int &, colortable_fields *, bool = false);
	string			format(Song *, Item, bool = false);
	color *			getcolor(Item, colortable_fields *);
	vector<Item> *		multiformat_item(string);
	Item			field_to_item(string);
	long			item_to_match(Item);
};


#endif /* _PMS_FIELD_H_ */

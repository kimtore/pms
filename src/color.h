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
 * color.h
 * Color class, holds information about color + attribute values for one color pair
 */

#ifndef _COLOR_H_
#define _COLOR_H_

class color
{
private:
	short			id;
	int			attr;
	bool			initialized;
	bool			isset;
public:
	static short		numpairs;

				color();
				color(int, int, int);
				~color();
	void			clean();
	int			pair();
	int			front;
	int			back;
	bool			set(int, int, int);
};

/*
 * All possible tag fields
 */
typedef struct
{
	color *			num;
	color *			file;
	color *			artist;
	color *			artistsort;
	color *			albumartist;
	color *			albumartistsort;
	color *			title;
	color *			album;
	color *			genre;
	color *			track;
	color *			trackshort;
	color *			time;
	color *			date;
	color *			year;
	color *			name;
	color *			composer;
	color *			performer;
	color *			disc;
	color *			comment;

} colortable_fields;

/*
 * All colored items in the topbar
 */
typedef struct
{
	colortable_fields	fields;

	color *			standard;
	color *			time_elapsed;
	color *			time_remaining;
	color *			progressbar;
	color *			progresspercentage;
	color *			librarysize; //songstotal
	color *			listsize;
	color *			queuesize;
	color *			livequeuesize; // songstotalinccurrent
	color *			playstate;
	color *			volume;

	color *			bitrate;
	color *			samplerate;
	color *			bits;
	color *			channels;

	color *			repeat;
	color *			random;
	color *			manualprogression;
	color *			mute;
	color *			repeatshort;
	color *			randomshort;
	color *			manualprogressionshort;
	color *			muteshort;


} colortable_topbar;


class Colortable
{
private:
	bool			isset;
	void			clear();

public:
				Colortable();
				~Colortable();

	void			defaults();

	/* other tables */
	colortable_fields	fields;
	colortable_topbar	topbar;

	/* main colors */
	color *			back;
	color *			standard;
	color *			border;
	color *			headers;
	color *			title;
	color *			status;
	color *			status_error;
	color *			position;

	/* list colors */
	color *			current;
	color *			cursor;
	color *			selection;
	color *			lastlist;
	color *			playinglist;

};


#endif /* _COLOR_H_ */

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

#ifndef _PMS_FIELDS_H_
#define _PMS_FIELDS_H_

typedef enum
{
	/* These fields come directly from MPD. */
	FIELD_DIRECTORY,
	FIELD_FILE,
	FIELD_POS,
	FIELD_ID,
	FIELD_TIME,
	FIELD_NAME,
	FIELD_ARTIST,
	FIELD_ARTISTSORT,
	FIELD_ALBUM,
	FIELD_TITLE,
	FIELD_TRACK,
	FIELD_DATE,
	FIELD_DISC,
	FIELD_GENRE,
	FIELD_ALBUMARTIST,
	FIELD_ALBUMARTISTSORT,

	/* Custom fields used only in PMS */
	FIELD_YEAR,
	FIELD_TRACKSHORT,

	/* These fields are mainly for use in the topbar. DO NOT include them in FIELD_COLUMN_VALUES. */
	FIELD_ELAPSED,
	FIELD_REMAINING,
	FIELD_MODES,
	FIELD_STATE,
	FIELD_QUEUESIZE,
	FIELD_QUEUELENGTH
}

field_t;

#define FIELD_COLUMN_VALUES 18
#define FIELD_TOTAL_VALUES 24

#endif /* _PMS_FIELDS_H_ */

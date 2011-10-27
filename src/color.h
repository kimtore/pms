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

#ifndef _PMS_COLOR_H_
#define _PMS_COLOR_H_

#include "field.h"
#include <string>
#include <vector>
using namespace std;

class Color
{
	private:

		short			id;

	public:

					Color();

		void			set(short nfront, short nback, int nattr);
		bool			set(string strcolor);
		string			getstrname();

		static short		color_count;

		int			pair;
		int			attr;
		short			front;
		short			back;
};

class Colortable
{
	private:
		
		short		dback;
		short		dfront;

	public:
				Colortable();
				~Colortable();

		void		load_defaults();

		vector<Color *>	table;

		/* Main colors */
		Color *		standard;
		Color *		topbar;
		Color *		statusbar;
		Color *		windowtitle;
		Color *		columnheader;
		Color *		console;
		Color *		error;
		Color *		readout;

		/* List colors */
		Color *		cursor;
		Color *		playing;
		Color *		selection;

		/* Field colors */
		Color *		field[FIELD_TOTAL_VALUES];

};


#endif /* _PMS_COLOR_H_ */

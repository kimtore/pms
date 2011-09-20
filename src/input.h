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

#ifndef _PMS_INPUT_H_
#define _PMS_INPUT_H_

#include <string>
using namespace std;

#define INPUT_NOINPUT -1
#define INPUT_BUFFERED 0
#define INPUT_RUN 1

#define INPUT_MODE_COMMAND 0
#define INPUT_MODE_INPUT 1
#define INPUT_MODE_SEARCH 2

class Input
{
	private:
		int		chbuf;
		int		mode;
		string		buffer;

	public:
		Input();

		/* Read next character from ncurses buffer */
		int		next();

		/* Setter and getter for mode */
		void		setmode(int nmode);
		int		getmode() { return mode; }
};

#endif /* _PMS_INPUT_H_ */

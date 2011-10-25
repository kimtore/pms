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

#ifndef _PMS_CONSOLE_H_
#define _PMS_CONSOLE_H_

#include <string>
#include <sys/time.h>
using namespace std;

#define MSG_LEVEL_ERR 0
#define MSG_LEVEL_WARN 1
#define MSG_LEVEL_INFO 2
#define MSG_LEVEL_DEBUG 3

class Logline
{
	public:
		Logline(int lvl, const char * ln);

		struct timeval tm;
		int level;
		string line;
};

/* Log a message to stderr */
void console_log(int level, const char * format, ...);

#define debug(_fmt, ...)	console_log(MSG_LEVEL_DEBUG, _fmt, __VA_ARGS__)
#define stinfo(_fmt, ...)	console_log(MSG_LEVEL_INFO, _fmt, __VA_ARGS__)
#define sterr(_fmt, ...)	console_log(MSG_LEVEL_ERR, _fmt, __VA_ARGS__)

#endif /* _PMS_CONSOLE_H_ */

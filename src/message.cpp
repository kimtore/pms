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
 * message.cpp - The Message class
 *
 */


#include <cstdarg>
#include <stdio.h>
#include "message.h"


Message::Message()
{
	clear();
}

void		Message::clear()
{
	code = 0;
	str.clear();
	time(&timestamp);
};

void		Message::assign(int c, const char * format, ...)
{
	va_list		ap;
	char		buffer[1024];

	va_start(ap, format);
	vsprintf(buffer, format, ap);
	va_end(ap);

	code = c;
	str = buffer;
	time(&timestamp);
};

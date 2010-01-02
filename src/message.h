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
 * message.h - The Message class
 *
 */

#ifndef _PMS_MESSAGE_H_
#define _PMS_MESSAGE_H_

#include <string>

using namespace std;


/*
 * Message verbosity levels
 */
enum
{
	MSG_STATUS = 0,
	MSG_CONSOLE,
	MSG_DEBUG
};

enum
{
	CERR_NONE = 0,
	CERR_NO_FILE,
	CERR_SYNTAX,
	CERR_UNKNOWN_COMMAND,
	CERR_EXCESS_ARGUMENTS,
	CERR_MISSING_IDENTIFIER,
	CERR_MISSING_VALUE,
	CERR_UNEXPECTED_TOKEN,
	CERR_INVALID_VALUE,
	CERR_INVALID_OPTION,
	CERR_INVALID_KEY,
	CERR_INVALID_COLOR,
	CERR_INVALID_COMMAND,
	CERR_INVALID_IDENTIFIER,
	CERR_INVALID_COLUMN,
	CERR_INVALID_TOPBAR_INDEX,
	CERR_INVALID_TOPBAR_POSITION
};

class Message
{
public:
				Message();

	void			clear();
	void			assign(int, const char *, ...);

	int			code;
	string			str;
	time_t			timestamp;
};


#endif /* _PMS_MESSAGE_H_ */

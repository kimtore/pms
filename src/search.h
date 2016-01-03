/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS	<<Practical Music Search>>
 * Copyright (C) 2006-2016  Kim Tore Jensen
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
 */


#ifndef _PMS_SEARCH_H_
#define _PMS_SEARCH_H_

#include "../config.h"
#include <string>

using namespace std;

#ifdef HAVE_REGEX
	/**
	 * Match a string against a regular expression, case insensitively.
	 *
	 * Returns true if the regular expression matches, false otherwise.
	 */
	bool
	match_regex(string * source, string * pattern);
#endif

/**
 * Match a string against another string, case insensitively.
 *
 * Returns true if the strings are identical, false otherwise.
 */
bool
match_exact(string * s1, string * s2);

/**
 * Match a string inside another string, case insensitively.
 *
 * Returns true if needle is found inside haystack.
 */
bool
match_inside(string * haystack, string * needle);

/**
 * Implementation of match_exact and match_inside.
 *
 * Returns true if needle is found inside haystack.
 */
static bool
match_run(string * haystack, string * needle, bool exact);

#endif /* _PMS_SEARCH_H_ */

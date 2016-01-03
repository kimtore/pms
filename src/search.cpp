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


#include "search.h"
#include <string>

using namespace std;

#ifdef HAVE_REGEX
#include <regex>
bool
match_regex(string * source, string * pattern)
{
	bool matched;
	regex reg;

	try
	{
		reg.assign(*pattern, std::regex_constants::icase);
		matched = regex_search(*source, reg);
	}
	catch (std::regex_error& err)
	{
		return false;
	}

	return matched;
}
#endif

bool
match_exact(string * s1, string * s2)
{
	return match_run(s1, s2, true);
}

bool
match_inside(string * haystack, string * needle)
{
	return match_run(haystack, needle, false);
}

static bool
match_run(string * haystack, string * needle, bool exact)
{
	bool			matched = exact;

	string::const_iterator	it_haystack;
	string::const_iterator	it_needle;

	for (it_haystack = haystack->begin(), it_needle = needle->begin(); it_haystack != haystack->end() && it_needle != needle->end(); it_haystack++)
	{
		/* exit if there aren't enough characters left to match the string */
		if (haystack->end() - it_haystack < needle->end() - it_needle) {
			return false;
		}

		/* check next character in needle with character in haystack */
		if (::toupper(*it_needle) == ::toupper(*it_haystack)) {
			/* matched a letter -- look for next letter */
			matched = true;
			it_needle++;
		} else if (exact) {
			/* didn't match a letter but need exact match */
			return false;
		} else {
			/* didn't match a letter -- start from first letter of needle */
			matched = false;
			it_needle = needle->begin();
		}
	}

	if (it_needle != needle->end()) {
		/* end of the haystack before getting to the end of the needle */
		return false;
	}

	if (exact && it_needle == needle->end() && it_haystack != haystack->end()) {
		/* need exact and got to the end of the needle but not the end of the
		 * haystack */
		return false;
	}

	return matched;
}


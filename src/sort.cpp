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

#include "songlist.h"
#include "fields.h"
#include "config.h"
#include "console.h"
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

extern Fieldtypes * fieldtypes;
extern Config * config;

static Field * sort_field;
static bool sort_numeric_compare;

static inline bool sort_cmp(Song * a, Song * b)
{
	static string sa;
	static string sb;

	if (sort_numeric_compare)
		return (atoi(a->f[sort_field->type].c_str()) < atoi(b->f[sort_field->type].c_str()));

	if (config->sort_case)
		return (a->f[sort_field->type] < b->f[sort_field->type]);

	sa = a->f[sort_field->type];
	sb = b->f[sort_field->type];
	
	std::transform(sa.begin(), sa.end(), sa.begin(), ::tolower);
	std::transform(sb.begin(), sb.end(), sb.begin(), ::tolower);

	return (sa < sb);
}

void Songlist::sort(string sortstr)
{
	size_t start = 0, end = 0;
	string s;
	vector<Song *> * source;

	source = &(searchresult ? searchresult->songs : songs);

	while (start < sortstr.size())
	{
		if ((end = sortstr.find(' ', start)) != string::npos)
			s = sortstr.substr(start, end - start);
		else
			s = sortstr.substr(start);

		if ((sort_field = fieldtypes->find(s)) != NULL)
		{
			switch(sort_field->type)
			{
				case FIELD_POS:
				case FIELD_ID:
				case FIELD_TIME:
				case FIELD_TRACK:
				case FIELD_DATE:
				case FIELD_YEAR:
				case FIELD_DISC:
					sort_numeric_compare = true;
					break;

				default:
					sort_numeric_compare = false;
			}

			if (start == 0)
				std::sort(source->begin(), source->end(), sort_cmp);
			else
				std::stable_sort(source->begin(), source->end(), sort_cmp);
		}
		else
			sterr("Invalid sort field `%s', ignoring.", s.c_str());

		if (end == string::npos)
			break;

		start = end + 1;
	}
}

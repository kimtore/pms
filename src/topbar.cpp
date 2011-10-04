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

#include <string>
#include <vector>
#include "topbar.h"
#include "field.h"
#include "console.h"
#include "song.h"
#include "config.h"
using namespace std;

extern Fieldtypes fieldtypes;
extern Config config;

Topbarchunk::Topbarchunk(string s, Color * c)
{
	str = s;
	color = c;
}

Topbarsegment::Topbarsegment()
{
	condition = 0;
	format = "";
}

Topbarsegment::~Topbarsegment()
{
	vector<Topbarchunk *>::iterator i;
	for (i = chunks.begin(); i != chunks.end(); ++i)
		delete *i;
}

unsigned int Topbarsegment::compile(Song * song)
{
	string str;
	size_t s = 0, e = 0;
	unsigned int strl = 0;
	vector<Topbarchunk *>::iterator i;
	vector<Field *>::iterator field;

	for (i = chunks.begin(); i != chunks.end(); ++i)
		delete *i;

	chunks.clear();
	if (!format.size())
		return 0;

	field = fields.begin();

	do
	{
		if ((e = format.find('$', s)) != string::npos)
		{
			/* our string */
			chunks.push_back(new Topbarchunk(format.substr(s, e - s), config.colors.topbar));
			chunks.push_back(new Topbarchunk((*field)->format(song), config.colors.field[(*field)->type]));
			++field;
			s = e + 1;
			continue;
		}
	}
	while(field != fields.end());

	if (s != string::npos)
		chunks.push_back(new Topbarchunk(format.substr(s), config.colors.topbar));
	
	for (i = chunks.begin(); i != chunks.end(); ++i)
		strl += (*i)->str.size();

	return strl;
}

Topbarline::~Topbarline()
{
	vector<Topbarsegment *>::iterator i;
	for (i = segments.begin(); i != segments.end(); ++i)
		delete *i;
}

void Topbar::clear()
{
	vector<Topbarline *>::iterator i;
	unsigned int t;

	for (t = 0; t < 3; ++t)
	{
		for (i = lines[t].begin(); i != lines[t].end(); ++i)
			delete *i;
		lines[t].clear();
	}
};

int Topbar::set(string format)
{
	string::iterator it;
	string working = "";
	string varname = "";
	bool bracket = false;
	bool var = false;
	unsigned int pos = 0; /* left/center/right */
	Field * field;
	Topbarline * line = new Topbarline;
	Topbarsegment * segment = new Topbarsegment;

	cached_format = format;

	clear();
	if (!format.size())
		return true;

	lines[pos].push_back(line);

	for (it = format.begin(); it != format.end(); ++it)
	{
		/* Variable handling */
		if (var)
		{
			if (*it < 'a' || *it > 'z')
			{
				--it;
				if ((field = fieldtypes.find(varname)) != NULL)
				{
					segment->fields.push_back(field);
					segment->format += '$';
					var = false;
					varname.clear();
					continue;
				}
				sterr("Topbar: unknown variable `%s' near %s", varname.c_str(), working.c_str());
				delete segment;
				return false;
			}

			varname += *it;
			working += *it;
			continue;
		}
		working += *it;

		/*
		 * Closing and opening of brackets
		 */
		if (*it == '{')
		{
			if (!bracket)
			{
				bracket = true;
				continue;
			}
			sterr("Topbar: unexpected `%c', expected identifier near %s", *it, working.c_str());
			delete segment;
			return false;
		}
		if (*it == '}')
		{
			if (bracket)
			{
				bracket = false;
				segment->src = working;
				line->segments.push_back(segment);
				segment = new Topbarsegment;
				line = new Topbarline;
				if (++pos > 2)
					pos = 0;
				lines[pos].push_back(line);
				working.clear();
				continue;
			}
		}
		if (!bracket)
		{
			sterr("Topbar: unexpected `%c', expected `{' near %s", *it, working.c_str());
			delete segment;
			return false;
		}

		if (*it == '$')
		{
			var = true;
			continue;
		}

		/* If all else fails, append to text */
		segment->format += *it;
	}

	if (bracket)
	{
		sterr("Topbar: unexpected end of segment, expected `}' near %s", working.c_str());
		delete segment;
		return false;
	}

	return true;
}

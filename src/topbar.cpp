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
#include "mpd.h"
using namespace std;

extern Fieldtypes fieldtypes;
extern Config config;
extern MPD mpd;

Topbarchunk::Topbarchunk(string s, Color * c)
{
	str = s;
	color = c;
}

Topbarsegment::Topbarsegment()
{
	condition.t = CONDITION_NONE;
	condition.f = CONDITION_NONE;
	condition.ctl = CONDITION_NONE;
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
	unsigned int current_condition = CONDITION_NONE;
	size_t s = 0, e = 0;
	unsigned int strl = 0;
	vector<Topbarchunk *>::iterator i;
	vector<Field *>::iterator field;

	for (i = chunks.begin(); i != chunks.end(); ++i)
		delete *i;

	chunks.clear();
	if (!format.size())
		return 0;

	/* What is our current condition? */
	if (mpd.status.state == MPD_STATE_PLAY)
		current_condition |= CONDITION_PLAYING;
	if (mpd.currentsong)
		current_condition |= CONDITION_SONG;
	if (mpd.is_connected())
		current_condition |= CONDITION_CONNECTED;

	/* Evaluate wanted conditions against current conditions */
	if ((condition.t & current_condition) != condition.t)
		return 0;
	if ((condition.f & ~current_condition) != condition.f)
		return 0;

	/* Iterate through fields and place them into chunks */
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
	bool conditional = false;
	bool invcondition;
	unsigned int condition = CONDITION_NONE;
	condition_t lastcondition = { CONDITION_NONE, CONDITION_NONE, CONDITION_NONE };
	vector<condition_t> conditions;
	unsigned int pos = 0; /* left/center/right */
	Field * field;
	Topbarline * line;
	Topbarsegment * segment;

	cached_format = format;

	clear();
	if (!format.size())
		return true;

	line = new Topbarline;
	segment = new Topbarsegment;
	lines[pos].push_back(line);

	for (it = format.begin(); it != format.end(); ++it)
	{
		/* Variable handling */
		if (var)
		{
			if (*it < 'a' || *it > 'z')
			{
				/* End of variable, store this field into segment */
				if ((field = fieldtypes.find(varname)) != NULL)
				{
					--it;
					segment->fields.push_back(field);
					segment->format += '$';
					var = false;
					varname.clear();
					continue;
				}

				/* Open conditional */
				else if ((*it == '(' && varname == "if") ||
						(*it == '{' && varname == "else"))
				{
					/* Store segment */
					segment->src = working;
					line->segments.push_back(segment);

					/* Create new segment */
					segment = new Topbarsegment;
					if (conditions.size())
						segment->condition = conditions.back();

					/* Reverse conditional "else" */
					if (*it == '{')
					{
						--it;

						if (invcondition)
							segment->condition.t |= lastcondition.ctl;
						else
							segment->condition.f |= lastcondition.ctl;

						segment->condition.ctl = condition;
						conditions.push_back(segment->condition);
					}

					/* Open a new segment based on conditional */
					var = false;
					conditional = true;
					working.clear();
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

		/* Conditional name handling */
		if (conditional && *it != '{')
		{
			if ((*it < 'a' || *it > 'z') && *it != '!')
			{
				/* Detect conditional name */
				if (*it != ')')
				{
					sterr("Topbar: unexpected `%c', expected `)' near %s", *it, working.c_str());
					delete segment;
					return false;
				}

				if ((invcondition = (varname.size() > 1 && varname[0] == '!')))
					varname = varname.substr(1);

				if (varname == "playing")
					condition = CONDITION_PLAYING;
				else if (varname == "song")
					condition = CONDITION_SONG;
				else if (varname == "connected")
					condition = CONDITION_CONNECTED;
				else
				{
					sterr("Topbar: unknown condition `%s' near %s", varname.c_str(), working.c_str());
					delete segment;
					return false;
				}

				/* Valid condition name, assign to current segment */
				if (!invcondition)
					segment->condition.t |= condition;
				else
					segment->condition.f |= condition;

				segment->condition.ctl = condition;
				conditions.push_back(segment->condition);
				varname.clear();
				working.clear();
				continue;
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
			if (conditional)
			{
				conditional = false;
				continue;
			}
			sterr("Topbar: unexpected `%c', expected identifier near %s", *it, working.c_str());
			delete segment;
			return false;
		}
		if (*it == '}')
		{
			if (conditional)
			{
				sterr("Topbar: unexpected `%c', expected `)' near %s", *it, working.c_str());
				delete segment;
				return false;
			}

			/* End conditional bracket */
			if (conditions.size())
			{
				lastcondition = conditions.back();
				conditions.pop_back();

				/* Store this segment */
				segment->src = working;
				line->segments.push_back(segment);

				/* Open a new segment based on previous conditional */
				segment = new Topbarsegment;
				if (conditions.size())
					segment->condition = conditions.back();

				working.clear();
				continue;
			}

			if (bracket)
			{
				bracket = false;
				segment->src = working;
				line->segments.push_back(segment);
				segment = new Topbarsegment;
				if (conditions.size())
					segment->condition = conditions.back();

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

	if (conditions.size())
	{
		sterr("Topbar: unexpected end of line, expected end of conditional `}' near %s", working.c_str());
		delete segment;
		return false;
	}

	if (conditional)
	{
		sterr("Topbar: unexpected end of line, expected end of expression `)' near %s", working.c_str());
		delete segment;
		return false;
	}

	if (bracket)
	{
		sterr("Topbar: unexpected end of line, expected end of segment `}' near %s", working.c_str());
		delete segment;
		return false;
	}

	return true;
}

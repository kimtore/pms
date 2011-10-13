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

#include "color.h"
#include "curses.h"
#include "field.h"
#include "console.h"

short Color::color_count = 0;

Colortable::Colortable()
{
	pair_content(-1, &dfront, &dback);

	standard = new Color(dfront, dback, 0);
	topbar = new Color(COLOR_WHITE, -1, 0);
	statusbar = new Color(COLOR_WHITE, -1, 0);
	windowtitle = new Color(COLOR_CYAN, -1, A_BOLD);
	columnheader = new Color(COLOR_WHITE, -1, 0);
	console = new Color(COLOR_WHITE, -1, 0);
	error = new Color(COLOR_WHITE, COLOR_RED, A_BOLD);
	readout = new Color(COLOR_WHITE, -1, 0);

	cursor = new Color(COLOR_BLACK, COLOR_WHITE, 0);
	playing = new Color(COLOR_BLACK, COLOR_YELLOW, 0);

	field[FIELD_DIRECTORY] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_FILE] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_POS] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_ID] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_TIME] = new Color(COLOR_MAGENTA, -1, 0);
	field[FIELD_NAME] = new Color(COLOR_WHITE, -1, A_BOLD);
	field[FIELD_ARTIST] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_ARTISTSORT] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_ALBUM] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_TITLE] = new Color(COLOR_WHITE, -1, A_BOLD);
	field[FIELD_TRACK] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_DATE] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_DISC] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_GENRE] = new Color(COLOR_WHITE, -1, 0);
	field[FIELD_ALBUMARTIST] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_ALBUMARTISTSORT] = new Color(COLOR_YELLOW, -1, 0);

	field[FIELD_YEAR] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_TRACKSHORT] = new Color(COLOR_CYAN, -1, 0);

	field[FIELD_ELAPSED] = new Color(COLOR_GREEN, -1, 0);
	field[FIELD_REMAINING] = new Color(COLOR_MAGENTA, -1, 0);
	field[FIELD_VOLUME] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_PROGRESSBAR] = new Color(COLOR_BLACK, -1, A_BOLD);
	field[FIELD_MODES] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_STATE] = new Color(COLOR_CYAN, -1, 0);
	field[FIELD_QUEUESIZE] = new Color(COLOR_YELLOW, -1, 0);
	field[FIELD_QUEUELENGTH] = new Color(COLOR_WHITE, -1, 0);
}

Colortable::~Colortable()
{
}

Color::Color(short nfront, short nback, int nattr)
{
	id = Color::color_count;
	set(nfront, nback, nattr);
	Color::color_count++;
}

void Color::set(short nfront, short nback, int nattr)
{
	front = nfront;
	back = nback;
	attr = nattr;
	init_pair(id, front, back);
	pair = COLOR_PAIR(id) | attr;
}

bool Color::set(string strcolor)
{
	vector<string>::iterator it;
	vector<string> cols;
	string t;
	size_t start = 0, end = 0;
	short nfront = -1;
	short nback = -1;
	int nattr = 0;
	short * cur = &nfront;

	while (start + 1 < strcolor.size())
	{
		if ((end = strcolor.find(' ', start)) != string::npos)
			t = strcolor.substr(start, end - start);
		else
			t = strcolor.substr(start);

		cols.push_back(t);

		if (end == string::npos)
			break;

		start = end + 1;
	}

	for (it = cols.begin(); it != cols.end(); ++it)
	{
		/* Attributes */
		if (*it == "bold")
			nattr |= A_BOLD;
		else if (*it == "reverse")
			nattr |= A_REVERSE;
		else if (cur == NULL)
		{
			sterr("Trailing characters near `%s'", it->c_str());
			return false;
		}
		else
		{
			/* Front colors only */
			if (cur == &nfront && (*it == "gray" || *it == "grey"))
			{
				*cur = COLOR_BLACK;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && (*it == "brightred" || *it == "lightred"))
			{
				*cur = COLOR_RED;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && (*it == "brightgreen" || *it == "lightgreen"))
			{
				*cur = COLOR_GREEN;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && *it == "yellow")
			{
				*cur = COLOR_YELLOW;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && (*it == "brightblue" || *it == "lightblue"))
			{
				*cur = COLOR_BLUE;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && (*it == "brightmagenta" || *it == "lightmagenta"))
			{
				*cur = COLOR_MAGENTA;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && (*it == "brightcyan" || *it == "lightcyan"))
			{
				*cur = COLOR_CYAN;
				nattr |= A_BOLD;
			}
			else if (cur == &nfront && *it == "white")
			{
				*cur = COLOR_WHITE;
				nattr |= A_BOLD;
			}

			/* Applies everywhere */
			else if (*it == "black")
				*cur = COLOR_BLACK;
			else if (*it == "red")
				*cur = COLOR_RED;
			else if (*it == "green")
				*cur = COLOR_GREEN;
			else if (*it == "brown")
				*cur = COLOR_YELLOW;
			else if (*it == "blue")
				*cur = COLOR_BLUE;
			else if (*it == "magenta")
				*cur = COLOR_MAGENTA;
			else if (*it == "cyan")
				*cur = COLOR_CYAN;
			else if (*it == "gray" || *it == "brightgray" || *it == "lightgray" || *it == "white")
				*cur = COLOR_WHITE;
			else
			{
				sterr("Invalid color `%s' for use in %s, ignoring.", it->c_str(), cur == &nfront ? "foreground" : "background");
				return false;
			}

			if (cur == &nfront)
				cur = &nback;
			else
				cur = NULL;
		}
	}

	set(nfront, nback, nattr);

	return true;
}

string Color::getstrname()
{
	string f;
	string b;
	string a;

	/*
	 * Foreground colors
	 */
	if (attr & A_BOLD)
	{
		if (front == COLOR_BLACK)
			f = "gray";
		else if (front == COLOR_YELLOW)
			f = "yellow";
		else if (front == COLOR_WHITE)
			f = "white";
	}
	if (f.empty())
	{
		if (front == COLOR_BLACK)
			f = "black";
		else if (front == COLOR_RED)
			f = "red";
		else if (front == COLOR_GREEN)
			f = "green";
		else if (front == COLOR_YELLOW)
			f = "brown";
		else if (front == COLOR_BLUE)
			f = "blue";
		else if (front == COLOR_MAGENTA)
			f = "magenta";
		else if (front == COLOR_CYAN)
			f = "cyan";
		else if (front == COLOR_WHITE)
			f = "brightgray";

		if (attr & A_BOLD)
			f = "bright" + f;
	}

	/*
	 * Background colors
	 */
	if (back == -1);
	else if (back == COLOR_BLACK)
		b = "black";
	else if (back == COLOR_RED)
		b = "red";
	else if (back == COLOR_GREEN)
		b = "green";
	else if (back == COLOR_YELLOW)
		b = "brown";
	else if (back == COLOR_BLUE)
		b = "blue";
	else if (back == COLOR_MAGENTA)
		b = "magenta";
	else if (back == COLOR_CYAN)
		b = "cyan";
	else if (back == COLOR_WHITE)
		b = "gray";

	/*
	 * Attributes
	 */
	if (attr & A_REVERSE)
		a = "reverse";

	if (b.size())
		f += " " + b;
	if (a.size())
		f += " " + a;

	return f;
}

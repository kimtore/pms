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

#include "curses.h"
#include <cstring>
#include <string>
#include <stdlib.h>
#include <math.h>

using namespace std;

Curses::Curses()
{
	if ((initscr()) == NULL)
	{
		ready = false;
		return;
	}

	raw();
	noecho();
	halfdelay(10);
	keypad(stdscr, true);
	curs_set(0);

	if (has_colors() && start_color())
	{
		use_default_colors();
		hascolors = true;
	}

	detect_dimensions();
	clear();
	refresh();

	ready = true;
}

Curses::~Curses()
{
	endwin();
}

void Curses::detect_dimensions()
{
	memset(&self, 0, sizeof self);
	memset(&topbar, 0, sizeof topbar);
	memset(&main, 0, sizeof main);
	memset(&statusbar, 0, sizeof statusbar);
	memset(&position, 0, sizeof position);

	self.right = COLS - 1;
	self.bottom = LINES - 1;

	topbar.top = 1;
	topbar.bottom = topbar.top;
	topbar.right = self.right;

	main.top = topbar.bottom + 1;
	main.bottom = self.bottom - 1;
	main.right = self.right;

	statusbar.top = self.bottom;
	statusbar.bottom = self.bottom;
	statusbar.right = self.right - 3;

	position.top = statusbar.top;
	position.bottom = statusbar.bottom;
	position.left = statusbar.right + 1;
	position.right = self.right;
}

void Curses::flush()
{
	refresh();
}

void Curses::clearline(Rect * rect, int line)
{
	Rect r;

	if (!rect)
		return;

	memcpy(&r, rect, sizeof r);
	if ((r.top += line) > r.bottom)
		return;
	
	r.bottom = r.top;

	wipe(&r);
}

void Curses::wipe(Rect * rect)
{
	int ix, yx;

	if (!rect)
		return;
	
	for (yx = rect->top; yx <= rect->bottom; yx++)
	{
		move(yx, rect->left);
		for (ix = rect->left; ix <= rect->right; ix++)
		{
			addch(' ');
		}
	}
}

void Curses::print(Rect * rect, int y, int x, const char * fmt, ...)
{
	va_list			ap;
	unsigned int		i = 0;
	double			f = 0;
	string			output = "";
	bool			parse = false;
	bool			attr = false;
	char			buf[1024];
	string			colorstr;
	int			pair = 0;
	unsigned int		maxlen;		// max allowed characters printed on screen
	unsigned int		printlen = 0;	// num characters printed on screen

	if (!rect)
		return;
	
	va_start(ap, fmt);

	//pair = c->pair();
	pair = 0;
	maxlen = rect->right - rect->left + 1;
	move(rect->top + y, rect->left + x);
	attron(pair);

	while(*fmt && printlen < maxlen)
	{
		if (*fmt == '%' && !parse)
		{
			if (*(fmt + 1) == '%')
			{
				fmt += 2;
				output = "%%";
				printw(output.c_str());
				continue;
			}
			parse = true;
			attr = true;
			++fmt;
		}

		if (parse)
		{
			switch(*fmt)
			{
				case '/':
				/* Turn off attribute, SGML style */
					attr = false;
					break;
				case 'B':
					if (attr)
						attron(A_BOLD);
					else
						attroff(A_BOLD);
					parse = false;
					break;
				case 'R':
					if (attr)
						attron(A_REVERSE);
					else
						attroff(A_REVERSE);
					parse = false;
					break;
				case 'd':
					parse = false;
					i = va_arg(ap, int);
					sprintf(buf, "%d", i);
					printw(buf);
					printlen += strlen(buf);
					i = 0;
					break;
				case 'f':
					parse = false;
					f = va_arg(ap, double);
					sprintf(buf, "%f", f);
					printw(buf);
					printlen += strlen(buf);
					break;
				case 's':
					parse = false;
					output = va_arg(ap, const char *);
					if (output.size() >= (maxlen - printlen))
					{
						output = output.substr(0, (maxlen - printlen - 1));
					}
					sprintf(buf, "%s", output.c_str());
					printw(buf);
					printlen += strlen(buf);
					break;
				case 0:
					parse = false;
					continue;
				default:
					/* Use colors? */
					i = atoi(fmt);
					if (i >= 0)
					{
						if (attr)
						{
							attroff(pair);
							attron(i);
						}
						else
						{
							attroff(i);
							attron(pair);
						}

						/* Skip characters */
						fmt += static_cast<int>(floor(log(i)) + 1);
					}
					parse = false;
					break;
			}
		}
		else
		{
			output = *fmt;
			printw(output.c_str());
			++printlen;
		}
		++fmt;
	}

	va_end(ap);
	attroff(pair);
}

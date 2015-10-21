/* vi:set ts=8 sts=8 sw=8 noet:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2015  Kim Tore Jensen
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
 * display.cpp - ncurses, display and window management
 *
 */

#include "display.h"
#include "config.h"
#include "pms.h"
#include <cstdio>
#include <sstream>

extern Pms *			pms;

Point::Point()
{
	x = 0;
	y = 0;
}

Point::Point(uint16_t x_, uint16_t y_)
{
	x = x_;
	y = y_;
}

Point &
Point::operator=(const Point & src)
{
	x = src.x;
	y = src.y;

	return *this;
}

BBox::BBox()
{
	tl.x = 0;
	tl.y = 0;
	br.x = 0;
	br.y = 0;
	window = NULL;
}

inline
uint16_t
BBox::top()
{
	return tl.y;
}

inline
uint16_t
BBox::bottom()
{
	return br.y;
}

inline
uint16_t
BBox::left()
{
	return tl.x;
}

inline
uint16_t
BBox::right()
{
	return br.x;
}

inline
uint16_t
BBox::width()
{
	return br.x - tl.x + 1;
}

inline
uint16_t
BBox::height()
{
	return br.y - tl.y + 1;
}

bool
BBox::clear(color * c)
{
	int16_t y = height() - 1;
	int16_t w = width();

	if (c && wattron(window, c->pair()) == ERR) {
		return false;
	}

	while (y >= 0) {
		if (mvwhline(window, y--, 0, ' ', w) == ERR) {
			return false;
		}
	}

	if (c && wattroff(window, c->pair()) == ERR) {
		return false;
	}

	return true;
}

bool
BBox::resize(const Point & tl_, const Point & br_)
{
	int rc;

	if (window != NULL) {
		rc = delwin(window);
		assert(rc != ERR);
		if (rc == ERR) {
			abort();
		}
		window = NULL;
	}

	tl = tl_;
	br = br_;

	window = newwin(height(), width(), tl.y, tl.x);
	assert(window != NULL);

	return (window != NULL);
}

bool
BBox::refresh()
{
	return (wrefresh(window) != ERR);
}

/*
 *
 * Display class
 *
 */

Display::Display(Control * n_comm)
{
	active_list = NULL;
}

Display::~Display()
{
	this->uninit();
}

/*
 * Switch mouse support on or off by setting the mouse mask
 */
mmask_t		Display::setmousemask()
{
	if (pms->options->mouse)
		mmask = mousemask(ALL_MOUSE_EVENTS, &oldmmask);
	else
		mmask = mousemask(0, &oldmmask);

	return mmask;
}

/*
 * Initialize ncurses
 */
bool
Display::init()
{
	/* Fetch most keys and turn off echoing */
	initscr();
	raw();
	noecho();
	keypad(stdscr, true);
	setmousemask();

	if (has_colors()) {
		start_color();
		use_default_colors();
	}

	/* Hide cursor */
	curs_set(0);

	resized();

	return true;
}

/*
 * Delete all windows and end ncurses mode
 */
void
Display::uninit()
{
	vector<List *>::iterator i;

	i = lists.begin();

	while (i != lists.end()) {
		delete *i;
		++i;
	}

	lists.clear();

	endwin();
}

/*
 * Resizes windows.
 */
void
Display::resized()
{
	vector<List *>::iterator iter;

	topbar.resize(Point(0, 0), Point(COLS - 1, pms->options->topbar_lines.size() - 1));
	titlebar.resize(Point(0, topbar.bottom() + 1), Point(COLS - 1, topbar.bottom() + 1));
	main_window.resize(Point(0, titlebar.bottom() + 1), Point(COLS - 1, LINES - 2));
	statusbar.resize(Point(0, main_window.bottom() + 1), Point(COLS - 4, main_window.bottom() + 1));
	position_readout.resize(Point(statusbar.right() + 1, statusbar.bottom()), Point(COLS - 1, statusbar.bottom()));

	iter = lists.begin();
	while (iter != lists.end()) {
		(*iter++)->set_column_size();
	}
}

/*
 * Flushes drawn output to screen for all windows on current screen.
 */
void
Display::refresh()
{
	topbar.refresh();
	titlebar.refresh();
	main_window.refresh();
	statusbar.refresh();
	position_readout.refresh();
}

bool
Display::add_list(List * list)
{
	lists.push_back(list);
	list->set_bounding_box(&main_window);
	return true;
}

bool
Display::activate_list(List * list)
{
	assert(list);

	if (list == active_list) {
		return false;
	}

	last_list = active_list;
	active_list = list;

	return true;
}

List *
Display::find(const char * title)
{
	vector<List *>::iterator i;

	i = lists.begin();
	while (i != lists.end()) {
		if (!strcmp(title, (*i)->title())) {
			return *i;
		}
		++i;
	}

	return NULL;
}

bool
Display::draw_topbar()
{
	Point				p;
	uint32_t			position;
	uint32_t			formatted_length;
	vector<Topbarline *>::iterator	iter;
	Song *				song;
	string				s;

	assert(topbar.height() == pms->options->topbar_lines.size());

	if (!topbar.height()) {
		return false;
	}

	topbar.clear(NULL);

	song = pms->cursong();

	iter = pms->options->topbar_lines.begin();

	while (iter != pms->options->topbar_lines.end()) {

		position = TOPBAR_FIELD_RIGHT + 1;

		while (position-- != TOPBAR_FIELD_LEFT) {

			/* Get a formatted string that can be printed in the topbar. */
			s = pms->formatter->format(song, (*iter)->strings[position], formatted_length, &(pms->options->colors->topbar.fields));

			if (formatted_length && s.size()) {
				if (position == TOPBAR_FIELD_LEFT) {
					p.x = 0;
				} else if (position == TOPBAR_FIELD_CENTER) {
					p.x = (topbar.width() / 2) - (formatted_length / 2);
				} else if (position == TOPBAR_FIELD_RIGHT) {
					p.x = (topbar.width() - formatted_length);
				}

				colprint(&topbar, p.y, p.x, pms->options->colors->topbar.standard, s.c_str());
			}
		}

		++iter;
		++p.y;
	}

	return true;
}

bool
Display::draw_titlebar()
{
	int		left, right;
	string		t;

	assert(active_list);

	titlebar.clear(NULL);

	t = active_list->title();
	if (!t.size()) {
		return false;
	}

	left = (titlebar.width() / 2) - (t.size() / 2);
	right = left + t.size();

	wattron(titlebar.window, pms->options->colors->border->pair());
	mvwaddch(titlebar.window, 0, left - 2, ACS_RTEE);
	mvwaddch(titlebar.window, 0, right + 1, ACS_LTEE);
	wattroff(titlebar.window, pms->options->colors->border->pair());

	wattron(titlebar.window, pms->options->colors->title->pair());
	mvwprintw(titlebar.window, 0, left - 1, " %s ", t.c_str());
	wattroff(titlebar.window, pms->options->colors->title->pair());

	return true;
}

bool
Display::draw_main_window()
{
	assert(active_list);
	return active_list->draw();
}

bool
Display::draw_position_readout()
{
	uint32_t	percent;
	char		buffer[16];

	if (active_list->size() < active_list->bbox->height()) {
		strcpy(buffer, "All");
	} else if (active_list->top_position() == active_list->min_top_position()) {
		strcpy(buffer, "Top");
	} else if (active_list->top_position() == active_list->max_top_position()) {
		strcpy(buffer, "Bot");
	} else {
		percent = 100 * active_list->top_position() / (active_list->size() - active_list->bbox->height() + 1);
		sprintf(buffer, "%2d%%%%", percent);
	}

	/* Clear window */
	position_readout.clear(NULL);

	/* Draw string */
	colprint(&position_readout, 0, 0, pms->options->colors->position, buffer);

	return true;
}

/*
 * Redraws all visible windows
 */
bool
Display::draw()
{
	/* FIXME: check if drawing is needed */
	draw_topbar();
	draw_titlebar();
	draw_main_window();
	draw_position_readout();
	/* FIXME: what about statusbar? */
	return true;
}

/*
 * Redraws all visible windows regardless of state
 */
void		Display::forcedraw()
{
	assert(false);
	/*
	topbar->draw();
	statusbar->draw();
	positionreadout->draw();
	if (curwin) curwin->draw();
	*/
}

/*
 * Set XTerm window title.
 *
 * The current xterm title exists under the WM_NAME property,
 * and can be retrieved with `xprop -notype -id $WINDOWID WM_NAME`.
 */
void
Display::set_xterm_title()
{
	unsigned int	reallen;
	string		title;
	ostringstream	oss;

	if (!pms->options->xtermtitle.size()) {
		return;
	}

	if (getenv("WINDOWID")) {
		title = pms->formatter->format(pms->cursong(), pms->options->xtermtitle, reallen, NULL, true);
		pms->log(MSG_DEBUG, 0, _("Set XTerm window title: '%s'\n"), title.c_str());

		oss << "\033]0;" << title << '\007';
		putp(oss.str().c_str());

		/* stdout is in line buffered mode be default and thus needs
		 * explicit flush to communicate with terminal successfully. */
		fflush(stdout);
	} else {
		pms->log(MSG_DEBUG, 0, _("Disabling XTerm window title: WINDOWID not found.\n"));
		pms->options->xtermtitle = "";
	}
}

Song *
Display::cursorsong()
{
	Songlist * songlist;

	if ((songlist = SONGLIST(active_list)) == NULL) {
		return NULL;
	}

	return songlist->cursorsong();
}




/*
 *
 * End of display class.
 *
 * Playlist column class
 *
 */
pms_column::pms_column(string n_title, Item n_type, unsigned int n_minlen)
{
	title	= n_title;
	type	= n_type;
	minlen	= n_minlen;
	abslen	= -1;
	median	= 0;
	items	= 0;
}

void		pms_column::addmedian(unsigned int n)
{
	++items;
	median += n;
	abslen = -1;
}

unsigned int	pms_column::len()
{
	if (abslen < 0)
	{
		if (items == 0)
			abslen = 0;
		else
			abslen = (median / items);
	}
	if ((unsigned int)abslen < minlen) {
		return minlen;
	}

	return (unsigned int)abslen;
}


/*
 * Prints formatted output onto a window. Borders are handled correctly.
 *
 * %s		= char *
 * %d		= int
 * %f		= double
 * %B %/B	= bold on/off
 * %R %/R	= reverse on/off
 * %0-n% %/0-n%	= color on/off
 *
 */
void colprint(BBox * bbox, int y, int x, color * c, const char *fmt, ...)
{
	va_list			ap;
	unsigned int		i = 0;
	double			f = 0;
	string			output = "";
	bool			parse = false;
	bool			attr = false;
	attr_t			attrval = 0;
	char			buf[1024];
	string			colorstr;
	int			colorint;
	int			pair = 0;
	unsigned int		maxlen;		// max allowed characters printed on screen
	unsigned int		printlen = 0;	// num characters printed on screen

	assert(bbox);

	va_start(ap, fmt);

	/* Check if string is out of range, and cuts if necessary */
	if (x < 0)
	{
		if (strlen(fmt) < abs(x))
			return;

		fmt += abs(x);
		x = 0;
	}

	if (c != NULL)
		pair = c->pair();

	wmove(bbox->window, y, x);
	wattron(bbox->window, pair);

	maxlen = bbox->width() - x + 1;

	while(*fmt && printlen < maxlen)
	{
		if (*fmt == '%' && !parse)
		{
			if (*(fmt + 1) == '%')
			{
				fmt += 2;
				output = "%%";
				wprintw(bbox->window, _(output.c_str()));
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
						wattron(bbox->window, A_BOLD);
					else
						wattroff(bbox->window, A_BOLD);
					parse = false;
					break;
				case 'R':
					if (attr)
						wattron(bbox->window, A_REVERSE);
					else
						wattroff(bbox->window, A_REVERSE);
					parse = false;
					break;
				case 'd':
					parse = false;
					i = va_arg(ap, int);
					sprintf(buf, "%d", i);
					wprintw(bbox->window, _(buf));
					printlen += strlen(buf);
					i = 0;
					break;
				case 'f':
					parse = false;
					f = va_arg(ap, double);
					sprintf(buf, "%f", f);
					wprintw(bbox->window, _(buf));
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
					wprintw(bbox->window, _(buf));
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
							wattroff(bbox->window, pair);
							wattron(bbox->window, i);
						}
						else
						{
							wattroff(bbox->window, i);
							wattron(bbox->window, pair);
						}

						/* Skip characters */
						colorint = static_cast<int>(i);
						colorstr = Pms::tostring(colorint);
						fmt += (colorstr.size());
					}
					parse = false;
					break;
			}
		}
		else
		{
			output = *fmt;
			wprintw(bbox->window, _(output.c_str()));
			++printlen;
		}
		++fmt;
	}

	va_end(ap);
	wattroff(bbox->window, pair);
}


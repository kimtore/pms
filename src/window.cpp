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

#include "window.h"
#include "curses.h"
#include "config.h"

extern Windowmanager wm;
extern Curses curses;
extern Config config;

void Window::draw()
{
	int i;

	if (!rect || !visible())
		return;

	for (i = 0; i <= rect->bottom - rect->top; i++)
		drawline(i);
}

void Window::clear()
{
	curses.wipe(rect, config.colors.standard);
}

unsigned int Window::height()
{
	if (!rect) return 0;
	return rect->bottom - rect->top;
}

Wmain::Wmain()
{
	position = 0;
	cursor = 0;
}

void Wmain::draw()
{
	Window::draw();
	wm.readout->draw();
}

void Wmain::scroll_window(int offset)
{
	int limit = static_cast<int>(content_size() - rect->bottom - rect->top + 1);

	if (limit < 0)
		limit = 0;

	offset = position + offset;

	if (offset < 0)
	{
		offset = 0;
		curses.bell();
	}
	if (offset > limit)
	{
		offset = limit;
		curses.bell();
	}
	
	if ((int)position == offset)
		return;

	position = offset;

	if (cursor < position)
		cursor = position;
	else if (cursor > position + height())
		cursor = position + height();
	
	draw();
	curses.flush();
}

void Wmain::set_position(unsigned int absolute)
{
	scroll_window(absolute - position);
}

void Wmain::move_cursor(int offset)
{
	offset = cursor + offset;

	if (offset < 0)
	{
		offset = 0;
		curses.bell();
	}
	else if (offset > (int)content_size() - 1)
	{
		offset = content_size() - 1;
		curses.bell();
	}

	if ((int)cursor == offset)
		return;

	cursor = offset;

	if (cursor < position)
		set_position(cursor);
	else if (cursor > position + height())
		set_position(cursor - height());
	
	draw();
	curses.flush();
}

void Wmain::set_cursor(unsigned int absolute)
{
	move_cursor(absolute - cursor);
}

bool Wmain::visible()
{
	return wm.active == this;
}

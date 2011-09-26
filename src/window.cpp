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

extern Windowmanager wm;
extern Curses curses;

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
	curses.wipe(rect);
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

void Wmain::scroll_window(int offset)
{
	int limit = static_cast<int>(content_size() - rect->bottom - rect->top + 1);

	offset = position + offset;

	if (offset < 0)
		offset = 0;
	if (offset > limit)
		offset = limit;
	
	position = offset;

	if (cursor < position)
		cursor = position;
	else if (cursor > position + height())
		cursor = position + height();
	
	wm.readout->draw();
	if (visible()) draw();
}

void Wmain::set_position(unsigned int absolute)
{
	position = absolute;
	scroll_window(0);
}

void Wmain::move_cursor(int offset)
{
	offset = cursor + offset;

	if (offset < 0)
		offset = 0;
	else if (offset > (int)content_size() - 1)
		offset = content_size() - 1;

	cursor = offset;

	if (cursor < position)
		set_position(cursor);
	else if (cursor > position + height())
		set_position(cursor - height());
	
	wm.readout->draw();
	if (visible()) draw();
}

void Wmain::set_cursor(unsigned int absolute)
{
	cursor = absolute;
	move_cursor(0);
}

bool Wmain::visible()
{
	return wm.active == this;
}

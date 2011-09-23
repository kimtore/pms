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

extern Windowmanager wm;

void Window::draw()
{
	int i;

	if (!rect || !visible())
		return;

	for (i = 0; i <= rect->bottom - rect->top; i++)
		drawline(i);
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

	if (offset > limit)
		offset = limit;
	if (offset < 0)
		offset = 0;
	
	position = offset;

	wm.readout->draw();
	if (visible()) draw();
}

void Wmain::set_position(unsigned int absolute)
{
	int limit = static_cast<int>(content_size() - rect->bottom - rect->top + 1);
	int abs = absolute;

	if (limit < 0)
		abs = 0;
	else if (abs > limit || abs < 0)
		abs = limit;

	position = abs;

	wm.readout->draw();
	if (visible()) draw();
}

void Wmain::move_cursor(int offset)
{
	cursor += offset;

	if (cursor < 0)
		cursor = 0;
	else if (cursor > content_size() - 1)
		cursor = content_size() - 1;
	
	wm.readout->draw();
	if (visible()) draw();
}

void Wmain::set_cursor(unsigned int absolute)
{
	set_position(absolute);
}

bool Wmain::visible()
{
	return wm.active == this;
}

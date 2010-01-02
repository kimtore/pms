/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2010  Kim Tore Jensen
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
 * mediator.h - track changes everywhere independent of systems
 *
 */


#include <vector>
#include <string>
#include "mediator.h"

using namespace std;


/*
 * Add an object to the changed list
 */
void		Mediator::add(string key)
{
	vector<string>::iterator	i;

	i = keys.begin();
	while (i != keys.end())
	{
		if (*i == key)
			return;
		++i;
	}

	keys.push_back(key);
}

/*
 * Determine if an object is in the list, and remove it afterwards
 */
bool		Mediator::changed(string key)
{
	vector<string>::iterator	i;

	i = keys.begin();
	while (i != keys.end())
	{
		if (*i == key)
		{
			keys.erase(i);
			return true;
		}
		++i;
	}

	return false;
}

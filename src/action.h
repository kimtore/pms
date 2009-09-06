/* vi:set ts=8 sts=8 sw=8:
 *
 * PMS  <<Practical Music Search>>
 * Copyright (C) 2006-2009  Kim Tore Jensen
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
 * action.h
 * 	Executes key-bound actions
 */

#ifndef _ACTION_H_
#define _ACTION_H_

#include "list.h"
#include "display.h"


bool		handle_command(pms_pending_keys);
bool		init_commandmap();

void		generr();
int		playnext(pms_play_mode, int);
song_t		gotonextentry(string, bool);
int		multiplay(long, int);
bool		setwin(pms_window *);
bool		makeselection(Songlist *, pms_pending_keys, string);
int		removesongs(Songlist *);
int		createwindow(string, pms_window *&, Songlist *&);

#endif /* _ACTION_H_ */

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
#include "message.h"


class Interface
{
private:
	Message *		msg;

public:
				Interface();
				~Interface();

	bool			check_events();

	//FIXME: REMOVE
	int			action;
	string			param;


	/*
	 * PMS specific stuff
	 */
	long			exec(string);
	long			version();
	long			clear_topbar(int);
	long			redraw();
	long			rehash();
	long			write_config(string);
	long			source(string);
	long			quit();
	long			shell(string);
	long			show_info();

	/*
	 * MPD administrative
	 */
	long			password(string);
	long			update_db(string);

	long			set_input_mode(int); //command mode, jump mode, search mode

	long			toggle();
	long			escape();

	long			next_of();
	long			prev_of();
	long			goto_random();
	long			goto_current();
	long			next_result();
	long			prev_result();

	long			text_updated();
	long			text_return();
	long			text_escape();

	/*
	 * Normal player actions
	 */
	long			setvolume(string);
	long			mute();
	long			crossfade(int);
	long			seek(int);



	long			play();	// play of type, too: playartist, playalbum, playrandom, etc
	long			add(); // play + add to, add all
	long			next(bool);
	long			prev();
	long			pause(bool);
	long			stop();
	long			remove();
	long			move();
	long			shuffle();
	long			clear();
	long			crop(); // crop playing and crop selection
	long			select(); //FIXME
	long			unselect();
	long			toggleselect();
	long			clearselection();

	long			cycle_playmode();
	long			cycle_repeatmode();

	long			create_playlist();
	long			save_playlist();
	long			delete_playlist();

	long			next_window();
	long			prev_window();
	long			change_window();
	long			last_window();

	long			show_bindings();
	long			activate_playlist();

	long			move_cursor(); // move up/down, pgup/pgdn, center, etc
	long			scroll_window();
};


bool		handle_command(pms_pending_keys);
bool		init_commandmap();

void		generr();
int		playnext(long, int);
song_t		gotonextentry(string, bool);
int		multiplay(long, int);
bool		setwin(pms_window *);
bool		makeselection(Songlist *, pms_pending_keys, string);
int		removesongs(Songlist *);
int		createwindow(string, pms_window *&, Songlist *&);

#endif /* _ACTION_H_ */

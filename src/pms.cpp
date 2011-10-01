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

#include "pms.h"
#include "console.h"
#include "curses.h"
#include "config.h"
#include "window.h"
#include "command.h"
#include "mpd.h"
#include "input.h"
#include <stdlib.h>

extern Config		config;
extern MPD		mpd;
extern Curses		curses;
extern Windowmanager	wm;
extern Input		input;
extern Commandlist 	commandlist;

int PMS::run_event(Inputevent * ev)
{
	if (!ev) return false;

	if (ev->result == INPUT_RESULT_RUN)
	switch(ev->action)
	{
		case ACT_MODE_INPUT:
			input.setmode(INPUT_MODE_INPUT);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_MODE_COMMAND:
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_RUN_CMD:
			run_cmd(ev->text);
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_QUIT:
			return quit();

		case ACT_RESIZE:
			curses.detect_dimensions();
			wm.playlist->update_column_length();
			wm.library->update_column_length();
			wm.draw();
			curses.flush();
			return true;

		case ACT_NEXT_WINDOW:
			return cycle_windows(ev->multiplier);

		case ACT_PREVIOUS_WINDOW:
			return cycle_windows(-ev->multiplier);

		case ACT_SCROLL_UP:
			return scroll_window(-ev->multiplier);

		case ACT_SCROLL_DOWN:
			return scroll_window(ev->multiplier);

		case ACT_CURSOR_UP:
			return move_cursor(-ev->multiplier);

		case ACT_CURSOR_DOWN:
			return move_cursor(ev->multiplier);

		case ACT_CURSOR_PGUP:
			return move_cursor_page(-ev->multiplier);

		case ACT_CURSOR_PGDOWN:
			return move_cursor_page(ev->multiplier);

		case ACT_CURSOR_TOP:
			return set_cursor_top();

		case ACT_CURSOR_BOTTOM:
			return set_cursor_bottom();

		case ACT_CURSOR_HOME:
			return set_cursor_home();

		case ACT_CURSOR_END:
			return set_cursor_end();

		case ACT_CURSOR_CURRENTSONG:
			return set_cursor_currentsong();

		case ACT_CONSUME:
			return set_consume(!mpd.status.consume);

		case ACT_CROSSFADE:
			return set_crossfade(ev->text);

		case ACT_RANDOM:
			return set_random(!mpd.status.random);

		case ACT_REPEAT:
			return set_repeat(!mpd.status.repeat);

		case ACT_SETVOL:
			return set_volume(ev->text);

		case ACT_SINGLE:
			return set_single(!mpd.status.single);

		case ACT_TOGGLEPLAY:
			return toggle_play();

		case ACT_PLAY:
			return play();

		default:
			return false;
	}

	else if (ev->result == INPUT_RESULT_BUFFERED)
	{
		wm.statusbar->draw();
		curses.flush();
	}


	return false;
}

int PMS::run_cmd(string cmd)
{
	Inputevent ev;
	Command * c;
	size_t i;

	if ((i = cmd.find(' ')) != string::npos)
	{
		cmd = cmd.substr(0, i);
		ev.text = cmd.size() > i ? cmd.substr(i + 1) : "";
	}

	c = commandlist.find(wm.context, cmd);
	if (!c)
	{
		sterr("Undefined command '%s'", cmd.c_str());
		return false;
	}

	ev.action = c->action;
	ev.context = wm.context;
	ev.result = INPUT_RESULT_RUN;
	return run_event(&ev);
}

int PMS::quit()
{
	config.quit = true;
	return true;
}

int PMS::scroll_window(int offset)
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->scroll_window(offset);
	return true;
}

int PMS::move_cursor(int offset)
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->move_cursor(offset);
	return true;
}

int PMS::move_cursor_page(int offset)
{
	bool beep;
	Wmain * window;
	window = WMAIN(wm.active);

	beep = config.use_bell;
	config.use_bell = false;
	window->move_cursor(offset * window->height());

	if (offset < 0)
		window->set_position(window->cursor - window->height());
	else
		window->set_position(window->cursor);

	config.use_bell = beep;
	return true;
}

int PMS::set_cursor_top()
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->set_cursor(window->position);
	return true;
}

int PMS::set_cursor_bottom()
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->set_cursor(window->position + window->height());
	return true;
}

int PMS::set_cursor_home()
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->set_cursor(0);
	return true;
}

int PMS::set_cursor_end()
{
	Wmain * window;
	window = WMAIN(wm.active);
	window->set_cursor(window->content_size() - 1);
	return true;
}

int PMS::set_cursor_currentsong()
{
	Wsonglist * window;
	Song * song;
	size_t pos;

	/* Silently ignore if there is no song playing. */
	if ((song = mpd.currentsong) == NULL)
		return false;

	/* Get current window */
	window = WSONGLIST(wm.active);
	if (window == NULL)
	{
		sterr("Current window is not a playlist, so cannot locate any songs here.", NULL);
		return false;
	}

	/* If the song has a position, and we are in the playlist, jump to that spot. */
	if (song->pos != -1 && window->songlist->playlist && song->pos < (int)window->songlist->size())
	{
		pos = song->pos;
	}

	/* Use song hash to look it up. */
	else if ((pos = window->songlist->find(song->fhash)) == string::npos)
	{
		sterr("Currently playing song is not in this songlist.", NULL);
		return false;
	}

	window->set_cursor(pos);
	return true;
}

int PMS::cycle_windows(int offset)
{
	wm.cycle(offset);
	return true;
}

int PMS::set_consume(bool consume)
{
	return mpd.set_consume(consume);
}

int PMS::set_crossfade(string crossfade)
{
	return mpd.set_crossfade(atoi(crossfade.c_str()));
}

int PMS::set_random(bool random)
{
	return mpd.set_random(random);
}

int PMS::set_repeat(bool repeat)
{
	return mpd.set_repeat(repeat);
}

int PMS::set_volume(string volume)
{
	if (volume.size() == 0)
		return false;
	
	if (volume[0] == '+')
		return mpd.set_volume(mpd.status.volume + atoi(volume.substr(1).c_str()));
	else if (volume[0] == '-')
		return mpd.set_volume(mpd.status.volume - atoi(volume.substr(1).c_str()));
	else
		return mpd.set_volume(atoi(volume.c_str()));
}

int PMS::set_single(bool single)
{
	return mpd.set_single(single);
}

int PMS::toggle_play()
{
	return mpd.pause(mpd.status.state == MPD_STATE_PLAY ? true : false);
}

int PMS::play()
{
	Song * s;
	if ((s = cursorsong()) == NULL)
		return false;

	if (s->id != -1)
		return (mpd.playid(s->id) == MPD_GETLINE_OK);
	else
		return (mpd.playid(mpd.addid(s->f[FIELD_FILE])) == MPD_GETLINE_OK);
}




Song * cursorsong()
{
	Wsonglist * win;
	win = WSONGLIST(wm.active);
	if (win == NULL)
		return NULL;
	return win->cursorsong();
}

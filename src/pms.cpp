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
#include "songlist.h"
#include "search.h"
#include <stdlib.h>

extern Config		config;
extern MPD		mpd;
extern Curses		curses;
extern Windowmanager	wm;
extern Input		input;
extern Commandlist 	commandlist;

int PMS::run_event(Inputevent * ev)
{
	static Inputevent lastev;
	Inputevent sev;

	if (!ev) return false;

	if (ev->result == INPUT_RESULT_RUN && ev->action != ACT_REPEATACTION && ev->action != ACT_RUN_CMD && ev != &lastev)
		lastev = *ev;

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

		case ACT_MODE_SEARCH:
			input.setmode(INPUT_MODE_SEARCH);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_REHASH:
			sev = *ev;
			config.source_default_config();
			lastev = sev;
			return true;

		case ACT_SOURCE:
			sev = *ev;
			config.source(ev->text);
			lastev = sev;
			return true;

		case ACT_RUN_CMD:
			run_cmd(ev->text, ev->multiplier);
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->draw();
			curses.flush();
			return true;

		case ACT_RUN_SEARCH:
			run_search(ev->text, ev->multiplier);
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->draw();
			wm.active->draw();
			curses.flush();
			return true;

		case ACT_REPEATACTION:
			if (ev->multiplier != 1)
				lastev.multiplier = ev->multiplier;
			return run_event(&lastev);

		case ACT_SET:
			return set_opt(ev);

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

		case ACT_ACTIVATE_SONGLIST:
			return activate_songlist();

		case ACT_ADD:
			return add(ev->multiplier);

		case ACT_REMOVE:
			return remove(ev->multiplier);

		case ACT_UPDATE:
			return update(ev->text);

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

		case ACT_CURSOR_RANDOM:
			return set_cursor_random();

		case ACT_CROSSFADE:
			return set_crossfade(ev->text);

		case ACT_PASSWORD:
			return set_password(ev->text);

		case ACT_TOGGLEPLAY:
			return toggle_play();

		case ACT_PLAY:
			return play();

		case ACT_STOP:
			return stop();

		case ACT_NEXT:
			return change_song(ev->multiplier);

		case ACT_PREVIOUS:
			return change_song(-ev->multiplier);

		case ACT_SEEK_FORWARD:
			return seek(ev->multiplier);

		case ACT_SEEK_BACK:
			return seek(-ev->multiplier);

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

int PMS::run_cmd(string cmd, unsigned int multiplier, bool batch)
{
	Inputevent ev;
	Command * c;
	size_t i;

	/* Strip whitespace from beginning */
	if ((i = cmd.find_first_not_of(' ')) > 0)
	{
		if (cmd.size() > i)
			cmd = cmd.substr(i);
		else
			return false;
	}

	/* Separate command and param */
	if ((i = cmd.find(' ')) != string::npos)
	{
		ev.text = cmd.size() > i ? cmd.substr(i + 1) : "";
		cmd = cmd.substr(0, i);
	}

	/* Ignore comments and empty lines */
	if (!cmd.size() || cmd[0] == '#')
		return false;

	c = commandlist.find(wm.context, cmd);
	if (!c)
	{
		sterr("Undefined command '%s'", cmd.c_str());
		return false;
	}

	ev.action = c->action;
	ev.context = wm.context;
	ev.result = INPUT_RESULT_RUN;
	ev.silent = batch;
	ev.multiplier = multiplier > 0 ? multiplier : 1;
	return run_event(&ev);
}

int PMS::run_search(string terms, unsigned int multiplier)
{
	Wsonglist * window;
	Song * song;
	size_t i;

	if ((window = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a playlist, so cannot locate any songs here.", NULL);
		return false;
	}

	song = window->cursorsong();

	if (terms.size() > 0)
		window->songlist->search(SEARCH_MODE_FILTER, config.search_field_mask, terms);
	else
		window->songlist->search(SEARCH_MODE_NONE);

	window->set_cursor(0);

	if (song && (i = window->songlist->sfind(song->fhash)) != string::npos)
		window->set_cursor(i);

	return true;
}

int PMS::set_opt(Inputevent * ev)
{
	option_t * opt;

	opt = config.readline(ev->text, !ev->silent);
	if (!opt)
		return false;

	if (opt->mask & OPT_CHANGE_MPD)
		mpd.apply_opts();
	if (opt->mask & OPT_CHANGE_DIMENSIONS)
		curses.detect_dimensions();

	if (opt->mask & OPT_CHANGE_REDRAW)
		wm.draw();
	else if (opt->mask & OPT_CHANGE_DRAWLIST)
		wm.active->draw();
	else
		return true;

	curses.flush();

	return true;
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
	else if ((pos = window->songlist->sfind(song->fhash)) == string::npos)
	{
		sterr("Currently playing song is not in this songlist.", NULL);
		return false;
	}

	window->set_cursor(pos);
	return true;
}

int PMS::set_cursor_random()
{
	Wsonglist * window;

	/* Get current window */
	window = WSONGLIST(wm.active);
	if (window == NULL)
	{
		sterr("Current window is not a playlist, so cannot locate any songs here.", NULL);
		return false;
	}

	if (window->songlist->size() == 0)
		return false;
	
	window->set_cursor(window->songlist->randpos());
	return true;
}

int PMS::cycle_windows(int offset)
{
	wm.cycle(offset);
	return true;
}

int PMS::activate_songlist()
{
	Wsonglist * win;
	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a playlist, and cannot be set as the primary playback list.", NULL);
		return false;
	}

	return mpd.activate_songlist(win->songlist);
}

int PMS::add(int count)
{
	bool status = true;
	int c = count;
	Wsonglist * win;
	vector<Song *>::iterator song;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a playlist. Cannot add any songs from here.", NULL);
		return false;
	}

	song = win->songlist->songs.begin() + win->cursor;
	while (c > 0 && song != win->songlist->songs.end())
	{
		status = status && mpd.addid((*song)->f[FIELD_FILE]);
		--c;
		++song;
	}
	count -= c;

	if (status && count >= 1)
	{
		if (count > 1)
			stinfo("%d songs added to playlist.", count);
		else if (count == 1)
			stinfo("`%s' added to playlist.", (*--song)->f[FIELD_TITLE].c_str());

		if (config.advance_cursor)
			win->move_cursor(count);
	}
	else
	{
		stinfo("Failed to add some songs to playlist.", NULL);
	}

	return status;
}

int PMS::remove(int count)
{
	Wsonglist * win;
	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a playlist. Cannot remove any songs from here.", NULL);
		return false;
	}
	if (win->songlist->readonly)
	{
		sterr("This playlist is read-only.", NULL);
		return false;
	}

	if (mpd.remove(win->songlist, win->cursor, count))
	{
		win->move_cursor(0);
		return true;
	}

	return false;
}

int PMS::update(string dir)
{
	return mpd.update(dir);
}

int PMS::set_crossfade(string crossfade)
{
	return mpd.set_crossfade(atoi(crossfade.c_str()));
}

int PMS::set_password(string password)
{
	int i;
	if ((i = mpd.set_password(password)) == true)
		mpd.apply_opts(); /* any option desynch should be fixed here. */
	return i;
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

int PMS::stop()
{
	return mpd.stop();
}

int PMS::change_song(int steps)
{
	int s;
	while (steps < 0 && (s = mpd.previous()))
		++steps;
	while (steps > 0 && (s = mpd.next()))
		--steps;

	return s;
}

int PMS::seek(int seconds)
{
	return mpd.seek(seconds);
}



Song * cursorsong()
{
	Wsonglist * win;
	win = WSONGLIST(wm.active);
	if (win == NULL)
		return NULL;
	return win->cursorsong();
}

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
#include "clipboard.h"
#include <stdlib.h>

extern Config		config;
extern MPD		mpd;
extern Curses		curses;
extern Windowmanager	wm;
extern Input		input;
extern Commandlist 	commandlist;
extern Keybindings 	keybindings;

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
			wm.statusbar->qdraw();
			return true;

		case ACT_MODE_COMMAND:
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->qdraw();
			return true;

		case ACT_MODE_SEARCH:
			input.setmode(INPUT_MODE_SEARCH);
			wm.statusbar->qdraw();
			return true;

		case ACT_MODE_LIVESEARCH:
			input.setmode(INPUT_MODE_LIVESEARCH);
			wm.statusbar->qdraw();
			return true;

		case ACT_VISUAL:
			return visual();

		case ACT_EXIT_LIVESEARCH:
			return livesearch(ev->text, true);

		case ACT_RESET_SEARCH:
			return resetsearch();

		case ACT_CONNECT:
			config.autoconnect = true;
			mpd.mpd_disconnect();
			return true;

		case ACT_REHASH:
			sev = *ev;
			config.load_default_config();
			config.source_default_config();
			lastev = sev;
			return true;

		case ACT_SOURCE:
			sev = *ev;
			config.source(ev->text);
			lastev = sev;
			return true;

		case ACT_MAP:
			return map_keys(ev->text);

		case ACT_UNMAP:
			return keybindings.remove(ev->text);

		case ACT_RUN_CMD:
			run_cmd(ev->text, ev->multiplier);
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->qdraw();
			return true;

		case ACT_RUN_SEARCH:
			run_search(ev->text, ev->multiplier);
			input.setmode(INPUT_MODE_COMMAND);
			wm.statusbar->qdraw();
			wm.active->qdraw();
			wm.topbar->qdraw();
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
			return true;

		case ACT_NEXT_WINDOW:
			return cycle_windows(ev->multiplier);

		case ACT_PREVIOUS_WINDOW:
			return cycle_windows(-ev->multiplier);

		case ACT_TOGGLE_WINDOW:
			return wm.toggle();

		case ACT_GOTO_WINDOW:
			return wm.go(ev->text);

		case ACT_GOTO_WINDOW_POS:
			return wm.go((ev->text.size() ? atoi(ev->text.c_str()) : ev->multiplier) - 1);

		case ACT_SORT:
			return sortlist(ev->text.size() ? ev->text : config.default_sort);

		case ACT_ACTIVATE_SONGLIST:
			return activate_songlist();

		case ACT_ADD:
			if (ev->text.size())
				return add(ev->text, ev->multiplier);
			return add(ev->multiplier);

		case ACT_REMOVE:
			return remove(ev->multiplier);

		case ACT_YANK:
			return yank(ev->multiplier);

		case ACT_PUT:
			return put(ev->multiplier);

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
		wm.topbar->qdraw();

		if (input.mode == INPUT_MODE_LIVESEARCH)
			return livesearch(ev->text);
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

int PMS::map_keys(string params)
{
	Command * cmd;
	size_t start = 0, end = 0;
	int context;
	vector<string> s;

	while (s.size() <= 3)
	{
		if ((end = params.find(' ', start)) == string::npos)
		{
			if (s.size() >= 2)
			{
				s.push_back(params.substr(start));
				break;
			}

			sterr("Syntax: map <context> <sequence> <command>", NULL);
			return false;
		}
		s.push_back(params.substr(start, end - start));
		start = end + 1;
	}

	if (s[0] == "console")
		context = CONTEXT_CONSOLE;
	else if (s[0] == "songlist")
		context = CONTEXT_SONGLIST;
	else if (s[0] == "list")
		context = CONTEXT_LIST;
	else if (s[0] == "all")
		context = CONTEXT_ALL;
	else
	{
		sterr("Invalid context `%s', expecteded one of `console', `songlist', `list', `all'", s[0].c_str());
		return false;
	}

	if ((cmd = commandlist.find(context, s[2])) == NULL)
	{
		sterr("Invalid command `%s' for context `%s'", s[2].c_str(), s[0].c_str());
		return false;
	}

	return (keybindings.add(context, cmd->action, s[1], (s.size() == 4 ? s[3] : "")) != NULL);

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
	if (opt->mask & OPT_CHANGE_PLAYMODE)
	{
		mpd.update_playstring();
		wm.statusbar->draw();
	}

	if (opt->mask & OPT_CHANGE_REDRAW)
	{
		wm.draw();
	}
	else
	{
		if (opt->mask & OPT_CHANGE_DRAWLIST)
			wm.active->draw();
		if (opt->mask & OPT_CHANGE_TOPBAR)
			wm.topbar->draw();

		if (!(opt->mask & OPT_CHANGE_DRAWLIST) && !(opt->mask & OPT_CHANGE_TOPBAR))
			return true;
	}

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
		sterr("Current window is not a song list, and cannot be set as the primary playback list.", NULL);
		return false;
	}

	return mpd.activate_songlist(win->songlist);
}

int PMS::sortlist(string sortstr)
{
	Wsonglist * win;
	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a song list, and cannot be sorted.", NULL);
		return false;
	}

	win->songlist->sort(sortstr);
	win->qdraw();
	return true;
}

int PMS::livesearch(string terms, bool exitsearch)
{
	Wsonglist * win;
	Song * song;
	song_t i;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a song list, and cannot be searched.", NULL);
		return false;
	}

	song = win->cursorsong();
	win->songlist->search(SEARCH_MODE_LIVE, config.search_field_mask, terms);
	if (exitsearch)
	{
		win->songlist->liveclear();
		input.setmode(INPUT_MODE_COMMAND);
		wm.statusbar->qdraw();
		curses.flush();
	}

	win->set_cursor(0);
	if (song && (i = win->songlist->sfind(song->fhash)) != string::npos)
		win->set_cursor(i);

	win->qdraw();
	return true;
}

int PMS::resetsearch()
{
	Wsonglist * win;
	Song * song;
	song_t i;

	if ((win = WSONGLIST(wm.active)) == NULL)
		return false;
	
	song = win->cursorsong();
	win->songlist->search(SEARCH_MODE_NONE);

	win->set_cursor(0);
	if (song && (i = win->songlist->sfind(song->fhash)) != string::npos)
		win->set_cursor(i);

	win->qdraw();
	return true;
}

int PMS::add(int count)
{
	bool status = true;
	size_t i;
	int added = 0;
	Wsonglist * win;
	selection_t sel;
	vector<Song *>::iterator it;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a playlist. Cannot add any songs from here.", NULL);
		return false;
	}

	sel = win->get_selection(count);
	if (sel->empty())
		return false;

	for (it = sel->begin(); it != sel->end(); ++it)
		if (mpd.addid((*it)->f[FIELD_FILE]))
			++added;

	if (added > 0)
	{
		if (added > 1)
			stinfo("%d songs added to playlist.", added);
		else
			stinfo("`%s' added to playlist.", (*--it)->f[FIELD_TITLE].c_str());

		if (config.advance_cursor)
			win->move_cursor(count);
	}
	else
	{
		stinfo("Failed to add some songs to playlist.", NULL);
	}

	win->songlist->clear_visual();

	return status;
}

int PMS::add(string uri, int count)
{
	bool status = true;

	while (count-- > 0)
		if (mpd.addid(uri) != -1)
			stinfo("`%s' added to playlist.", uri.c_str());
		else
			status = false;
	
	return status;
}

int PMS::remove(int count)
{
	Wsonglist * win;
	selection_t sel;
	size_t start;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Current window is not a song list. Cannot remove any songs from here.", NULL);
		return false;
	}
	if (win->songlist->readonly)
	{
		sterr("This song list is read-only.", NULL);
		return false;
	}

	sel = win->get_selection(count);

	if (sel->empty())
		return false;

	if ((count = mpd.remove(win->songlist, sel)) > 0)
	{
		win->songlist->visual_pos(&start, NULL);
		if (start == -1)
			start = win->cursor;
		if (start >= win->songlist->size() - sel->size())
			--start;
		win->set_cursor(start);
		clipboard.set(sel);
		win->songlist->clear_visual();
		return true;
	}

	return false;
}

int PMS::visual()
{
	Wsonglist * win;
	int pos;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Visual mode can only be used in songlists.", NULL);
		return false;
	}

	if (win->songlist->visual_start == -1)
		pos = win->cursor;
	else
		pos = -1;

	win->songlist->visual_start = pos;
	win->songlist->visual_stop = pos;
	win->qdraw();

	return true;
}

int PMS::yank(int count)
{
	Wsonglist * win;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Only songs can be yanked, and this is not a song list.", NULL);
		return false;
	}

	clipboard.set(win->get_selection(count));
	win->songlist->clear_visual();
	win->qdraw();
	stinfo("%d songs yanked to clipboard.", clipboard.size());

	if (config.advance_cursor)
		win->move_cursor(count);

	return true;
}

int PMS::put(int count)
{
	Wsonglist * win;

	if ((win = WSONGLIST(wm.active)) == NULL)
	{
		sterr("Can only put songs into a song list.", NULL);
		return false;
	}

	if (clipboard.empty())
	{
		sterr("Clipboard is empty.", NULL);
		return false;
	}

	if (win->songlist->readonly)
	{
		sterr("Cannot put songs into a read-only song list.", NULL);
		return false;
	}

	if (mpd.put(win->songlist, win->cursor + 1, &clipboard.songs))
	{
		stinfo("Put %d songs into song list.", clipboard.size());
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
	Song * song;

	if (mpd.status.single && mpd.status.repeat && mpd.currentsong)
		return mpd.playid(mpd.currentsong->id);

	if ((song = mpd.next_song_in_line(steps)) == NULL)
	{
		sterr("There is no candidate for %s song, try switching lists or setting repeat mode.", string(steps > 0 ? "next" : "previous").c_str());
		return false;
	}

	if (song->id != -1)
		return mpd.playid(song->id);
	else
		return mpd.playid(mpd.addid(song->f[FIELD_FILE]));
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

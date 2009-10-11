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
 * action.cpp
 * 	Executes key-bound actions
 */

#include "action.h"
#include "command.h"
#include "display.h"
#include "config.h"
#include "input.h"
#include "i18n.h"
#include "pms.h"

extern Pms *			pms;



Interface::Interface()
{
	msg = new Message();
}

Interface::~Interface()
{
	delete msg;
}


/*
 * Handle any events.
 * This is a frontend to all the other functions in this class.
 */
bool		Interface::check_events()
{
	msg->clear();
	action = pms->input->getpending();

	switch(action)
	{
		default:
		case PEND_NONE:
			return false;

		/*
		 * PMS specific stuff
		 */

		case PEND_EXEC:
			exec(param);
			break;

		case PEND_VERSION:
			version();
			break;

		case PEND_CLEAR_TOPBAR:
			clear_topbar(atoi(param.c_str()));
			break;

		case PEND_REDRAW:
			redraw();
			break;

		case PEND_REHASH:
			rehash();
			break;

		case PEND_WRITE_CONFIG:
			write_config(param);
			break;

		case PEND_SOURCE:
			source(param);
			break;

		case PEND_QUIT:
			quit();
			break;
		
		case PEND_SHELL:
			shell(param);
			break;

		case PEND_SHOW_INFO:
			show_info();
			break;

		/*
		 * MPD admin
		 */
		case PEND_PASSWORD:
			password(param);
			break;

		case PEND_UPDATE_DB:
			update_db(param);
			break;

		/*
		 * Normal player actions
		 */
		case PEND_NEXT:
			next(false);
			break;

		case PEND_REALLY_NEXT:
			next(true);
			break;

		case PEND_PREV:
			prev();
			break;

		case PEND_VOLUME:
			setvolume(param);
			break;

		case PEND_MUTE:
			mute();
			break;

		case PEND_CROSSFADE:
			crossfade(atoi(param.c_str()));
			break;

		case PEND_SEEK:
			seek(atoi(param.c_str()));
			break;
	}

	if (msg->code != 0 && msg->str.size() > 0)
	{
		pms->putlog(msg);
		msg->clear();
	}

	return true;
}

/*
 * Execute an input string from the command line
 */
long		Interface::exec(string s)
{
	if (pms->input->run(s, *msg))
	{
		pms->drawstatus();
		handle_command(pms->input->getpending()); //FIXME

		return STOK;
	}
	else if (pms->input->text.substr(0, 1) == "!")
	{
		return shell(pms->input->text.substr(1));
	}
	else if (!pms->config->readline(s))
	{
		if (pms->msg->code == CERR_NONE)
			pms->log(MSG_STATUS, STOK, "  %s", pms->msg->str.c_str());
		else
			pms->log(MSG_STATUS, STERR, _("Error %d: %s"), pms->msg->code, pms->msg->str.c_str());

		return pms->msg->code;
	}
}

/*
 * Print program name and version.
 */
long		Interface::version()
{
	pms->log(MSG_STATUS, STOK, "%s %s", PMS_NAME, PACKAGE_VERSION);
	return STOK;
}

/*
 * Clear out a line from the topbar.
 * If line is 0, clear everything.
 */
long		Interface::clear_topbar(int line = 0)
{
	if (pms->options->topbar.size() == 0)
		return STOK;

	if (line == 0)
	{
		pms->options->topbar.clear();
		pms->log(MSG_STATUS, STOK, _("Cleared out the entire topbar."));
		return STOK;
	}

	if (line < 1 || line > pms->options->topbar.size())
	{
		pms->log(MSG_STATUS, STERR, _("Out of range, acceptable range is 1-%d."), pms->options->topbar.size());
		return STERR;
	}
	pms->options->topbar.erase(pms->options->topbar.begin() + line - 1);
	pms->log(MSG_STATUS, STOK, _("Removed line %d from topbar."), line);

	pms->mediator->add("redraw.topbar");

	return STOK;
}

/*
 * Redraw everything
 */
long		Interface::redraw()
{
	pms->mediator->add("redraw");
	return STOK;
}

/*
 * Re-read the configuration file
 */
long		Interface::rehash()
{
	pms->options->reset();
	msg->code = source(pms->options->get_string("configfile"));

	if (msg->code == 0)
		pms->log(MSG_STATUS, STOK, _("Reloaded configuration file."));

	return msg->code;
}

/*
 * Save the current configuration to a file
 */
long		Interface::write_config(string file)
{
	if (file.size() == 0)
		file = pms->options->get_string("configfile");

	return STERR;
	//FIXME: implement this
}

/*
 * Read a configuration file or script
 */
long		Interface::source(string file)
{
	if (pms->config->source(file))
	{
		pms->log(MSG_STATUS, STOK, _("Read configuration file: %s"), file.c_str());
		return STOK;
	}
	else
	{
		pms->log(MSG_STATUS, STERR, _("Error reading %s: %s"), file.c_str(), pms->msg->str.c_str());
		return STERR;
	}
}

/*
 * Quit PMS.
 */
long		Interface::quit()
{
	pms->shutdown();
	return STOK;
}

/*
 * Run a shell command
 */
long		Interface::shell(string command)
{
	pms->run_shell(command);
	pms->drawstatus();
	return pms->msg->code;
}

/*
 * Put song info into the console
 */
long		Interface::show_info()
{
	Song *		song;

	song = pms->disp->cursorsong();
	if (song == NULL)
	{
		pms->log(MSG_STATUS, STERR, _("No info could be retrieved."));
		return STERR;
	}

	pms->log(MSG_STATUS, STOK, "%s%s", pms->options->get_string("libraryroot").c_str(), song->file.c_str());
	pms->log(MSG_CONSOLE, STOK, _("--- song info ---\n"));
	pms->log(MSG_CONSOLE, STOK, "file\t\t = %s%s\n", pms->options->get_string("libraryroot").c_str(), song->file.c_str());
	pms->log(MSG_CONSOLE, STOK, "artist\t\t = %s\n", song->artist.c_str());
	pms->log(MSG_CONSOLE, STOK, "albumartist\t = %s\n", song->albumartist.c_str());
	pms->log(MSG_CONSOLE, STOK, "albumartistsort\t = %s\n", song->albumartistsort.c_str());
	pms->log(MSG_CONSOLE, STOK, "date\t\t = %s\n", song->date.c_str());
	pms->log(MSG_CONSOLE, STOK, "year\t\t = %s\n", song->year.c_str());
	pms->log(MSG_CONSOLE, STOK, "artistsort\t = %s\n", song->artistsort.c_str());
	pms->log(MSG_CONSOLE, STOK, "title\t\t = %s\n", song->title.c_str());
	pms->log(MSG_CONSOLE, STOK, "album\t\t = %s\n", song->album.c_str());
	pms->log(MSG_CONSOLE, STOK, "track\t\t = %s\n", song->track.c_str());
	pms->log(MSG_CONSOLE, STOK, "disc\t\t = %s\n", song->disc.c_str());
	pms->log(MSG_CONSOLE, STOK, _("--- end of info ---\n"));

	return STOK;
}

/*
 * Send a password to MPD
 */
long		Interface::password(string pass)
{
	if (pass.size() == 0)
	{
		pms->log(MSG_STATUS, STERR, _("You have to specify a password."));
		return STERR;
	}
	if (pms->comm->sendpassword(pass))
	{
		pms->log(MSG_STATUS, STOK, _("Password accepted by mpd."));
		pms->options->set_string("password", pms->input->param);
		return STOK;
	}
	else
	{
		generr();
		return STERR;
	}
}

/*
 * Update mpd's library.
 * If location is empty, update everything.
 * location 'this' means the selected (cursor) song.
 * location 'current' means the currently playing song.
 * location 'thisdir' means the selected (cursor) song's directory.
 * location 'currentdir' means the currently playing song's directory.
 */
long		Interface::update_db(string location)
{
	string		libroot;

	libroot = pms->options->get_string("libraryroot");

	if (location.size() > 0)
	{
		if (location == "this" && pms->disp && pms->disp->cursorsong())
		{
			location = pms->disp->cursorsong()->file;
			pms->log(MSG_DEBUG, STOK, "Encountered location 'this', translating to %s\n", location.c_str());
		}
		else if (location == "current" && pms->cursong())
		{
			location = pms->cursong()->file;
			pms->log(MSG_DEBUG, STOK, "Encountered location 'current', translating to %s\n", location.c_str());
		}
		else if (location == "thisdir" && pms->disp && pms->disp->cursorsong())
		{
			location = pms->disp->cursorsong()->dirname();
			pms->log(MSG_DEBUG, STOK, "Encountered location 'thisdir', translating to %s\n", location.c_str());
		}
		else if (location == "currentdir" && pms->cursong())
		{
			location = pms->cursong()->dirname();
			pms->log(MSG_DEBUG, STOK, "Encountered location 'currentdir', translating to %s\n", location.c_str());
		}
		else if (libroot.size() > 0)
		{
			if (location.substr(0, libroot.size()) == libroot && location.size() > libroot.size())
			{
				location = location.substr(libroot.size());
				pms->log(MSG_DEBUG, STOK, "Encountered library root in update parameter, stripping to %s\n", location.c_str());
			}
		}
	}
	else
	{
		location = "/";
	}

	if (pms->comm->rescandb(location))
	{
		if (location == "/")
			pms->log(MSG_STATUS, STOK, _("Scanning entire library for changes..."));
		else
			pms->log(MSG_STATUS, STOK, _("Scanning '%s' for changes..."), location.c_str());

		return STOK;
	}
	else if (pms->comm->status()->db_updating)
	{
		pms->log(MSG_STATUS, STERR, _("A library update is already in progress. Please wait for it to finish first."));
	}
	else
	{
		generr();
	}

	return STERR;
}


/*
 *************************** NORMAL PLAYER ACTIONS **************************
 */


long		Interface::play()
{
}

long		Interface::add()
{
}

/*
 * Skip to the next song in line.
 */
long		Interface::next(bool ignore_playmode = false)
{
	if (playnext(ignore_playmode ? PLAYMODE_LINEAR : pms->options->get_long("playmode"), true) == MPD_SONG_NO_ID)
	{
		pms->log(MSG_STATUS, STERR, _("You have reached the end of the list."));
		return STERR;
	}
	pms->drawstatus();
	return STOK;
}

long		Interface::prev()
{
	Song *		cs;
	song_t		i;

	cs = pms->cursong();
	if (cs == NULL)
	{
		if (pms->comm->playlist()->size() == 0)
		{
			pms->log(MSG_STATUS, STERR, _("Can't skip backwards - there are no songs in the playlist."));
			return STERR;
		}
		i = pms->comm->activelist()->size();
	}
	else
	{
		if (cs->pos <= 0)
		{
			if (pms->options->get_long("repeatmode") == REPEAT_LIST)
			{
				i = pms->comm->playlist()->size();
			}
			else
			{
				pms->log(MSG_STATUS, STERR, _("No previous song."));
				return STERR;
			}
		}
		else
		{
			i = cs->pos;
		}
	}

	--i;
	if (i < 0 || i >= pms->comm->playlist()->size())
	{
		pms->log(MSG_CONSOLE, STERR, _("Previous song: out of range.\n"));
		return STERR;
	}

	cs = pms->comm->playlist()->songs[i];
	pms->comm->playid(cs->id);
	pms->drawstatus();

	return STOK;
}

long		Interface::pause()
{
}

long		Interface::toggleplay()
{
}

long		Interface::stop()
{
}

/*
 * Set or adjust volume.
 * A pure integer value means set volume to this value.
 * +/- before the integer means adjust volume by this percentage.
 */
long		Interface::setvolume(string vol)
{
	bool		ok;

	if (vol.size() == 0)
	{
		pms->log(MSG_STATUS, STERR, _("Unexpected end of line, expected sign or integer value."));
		return STERR;
	}
	if (vol[0] != '+' && vol[0] != '-')
	{
		ok = pms->comm->setvolume(atoi(vol.c_str()));
	}
	else
	{
		if (vol.size() == 1)
		{
			pms->log(MSG_STATUS, STERR, _("Unexpected end of line, expected integer value."));
			return STERR;
		}
		ok = pms->comm->volume(atoi(vol.c_str()));
	}
	if (ok)
	{
		pms->log(MSG_STATUS, STOK, _("Volume: %d%%%%"), pms->comm->status()->volume);
		return STOK;
	}
	else
	{
		generr();
		return STERR;
	}
}

/*
 * Toggle muted status
 */
long		Interface::mute()
{
	if (!pms->comm->mute())
	{
		generr();
		return STERR;
	}

	if (pms->comm->muted())
		pms->log(MSG_STATUS, STOK, "Mute is on, from %d%%%%", pms->comm->mvolume());
	else
		pms->log(MSG_STATUS, STOK, "Mute is off, volume=%d%%%%", pms->comm->status()->volume);

	return STOK;
}

/*
 * Set crossfade time
 */
long		Interface::crossfade(int seconds)
{
	if (seconds == 0)
	{
		seconds = pms->comm->crossfade();
		if (seconds == 0)
			pms->log(MSG_STATUS, STOK, "Crossfade switched off."); 
		else if (seconds > 0)
			pms->log(MSG_STATUS, STOK, "Crossfade switched on and is set to %d seconds.", seconds);
	}
	else
	{
		seconds = pms->comm->crossfade(seconds);
		if (seconds >= 0)
			pms->log(MSG_STATUS, STOK, "Crossfade set to %d seconds.", seconds);
	}
	if (seconds == -1)
	{
		generr();
		return STERR;
	}
	return STOK;
}

/*
 * Seek in the current stream
 */
long		Interface::seek(int seconds)
{
	if (!pms->cursong() || pms->comm->status()->state < MPD_STATUS_STATE_PLAY)
	{
		pms->log(MSG_STATUS, STERR, _("Can't seek when player is stopped."));
		return STERR;
	}

	if (seconds == 0)
		return STERR;

	/* Overflow handling */
	if (pms->comm->status()->time_elapsed + seconds >= pms->comm->status()->time_total)
	{
		/* Skip forward or loop if repeat-one */
		if (pms->options->get_long("repeatmode") == REPEAT_ONE)
			pms->comm->playid(pms->cursong()->id);
		else
			playnext(pms->options->get_long("playmode"), true);
	}
	else if (pms->comm->status()->time_elapsed + seconds < 0)
	{
		/* Skip backwards */
		if (pms->options->get_long("repeatmode") == REPEAT_ONE)
		{
			if (pms->comm->seek(pms->cursong()->time + seconds))
				return STOK;
		}
		else
		{
			if (prev() == STOK)
			{
				pms->comm->update(true);
				if (pms->comm->seek(pms->cursong()->time + seconds))
					return STOK;
			}
			else
			{
				stop();
			}
		}
	}
	else
	{
		/* Normal seek */
		if (pms->comm->seek(seconds))
			return STOK;
	}
	generr();
	return STERR;
}




/*
 * Executes actions. This is the last bit of the user interface.
 */
bool		handle_command(pms_pending_keys action)
{
	Message		err;
	Song *		song = NULL;
	Songlist *	list = NULL;
	Songlist *	dlist = NULL;
	pms_window *	win = pms->disp->actwin();
	pms_window *	tmpwin;
	Input_mode	mode;
	int		i = 0;
	long		l = 0;
	song_t		sn = 0;
	string		s;

	{
		pms->interface->action = action;
		pms->interface->param = pms->input->param;

		if (pms->interface->check_events())
			return true;
	}

	if (win) list = win->plist();

	switch(action)
	{
		case PEND_NONE:
			return false;

		case PEND_DELETE:
			if (!list) return false;
			i = removesongs(list);
			if (i <= 0)
			{
				pms->log(MSG_STATUS, STERR, "No songs removed.");
				return false;
			}

			win->wantdraw = true;
			pms->log(MSG_STATUS, STOK, "Removed %d %s.", i, (i == 1 ? "song" : "songs"));
			break;

		case PEND_MOVE_DOWN:
		case PEND_MOVE_UP:
			if (pms->input->mode() == INPUT_COMMAND || pms->input->mode() == INPUT_JUMP)
			{
				pms->input->goprev();
				pms->drawstatus();
				break;
			}
			pms->disp->movecursor(action == PEND_MOVE_DOWN ? 1 : -1);
			if (win) win->wantdraw = true;
			break;

		case PEND_MOVE_HALFPGDN:
		case PEND_MOVE_HALFPGUP:
			if (!win) break;
			//vim seems to integer divide number of rows visible by 
			//2 unless only one row is visible
			i = (win->bheight() - 1) / 2;
			if (i < 1)
				i = 1;
			pms->disp->scrollwin(i * (action == PEND_MOVE_HALFPGDN ? 1 : -1));
			if (win) win->wantdraw = true;
			break;

		case PEND_MOVE_PGDN:
		case PEND_MOVE_PGUP:
			if (!win) break;
			if (win->bheight() - 1 > 4)
			{
				//more than four lines visible: vim leaves two 
				//previously visible lines visible
				i = win->bheight() - 1 - 2;
			}
			else if (win->bheight() - 1 == 4)
			{
				//four lines visible: vim leaves one previously 
				//visible line visible
				i = 3;
			}
			else
			{
				//three or fewer lines visible: vim moves a full 
				//page
				i = win->bheight() - 1;
			}
			pms->disp->scrollwin(i * (action == PEND_MOVE_PGDN ? 1 : -1));
			if (win) win->wantdraw = true;
			break;

		case PEND_SCROLL_DOWN:
		case PEND_SCROLL_UP:
			if (!win) break;
			pms->disp->scrollwin(action == PEND_SCROLL_DOWN ? 1 : -1);
			if (win) win->wantdraw = true;
			break;

		case PEND_CENTER_CURSOR:
			if (!win) break;
			if (pms->options->get_long("scroll_mode") != SCROLL_NORMAL) break;
			pms->disp->scrollwin(win->scursor() - win->cursordrawstart() - (win->bheight() - 1) / 2);
			break;

		case PEND_MOVE_HOME:
			pms->disp->setcursor(0);
			if (win) win->wantdraw = true;
			break;

		case PEND_MOVE_END:
			pms->disp->setcursor(pms->disp->actwin()->size() - 1);
			if (win) win->wantdraw = true;
			break;

		case PEND_GOTO_CURRENT:
			if (!pms->cursong() || !list) return false;
			if (!list->gotocurrent())
			{
				pms->log(MSG_STATUS, STERR, "Currently playing song is not here.");
				return false;
			}

			win->wantdraw = true;
			break;

		case PEND_STOP:
			if (!pms->comm->stop())
			{
				generr();
				return false;
			}

			pms->drawstatus();
			break;

		case PEND_PLAY:
			if (!list) return false;
			song = list->cursorsong();
			if (song == NULL) return false;
			pms->log(MSG_DEBUG, 0, "Playing song with id=%d pos=%d filename=%s\n", song->id, song->pos, song->file.c_str());
			i = song->id;
			if (i == MPD_SONG_NO_ID)
			{
				i = pms->comm->add(pms->comm->playlist(), song);
				if (i == MPD_SONG_NO_ID)
				{
					generr();
					return false;
				}
			}
			if (!pms->comm->playid(i))
			{
				generr();
				return false;
			}

			pms->drawstatus();
			break;

		case PEND_PLAYALBUM:
			multiplay(MATCH_ALBUM, 0);
			break;

		case PEND_PLAYARTIST:
			multiplay(MATCH_ARTIST, 0);
			break;

		case PEND_ADDALBUM:
			multiplay(MATCH_ALBUM, 1);
			break;

		case PEND_ADDARTIST:
			multiplay(MATCH_ARTIST, 1);
			break;

		case PEND_ADDALL:
			multiplay(MATCH_ALL, 1);
			break;

		case PEND_ADD:
		case PEND_ADDTO:
			if (!pms->disp->cursorsong() && !win->current())
			{
				pms->log(MSG_STATUS, STERR, "This is not a song.");
				break;
			}

			dlist = pms->comm->playlist();

			/* Add list to list */
			if (win && win->type() == WIN_ROLE_WINDOWLIST && win->current() && pms->input->win == NULL && action == PEND_ADD)
			{
				pms->log(MSG_DEBUG, 0, "Adding list to list.\n");
				list = win->current()->plist();
				pms->comm->add(list, pms->comm->playlist());
				pms->log(MSG_STATUS, STOK, "%d songs from %s appended to playlist.", list->size(), list->filename.c_str());
				setwin(pms->disp->findwlist(pms->comm->playlist()));
				break;
			}

			/* Addto spawns windowlist */
			if (action == PEND_ADDTO)
			{
				/* Need additional window param */
				if (pms->input->win == NULL)
				{
					/* Add list from windowlist not supported - TODO */
					if (win->type() == WIN_ROLE_WINDOWLIST)
					{
						pms->log(MSG_STATUS, STERR, "Not supported. Please select the songs you want to add before using the add-to command.");
						break;
					}
					pms->log(MSG_DEBUG, 0, "Storing window parameters: win=%p\n", win);
					pms->input->winstore(win);
					setwin(pms->disp->create_windowlist());
					break;
				}
				if (win->current())
					dlist = win->current()->plist();

				pms->log(MSG_DEBUG, 0, "Returned window parameters: win=%p, destination list=%p\n", win->current(), dlist);

				if (dlist == NULL)
				{
					pms->input->winclear();
					break;
				}

				if (pms->input->win)
					list = pms->input->win->plist();
				else
					list = NULL;
			}

			if (!list || !dlist) break;

			if (dlist == pms->comm->playlist())
				s = "playlist";
			else if (dlist == pms->comm->library())
				s = "library";
			else
				s = dlist->filename;

			/* Add arbitrary file or stream */
			if (pms->input->param.size() > 0)
			{
				song = new Song(pms->input->param);
				if (pms->comm->add(dlist, song) != MPD_SONG_NO_ID)
					pms->log(MSG_STATUS, STOK, "Added '%s' to %s.", pms->input->param.c_str(), s.c_str());
				else
					generr();

				delete song;
				break;
			}

			/* Add selected song(s) */
			song = list->popnextselected();
			while (song != NULL)
			{
				pms->log(MSG_DEBUG, 0, "Adding song at %p with id=%d pos=%d filename=%s\n", song, song->id, song->pos, song->file.c_str());
				if (pms->comm->add(dlist, song) != MPD_SONG_NO_ID)
					++i;
				else
					generr();

				song = list->popnextselected();
			}
			if (i == 0)
			{
				generr();
			}
			else
			{
				if (i == 1 && pms->options->get_bool("nextafteraction"))
					pms->disp->movecursor(1);
				pms->log(MSG_STATUS, STOK, _("Added %d %s to %s."), i, (i == 1 ? "song" : "songs"), s.c_str());
			}

			win->wantdraw = true;
			break;


		case PEND_PLAYRANDOM:
		case PEND_ADDRANDOM:
			if (!list) list = pms->comm->library();
			song = list->randsong();
			i = pms->comm->add(pms->comm->playlist(), song);
			if (i == MPD_SONG_NO_NUM)
				break;
			if (action == PEND_PLAYRANDOM)
				pms->comm->playid(i);
			if (win)
				win->wantdraw = true;
			pms->drawstatus();
			break;

		case PEND_GOTORANDOM:
			if (!list)
			{
				pms->log(MSG_STATUS, STERR, _("This command can only be run within a playlist."));
				break;
			}
			song = list->randsong(&sn);
			if (song == NULL) break;
			pms->disp->setcursor(sn);
			break;

		case PEND_MOVEITEMS:
			if (!list || !win)
			{
				pms->log(MSG_STATUS, STERR, _("You can't move anything else than songs."));
				break;
			}
			i = pms->comm->move(list, atoi(pms->input->param.c_str()));
			if (i == 0)
			{
				pms->log(MSG_STATUS, STERR, _("Can't move."));
				break;
			}
			else if (i == 1)
			{
				pms->disp->movecursor(atoi(pms->input->param.c_str()));
			}
			else
			{
				win->wantdraw = true;
			}
			pms->drawstatus();
			break;

		case PEND_PAUSE:
			pms->comm->pause(false);
			pms->drawstatus();
			break;

		case PEND_TOGGLEPLAY:
			pms->comm->pause(true);
			pms->drawstatus();
			break;

		case PEND_SHUFFLE:
			if (pms->comm->shuffle() != 0) break;
			pms->log(MSG_STATUS, STOK, "Playlist shuffled.");
			break;

		case PEND_REPEAT:
			switch(pms->options->get_long("repeatmode"))
			{
				case REPEAT_NONE:
					pms->options->set("repeatmode", "single");
					break;
				case REPEAT_ONE:
					pms->options->set("repeatmode", "yes");
					break;
				default:
				case REPEAT_LIST:
					pms->options->set("repeatmode", "no");
					break;
			}

			pms->log(MSG_DEBUG, 0, "Repeatmode set to %d\n", pms->options->get_long("repeatmode"));

			/* Have MPD manage repeat inside playlist */
			pms->comm->repeat(pms->options->get_long("repeatmode") == REPEAT_LIST && pms->comm->activelist() == pms->comm->playlist());

			pms->drawstatus();
			break;

		case PEND_CLEAR:
			if (!pms->comm->clear(list)) break;
			pms->log(MSG_STATUS, STOK, "Playlist cleared.");
			break;

		case PEND_CROP:
		case PEND_CROPSELECTION:
			if (!list) break;

			if (!pms->comm->crop(list, (action == PEND_CROP ? 0 : 1)))
				pms->log(MSG_STATUS, STERR, "Could not find playing song here.");
			else
				pms->log(MSG_STATUS, STOK, "Playlist cropped.");
			break;



		/* Command-mode + searching*/
		case PEND_COMMANDMODE:
			pms->input->mode(INPUT_COMMAND);
			pms->drawstatus();
			break;

		case PEND_TEXT_UPDATED:
			pms->drawstatus();
			if (pms->input->mode() == INPUT_JUMP)
			{
				if (!win) break;
				i = win->scursor();
				if ((unsigned int)i >= win->size()) i = 0;
				win->jumpto(pms->input->text, i);
			}
			break;

		case PEND_TEXT_RETURN:

			mode = pms->input->mode();
			pms->input->savehistory();

			if (win && win->type() == WIN_ROLE_WINDOWLIST)
				pms->input->mode(INPUT_LIST);
			else
				pms->input->mode(INPUT_NORMAL);

			if (mode == INPUT_COMMAND)
			{
				pms->interface->exec(pms->input->text);
			}
			else if (mode == INPUT_JUMP)
			{
				pms->input->searchterm = pms->input->text;

				if (win->posof_jump(pms->input->text, 0) == -1)
					pms->log(MSG_STATUS, STERR, _("Pattern not found: %s"), pms->input->text.c_str());

				//else do nothing so the search command is left visible
			}
			else
			{
				pms->drawstatus();
			}

			break;

		/* Special case for list */
		case PEND_RETURN:
			if (win && win->type() == WIN_ROLE_WINDOWLIST)
			{
				win = win->current();

				/* Stored mode - return to original call */
				if (pms->input->winpop())
				{
					handle_command(pms->input->getpending());
					if (pms->options->get_bool("addtoreturns"))
					{
						setwin(pms->input->win);
						pms->disp->lastwin = win;
					}
					else
					{
						setwin(win);
						pms->disp->lastwin = pms->input->win;
					}
					pms->input->winclear();
				}
				else
				{
					/* Windowlist mode, switch to new window */
					if (!setwin(win))
						pms->log(MSG_STATUS, STERR, _("Can't change window."));
				}
			}
			break;

		case PEND_TEXT_ESCAPE:
		case PEND_RETURN_ESCAPE:
			if (!win) break;
			switch(win->type())
			{
				case WIN_ROLE_BINDLIST:
					pms->input->mode(INPUT_NORMAL);
					setwin(pms->disp->lastwin);
					break;
				case WIN_ROLE_WINDOWLIST:
					pms->input->winclear();
					setwin(win->lastwin());
					break;
				default:
					pms->input->mode(INPUT_NORMAL);
					pms->drawstatus();
					break;
			}
			
			break;

		/* Searching */
		case PEND_JUMPNEXT:
			if (!win || win->type() != WIN_ROLE_PLAYLIST)
			{
				pms->log(MSG_STATUS, STERR, _("Can't search within this window."));
				break;
			}
			i = win->plist()->cursor() + 1;
			if ((unsigned int)i > win->plist()->end()) i = 0;
			if (win->jumpto(pms->input->searchterm, i))
			{
				pms->log(MSG_STATUS, STOK, "/%s", pms->input->searchterm.c_str());
			}
			else
			{
				pms->log(MSG_STATUS, STERR, "Pattern not found: %s", pms->input->searchterm.c_str());
			}
			break;

		case PEND_JUMPPREV:
			if (!win || win->type() != WIN_ROLE_PLAYLIST)
			{
				pms->log(MSG_STATUS, STERR, _("Can't search within this window."));
				break;
			}
			pms->log(MSG_STATUS, STOK, "?%s", pms->input->searchterm.c_str());
			if (win->jumpto(pms->input->searchterm, win->plist()->cursor(), true))
			{
				pms->log(MSG_STATUS, STOK, "?%s", pms->input->searchterm.c_str());
			}
			else
			{
				pms->log(MSG_STATUS, STERR, "Pattern not found: %s", pms->input->searchterm.c_str());
			}
			break;

		case PEND_JUMPMODE:
			if (!win || win->type() != WIN_ROLE_PLAYLIST)
			{
				pms->log(MSG_STATUS, STERR, _("Can't search within this window."));
				break;
			}
			pms->input->mode(INPUT_JUMP);
			pms->drawstatus();
			break;

		case PEND_PREVOF:
		case PEND_NEXTOF:
			if (pms->input->param.size() == 0)
			{
				pms->log(MSG_STATUS, STERR, _("This command has to be run with a field argument."));
				return false;
			}

			if (list == NULL)
			{
				pms->log(MSG_STATUS, STERR, _("This command has to be run within a playlist."));
				return false;
			}

			if (action == PEND_NEXTOF)
				sn = list->nextof(pms->input->param);
			else
				sn = list->prevof(pms->input->param);

			if (sn != MATCH_FAILED && win != NULL)
				win->setcursor(sn);
			else
				pms->log(MSG_STATUS, STERR, _("Could not find another entry of type '%s'."), pms->input->param.c_str());

			break;

		/* Window control */
		case PEND_CREATEPLAYLIST:
		case PEND_SAVEPLAYLIST:
			tmpwin = win;
			i = createwindow(pms->input->param, win, list);

			switch(i)
			{
			/* Created both playlist and window */
			case 0:
				win->setplist(list);
				if (action == PEND_SAVEPLAYLIST)
					list->set(pms->comm->playlist());
				else
				{
					pms->comm->clear(list);

					/* In case "create" was in reply to addto or something else */
					if (tmpwin && tmpwin->type() == WIN_ROLE_WINDOWLIST && pms->input->winpop())
					{
						tmpwin->setcursor(tmpwin->size());
						handle_command(pms->input->getpending());
						if (pms->options->get_bool("addtoreturns"))
						{
							setwin(pms->input->win);
							pms->disp->lastwin = win;
						}
						else
						{
							setwin(win);
							pms->disp->lastwin = pms->input->win;
						}
						pms->input->winclear();
						break;
					}
				}
				setwin(win);
				break;
			/* Already exists */
			case 1:
				setwin(win);
				s = "\"%s\" already exists.";
				pms->log(MSG_STATUS, STERR, s.c_str(), pms->input->param.c_str());
				break;
			/* No parameter */
			case -1:
				pms->input->mode(INPUT_COMMAND);
				pms->input->text = "create " + pms->input->param;
				pms->drawstatus();
				break;
			case -2:
				generr();
				break;
			case -3:
			default:
				pms->log(MSG_STATUS, STERR, "Internal error: can't create a window.");
				pms->log(MSG_DEBUG, 0, "Window creation failed in PEND_CREATEPLAYLIST, win=%p list=%p\n", win, list);
				break;
			case -4:
				pms->log(MSG_STATUS, STERR, "Internal error: can't find the right window.");
				pms->log(MSG_DEBUG, 0, "Window search failed in PEND_CREATEPLAYLIST, win=%p list=%p\n", win, list);
			}
			break;

		/* Delete a playlist */
		case PEND_DELETEPLAYLIST:

			if (pms->input->param.size() > 0)
			{
				s = pms->input->param;
				list = pms->comm->findplaylist(s);
			}
			else
			{
				/* In case of windowlist */
				if (!list)
				{
					win = win->current();
					if (!win) break;
					list = win->plist();
					if (!list) break;
				}

				if (list->filename.size() == 0)
				{
					pms->log(MSG_STATUS, STERR, "You can't remove a pre-defined playlist.");
					break;
				}

				s = list->filename;
			}

			win = pms->disp->findwlist(list);
			if (pms->comm->deleteplaylist(s))
			{
				pms->disp->delete_window(win);
				pms->log(MSG_STATUS, STOK, "Deleted playlist '%s'.", s.c_str());
			}
			else
			{
				generr();
			}

			win = pms->disp->actwin();
			if (win)
			{
				/* Update cursor position if this was a windowlist */
				win->current();
				win->wantdraw = true;
			}

			break;
			

		case PEND_NEXTWIN:
			if (!setwin(pms->disp->nextwindow()))
				pms->log(MSG_STATUS, STERR, "There is no next window.");
			break;

		case PEND_PREVWIN:
			if (!setwin(pms->disp->prevwindow()))
				pms->log(MSG_STATUS, STERR, "There is no previous window.");
			break;

		case PEND_CHANGEWIN:
			if (pms->input->param == "playlist")
				win = pms->disp->findwlist(pms->comm->playlist());
			else if (pms->input->param == "library")
				win = pms->disp->findwlist(pms->comm->library());
			else if (pms->input->param == "windowlist")
				win = pms->disp->create_windowlist();
//TODO: add this for 0.40.7
//			else if (pms->input->param == "directorylist")
//				win = pms->disp->create_directorylist();
			else
			{
				win = pms->disp->findwlist(pms->comm->findplaylist(pms->input->param));
				if (!win)
				{
					pms->log(MSG_STATUS, STERR, "Change window: invalid parameter '%s'", pms->input->param.c_str());
					break;
				}
			}

			if (win)
				setwin(win);

			break;

		case PEND_LASTWIN:
			if (win && win->type() == WIN_ROLE_WINDOWLIST)
			{
				win->switchlastwin();
				break;
			}
			setwin(pms->disp->lastwin);
			break;

		case PEND_SHOWBIND:
			if (win)
			{
				if (win->type() == WIN_ROLE_BINDLIST)
					break;
			}

			win = pms->disp->create_bindlist();
			if (!win)
				pms->log(MSG_STATUS, STERR, "Can not show the list of key pms->bindings.");
			else
				setwin(win);

			break;

		/*
		 * Specifies which playlist should be played through
		 */
		case PEND_ACTIVATELIST:
			if (!win) break;

			if (pms->input->param.size() == 0)
			{
				/* Inside windowlist window, select window from cursor - else use active window */
				if (list == NULL)
				{
					win = win->current();
					if (win == NULL)
						break;
					list = win->plist();
				}
			}
			else
			{
				/* Use parameter as list */
				list = pms->comm->findplaylist(pms->input->param);
			}

			if (list == NULL)
			{
				pms->log(MSG_STATUS, STERR, "Invalid playlist name.");
				break;
			}

			if (pms->comm->activatelist(list))
			{
				pms->drawstatus();
				win = pms->disp->actwin();
				if (win && win->type() == WIN_ROLE_WINDOWLIST)
					win->wantdraw = true;
			}
			else
				pms->log(MSG_STATUS, STERR, "Can not activate playlist '%s'.", list->filename.c_str());

			break;

		/* Selection */
		case PEND_TOGGLESELECT:
		case PEND_SELECT:
		case PEND_UNSELECT:
		case PEND_CLEARSELECTION:
			if (!win) break;
			if (!win->plist())
			{
				if (win->type() != WIN_ROLE_WINDOWLIST)
					break;
				win->toggleselect();
			}
			if (makeselection(win->plist(), action, pms->input->param) && (pms->input->param.size() == 0)
					&& pms->options->get_bool("nextafteraction") && action != PEND_CLEARSELECTION)
				win->movecursor(1);
			win->wantdraw = true;
			break;

		/* Other */


		/* Cycle through between linear play, random and play single song */
		case PEND_CYCLE_PLAYMODE:
			switch(pms->options->get_long("playmode"))
			{
				default:
				case PLAYMODE_MANUAL:
					pms->options->set_long("playmode", PLAYMODE_LINEAR);
					break;
				case PLAYMODE_LINEAR:
					pms->options->set_long("playmode", PLAYMODE_RANDOM);
					break;
				case PLAYMODE_RANDOM:
					pms->options->set_long("playmode", PLAYMODE_MANUAL);
					break;
			}

			/* Have MPD manage random inside playlist */
			pms->comm->random(pms->options->get_long("playmode") == PLAYMODE_RANDOM && pms->comm->activelist() == pms->comm->playlist());

			pms->drawstatus();
			break;

		case PEND_RESIZE:
			pms->disp->resized();
			pms->disp->forcedraw();
			pms->drawstatus();
			break;

		default:
			return false;
	}

	return true;
}

/*
 * Reports a generic error onto the statusbar
 */
void		generr()
{
	pms->log(MSG_STATUS, STERR, "%s", pms->comm->err());
}

/*
 * Adds or enqueues the next song based on play mode
 */
int		playnext(long mode, int playnow)
{
	Song *		song;
	int		i;

	if (mode == PLAYMODE_LINEAR)
	{
		if (!pms->cursong() || (int)pms->comm->playlist()->end() != pms->cursong()->pos)
			song = pms->comm->playlist()->nextsong();
		else
			song = pms->comm->activelist()->nextsong();

		if (!song) return MPD_SONG_NO_ID;

		if (song->id == MPD_SONG_NO_NUM)
			i = pms->comm->add(pms->comm->playlist(), song);
		else
			i = song->id;
	}
	else if (mode == PLAYMODE_RANDOM)
	{
		if (pms->cursong() && static_cast<int>(pms->comm->playlist()->end()) != pms->cursong()->pos)
		{
			song = pms->comm->playlist()->nextsong();
			if (!song) return MPD_SONG_NO_ID;
			i = song->id;
		}
		else
		{
			song = pms->comm->activelist()->randsong();
			if (!song) return MPD_SONG_NO_ID;
			i = pms->comm->add(pms->comm->playlist(), song);
		}
	}
	else	return MPD_SONG_NO_ID;

	if (i == MPD_SONG_NO_NUM)
		return MPD_SONG_NO_ID;

	if (playnow == true)
		pms->comm->playid(i);

	return i;
}

/*
 * Plays or adds all of one type
 */
int		multiplay(long mode, int playmode)
{
	Songlist *		list;
	Song *			song;
	pms_window *		win;
	int			i = MATCH_FAILED;
	int			listend;
	int			first = -1;
	string			pattern;
	string			pmode;

	win = pms->disp->actwin();
	if (!win || !win->plist()) return false;
	list = win->plist();
	if (list == pms->comm->playlist()) return false;
	song = win->plist()->cursorsong();
	if (song == NULL) return false;

	pmode = (playmode == 0 ? _("Playing") : _("Adding"));

	switch(mode)
	{
		case MATCH_ARTIST:
			if (!song->artist.size()) return false;
			pattern = song->artist;
			pms->log(MSG_STATUS, STOK, _("%s all songs by %s"), pmode.c_str(), song->artist.c_str());
			i = 0;
			break;

		case MATCH_ALBUM:
			if (!song->album.size()) return false;
			pattern = song->album;

			if (pms->comm->playlist()->match(pattern, pms->comm->playlist()->end(), pms->comm->playlist()->end(), mode | MATCH_EXACT) == MATCH_FAILED)
			{
				//last track of the current playlist is not part of this album
				i = 0;
				pms->log(MSG_STATUS, STOK, _("%s album '%s' by %s"), pmode.c_str(), song->album.c_str(), song->artist.c_str());
			}
			else
			{
				//last track of the current playlist is part of this album
				//get last track of the album
				song = list->songs[list->match(pattern, 0, list->end(), mode | MATCH_EXACT | MATCH_REVERSE)];
				if (pms->comm->playlist()->match(song->file, pms->comm->playlist()->end(), pms->comm->playlist()->end(), MATCH_FILE | MATCH_EXACT) != MATCH_FAILED)
				{
					//last track of playlist matches last track of album
					i = 0;
					pms->log(MSG_STATUS, STOK, _("%s album '%s' by %s"), pmode.c_str(), song->album.c_str(), song->artist.c_str());
				}
				else
				{
					//find position in the library of the playlist's last track, 
					//start adding from the one after that
					i = list->match(pms->comm->playlist()->songs[pms->comm->playlist()->end()]->file, 0, list->end(), MATCH_FILE | MATCH_EXACT) + 1;
					pms->log(MSG_STATUS, STOK, _("%s remainder of album '%s' by %s"), pmode.c_str(), song->album.c_str(), song->artist.c_str());
				}
			}
			break;

		case MATCH_ALL:
			pms->log(MSG_STATUS, STOK, _("%s all songs on the current list"), pmode.c_str());
			pattern = "";
			i = 0;
			break;

		default:
			return false;
	}

	listend = static_cast<int>(list->end());

	pms->comm->list_start();
	while (true)
	{
		i = list->match(pattern, i, list->end(), mode | MATCH_EXACT);
		if (i == MATCH_FAILED) break;
		if (first == -1)
			first = pms->comm->playlist()->size();
		pms->comm->add(pms->comm->playlist(), list->songs[i]);
		if (++i > listend) break;
	}
	if (!pms->comm->list_end())
		return false;

	if (first != -1 && playmode == 0)
		pms->comm->playpos(first);

	return true;
}

/*
 * Changes to a window, and sets appropriate input mode
 */
bool		setwin(pms_window * win)
{
	if (!win) return false;
	if (!pms->disp->activate(win)) return false;

	if (win->type() == WIN_ROLE_WINDOWLIST || win->type() == WIN_ROLE_BINDLIST)
		pms->input->mode(INPUT_LIST);
	else
		pms->input->mode(INPUT_NORMAL);

	if (pms->options->get_bool("followwindow"))
	{
		pms->comm->activatelist(win->plist());
		pms->drawstatus();
	}
	return true;
}

/*
 * Perform select, unselect or toggle select on one or more entries
 */
bool		makeselection(Songlist * list, pms_pending_keys action, string param)
{
	Song *		song;
	int		mode;
	int		i = -1;

	if (!list)	return false;

	switch(action)
	{
		case PEND_SELECT:
			mode = 1;
			break;
		case PEND_UNSELECT:
			mode = 0;
			break;
		case PEND_TOGGLESELECT:
			mode = 2;
			break;
		case PEND_CLEARSELECTION:
			for(i = 0; i <= list->end(); i++)
				list->selectsong(list->songs[i], 0);
			return true;
		default:
			return false;
	}

	/* Perform only on one object */
	if (param.size() == 0)
	{
		song = list->cursorsong();
		if (!song)
			return false;
		if (mode == 2)
			list->selectsong(song, !song->selected);
		else
			list->selectsong(song, mode);
		return true;
	}

	/* Perform on range of objects */
	i = list->match(param, 0, list->end(), MATCH_ALL);
	if (i == MATCH_FAILED)
		return false;
	while (i != MATCH_FAILED)
	{
		song = list->songs[i];
		if (!song)
			continue;
		if (mode == 2)
			list->selectsong(song, !song->selected);
		else
			list->selectsong(song, mode);

		if ((unsigned int)i == list->end())
			break;

		i = list->match(param, ++i, list->end(), MATCH_ALL);
	}

	return true;
}

/*
 * Remove selected songs in a list
 */
int		removesongs(Songlist * list)
{
	int				count = 0;
	Song *				song;
	vector<Song *>			songs;
	vector<Song *>::iterator	i;

	if (!list) return 0;

	song = list->popnextselected();
	while (song != NULL)
	{
		songs.push_back(song);
		song = list->popnextselected();
	}

	i = songs.begin();
	while (i != songs.end())
	{
		if (pms->comm->remove(list, *i))
			++count;
		++i;
	}

	return count;
}

/*
 * Create a playlist and connect a window to it, returns the window if successful
 */
int		createwindow(string param, pms_window *& win, Songlist *& list)
{
	win = NULL;
	list = NULL;

	if (param.size() == 0)
		return -1;

	list = pms->comm->findplaylist(param);

	if (list)
	{
		win = pms->disp->findwlist(list);
		if (win)
		{
			return 1;
		}
		return -4;
	}

	list = pms->comm->newplaylist(param);
	
	if (!list)
		return -2;

	win = static_cast<pms_window *> (pms->disp->create_playlist());

	if (!win)
		return -3;

	return 0;
}



/*
 * Defines all commands which can be used
 */
bool init_commandmap()
{
	pms->commands = new Commandmap();
	if (pms->commands == NULL)	return false;
	pms->bindings = new Bindings(pms->commands);
	if (pms->bindings == NULL)
	{
		delete pms->commands;
		return false;
	}

	/* Misc stuff */
	pms->commands->add("!", "Run a shell command", PEND_SHELL);
	pms->commands->add("command-mode", "Switch to command mode", PEND_COMMANDMODE);
	pms->commands->add("info", "Show file information in console", PEND_SHOW_INFO);
	pms->commands->add("password", "Send a password to the server", PEND_PASSWORD);
	pms->commands->add("source", "Read a script or configuration file", PEND_SOURCE);
	pms->commands->add("write-config", "Save configuration file", PEND_WRITE_CONFIG);
	pms->commands->add("rehash", "Read configuration file", PEND_REHASH);
	pms->commands->add("redraw", "Force screen redraw", PEND_REDRAW);
	pms->commands->add("version", "Show program version", PEND_VERSION);
	pms->commands->add("v", "Show program version", PEND_VERSION);
	pms->commands->add("clear-topbar", "Remove all contents in topbar", PEND_CLEAR_TOPBAR);
	pms->commands->add("quit", "Quit program", PEND_QUIT);
	pms->commands->add("q", "Quit program", PEND_QUIT);

	/* Playlist management */
	pms->commands->add("update", "Update MPD music library", PEND_UPDATE_DB);
	pms->commands->add("create", "Create an empty playlist", PEND_CREATEPLAYLIST);
	pms->commands->add("save", "Save the playlist as a new playlist", PEND_SAVEPLAYLIST);
	pms->commands->add("delete-list", "Delete the current playlist", PEND_DELETEPLAYLIST);
	pms->commands->add("select", "Add song under cursor to selection", PEND_SELECT);
	pms->commands->add("unselect", "Remove song under cursor to selection", PEND_UNSELECT);
	pms->commands->add("clear-selection", "Clear selection", PEND_CLEARSELECTION);
	pms->commands->add("toggle-select", "Toggle selection", PEND_TOGGLESELECT);
	pms->commands->add("remove", "Remove selected song(s) from list", PEND_DELETE);
	pms->commands->add("move", "Move songs by offset N", PEND_MOVEITEMS);

	/* Searching */
	pms->commands->add("next-result", "Jump to next result", PEND_JUMPNEXT);
	pms->commands->add("prev-result", "Jump to previous result", PEND_JUMPPREV);
	pms->commands->add("quick-find", "Go to quicksearch mode", PEND_JUMPMODE);
	pms->commands->add("next-of", "Jump to next of given field", PEND_NEXTOF);
	pms->commands->add("prev-of", "Jump to previous of given field", PEND_PREVOF);

	/* Playback */
	pms->commands->add("play", "Play song under cursor", PEND_PLAY);
	pms->commands->add("play-album", "Play entire album of song under cursor", PEND_PLAYALBUM);
	pms->commands->add("play-artist", "Play all songs from artist of song under cursor", PEND_PLAYARTIST);
	pms->commands->add("play-random", "Play a random song", PEND_PLAYRANDOM);
	pms->commands->add("add", "Add selected song(s) to playlist", PEND_ADD);
	pms->commands->add("add-to", "Add selected song(s) to a named playlist", PEND_ADDTO);
	pms->commands->add("add-album", "Add entire album of song under cursor to playlist", PEND_ADDALBUM);
	pms->commands->add("add-artist", "Add all songs from artist of song under cursor to playlist", PEND_ADDARTIST);
	pms->commands->add("add-random", "Add a random song to playlist", PEND_ADDRANDOM);
	pms->commands->add("add-all", "Add all songs from the currently visible list to playlist", PEND_ADDALL);
	pms->commands->add("next", "Next song by play mode", PEND_NEXT);
	pms->commands->add("really-next", "Next song in list regardless of play mode", PEND_REALLY_NEXT);
	pms->commands->add("prev", "Previous song", PEND_PREV);
	pms->commands->add("pause", "Pause or play", PEND_PAUSE);
	pms->commands->add("toggle-play", "Toggle playback, play if stopped", PEND_TOGGLEPLAY);
	pms->commands->add("stop", "Stop", PEND_STOP);
	pms->commands->add("shuffle", "Shuffle playlist", PEND_SHUFFLE);
	pms->commands->add("clear", "Clear the list", PEND_CLEAR);
	pms->commands->add("crop", "Crops list to currently playing song", PEND_CROP);
	pms->commands->add("cropsel", "Crops list to selected songs", PEND_CROPSELECTION);
	pms->commands->add("repeat", "Toggle repeat mode", PEND_REPEAT);
	pms->commands->add("volume", "Increase or decrease volume", PEND_VOLUME);
	pms->commands->add("mute", "Toggle mute", PEND_MUTE);
	pms->commands->add("crossfade", "Set crossfade time", PEND_CROSSFADE);
	pms->commands->add("seek", "Seek in stream", PEND_SEEK);
	pms->commands->add("playmode", "Cycle play mode", PEND_CYCLE_PLAYMODE);

	/* Movement */
	pms->commands->add("next-window", "Go to next playlist", PEND_NEXTWIN);
	pms->commands->add("prev-window", "Go to previous playlist", PEND_PREVWIN);
	pms->commands->add("change-window", "Go to named playlist", PEND_CHANGEWIN);
	pms->commands->add("last-window", "Switch to last window used", PEND_LASTWIN);
	pms->commands->add("help", "Show current key pms->bindings", PEND_SHOWBIND);
	pms->commands->add("goto-random", "Set the cursor to a random song.", PEND_GOTORANDOM);
	pms->commands->add("goto-current", "Find current song", PEND_GOTO_CURRENT);
	pms->commands->add("activate-list", "Activate list for playback", PEND_ACTIVATELIST);
	pms->commands->add("move-up", "Move cursor up", PEND_MOVE_UP);
	pms->commands->add("move-down", "Move cursor down", PEND_MOVE_DOWN);
	pms->commands->add("move-halfpgup", "Move cursor one half page up", PEND_MOVE_HALFPGUP);
	pms->commands->add("move-halfpgdn", "Move cursor one half page down", PEND_MOVE_HALFPGDN);
	pms->commands->add("move-pgup", "Move cursor one page up", PEND_MOVE_PGUP);
	pms->commands->add("move-pgdn", "Move cursor one page down", PEND_MOVE_PGDN);
	pms->commands->add("move-home", "Move cursor to start of list", PEND_MOVE_HOME);
	pms->commands->add("move-end", "Move cursor to end of list", PEND_MOVE_END);
	pms->commands->add("scroll-up", "Scroll up one line", PEND_SCROLL_UP);
	pms->commands->add("scroll-down", "Scroll down one line", PEND_SCROLL_DOWN);
	pms->commands->add("center-cursor", "Center the cursor", PEND_CENTER_CURSOR);
	pms->commands->add("centre-cursor", "Centre the cursor", PEND_CENTER_CURSOR);

	return true;
}



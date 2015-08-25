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
 * pms.cpp - The PMS main class
 *
 */

#include "pms.h"

using namespace std;

Pms *		pms;


/*
 * 1..2..3..
 */
int main(int argc, char *argv[])
{
	int		exitcode;

	pms = new Pms(argc, argv);
	if (!pms)
	{
		printf("Not enough memory, aborting.\n");
		return PMS_EXIT_LOMEM;
	}

	exitcode = pms->init();
	if (exitcode == 0)
	{
		exitcode = pms->main();
	}
	delete pms;
	return exitcode;	
}

/*
 * Init
 */
Pms::Pms(int c, char **v)
{
	argc = c;
	argv = v;
	disp = NULL;
}

/*
 * Unit
 */
Pms::~Pms()
{
}

/*
 * Connection and main loop
 */
int			Pms::main()
{
	string			t_str;
	pms_pending_keys	pending = PEND_NONE;
	char			pass[512] = "";
	bool			statechanged = false;
	bool			songchanged = false;
	pms_window *		win = NULL;
	time_t			timer = 0;

	/* Connection */
	printf(_("Connecting to host %s, port %ld..."), options->get_string("hostname").c_str(), options->get_long("port"));

	if (conn->connect() != 0)
	{
		printf(_("failed.\n"));
		printf("%s\n", conn->errorstr().c_str());

		return PMS_EXIT_CANTCONNECT;
	}

	printf(_("connected.\n"));

	/* Password? */
	if (options->get_string("password").size() > 0)
	{
		printf(_("Sending password..."));
		if (comm->sendpassword(options->get_string("password")))
			printf(_("password accepted.\n"));
		else
			printf(_("wrong password.\n"));
	}

	comm->get_available_commands();
	if (!(comm->authlevel() & AUTH_READ))
	{
		printf(_("This mpd server requires a password.\n"));
		while(true)
		{
			printf(_("Password: "));

			fgets(pass, 512, stdin) ? 1 : 0; //ternary here is a hack to get rid of a warn_unused_result warning
			if (pass[strlen(pass)-1] == '\n')
				pass[strlen(pass)-1] = '\0';

			options->set_string("password", pass);

			comm->sendpassword(pass);
			comm->get_available_commands();
			if (!(comm->authlevel() & AUTH_READ))
				printf(_("Wrong password, try again.\n"));
			else
				break;
		}
	}

	printf(_("Successfully logged in.\n"));

	/* Update lists */
	printf(_("Retrieving library and all playlists..."));
	comm->update(true);
	printf(_("done.\n"));

	comm->has_new_library();
	comm->has_new_playlist();
	printf(_("Sorting library..."));
	comm->library()->sort(options->get_string("sort"));
	printf(_("done.\n"));

	/* Center attention to current song */
	comm->library()->gotocurrent();
	comm->playlist()->gotocurrent();

	_shutdown = false;
	if (!disp->init())
	{
		printf(_("Can't initialize display!\n"));
		return PMS_EXIT_NODISPLAY;
	}

	/* Workaround for buggy ncurses clearing the screen on first getch() */
	getch();

	/* Set up library and playlist windows */
	playlist = disp->create_playlist();
	library = disp->create_playlist();
//	dirlist = disp->create_directorylist();
	if (!playlist || !library)
	{
		delete disp;
		printf(_("Can't initialize windows!\n"));
		return PMS_EXIT_NOWINDOWS;
	}
	playlist->settitle(_("Playlist"));
	library->settitle(_("Library"));
	playlist->list = comm->playlist();
	library->list = comm->library();

	resetstatus(-1);
	drawstatus();

	playlist->set_column_size();
	library->set_column_size();

	connect_window_list();

	/* Focus startup list */
	comm->activatelist(comm->playlist());
	t_str = options->get_string("startuplist");
	if (t_str == "library")
	{
		comm->activatelist(comm->library());
	}
	else if (t_str.size() > 0 && t_str != "playlist")
	{
		comm->activatelist(comm->findplaylist(t_str));
	}
	disp->activate(disp->findwlist(comm->activelist()));

	disp->forcedraw();
	disp->refresh();

	/*
	 * Main loop
	 */
	do
	{
		/* Has to have valid connection. */
		if (!conn->connected() || !comm->update(false) == -1)
		{
			if (timer == 0)
			{
				log(MSG_STATUS, STERR, "Disconnected from mpd: %s", comm->err());
			}
			if (difftime(time(NULL), timer) >= options->get_long("reconnectdelay"))
			{
				if (timer != 0)
					log(MSG_STATUS, STOK, _("Attempting reconnect..."));

				if (conn->connect() != 0)
				{
					if (timer != 0)
						log(MSG_STATUS, STERR, conn->errorstr().c_str());

					time(&timer);
					continue;
				}
				else
				{
					log(MSG_STATUS, STOK, _("Reconnected successfully."));
					comm->clearerror();
					timer = 0;
				}
			}
		}

		/* Get updated info about state and playlists */
		if (comm->has_new_library())
		{
			log(MSG_STATUS, STOK, _("Library updated."));
			if (disp->actwin())
				disp->actwin()->wantdraw = true;
			library->list->sort(options->get_string("sort"));
			library->set_column_size();
			connect_window_list();
		}
		if (comm->has_new_playlist())
		{
			if (disp->actwin())
				disp->actwin()->wantdraw = true;
			playlist->set_column_size();
		}

		/* Progress to next song? */
		progress_nextsong();

		/* Any pending keystrokes? */
		if (input->get_keystroke())
		{
			pending = input->dispatch();
			if (pending != PEND_NONE)
			{
				handle_command(pending);
				comm->update(true);
			}
		}

		songchanged = comm->song_changed();
		statechanged = comm->state_changed();
		if (songchanged)
		{
			/* Cursor follows playback if song changed */
			if (options->get_bool("followplayback"))
			{
				win = disp->findwlist(comm->activelist());
				if (win)
				{
					setwin(win);
					win->gotocurrent();
				}
			}
		}

		if (statechanged)
		{
			/* Shell command when song finishes */
			if (options->get_string("onplaylistfinish").size() > 0 && cursong() && cursong()->pos == comm->playlist()->end())
			{
				/* If a manual stop was issued, don't do anything */
				if (comm->status()->state == MPD_STATUS_STATE_STOP && pending != PEND_STOP)
				{
					/* soak up return value to suppress 
					 * warning */
					int code = system(options->get_string("onplaylistfinish").c_str());
				}
			}
		}


		/* Reset status */
		if (resetstatus(0) >= options->get_long("resetstatus") || songchanged || statechanged)
			drawstatus();

		/* Draw XTerm window title */
		disp->set_xterm_title();

		/* Check out mediator events */
		/* FIXME: add these into their appropriate places */
		if (mediator->changed("setting.sort"))
			comm->library()->sort(options->get_string("sort"));
		else if (mediator->changed("setting.ignorecase"))
			comm->library()->sort(options->get_string("sort"));
		else if (mediator->changed("setting.columns"))
			disp->actwin()->set_column_size();
		else if (mediator->changed("setting.mouse"))
			disp->setmousemask();
		else if (mediator->changed("redraw.topbar"))
			disp->resized();
		else if (mediator->changed("topbarvisible"))
			disp->resized();
		else if (mediator->changed("topbarborders"))
			disp->resized();
		else if (mediator->changed("topbarspace"))
			disp->resized();
		else if (mediator->changed("columnspace"))
			disp->resized();
		else if (mediator->changed("setting.topbarclear"))
		{
			if (options->get_bool("topbarclear"))
				options->topbar.clear();
		}

		/* Draw */
		disp->topbar->wantdraw = true;
		if (mediator->changed("redraw"))
			disp->forcedraw();
		else
			disp->draw();
		disp->refresh();

	}
	while (!_shutdown);

	log(MSG_CONSOLE, STOK, _("Shutting down program.\n"));

	delete disp;
	delete comm;
	delete conn;

	/* Unclutter the prompt */
	printf("\n");

	return PMS_EXIT_SUCCESS;
}

/*
 * Set up neccessary variables
 */
int			Pms::init()
{
	string			str;
	vector<string> *	tok;

	int			exitcode = PMS_EXIT_SUCCESS;
	char *			host;
	char *			port;
	char *			password;
	const char *		charset = NULL;
	
	/* Internal pointers */
	msg = new Message();
	mediator = new Mediator();
	interface = new Interface();
	formatter = new Formatter();

	/* Setup locales and internationalization */
	setlocale(LC_ALL, "");
	setlocale(LC_CTYPE, "");
	g_get_charset(&charset);
	bindtextdomain(GETTEXT_PACKAGE, LOCALE_DIR);
	bind_textdomain_codeset(GETTEXT_PACKAGE, charset);
	textdomain(GETTEXT_PACKAGE);

	/* Print program header */
	printf("%s v%s\n%s\n", PMS_NAME, PACKAGE_VERSION, PMS_COPYRIGHT);

	/* Read important environment variables */
	host = getenv("MPD_HOST");
	port = getenv("MPD_PORT");
	password = getenv("MPD_PASSWORD");

	/* Set up field types */
	fieldtypes = new Fieldtypes();
	fieldtypes->add("num", _("#"), FIELD_NUM, 0, NULL);
	fieldtypes->add("file", _("Filename"), FIELD_FILE, 0, sort_compare_file);
	fieldtypes->add("artist", _("Artist"), FIELD_ARTIST, 0, sort_compare_artist);
	fieldtypes->add("artistsort", _("Artist sort name"), FIELD_ARTISTSORT, 0, sort_compare_artistsort);
	fieldtypes->add("albumartist", _("Album artist"), FIELD_ALBUMARTIST, 0, sort_compare_albumartist);
	fieldtypes->add("albumartistsort", _("Album artist sort name"), FIELD_ALBUMARTISTSORT, 0, sort_compare_albumartistsort);
	fieldtypes->add("title", _("Title"), FIELD_TITLE, 0, sort_compare_title);
	fieldtypes->add("album", _("Album"), FIELD_ALBUM, 0, sort_compare_album);
	fieldtypes->add("track", _("Track"), FIELD_TRACK, 6, sort_compare_track);
	fieldtypes->add("trackshort", _("No"), FIELD_TRACKSHORT, 3, sort_compare_track);
	fieldtypes->add("length", _("Length"), FIELD_TIME, 7, sort_compare_length);
	fieldtypes->add("date", _("Date"), FIELD_DATE, 11, sort_compare_date);
	fieldtypes->add("year", _("Year"), FIELD_YEAR, 5, sort_compare_year);
	fieldtypes->add("name", _("Name"), FIELD_NAME, 0, sort_compare_name);
	fieldtypes->add("genre", _("Genre"), FIELD_GENRE, 0, sort_compare_genre);
	fieldtypes->add("composer", _("Composer"), FIELD_COMPOSER, 0, sort_compare_composer);
	fieldtypes->add("performer", _("Performer"), FIELD_PERFORMER, 0, sort_compare_performer);
	fieldtypes->add("disc", _("Disc"), FIELD_DISC, 5, sort_compare_disc);
	fieldtypes->add("comment", _("Comment"), FIELD_COMMENT, 0, sort_compare_comment);

	/* Set up default bindings */
	if (!init_commandmap())
	{
		return PMS_EXIT_NOCOMMAND;
	}
	options = new Options();
	init_default_keymap();

	/* Our configuration */
	config = new Configurator(options, bindings);

	/* Some default options */
	options->set_string("hostname", (host ? host : "127.0.0.1"));
	if (!password && host)
	{
		tok = splitstr(host, "@");
		if (tok->size() == 2)
		{
			options->set_string("hostname", (*tok)[0]);
			options->set_string("password", (*tok)[1]);
		}
		delete tok;
	}
	if (options->get_string("password").size() == 0)
	{
		options->set_string("password", (password ? password : ""));
	}
	options->set_long("port", (port ? atoi(port) : 6600));

	if (options->get_long("port") <= 0 || options->get_long("port") > 65535)
	{
		printf(_("Error: port number in environment variable MPD_PORT must be from 1-65535\n"));
		return PMS_EXIT_BADARGS;
	}

	/* Parse command-line */
	if (parse_args(argc, argv) == false)
	{
		return PMS_EXIT_BADARGS;
	}

	if (!config->loadconfigs())
		return PMS_EXIT_CONFIGERR;

	/* Seed random number generator */
	srand(time(NULL));

	/* Setup some important stuff */
	conn	= new Connection(options->get_string("hostname"), options->get_long("port"), options->get_long("mpd_timeout"));
	comm	= new Control(conn);
	disp	= new Display(comm);
	input	= new Input();
	if (!conn || !comm || !disp || !input)
		return PMS_EXIT_LOMEM;

	/* Initialization finished */
	return PMS_EXIT_SUCCESS;
}






/*
 * Converts long to string
 */
string			Pms::tostring(long number)
{
	ostringstream s;
	s << number;
	return s.str();
}

/*
 * Converts size_t to string
 */
string			Pms::tostring(size_t number)
{
	ostringstream s;
	s << number;
	return s.str();
}

/*
 * Converts int to string
 */
string			Pms::tostring(int number)
{
	ostringstream s;
	s << number;
	return s.str();
}

/*
 * Split a string into tokens
 */
vector<string> *	Pms::splitstr(string str, string delimiter)
{
	vector<string> *	tokens = new vector<string>;

	string::size_type last	= str.find_first_not_of(delimiter, 0);
	string::size_type pos	= str.find_first_of(delimiter, last);

	while (string::npos != pos || string::npos != last)
	{
		tokens->push_back(str.substr(last, pos - last));
		last = str.find_first_not_of(delimiter, pos);
		pos = str.find_first_of(delimiter, last);
	}
	
	return tokens;
}

/*
 * Join tokens into a string
 */
string			Pms::joinstr(vector<string> * source, vector<string>::iterator start, vector<string>::iterator end, string delimiter)
{
	string			dest = "";

	while (start != source->end())
	{
		dest += *start;

		if (start == end)
			break;

		if (++start != end)
		{
			dest += delimiter;
		}
	}
	
	return dest;
}

/*
 * Formats seconds into the format Dd H:MM:SS.
 */
string			Pms::timeformat(int seconds)
{
	static const int	day	= (60 * 60 * 24);
	static const int	hour	= (60 * 60);
	static const int	minute	= 60;

	int		i;
	string		s = "";

	/* No time */
	if (seconds < 0)
	{
		s = "--:--";
		return s;
	}

	/* days */
	if (seconds >= day)
	{
		i = seconds / day;
		s = Pms::tostring(i) + "d ";
		seconds %= day;
	}

	/* hours */
	if (seconds >= hour)
	{
		i = seconds / hour;
		s += zeropad(i, 1) + ":";
		seconds %= hour;
	}

	/* minutes */
	i = seconds / minute;
	s = s + zeropad(i, 2) + ":";
	seconds %= minute;

	/* seconds */
	s += zeropad(seconds, 2);

	return s;
}

/*
 * Return "song" or "songs" based on plural or not
 */
string			Pms::pluralformat(unsigned int i)
{
	if (i == 1)
		return _("song");
	else
		return _("songs");
}
/*
 * Pad integer with zeroes up to target length
 */
string			Pms::zeropad(int i, unsigned int target)
{
	string s;
	s = Pms::tostring(i);
	while(s.size() < target)
		s = '0' + s;
	return s;
}

/*
 * Replaces % with %%
 */
string			Pms::formtext(string text)
{
	string::const_iterator	i;
	string			nutext;

	i = text.begin();
	nutext.clear();

	while (i != text.end())
	{
		nutext += *i;
		if (*i == '%')
			nutext += *i;
		++i;
	}

	return nutext;
}

/*
 * Return true if the terminal supports Unicode
 */
bool			Pms::unicode()
{
	const char *		charset = NULL;

	g_get_charset(&charset);
	return strcmp(charset, "UTF-8") == 0;
}










/*
 * Run a shell command
 *
 * FIXME: perhaps this command should be within Interface class?
 * TODO: add %artist% tags through the field pattern parser: meaning %file% -> filename, not % -> filename
 *	...but current implementation is nice and vim-like
 */
bool			Pms::run_shell(string cmd)
{
	string				search;
	string				replace;
	string::size_type		pos;
	int				i;
	Songlist *			list;
	char				c;

	msg->clear();

	/*
	 * %: path to current song, not enclosed in quotes
	 */
	if (cursong())
	{
		search = "%";
		replace = options->get_string("libraryroot");
		replace += cursong()->file;
		pos = 0;
		while ((pos = cmd.find(search, pos)) != string::npos)
		{
			if (pos == 0 || cmd[pos - 1] != '\\')
				cmd.replace(pos, search.size(), replace);
			pos++;
		}
	}

	/*
	 * ##: path to each song in selection (or each song on the current 
	 * playlist if there is no selection), each enclosed with doublequotes 
	 * and separated by spaces
	 */
	list = disp->actwin()->plist();
	search = "##";
	if (cmd.find(search, 0) != string::npos && list && list->size())
	{
		replace = "";
		for (i = 0; i < list->size(); i++)
		{
			if (!list->selection.size || list->song(i)->selected)
			{
				replace += options->get_string("libraryroot");
				replace += list->song(i)->file;
				replace += "\" \"";
			}
		}
		if (replace.size() > 0)
		{
			replace = "\"" + replace.substr(0, replace.size() - 2);
			pos = 0;
			while ((pos = cmd.find(search, pos)) != string::npos)
			{
				if (pos == 0 || cmd[pos - 1] != '\\')
					cmd.replace(pos, search.size(), replace);
				pos++;
			}
		}
	}

	/*
	 * #: path to song the cursor is on, not enclosed in quotes
	 */
	if (disp->cursorsong())
	{
		search = "#";
		replace = options->get_string("libraryroot");
		replace += disp->cursorsong()->file;
		pos = 0;
		while ((pos = cmd.find(search, pos)) != string::npos)
		{
			if (pos == 0 || cmd[pos - 1] != '\\')
				cmd.replace(pos, search.size(), replace);
			pos++;
		}
	}

	//pms->log(MSG_DEBUG, 0, "running shell command '%s'\n", cmd.c_str());
	endwin();

	msg->code = system(cmd.c_str());
	msg->code = WEXITSTATUS(msg->code);

	pms->log(MSG_DEBUG, 0, "Shell returned %d\n", msg->code);
	if (msg->code != 0)
		printf(_("\nShell returned %d\n"), msg->code);

	printf(_("\nPress ENTER to continue"));
	fflush(stdout);
	{
		/* soak up return value to suppress warning */
		int key = scanf("%c", &c);
	}

	reset_prog_mode();
	refresh();

	return true;
}

/*
 * Returns the currently playing song
 */
Song *			Pms::cursong()
{
	if (!comm) return NULL;
	return comm->song();
}

/* 
 * Reset status to its original state
 */
void			Pms::drawstatus()
{
	if (input->mode() == INPUT_JUMP)
		log(MSG_STATUS, STOK, "/%s", formtext(input->text).c_str());
	else if (input->mode() == INPUT_FILTER)
		log(MSG_STATUS, STOK, ":g/%s", formtext(input->text).c_str());
	else if (input->mode() == INPUT_COMMAND)
		log(MSG_STATUS, STOK, ":%s", formtext(input->text).c_str());
	else
		log(MSG_STATUS, STOK, "%s", playstring().c_str());

	resetstatus(-1);
}

/*
 * Measures time from last statusbar text
 */
int			Pms::resetstatus(int set)
{
	static time_t 		stored = time(NULL);
	static time_t 		now = time(NULL);

	if (set == 1)
		time(&stored);
	else if (set == -1)
		stored = 0;

	if (stored == 0)
		return 0;

	if (time(&now) == -1)
		return -1;

	return (static_cast<int>(difftime(now, stored)));
}

/*
 * Return a textual description on how song progression works
 */
string			Pms::playstring()
{
	string		s;
	string		list = "<unknown>";
	bool		is_last;

	long		playmode = options->get_long("playmode");
	long		repeatmode = options->get_long("repeat");

	if (!comm->status() || !conn->connected())
	{
		s = "Not connected.";
		return s;
	}

	if (comm->status()->state == MPD_STATUS_STATE_STOP || !cursong())
	{
		s = "Stopped.";
		return s;
	}
	else if (comm->status()->state == MPD_STATUS_STATE_PAUSE)
	{
		s = "Paused...";
		return s;
	}

	if (comm->activelist())
		list = comm->activelist()->filename;

	if (list.size() == 0)
	{
		if (comm->activelist() == comm->library())
			list = "library";
		else if (comm->activelist() == comm->playlist())
			list = "playlist";
	}

	s = "Playing ";

	if (playmode == PLAYMODE_MANUAL)
	{
		s += "this song, then stopping.";
		return s;
	}

	if (repeatmode == REPEAT_ONE)
	{
		s += "the same song indefinitely.";
		return s;
	}

	is_last = (cursong()->pos == static_cast<int>(comm->playlist()->end()));

	if (!is_last && !(comm->activelist() == comm->playlist() && repeatmode == REPEAT_LIST))
	{
		s += "through playlist, then ";
	}

	if (playmode == PLAYMODE_RANDOM)
	{
		s += "random songs from " + list + ".";
		return s;
	}

	if (repeatmode == REPEAT_LIST)
	{
		s += "songs from " + list + " repeatedly.";
		return s;
	}

	if (repeatmode == REPEAT_NONE)
	{
		if (comm->activelist() == comm->playlist())
		{
			if (options->get_bool("followcursor"))
			{
				if (is_last)
					s += "this song, then ";

				s += "following cursor.";
				return s;
			}

			if (is_last)
				s += "this song, then stopping.";
			else
				s += "stopping.";

			return s;
		}
		else
		{
			s += "songs from " + list + ".";
			return s;
		}
	}

	return s;
}

/*
 * Put an arbitrary message into the message log
 */
void			Pms::putlog(Message * m)
{
	if (m->code == 0 && m->str.size() == 0)
		return;

	log(MSG_CONSOLE, m->code, m->str.c_str());
}

/*
 * Log a message.
 * Verbosity levels:
 *  0 = statusbar
 *  1 = console
 *  2 = debug
 */
void			Pms::log(int verbosity, long code, const char * format, ...)
{
	long		loglines;
	va_list		ap;
	char		buffer[1024];
	char		tbuffer[20];
	string		level;
	Message *	m;
	tm *		timeinfo;
	color *		pair;

	if (verbosity >= MSG_DEBUG && !pms->options->get_bool("debug"))
		return;

	m = new Message();
	if (m == NULL)
		return;

	va_start(ap, format);
	vsprintf(buffer, format, ap);
	va_end(ap);

	m->str = buffer;
	m->code = code;

	if (verbosity == MSG_STATUS)
	{
		m->str += "\n";

		if (code == STOK)
			pair = options->colors->status;
		else
			pair = options->colors->status_error;

		disp->statusbar->clear(false, pair);
		colprint(disp->statusbar, 0, 0, pair, "%s", buffer);
		resetstatus(1);
	}

	if (verbosity <= MSG_DEBUG && pms->options->get_bool("debug"))
	{
		timeinfo = localtime(&(m->timestamp));
		strftime(tbuffer, 20, "%Y-%m-%d %H:%M:%S", timeinfo);
		if (verbosity == MSG_STATUS)
			level = "status";
		else if (verbosity == MSG_CONSOLE)
			level = "console";
		else if (verbosity == MSG_DEBUG)
			level = "debug";
		fprintf(stderr, "%s /%s/ %s", tbuffer, level.c_str(), m->str.c_str());
	}

	if (!disp && verbosity < MSG_DEBUG)
	{
		printf("%s", buffer);
	}

	msglog.push_back(m);
	loglines = options->get_long("msg_buffer_size");
	if (loglines > 0 && msglog.size() > loglines)
		msglog.erase(msglog.begin());
}

/*
 * Checks if time is right for song progression, and takes necessary action
 */
bool			Pms::progress_nextsong()
{
	static song_t		lastid = MPD_SONG_NO_ID;
	static Song *		lastcursor = NULL;
	Songlist *		list = NULL;
	unsigned int		remaining;

	long			repeatmode;
	long			playmode;

	if (!cursong())		return false;

	if (comm->status()->state != MPD_STATUS_STATE_PLAY)
		return false;

	remaining = (comm->status()->time_total - comm->status()->time_elapsed - comm->status()->crossfade);

	repeatmode = options->get_long("repeat");
	playmode = options->get_long("playmode");

	/* Too early */
	if (remaining > options->get_long("nextinterval") || lastid == cursong()->id)
		return false;

	/* No auto-progression, even when in the middle of playlist */
	if (playmode == PLAYMODE_MANUAL)
	{
		if (remaining <= options->get_long("stopdelay"))
		{
			pms->log(MSG_DEBUG, 0, "Manual playmode, stopping playback.\n");
			comm->stop();
			lastid = cursong()->id;
			return true;
		}
		else
		{
			return false;
		}
	}
	/* Defeat desync with server */
	lastid = cursong()->id;

	/* Normal progression: reached end of playlist */
	if (comm->status()->song == static_cast<int>(playlist->list->end()))
	{
		/* List to play from */
		list = comm->activelist();
		if (!list) return false;

		if (list == comm->playlist())
		{
			/* Let MPD handle repeating of the playlist itself */
			if (repeatmode == REPEAT_LIST)
				return false;

			/* Let MPD handle random songs from playlist */
			if (playmode == PLAYMODE_RANDOM)
				return false;
		}

		pms->log(MSG_DEBUG, 0, "Auto-progressing to next song.\n");

		/* Playback follows cursor */
		if (options->get_bool("followcursor") && lastcursor != disp->cursorsong() && disp->cursorsong()->file != cursong()->file)
		{
			pms->log(MSG_DEBUG, 0, "Playback follows cursor: last cursor=%p, now cursor=%p.\n", lastcursor, disp->cursorsong());
			lastcursor = disp->cursorsong();
			lastid = comm->add(comm->playlist(), lastcursor);
		}

		/* Normal song progression */
		lastid = playnext(playmode, false);
	}

	if (lastcursor == NULL)
	{
		lastcursor = disp->cursorsong();
	}

	return (lastid != MPD_SONG_NO_ID);
}

/*
 * Create new windows for each custom playlist
 */
bool			Pms::connect_window_list()
{
	bool				ok = true;
	pms_window *			win;
	vector<Songlist *>::iterator	i;

	i = comm->playlists.begin();
	while (i != comm->playlists.end())
	{
		if (disp->findwlist(*i) == NULL)
		{
			win = disp->create_playlist();
			if (win)
				win->setplist(*i);
			else
				ok = false;
		}
		++i;
	}

	return ok;
}

/*
 * Default key bindings
 */
void			Pms::init_default_keymap()
{
	bindings->clear();

	/* Movement */
	bindings->add("up", "move-up");
	bindings->add("down", "move-down");
	bindings->add("pageup", "move-pgup");
	bindings->add("pagedown", "move-pgdn");
	bindings->add("^B", "move-pgup");
	bindings->add("^F", "move-pgdn");
	bindings->add("^U", "move-halfpgup");
	bindings->add("^D", "move-halfpgdn");
	bindings->add("^Y", "scroll-up");
	bindings->add("^E", "scroll-down");
	bindings->add("z", "center-cursor");
	bindings->add("home", "move-home");
	bindings->add("end", "move-end");
	bindings->add("g", "goto-current");
	bindings->add("R", "goto-random");
	bindings->add("j", "move-down");
	bindings->add("k", "move-up");
	bindings->add("t", "prev-window");
	bindings->add("T", "next-window");
	bindings->add("(", "prev-of album");
	bindings->add(")", "next-of album");
	bindings->add("{", "prev-of artist");
	bindings->add("}", "next-of artist");
	bindings->add("1", "change-window playlist");
	bindings->add("2", "change-window library");
	bindings->add("w", "change-window windowlist");
	// TODO: add this for a later version
	//bindings->add("W", "change-window directorylist");
	bindings->add("tab", "last-window");

	/* Searching */
	bindings->add("/", "quick-find");
	bindings->add("n", "next-result");
	bindings->add("N", "prev-result");

	/* Playlist management */
	bindings->add("a", "add");
	bindings->add("A", "add-to");
	bindings->add("b", "add-album");
	bindings->add("B", "play-album");
	bindings->add("delete", "remove");
	bindings->add("c", "cropsel");
	bindings->add("C", "crop");
	bindings->add("insert", "toggle-select");
	bindings->add("F12", "activate-list");
	bindings->add("^X", "delete-list");
	bindings->add("J", "move 1");
	bindings->add("K", "move -1");

	/* Controls */
	bindings->add("return", "play");
	bindings->add("kpenter", "play");
	bindings->add("backspace", "stop");
	bindings->add("p", "pause");
	bindings->add("space", "toggle-play");
	bindings->add("l", "next");
	bindings->add("h", "prev");
	bindings->add("M", "mute");
	bindings->add("m", "playmode");
	bindings->add("r", "repeat");
	bindings->add("+", "volume +5");
	bindings->add("-", "volume -5");
	bindings->add("left", "seek -5");
	bindings->add("right", "seek 5");

	/* Maintenance */
	bindings->add("f", "toggle followcursor");
	bindings->add("F", "toggle followplayback");
	bindings->add("^F", "toggle followwindow");
	bindings->add(":", "command-mode");
	bindings->add("u", "update-library");
	bindings->add("v", "version");
	bindings->add("q", "quit");
	bindings->add("F1", "help");
	bindings->add("^L", "redraw");
}










/*
 * Print the version string
 */
void			Pms::print_version()
{
   	printf("Uses libmpdclient (c) 2003-2007 by Warren Dukes (warren.dukes@gmail.com)\n");
	printf("This program is licensed under the GNU General Public License 3.\n");
}

/*
 * Print switch usage
 */
void			Pms::print_usage()
{
	printf("Usage:\n");
	printf("  -%s\t\t\t%s\n", "v", "print version and exit");
	printf("  -%s\t\t%s\n", "? --help", "display command-line options");
	printf("  -%s\t\t\t%s\n", "d", "turn on debugging to stderr");
	printf("  -%s\t\t%s\n", "c <filename>", "use an alternative config file");
	printf("  -%s\t\t%s\n", "h <hostname>", "connect to this MPD server");
	printf("  -%s\t\t%s\n", "p <port>", "connect to this port");
	printf("  -%s\t\t%s\n", "P <password>", "give this password to MPD server");
}

/*
 * Helper function, prints an error
 */
bool			Pms::require_arg(char c)
{
	printf("Error: option '%c' requires an argument.\n", c);
	print_usage();
	return false;
}

/*
 * Parse command-line arguments
 */
bool			Pms::parse_args(int argc, char * argv[])
{
	int			argn;
	string			value = "";
	string			arg = "";
	bool			switched = false;
	string			s;
	string::iterator	i;

	if (argc <= 1)
		return true;

	for (argn = 1; argn < argc; argn++)
	{
		s = argv[argn];

		if (s == "--help")
		{
			print_usage();
			return false;
		}

		i = s.begin();

		while (i != s.end())
		{
			if (!switched)
				if (*i != '-')
					return false;

			switch (*i)
			{
				case 'd':
					options->set_bool("debug", true);
					break;
				case 'v':
					print_version();
					return false;
				case '?':
					print_usage();
					return false;
				case 'c':
					if (++argn >= argc)
						return require_arg(*i);
					options->set_string("configfile", argv[argn]);
					break;
				case 'h':
					if (++argn >= argc)
						return require_arg(*i);
					options->set_string("hostname", argv[argn]);
					break;
				case 'p':
					if (++argn >= argc)
						return require_arg(*i);
					options->set_long("port", atoi(argv[argn]));
					if (options->get_long("port") <= 0 || options->get_long("port") > 65535)
					{
						printf(_("Error: port number must be from 1-65535\n"));
						return false;
					}
					break;
				case 'P':
					if (++argn >= argc)
						return require_arg(*i);
					options->set_string("password", argv[argn]);
					break;
				case '-':
					if (switched)
					{
						print_usage();
						return false;
					}
					switched = true;
					break;
				default:
					printf(_("Error: unknown option '%c'\n"), *i);
					print_usage();
					return false;
			}
			++i;
		}
	
		switched = false;
	}

	return true;
}


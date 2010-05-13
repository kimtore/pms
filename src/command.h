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
 */

#ifndef _PMS_COMMAND_H_
#define _PMS_COMMAND_H_

#include <string>
#include <time.h>
#include "libmpdclient.h"
#include "conn.h"
#include "list.h"

using namespace std;


/*
 * Which permissions we have
 */
typedef enum
{
	AUTH_NONE	= 0,
	AUTH_READ	= 1 << 0,
	AUTH_ADD	= 1 << 1,
	AUTH_CONTROL	= 1 << 2,
	AUTH_ADMIN	= 1 << 3
}
Mpd_authlevel;


/*
 * A list of all of mpd's commands.
 */
typedef struct
{
	bool		add;
	bool		addid;
	bool		clear;
	bool		clearerror;
	bool		close;
	bool		commands;
	bool		count;
	bool		crossfade;
	bool		currentsong;
	bool		delete_;
	bool		deleteid;
	bool		disableoutput;
	bool		enableoutput;
	bool		find;
	bool		idle;
	bool		kill;
	bool		list;
	bool		listall;
	bool		listallinfo;
	bool		listplaylist;
	bool		listplaylistinfo;
	bool		listplaylists;
	bool		load;
	bool		lsinfo;
	bool		move;
	bool		moveid;
	bool		next;
	bool		notcommands;
	bool		outputs;
	bool		password;
	bool		pause;
	bool		ping;
	bool		play;
	bool		playid;
	bool		playlist;
	bool		playlistadd;
	bool		playlistclear;
	bool		playlistdelete;
	bool		playlistfind;
	bool		playlistid;
	bool		playlistinfo;
	bool		playlistmove;
	bool		playlistsearch;
	bool		plchanges;
	bool		plchangesposid;
	bool		previous;
	bool		random;
	bool		rename;
	bool		repeat;
	bool		rm;
	bool		save;
	bool		filter;
	bool		seek;
	bool		seekid;
	bool		setvol;
	bool		shuffle;
	bool		single;
	bool		stats;
	bool		status;
	bool		stop;
	bool		swap;
	bool		swapid;
	bool		tagtypes;
	bool		update;
	bool		urlhandlers;
	bool		volume;
}
Mpd_allowed_commands;


/*
 * mpd's "status" and "stats" information
 */
class Mpd_status
{
private:
	mpd_Stats *	stats;
	mpd_Status *	status;
public:
			Mpd_status();
			~Mpd_status();

	bool		alive() const;
	void		assign_status(mpd_Status *);
	void		assign_stats(mpd_Stats *);

	bool		muted;
	int		volume;
	bool		repeat;
	bool		single;
	bool		random;
	int		playlist_length;
	long long	playlist;
	long long	storedplaylist;
	int		state;
	int		crossfade;
	song_t		song;
	song_t		songid;
	int		time_elapsed;
	int		time_total;
	bool		db_updating;
	int		error;
	string		errstr;

	/* Audio decoded properties */
	int		bitrate;
	unsigned int	samplerate;
	int		bits;
	int		channels;

	/* Stats */
	song_t		artists_count;
	song_t		albums_count;
	song_t		songs_count;
	unsigned long	uptime;
	unsigned long	db_update_time;
	unsigned long	playtime;
	unsigned long	db_playtime;

	/* Cache to detect changes */
	int		last_state;
	long long	last_playlist;
	unsigned long	last_db_update_time;
	bool		last_db_updating;
	int		update_job_id;
};


/*
 * Directory class, holds information on songs in current directory
 */
class Directory
{
private:
	Directory *			parent_;
	string				name_;
public:
					Directory(Directory *, string);
					~Directory();
	
	int				cursor;
	vector<Song *>			songs;
	vector<Directory *>		children;
	Directory *			add(string);
	string				name() { return (name_.size() == 0 ? "/" : name_); };
	Directory *			parent() { return parent_; };
	string				path();

//	void				debug_tree();
};


/*
 * Interface with mpd server
 */
class Control
{
private:
	Connection *		conn;
	Mpd_status *		st;
	Mpd_allowed_commands	commands;

	mpd_Stats		*_stats;
	Song			*_song;
	Songlist		*_playlist;
	Songlist		*_library;
	Songlist		*_active;

	long long		last_playlist_version;
	int			last_song;
	int			oldsong;
	bool			_has_new_playlist;
	bool			_has_new_library;
	int			command_mode;
	int			mutevolume;
	int			crossfadetime;

	/* Update interval timer */
	time_t			mytime[2];
	int			usetime;

	int			get_current_playing();
	int			get_stats();
	void			retrieve_lists(vector<Songlist *> &);
	unsigned int		update_playlists();
	void			update_playlist();
	void			update_library();
	bool			finish();

public:
				Control(Connection *);
				~Control();

	vector<Songlist *>	playlists;
	Directory *		rootdir;

	bool			alive();
	const char *		err();		// Reports errors from mpd server

	/* Server management */
	int			authlevel();
	void			get_available_commands();
	bool			rescandb(string = "/");
	bool			sendpassword(string);
	void			clearerror();

	/* Set/end command list mode */
	bool			list_start();
	bool			list_end();

	/* List management */
	song_t			add(Songlist *, Song *);
	song_t			add(Songlist * source, Songlist * dest);
	int			remove(Songlist *, Song *);
	int			prune(Songlist *, Songlist *);

	/* Play controls */
	bool			play();
	bool			playid(song_t);
	bool			playpos(song_t);
	bool			pause(bool);
	bool			stop();

	/* Player management */
	bool			shuffle();
	bool			seek(int);
	bool			random(int);
	bool			repeat(bool);
	bool			single(bool);
	bool			setvolume(int);
	bool			volume(int);
	bool			mute();
	bool			muted();
	int			mvolume() { return mutevolume; };

	/* List management */
	Songlist *	findplaylist(string filename);
	Songlist *	newplaylist(string filename);
	bool		deleteplaylist(string filename);
	Songlist *	activelist();
	bool		activatelist(Songlist *);
	int		clear(Songlist *);
	bool		crop(Songlist *, int);
	unsigned int	move(Songlist *, int offset);

	int		crossfade();
	int		crossfade(int);



	Mpd_status	*status() { return st; };
	mpd_Stats	*stats() { return _stats; };
	Song		*song() { return _song; };
	Songlist	*playlist() { return _playlist; };
	Songlist	*library() { return _library; };

	bool		increment();
	int		update(bool);
	bool		get_status();
	bool		has_new_playlist();
	bool		has_new_library();
	bool		song_changed();
	bool		state_changed();
	Songlist *	plist(int);
};
 
#endif /* _PMS_COMMAND_H_ */

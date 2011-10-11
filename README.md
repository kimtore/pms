# PMS
## Practical Music Search
Practical Music Search is an ncurses-based client for
[MPD](http://mpd.wikia.com/wiki/Music_Player_Daemon_Wiki). It has a command line
interface much like [Vim](http://www.vim.org), and supports custom colors,
layouts, and key bindings. [PMS](https://github.com/ambientsound/pms) aims to
be accessible and highly configurable.

### Installing
    $ cmake . && make && sudo make install

### Dependencies
You need to have [MPD](http://mpd.wikia.com/wiki/Music_Player_Daemon_Wiki)
installed and working before using [PMS](https://github.com/ambientsound/pms),
but not neccessarily on the same machine.

This client works best with recent
[MPD](http://mpd.wikia.com/wiki/Music_Player_Daemon_Wiki) versions (>= 0.15.0).

[PMS](https://github.com/ambientsound/pms) depends on the following libraries:

- ncurses (>= 5.0)
- glib2 (>= 2.0)
- boost_regex (>= 1.36.0) to enable regular expression searches

### Copying
Copyright (C) 2006-2011 [Kim Tore Jensen](kimtjen@gmail.com)

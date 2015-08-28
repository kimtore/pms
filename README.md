# Practical Music Search

Practical Music Search is a ncurses-based client for MPD. It has a command line
interface much like Vim, and supports custom colors, layouts, and key bindings.
PMS aims to be accessible and highly configurable.


## Compiling

PMS is a client for the [Music Player Daemon](http://musicpd.org). You need to 
have MPD installed and working before using PMS, but not neccessarily on the
same machine.

The client only works with MPD versions >= 0.15.0.

You'll need `glib >= 2.0`, `ncurses >= 5.0`, and `libmpdclient >= 2.5` to build
PMS. Installing `boost_regex >= 1.36.0` will enable regular expression searches.
In addition, if building from Git, you'll need the `intltool` package. On
Debian-based systems, you can install them by running:

```
sudo apt-get install build-essential intltool libncursesw5-dev libglib2.0-dev libmpdclient-dev
```

To build from a release tarball, run:

```
./configure && make
```

From the Git tree, run:

```
./rebuild.sh
```

Then, you may install PMS by running `sudo make install`.


## Configuration

Consult the man page for configuration options.

There are example configuration files in the `examples` directory.

Hint: type `:help` from within PMS to show a list of all current keyboard
bindings.


## Bugs, feature requests, etc.

There are many bugs. Please report them if you discover them.

Please use the [issue tracker](https://github.com/ambientsound/pms/issues) to
report bugs, or send them to the author's e-mail address.


## Authors

Copyright (c) 2006-2015 Kim Tore Jensen <kimtjen@gmail.com>.

* Kim Tore Jensen <<kimtjen@gmail.com>>
* Bart Nagel <<bart@tremby.net>>

The source code and latest version can be found at Github:
<https://github.com/ambientsound/pms>.

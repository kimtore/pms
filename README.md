# Practical Music Search

Practical Music Search is a ncurses-based client for MPD. It has a command line
interface much like Vim, and supports custom colors, layouts, and key bindings.
PMS aims to be accessible and highly configurable.


## Deprecation warning

**The 0.42.x branch is based on the same code as PMS 0.42, the latest official release from May 2010. This branch has been discontinued, and work is being focused on the Go version in the master branch. You are welcome to contribute to the legacy code, but please note that no new releases will be made from this branch.**


## Compiling

PMS is a client for the [Music Player Daemon](http://musicpd.org). You need to 
have MPD installed and working before using PMS, but not neccessarily on the
same machine.

The client only works with MPD versions >= 0.15.0.

You'll need `ncurses >= 5.0` and `libmpdclient >= 2.5` to build PMS.
If your c++ compiler supports c++11's regex (like `gcc-c++ >= 4.9`),
it will enable regular expression searches.

To install the dependencies on Debian-based systems you may run:
```
sudo apt-get install build-essential libncursesw5-dev libmpdclient-dev
```

Pandoc is required to build the manpage (it is created by default, unless pandoc is missing).

PMS can be build with CMake or GNU autotools.

### Building with CMake

The following commands are required to build and install PMS with CMake:
```
cmake .
make
sudo make install
```

### Building with autotools

If building from Git with autotools, you'll need the `intltool` package as well.

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

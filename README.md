# Practical Music Search

[![Build Status](https://travis-ci.org/ambientsound/pms.svg?branch=go)](https://travis-ci.org/ambientsound/pms)

Practical Music Search is an interactive console-based client for the [Music Player Daemon](https://www.musicpd.org/), written in Go. The interface is similar to Vim, and features lightning fast full-text searches, sorting, custom colors, configurable layouts, visual selection, keyboard shortcuts, and much more.

This software was previously written in C++. The master branch now contains a rewrite, currently implemented in Go. This branch is still somewhat experimental, and much of the old functionality is missing (see [Roadmap](#roadmap)), but it is under rapid development, and is usable for basic tasks and a little more. The full-text search is very fast; consider giving it a try!


## Running

You are assumed to have a working [Go development environment](https://golang.org/doc/install). PMS requires Go >= 1.8.

To install the application and dependencies, and run PMS, assuming you have $GOBIN in your path:

```
git clone https://github.com/ambientsound/pms $GOPATH/src/github.com/ambientsound/pms
cd $GOPATH/src/github.com/ambientsound/pms
make
pms
```

If PMS crashes, it's very helpful with debugging information. This can be obtained using the debug log:

```
pms --debug /tmp/pms.log
```


## Requirements

PMS wants to build a search index from MPD's database. In order to be truly practical, PMS must support fuzzy matching, scoring, and sub-millisecond full-text searches. This is accomplished by leveraging [Bleve](https://github.com/blevesearch/bleve), a full-text search and indexing library.

A full-text search index takes up both space and memory. For a library of about 30 000 songs, you should expect using about 500MB of disk space and around 1GB of RAM.

PMS is multithreaded and will benefit from multicore CPUs.


## Configuring

### MPD

MPD's output buffer size bust be set to something reasonably high so that the `listallinfo` command does not overflow MPD's send buffer.

```
cat >>/etc/mpd.conf<<<EOF
max_output_buffer_size "262144"
EOF
```

### PMS

PMS will honor the `MPD_HOST` and `MPD_PORT` variables.

See `pms --help` for configuration options. Configuration files are not implemented yet, but the configuration can be changed while running the program.

The default configuration can be found in [options/defaults.go](options/defaults.go).


## Roadmap

The current goal of the Go implementation is to implement most of the features found in the 0.42 branch.

This functionality is not implemented yet:

* Basic player controls (~~play~~, ~~add~~, ~~pause~~, ~~stop~~, ~~next~~, ~~prev~~, ~~volume~~, consume, repeat, single, random).
* Automatic add to queue when queue is nearing end.
* Customizable topbar.
* ~~Customizable colors~~.
* Tab completion.
* Reading configuration files.
* ~~Multiple selection~~, copy, and paste.
* Remote playlist management.
* ...and probably more.


## Contributing

There are probably many bugs, and much expected functionality is missing. You're welcome to report any bugs or feature requests by using the Github issue tracker. Code contributions are more than welcome, please submit a merge request on Github.

For general discussion about the project, or to contact the project devs, you can use the IRC channel `#pms` on Freenode.

This project adheres to the [Contributor Covenant Code of Conduct](code_of_conduct.md). By participating, you are expected to uphold this code.


## Authors

Copyright (c) 2006-2017 Kim Tore Jensen <<kimtjen@gmail.com>>.

* Kim Tore Jensen <<kimtjen@gmail.com>>
* Bart Nagel <<bart@tremby.net>>

The source code and latest version can be found at Github:
<https://github.com/ambientsound/pms>.

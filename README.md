# Practical Music Search

~~Practical Music Search is a ncurses-based client for MPD. It has a command line interface much like Vim, and supports custom colors, layouts, and key bindings. PMS aims to be accessible and highly configurable.~~

This is an experimental branch of PMS, re-implemented in Go. Much is missing yet, but the UI is usable and the connection to MPD seems stable enough. The full-text search is very fast, and definitely worth a try.


## Running

You are assumed to have a working [Go development environment](https://golang.org/doc/install).

Install the dependencies:

```
go get github.com/blevesearch/bleve  # Bleve search index
go get github.com/fhs/gompd/mpd      # MPD client library
go get github.com/jessevdk/go-flags  # POSIX-style command-line flags
```

Then install PMS itself. Assuming that you have `$GOPATH/bin` in your path,

```
go get github.com/ambientsound/pms
pms 2>/tmp/pms.log  # or /dev/null
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

Also see `pms --help` for configuration options.


## Contributing

There are many bugs and much of expected functionality is missing. You're welcome to contribute code, just send a merge request on Github.

IRC channel `#pms` on Freenode for open discussion.


## Authors

Copyright (c) 2006-2017 Kim Tore Jensen <kimtjen@gmail.com>.

* Kim Tore Jensen <<kimtjen@gmail.com>>
* Bart Nagel <<bart@tremby.net>>

The source code and latest version can be found at Github:
<https://github.com/ambientsound/pms>.

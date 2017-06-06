# Practical Music Search

[![Build Status](https://travis-ci.org/ambientsound/pms.svg?branch=go)](https://travis-ci.org/ambientsound/pms)
[![Go Report Card](https://goreportcard.com/badge/github.com/ambientsound/pms)](https://goreportcard.com/report/github.com/ambientsound/pms)
[![codecov](https://codecov.io/gh/ambientsound/pms/branch/master/graph/badge.svg)](https://codecov.io/gh/ambientsound/pms/branch/master)
[![License](https://img.shields.io/github/license/ambientsound/pms.svg)](LICENSE)

Practical Music Search is an interactive console client for the [Music Player Daemon](https://www.musicpd.org/), written in Go. Its interface is similar to Vim, and aims to be fast, configurable, and practical.

PMS has many features that involve sorting, searching, and navigating. It’s designed to let you navigate your music collection in an effective way. Some of the currently implemented features are:

* Vim-style look and feel!
* Can be configured to consume a very small amount of screen space.
* MPD player controls: play, add, pause, stop, next, prev, volume.
* Highly customizable top bar, tag headers, text styles, colors, and keyboard bindings.
* Fast library search, featuring UTF-8 normalization, fuzzy search, and scoring.
* Selecting songs, by _visual mode_, manual selection, and specific tags.
* Many forms of tracklist manipulation, such as cut, copy, paste, filter, and sort.
* Config files, tab completion, history, and much more!


## Documentation

[Documentation](doc/README.md) is available in the project repository.


## Project status

PMS is in _alpha_ stage, and is a work in progress. A more feature complete, but older, version can be found in the `0.42.x` branch.

This software was previously written in C++. The master branch now contains a rewrite, currently implemented in Go.
The current goal of the Go implementation is to implement most of the features found in the 0.42 branch.

This functionality is present in the `0.42.x` branch, but missing in master:

* Controlling player flags: consume, repeat, single, random.
* Automatically add songs to the queue when it is nearing end.
* Remote playlist management.
* ...and probably more.


## Getting started

You’re assumed to have a working [Go development environment](https://golang.org/doc/install). PMS requires Go >= 1.8.

To install the application and dependencies, and run PMS, assuming you have `$GOBIN` in your path:

```
go get -u github.com/ambientsound/pms
cd $GOPATH/src/github.com/ambientsound/pms
make
pms
```

If PMS crashes, and you want to report a bug, please include the debug log:

```
pms --debug /tmp/pms.log 2>>/tmp/pms.log
```


## Requirements

PMS wants to build a search index from MPD's database. To be truly practical, PMS must support fuzzy matching, scoring, and sub-millisecond full-text searches. This is accomplished by using [Bleve](https://github.com/blevesearch/bleve), a full-text search and indexing library.

A full-text search index takes up both space and memory. For a library of about 30 000 songs, you should expect using about 500 MB of disk space and around 1 GB of RAM.

PMS is multithreaded and benefits from multicore CPUs.


## Contributing

There are bugs, and much of the expected functionality is missing. Code contributions are warmly received through merge requests on Github. You're also welcome to report any bugs or feature requests by using the Github issue tracker.

For general discussion about the project, or to contact the project devs, you can use the IRC channel `#pms` on Freenode.

This project adheres to the [Contributor Covenant Code of Conduct](code_of_conduct.md). By participating, you are expected to uphold this code.


## Authors

Copyright (c) 2006-2017 Kim Tore Jensen <<kimtjen@gmail.com>>.

* Kim Tore Jensen <<kimtjen@gmail.com>>
* Bart Nagel <<bart@tremby.net>>

The source code and latest version can be found at Github:
<https://github.com/ambientsound/pms>.

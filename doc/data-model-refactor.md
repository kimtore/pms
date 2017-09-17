# Refactoring how and where data is stored

## Background

There are many different kinds of data stored within PMS. At the time of
writing, they are scattered throughout different parts of the program. As the
application grows more complex, the internal data representation will need to
be consolidated.

## Goals

* Create a single entry point for accessing global data.
* Thread-safe access and manipulation.

## Which data is stored?

The following data collections are shown in list views, directly to the user:

* Track lists
  * Queue
  * Library (should be read/write, but undoable)
  * Ephemeral lists (search results, etc.)
  * Remote playlists (*not implemented*)
* Help screen (*not implemented*)
* Outputs (*not implemented*)
* File browser (*not implemented*)
* Album browser (*not implemented*)

Other collections, which are not shown in a list view, but should still be accessible from components:

* Clipboards
* Current song
* Keyboard bindings
* Library (canonical read-only copy)
* MPD statistics
* Options
* Player status
* Search history
* Search index (bound to the library, and might possibly also have temporary in-memory indices for other track lists)
* Tab completion history
* Track list editing history (*not implemented*, also known as undo/redo)

## Components using data

* Commands
* Top bar fragments
* Command-line input

# Introduction

## Configuration

By default, PMS tries to read your configuration from
`$HOME/.config/pms/pms.conf`.
If you defined paths in either `$XDG_CONFIG_DIRS` or `$XDG_CONFIG_HOME`, PMS will look for `pms.conf` there.

```
# Sample PMS configuration file.
# All whitespace, newlines, and text after a hash sign will be ignored.

# The 'center' option will make sure the cursor is always centered on screen.
set center

# Some custom keyboard bindings.
bind <Alt-Left> cursor prevOf year    # jump to previous year.
bind <Alt-Right> cursor nextOf year   # jump to next year.

# Pink statusbar.
style statusbar black darkmagenta

# Minimalistic topbar.
set topbar="Now playing: ${tag|artist} \\- ${tag|title} (${elapsed})"
```

See the [default configuration](/options/defaults.go).


## Basic movement

The default bindings for movement are similar to vim.

`j` and `k` move down and up,
`gt` and `gT` (or just `t` and `T`) move forward and back between lists.
Use `<Ctrl-F>` and `<Ctrl-B>` to move a page down or up,
or `<Ctrl-D>` and `<Ctrl-U`> for half a page at a time.
`gg` and `G` go to the very top and bottom of the list,
while `H`, `M`, and `L` go to the top, middle, and bottom of the current viewport.

You can also move quickly from album to album using `b` and `e`,
which are examples of [`cursor prevOf` and `cursor nextOf` commands](commands.md#move-the-cursor-and-viewport).


## Adding tracks to the playlist

A highlighted track (or selection of tracks) can be added to the current list with `a`,
or added and played with `<Enter>`.
(When the current playlist is focused, `<Enter>` will just play, rather than also adding a duplicate.)

`x`, meanwhile, will delete the highlighted track from the list.


## Searching for tracks

PMS employs a very fast and powerful search engine called _Bleve_.
The following is an example on how to do a search in PMS:

To start a search, type `/` (or `:inputmode search`).
The tracklist will be cleared, and a slash will appear in the statusline.
Type at least two characters to start searching.
The tracklist will update itself as you type.

Search results will be sorted by match score.
If you want to sort your search result, press `<Ctrl-S>` (or type `:sort`) to sort by the default sort parameters.

To drill down into the search, highlight a song,
then press `<Ctrl-J>` (or type `:isolate artist`) to show all tracks with the same artist,
or `<Ctrl-T>` (`:isolate albumartist album`) to show all tracks in the same album.

To select tracks, type `m` (`:select toggle`) to mark one at a time,
or use the visual selection by typing `v` (`:select visual`).
You could also type `&` (`:select nearby albumartist album`) to select the entire album.
Press `a` (`:add`) to add the selected songs to the queue,
or `<Enter>` (`:play selection`) to play them immediately.


## Known issues

If having connection problems, you might be hitting a buffer limit in MPD.
It may help to configure your MPD server according to [configuring PMS and MPD](mpd.md).

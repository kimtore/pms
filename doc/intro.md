# Introduction

## Configuration

By default, PMS tries to read your configuration from
`$HOME/.config/pms/pms.conf`. If you defined paths in either `$XDG_CONFIG_DIRS`
or `$XDG_CONFIG_HOME`, PMS will look for `pms.conf` there.

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


## Searching for tracks

PMS employs a very fast and powerful search engine called _Bleve_. The
following is an example on how to do a search in PMS:

To start a search, type `/`. The tracklist will be cleared, and a slash will
appear in the statusline. Type at least two characters to start searching. The
tracklist will update itself as you type.

Search results will be sorted by match score. If you want to sort your search
result, press `<Ctrl-S>` to sort by the default sort parameters.

To drill down into the search, select a song, then press `<Ctrl-J>` to show all
tracks with the same artist, or `<Ctrl-T>` to show all tracks in the same
album.

To select tracks, type `m`, or use the visual selection by typing `v`. You
could also type `&` to select the entire album. Press `a` to add the selected
songs to the queue, or `<Enter>` to play them immediately.


## Known issues

If having connection problems, you might be hitting a buffer limit in MPD.
Please configure your MPD server according to [configuring PMS and
MPD](mpd.md).

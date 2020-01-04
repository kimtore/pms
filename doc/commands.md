# Commands

Practical Music Search is driven by _commands_.
Commands are strings of text, and can be entered in the [_multibar_](#switching-input-modes) or in configuration files.

Below is a list of commands recognized by PMS, along with their parameters and description.

Literal text is spelled out normally. Placeholders are enclosed in `<angle brackets>`, and optional parameters are enclosed in `[square brackets]`.


## Move the cursor and viewport

These commands are split into `cursor` and `viewport` namespaces.
`cursor` commands primarily move the cursor, and `viewport` commands primarily move the viewport.

* `cursor current`

  Move the cursor to MPD's current track.

* `cursor up`  
  `cursor down`

  Move the cursor up or down one row.

* `cursor home`  
  `cursor end`

  Move the cursor to the very first or last track in the current list.

* `cursor high`  
  `cursor middle`  
  `cursor low`

  Move the cursor to the top, middle, or bottom of the current viewport.

* `cursor nextOf <tag> [<tag> [...]]`

  Move the cursor down to the next track on the list
  where any of the given tags do not match the corresponding tag of the current track.

  For example, `cursor nextOf album musicbrainz_albumid` will move to the first track of the next album;
  more specifically, the first track where either the album title or the Musicbrainz album ID does not match the corresponding tag on the current track.

* `cursor prevOf <tag> [<tag> [...]]`

  Move the cursor up to the last track on the list in sequence
  where all of the given tags match the corresponding tags of the current track.
  If the current song *was* the last match, continue searching upwards until the tags differ again.

  For example, `cursor prevOf artist musicbrainz_artistid` will move the cursor up the list
  to the first track by the current artist;
  more specifically, to the final track encountered in sequence where both of these tags match.
  If used on the top track by an artist, the cursir will move up further to the first track of the previous artist.

* `cursor random`

  Move the cursor to a random position in the current list.

* `cursor +<N>`  
  `cursor -<N>`

  Move the cursor by a particular number of rows.
  Negative numbers move the cursors up, and positive numbers move the cursor down.

* `cursor <N>`

  Move the cursor to an absolute position in the current list, where `0` is the very first track.

* `viewport up`  
  `viewport down`

  Move the viewport up or down one row;
  leave the cursor on its current song if possible.


* `viewport halfpageup`  
  `viewport halfpgup`  
  `viewport halfpaged[ow]n`  
  `viewport halfpgdn`

  Move the viewport up or down a number of rows equal to half the viewport height.
  Independently, move the cursor up or down the same number.

* `viewport pageup`  
  `viewport pgup`  
  `viewport paged[ow]n`  
  `viewport pgdn`

  Move the viewport up or down one full page (actually slightly less in most cases),
  leaving the cursor on its current song where possible.

* `viewport high`  
  `viewport low`

  Move the viewport up as high or as low as possible while leaving the cursor in view,
  still pointing to the same song.
  (When `center` is set the cursor will not end up pointing to the same song.)

* `viewport middle`

  Move the viewport so that the cursor is in the middle of the viewport,
  still pointing to the same song.

## Manipulating lists

These commands switch between, create, and edit lists.

* `list next`  
  `list prev`

  Switch to the next or previous list.

* `list <N>`

  Switch to the list with the given index.

* `list duplicate`

  Duplicate the current list.

* `list remove`

  Remove the currently visible list, if possible.

* `isolate <tag> [<tag> [...]]`

  Search for tracks with similar tags to the current [selection](#selecting-tracks), and create a new tracklist with the results.
  The tracklist is sorted by the default sort criteria.

  See also [`inputmode search`](#switching-input-modes) for another way to create new lists.

* `sort [<tag> [...]]`

  Sort the current tracklist by the tags specified in the `sort` option if no tags are given, or otherwise by the specified tags.
  The most significant sort criterion is specified last.

  The first sort is performed as an unstable sort, while the remainder use a stable sorting algorithm.

### Adding, removing, and moving tracks

* `add [<uri> [...]]`

  Add one or more files or URIs to the queue.
  If no parameters are given, the current [selection](#selecting-tracks) is assumed.

  See also [`play cursor` and `play selection`](#controlling-playback).

* `yank`  
  `copy`

  Replace the clipboard contents with the currently selected tracks.

* `cut`

  Remove the current [selection](#selecting-tracks) from the tracklist, and replace the clipboard contents with the removed tracks.

* `paste [after]`  
  `paste before`

  Insert the contents of the clipboard after (this is default) or before the cursor position.


## Selecting tracks

The `select` commands allow the tracklist selection to be manipulated.

* `select toggle`

  Toggle selection status for the track under the cursor.

  When used from visual mode, all tracks currently in the visual selection will have their manual selection status toggled, and visual mode is switched off.

* `select visual`

  Toggle _visual mode_ and anchor the selection on the track under the cursor.

* `select nearby <tag> [<tag> [...]]`

  Set the visual selection to nearby tracks with the same specified tags as the track under the cursor.
  If there is already a visual selection, it will be cleared instead.


## Controlling playback

* `prev[ious]`

  Skip back to the previous track.

* `next`

  Skip to the next track.

* `pause`

  Pause or resume playback.

* `play`

  Play MPD's current song, or start from the beginning of the queue if there is none.

* `play cursor`

  Add the song under the cursor to the queue if necessary, and start playing.

* `play selection`

  Add the entire [selection](#selecting-tracks) to the queue, and start playing from the first selected song.
  If there is no selection, fall back to the song under the cursor.

* `seek +<N>`  
  `seek -<N>`

  Seek relatively by a given number of seconds.

* `seek <N>`

  Seek to a particular point in the song, measured in seconds.

* `stop`

  Stop playback.

* `single [toggle]`  
  `single on`  
  `single off`

  Toggle MPD's single mode playback style, or switch it on or off.

### Controlling the volume

These commands control the volume. The volume range is from 0 to 100.

* `volume <N>`

  Set the volume to an absolute value.

* `volume +<N>`  
  `volume -<N>`

  Adjust the volume by a relative value.

* `volume mute`

  Toggle mute.


## Switching input modes

* `inputmode normal`

  Switch to normal mode, where key bindings take effect.

* `inputmode input`

  Switch to input mode: focus the multibar, where commands can be typed in.

* `inputmode search`

  Switch to search mode, where searches execute as you type.

  When `<Enter>` is pressed from search mode, the result is a new list containing the current search results.


## Customizing PMS

### Setting global options

The command `set` or its shorthand `se` can be used to change global program options at runtime.
The [list of available options](options.md) is documented elsewhere.

* `set <option>=<value>`

  Set a non-boolean option to a particular value.

* `set <option>`  
  `set no<option>`

  Switch a boolean option on or off.

* `set inv<option>`  
  `set <option>!`

  Toggle a boolean option.

* `set <option>?`

  Query the current value of an option.

### Setting key sequences

These commands bind and unbind key sequences to commands.

A _key sequence_ can have any number of elements, each of which is any of:

* a letter
* a "special" key enclosed in angle brackets, such as `<space>` or `<f1>`
* a key with modifiers, enclosed in angle brackets, such as `<Ctrl-X>`, `<Alt-A>`, or `<Shift-Escape>`

Modifier keys are `Ctrl`, `Alt`, `Meta`, and `Shift`.
The first letter of these four words can also be used as a shorthand.

Special keys are too numerous to list here, but can be found in the [complete list of special keys](/keysequence/names.go).

Regular keys such as letters, numbers, symbols, unicode characters, etc. will never have the `Shift` modifier key.
Generally, terminal applications have far less insight into keyboard activity than graphical applications,
and therefore you should avoid depending too much on availability of modifiers or any specific keys.

Contexts are a way to make key bindings context sensitive. Choose between `global`, `library`, `tracklist`, and `windows`.
You can bind a key sequence to multiple contexts. The local context takes precedence, so a sequence bound to
the `tracklist` context will always be attempted before `global`.

* `bind <context> <key sequence> <command>`

  Configure a specific keyboard input sequence to execute a command.

* `unbind <context> <key sequence>`

  Unbind a key sequence.

### Setting styles

* `style <name> [<foreground> [<background>]] [bold] [underline] [reverse] [blink]`

  Specify the style of a UI item.
  See the [styling guide](styling.md#text-style) for details.

  The keywords `bold`, `underline`, `reverse`, and `blink` can be specified literally.
  Any keyword order is accepted, but the background color, if specified, must come after the foreground color.


## Miscellaneous

* `print [<tag> [...]]`

  Show the contents of the given tag(s) for the track under the cursor.

* `q[uit]`

  Exit the program. Any unsaved changes will be lost.

* `redraw`

  Force a screen redraw. Useful if rendering has gone wrong.

* `show logs`  
  `show music`

  Switch between different views. `music` will enable you to cycle between song lists,
  while `logs` shows an event log of what's happened so far in your session.

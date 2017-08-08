# Commands

Practical Music Search is driven by _commands_. Commands are sentences of text, and can be entered in the _multibar_ or in configuration files.

Below is a list of commands recognized by PMS, along with their parameters and description.

Literal text is spelled out normally, along with `<required>` parameters and `[optional]` parameters.


## add

```
add [uri] [uri...]
```

Adds one or more files or URI's to the queue. If no parameters are given, the current tracklist selection is assumed.

## bind

```
bind <key sequence> <command>
```

Configures a specific keyboard input sequence to execute a command.

The _key sequence_ may be either
* a string of letters,
* a string of `<special>` keys, such as `<space>` or `<f1>`
* a string of keys with modifiers, such as `<Ctrl-X>`, `<Alt-A>`, or `<Shift-Escape>`,
* or a combination of all three, such as `<Ctrl-X><Delete>quit`.

Modifier keys are `Ctrl`, `Alt`, `Meta`, and `Shift`. The first letter of these four words can also be used as a shorthand.

Special keys are too many to list here, but can be found in the [complete list of special keys](keysequence/names.go).

Regular keys such as letters, numbers, symbols, unicode characters, etc. will never have the `Shift` modifier key. Generally, terminal applications have far less visibility into keyboard
activity than graphical applications. Hence, you should avoid depending overly much on availability of modifiers, or the availability of any specific keys.

## copy

`copy` is an alias for [yank](#yank).

## cursor

```
cursor current
cursor down
cursor end
cursor home
cursor nextOf [tag] [tag...]
cursor pagedn
cursor pagedown
cursor pageup
cursor pgdn
cursor pgup
cursor prevOf [tag] [tag...]
cursor random
cursor up
cursor <+N>
cursor <-N>
cursor <N>
```

Moves the cursor. `current` moves to the currently playing song. `random` picks a random song from the tracklist.

The `nextOf` and `prevOf` commands locates the next or previous track that has different tags from the track under the cursor. Multiple tags are allowed.

The cursor can also be set to a relative or absolute position using the three last forms.

## cut

```
cut
```

Removes the current selection from the tracklist, and replace the clipboard contents with the removed tracks.

## inputmode

```
inputmode normal
inputmode input
inputmode search
```

Switches between modes. Normal mode is where key bindings take effect. Input mode allows commands to be typed in, while search mode executes searches while you type.

## isolate

```
isolate <tag> [tag...]
```

Searches for songs with similar tags as the current selection, and creates a new tracklist. The tracklist is sorted by the default sort criteria.

## list
## next
## pause
## paste

```
paste
paste before
paste after
```

Insert the contents of the clipboard before or after the cursor position. If no position is given, `after` is assumed.

## play

```
play
play cursor
play selection
```

Start playing. Without any parameters, `play` will resume playing MPD's current song, or start from the beginning of the queue.

If invoked with the `cursor` argument, `play` will add the song under the cursor to the queue if necessary, and start playing.

The `selection` argument is like `cursor`, but adds the entire selection to the queue, and starts playing from the first selected song. If there is no selection, fall back to the cursor.

## prev
## previous
## print
## q

`q` is an alias for [quit](#quit).

## quit

```
quit
```

Exit the program. Any unsaved changes will be lost.

## redraw

```
redraw
```

Force a screen redraw. Useful if rendering has gone wrong.

## se

`se` is an alias for [set](#set).

## select

```
select toggle
select visual
select nearby <tag> [tag...]
```

Manipulate tracklist selection.

`visual` will toggle _visual mode_ and anchor the selection on the track under the cursor.

`toggle` will toggle selection status for the track under the cursor. If in visual mode when using `toggle`, the visual selection will be converted to manual selection, and visual mode switched off.

`nearby` will set the visual selection to nearby tracks having the same specified _tags_ as the track under the cursor. If there is already a visual selection, it will be cleared instead.

## seek

```
seek <+N>
seek <-N>
seek <N>
```

Skip forwards or backwards in the current track. It is also possible to seek to an absolute position.

## set

```
set option=value
set option
set nooption
set invoption
set option?
set option!
```

Change global program options. Boolean option are enabled with `set option` and disabled with `set nooption`. The value can be flipped using `set invoption` or `set option!`. Values can be queried using `set option?`.

### single

```
single on
single off
single toggle
```

`single` is an alias for `single toggle`.

Turn MPD's single mode playback style on or off.

## sort

```
sort
sort <tag> [tag...]
```

Sort the current tracklist by the specified tags. The most significant sort criteria is specified last. The first sort is performed as an unstable sort, while the remainder is a stable sort.

If no tags are given, the tracklist is sorted by the tags specified in the `sort` option.

## stop

```
stop
```

Stops playback.

## style

```
style <name> [foreground [background]] [bold] [underline] [reverse] [blink]
```

Specify the style of an UI item. Please see the [styling guide](styling.md#text-style) for details.

The keywords `bold`, `underline`, `reverse`, and `blink` can be specified literally. Any keyword order is accepted, but the foreground color must come prior to the background color, if specified.

## unbind

```
unbind <key sequence>
```

Unbind a given sequence.

See [bind](#bind).

## volume

```
volume <+N>
volume <-N>
volume <N>
volume mute
```

Set or mute the volume, either to an absolute or relative value. The volume range is `0-100`.

## yank

```
yank
```

Replaces the clipboard contents with the currently selected tracks.

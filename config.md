# Configuring PMS

This document aims to document all configuration options, commands, variables, and functionality.

Literal text is spelled out normally, along with `<required>` parameters and `[optional]` parameters.

## Options

### columns

```
set columns=artist,track,title,album,year,time
```

Defines which tags should be shown in the tracklist.

### sort

```
set sort=file,track,disc,album,year,albumartistsort
```

Set default sort order, for when using `sort` without any parameters.

### topbar

```
set topbar="|$shortname $version||;${tag|artist} - ${tag|title}||${tag|album}, ${tag|year};$volume $mode $elapsed ${state} $time;|[${list|index}/${list|total}] ${list|title}||;;"
```

Define the layout and visible items in the _top bar_.

The top bar is an informational area where information about the current state of both MPD and PMS is shown.

#### Special characters

* `|` divides the line into one more piece.
* `;` starts a new line.
* `$` selects a variable. You can use either `${variable|parameter}` or `$variable`.
* `\` escapes the special meaning of the next character, so it will be printed literally.

#### Variables

* `${elapsed}` is the time elapsed in the current track.
* `${list}`
    * `${list|index}` is the numeric index of the current tracklist.
    * `${list|title}` is the title of the current tracklist.
    * `${list|total}` is the total number of tracklists.
* `${mode}` is the status of the player switches `consume`, `random`, `single`, and `repeat`, printed as four characters (`czsr`).
* `${shortname} is the short name of this program.`
* `${state}`
    * `${state}` is the current player state, represented as a two-character ASCII art.
    * `${state|unicode}` is the current player state, represented with unicode symbols.
* `${tag|<tag>}` is a specific tag in the currently playing song.
* `${time}` is the total length of the current track.
* `${version}` is the git version this program was compiled from.
* `${volume}` is the current volume, or `MUTE` if the volume is zero.

#### Examples

One line, showing player state, artist, and track title.

```
set topbar="${state|unicode} ${tag|artist} - ${tag|title}"
```

One line, with left- and right-justified text.

```
set topbar="this text is to the left|this text is to the right"
```

Three lines containing center-justified text on the middle line.

```
set topbar=";|this text is in the center||;;"
```

## Commands

### add

```
add [uri] [uri...]
```

Adds one or more URI's to MPD's queue. If no URI's are specified, the current tracklist selection is assumed.

### bind

```
bind <keysequence> <command>
```

Binds a key sequence to execute a command. The key sequence is a combination of literal keys and special keys. Special keys are represented by wrapping them in angle brackets, e.g. `<F1>` or `<C-c>`.

For reference, please see the [complete list of special keys](input/parser/keynames.go).

### cursor

```
cursor current
cursor down
cursor end
cursor home
cursor next-of [tag] [tag...]
cursor pagedn
cursor pagedown
cursor pageup
cursor pgdn
cursor pgup
cursor prev-of [tag] [tag...]
cursor random
cursor up
cursor <+N>
cursor <-N>
cursor <N>
```

Moves the cursor. `current` moves to the currently playing song. `random` picks a random song from the tracklist.

The `next-of` and `prev-of` commands locates the next or previous track that has different tags from the track under the cursor. Multiple tags are allowed.

The cursor can also be set to a relative or absolute position using the three last forms.

### inputmode

```
inputmode normal
inputmode input
inputmode search
```

Switches between modes. Normal mode is where key bindings take effect. Input mode allows commands to be typed in, while search mode executes searches while you type.

### isolate
### list
### next
### pause
### play
### prev
### previous
### print
### q
### quit
### redraw
### remove
### se
### select

```
select toggle
select visual
```

Manipulate tracklist selection. `visual` will toggle _visual mode_ and anchor the selection on the track under the cursor.

`toggle` will toggle selection status for the track under the cursor. If in visual mode when using `toggle`, the visual selection will be converted to manual selection, and visual mode switched off.

### seek

```
seek <+N>
seek <-N>
seek <N>
```

Skip forwards or backwards in the current track. It is also possible to seek to an absolute position.

### set

```
set option=value
set option
set nooption
set invoption
set option?
set option!
```

Change global program options. Boolean option are enabled with `set option` and disabled with `set nooption`. The value can be flipped using `set invoption` or `set option!`. Values can be queried using `set option?`.

### sort
### stop

```
stop
```

Stops playback.

### style
### volume

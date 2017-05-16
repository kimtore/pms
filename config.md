# Configuring PMS

This document aims to document all configuration options, commands, variables, and functionality.

Literal text is spelled out normally, along with `<required>` parameters and `[optional]` parameters.

## add

```
add [uri] [uri...]
```

Adds one or more URI's to MPD's queue. If no URI's are specified, the current tracklist selection is assumed.

## bind

```
bind <keysequence> <command>
```

Binds a key sequence to execute a command. The key sequence is a combination of literal keys and special keys. Special keys are represented by wrapping them in angle brackets, e.g. `<F1>` or `<C-c>`.

For reference, please see the [complete list of special keys](input/parser/keynames.go).

## cursor

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

## inputmode

```
inputmode normal
inputmode input
inputmode search
```

Switches between modes. Normal mode is where key bindings take effect. Input mode allows commands to be typed in, while search mode executes searches while you type.

## isolate
## list
## next
## pause
## play
## prev
## previous
## print
## q
## quit
## redraw
## remove
## se
## select

```
select toggle
select visual
```

Manipulate tracklist selection. `visual` will toggle _visual mode_ and anchor the selection on the track under the cursor.

`toggle` will toggle selection status for the track under the cursor. If in visual mode when using `toggle`, the visual selection will be converted to manual selection, and visual mode switched off.

## set
## sort
## stop
## style
## volume

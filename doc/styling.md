# Styling guide

## Text style

All UI elements in PMS can be styled with custom colors and text attributes. Styles are set with the `style` command, please see the [commands documentation](commands.md#style) for details.

Colors are either literal names, such as `red`, `yellow`, or `green`, or hexadecimal color values such as `@ff0077` (the hash sign `#` is used for comments; the `@` character is used instead). A full list of supported color names can be found in the [tcell documentation](https://github.com/gdamore/tcell/blob/master/color.go#L820).

### Tags

Any _tag_, such as `artist`, `album`, `date`, `time`, etc. can be styled. These tags are not included in this list. Their styles apply both to the tracklist and in the top bar. All tag names are in lowercase.

### Tracklist

* `allTagsMissing` - song style in the tracklist when all the _essential_ tags are missing (`artist`, `album`, and `title`).
* `currentSong` - color of the entire line in the tracklist, highlighting the currently playing song.
* `cursor` - color of the entire line in the tracklist, highlighting the cursor position.
* `header` - column headers showing tag names.
* `mostTagsMissing` - song style in the tracklist when _most_ tags are missing (`artist` and `title`).
* `selection` - line color of selected songs.

### Top bar

Please see [top bar](#top-bar) for corresponding variables.

* `elapsed` - corresponds to `${elapsed}`.
* `listIndex` - corresponds to `${list|index}`.
* `listTitle` - corresponds to `${list|title}`.
* `listTotal` - corresponds to `${list|total}`.
* `mute` - the color of the `${volume}` widget when the volume is zero.
* `shortName` - corresponds to `${elapsed}`.
* `state` - corresponds to `${elapsed}`.
* `switches` - corresponds to `${elapsed}`.
* `tagMissing` - the color used for tag styling when the tag is missing.
* `topbar` - the default color of the top bar text and whitespace.
* `version` - corresponds to `${elapsed}`.
* `volume` - corresponds to `${volume}`.

### Statusbar

* `commandText` - text color when writing text in command input mode.
* `errorText` - text color of error messages in the status bar.
* `readout` - position readout at the bottom right.
* `searchText` - text color when searching.
* `sequenceText` - text color of uncompleted keyboard bindings.
* `statusbar` - normal text in the status bar.
* `visualText` - text color of the `-- VISUAL --` text when selecting songs in visual mode.


## Top bar

The top bar is an informational area where information about the current state
of both MPD and PMS is shown. The height and contents is very flexible, and can
be made to look exactly as you want to.

### Examples

Borders are added for emphasis, and are not part of the rendered text.

One line, showing player state, artist, track title, and time elapsed.

```
set topbar="${state|unicode} ${tag|artist} - ${tag|title} (${elapsed})"
```
→
```
.---------------------------------------------------------------------------.
|▶ Madrugada - Black Mambo (02:41)                                          |
'---------------------------------------------------------------------------'
```

One line, with left- and right-justified text.

```
set topbar="this text is to the left|this text is to the right"
```
→
```
.---------------------------------------------------------------------------.
|this text is to the left                          this text is to the right|
'---------------------------------------------------------------------------'
```

Three lines containing center-justified text on the middle line.

```
set topbar=";|the ${tag|artist} is in the center||;;"
```
→
```
.---------------------------------------------------------------------------.
|                                                                           |
|                       the Spoonbill is in the center                      |
|                                                                           |
'---------------------------------------------------------------------------'
```

### Variables

* `${elapsed}` is the time elapsed in the current track.
* `${list}`
    * `${list|index}` is the numeric index of the current tracklist.
    * `${list|title}` is the title of the current tracklist.
    * `${list|total}` is the total number of tracklists.
* `${mode}` is the status of the player switches `consume`, `random`, `single`, and `repeat`, printed as four characters (`czsr`).
* `${shortname}` is the short name of this program.
* `${state}`
    * `${state}` is the current player state, represented as a two-character ASCII art.
    * `${state|unicode}` is the current player state, represented with unicode symbols.
* `${tag|<tag>}` is a specific tag in the currently playing song.
* `${time}` is the total length of the current track.
* `${version}` is the git version this program was compiled from.
* `${volume}` is the current volume, or `MUTE` if the volume is zero.

### Special characters

* `|` divides the line into one more piece.
* `;` starts a new line.
* `$` selects a variable. You can use either `${variable|parameter}` or `$variable`.
* `\` escapes the special meaning of the next character, so it will be printed literally.

# Styling guide

## Text style

All UI elements in PMS can be styled with custom colors and text attributes.
Styles are set with the `style` command; please see the [commands documentation](commands.md#setting-styles) for details.

Colors are either literal names, such as `red`, `yellow`, or `green`, or hexadecimal color values such as `@ff0077`
(the hash sign `#` is used for comments; the `@` character is used instead).
A full list of supported color names can be found in the [tcell documentation](https://github.com/gdamore/tcell/blob/master/color.go#L820).

### Tags

Any _tag_, such as `artist`, `album`, `date`, `time`, etc. can be styled.
These tags are not included in this list.
Their styles apply both to the tracklist and in the top bar.
All tag names are in lowercase.

### Tracklist

* `header`

  Column headers showing tag names.

* `allTagsMissing`

  Song style in the tracklist when all the "essential" tags are missing (`artist`, `album`, and `title`).

* `mostTagsMissing`

  Song style in the tracklist when _most_ tags are missing (`artist` and `title`).

* `currentSong`

  Color of the entire line in the tracklist, highlighting the currently playing song.

* `selection`

  Line color of selected songs.

* `cursor`

  Color of the entire line in the tracklist, highlighting the cursor position.

### Log console

* `logLevel`

  Log level, such as DEBUG, INFO, or ERROR.

* `logMessage`

  Log messages.

* `timestamp`

  Log timestamps

### Top bar

See [below](#top-bar-variables) for corresponding variables.

* `elapsedPercentage`

  Corresponds to `${elapsed|percentage}`.

* `elapsedTime`

  Corresponds to `${elapsed}`.

* `listIndex`

  Corresponds to `${list|index}`.

* `listTitle`

  Corresponds to `${list|title}`.

* `listTotal`

  Corresponds to `${list|total}`.

* `mute`

  The color of the `${volume}` widget when the volume is zero.

* `shortName`

  Corresponds to `${shortname}`.

* `state`

  Corresponds to `${state}` and `${state|unicode}`.

* `switches`

  Corresponds to `${mode}`.

* `tagMissing`

  The color used for tag styling when the tag is missing.

* `topbar`

  The default color of the top bar text and whitespace.

* `version`

  Corresponds to `${version}`.

* `volume`

  Corresponds to `${volume}`.

### Statusbar

* `commandText`

  Text color when writing text in command input mode.

* `errorText`

  Text color of error messages in the status bar.

* `readout`

  Position readout at the bottom right.

* `searchText`

  Text color when searching.

* `sequenceText`

  Text color of uncompleted keyboard bindings.

* `statusbar`

  Normal text in the status bar.

* `visualText`

  Text color of the `-- VISUAL --` text when selecting songs in visual mode.


## Top bar

The top bar is an informational area where information about the current state of both MPD and PMS is shown.
The height and contents are very flexible, and can be made to look exactly as you want them to.

### Examples

Borders are added for emphasis, and are not part of the rendered text.

* One line, showing player state, artist, track title, and time elapsed.

  ```
  set topbar="${state|unicode} ${tag|artist} - ${tag|title} (${elapsed})"
  ```
  →
  ```
  .---------------------------------------------------------------------------.
  |▶ Madrugada - Black Mambo (02:41)                                          |
  '---------------------------------------------------------------------------'
  ```

* One line, with left- and right-justified text.

  ```
  set topbar="this text is to the left|this text is to the right"
  ```
  →
  ```
  .---------------------------------------------------------------------------.
  |this text is to the left                          this text is to the right|
  '---------------------------------------------------------------------------'
  ```

* Three lines containing center-justified text on the middle line.

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

### Top bar variables

#### Playback state

* `${state}`

  The current player state, represented as a two-character ASCII art.

* `${state|unicode}`

  The current player state, represented with unicode symbols.

* `${elapsed}`

  The time elapsed in the current track.

* `${elapsed|percentage}`

  The time elapsed in the current track as a percentage of its total length.

* `${time}`

  The total length of the current track.

* `${mode}`

  The status of the player switches `consume`, `random`, `single`, and `repeat`, printed as four characters (`czsr`).

* `${volume}`

  The current volume, or `MUTE` if the volume is zero.

* `${tag|<tag>}`

  A specific tag of the currently playing song, such as `${tag|artist}`.

#### Information about the current list

* `${list|index}`

  The numeric index of the current tracklist.

* `${list|title}`

  The title of the current tracklist.

* `${list|total}`

  The total number of tracklists.

#### Miscellaneous

* `${shortname}`

  The short name of this program.

* `${version}`

  The git version this program was compiled from.

### Special characters

* `|` divides the line into one more piece.
* `;` starts a new line.
* `$` selects a variable. You can use either `${variable|parameter}` or `$variable`.
* `\` escapes the special meaning of the next character, so it will be printed literally.

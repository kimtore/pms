# Options

## center

```
set center
set nocenter
```

If set, the viewport is automatically moved so that the cursor stays in the center, if possible.


## columns

```
set columns=artist,track,title,album,year,time
```

Defines which tags should be shown in the tracklist.


## sort

```
set sort=file,track,disc,album,year,albumartistsort
```

Set default sort order, for when using `sort` without any parameters.


## topbar

```
set topbar="|$shortname $version||;${tag|artist} - ${tag|title}||${tag|album}, ${tag|year};$volume $mode $elapsed ${state} $time;|[${list|index}/${list|total}] ${list|title}||;;"
```

Define the layout and visible items in the _top bar_. Please see the [styling guide](styling.md#top-bar) for information on how to configure the top bar.

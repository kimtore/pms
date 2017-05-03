package options

// AddDefaultOptions adds internal options that can be set by the user through
// the command-line interface.
func (o *Options) AddDefaultOptions() {
	o.Add(NewStringOption("columns"))
	o.Add(NewStringOption("sort"))
}

// Defaults is the default, internal configuration file.
const Defaults string = `
set columns=artist,track,title,album,year,time
set sort=file,track,disc,album,year,albumartistsort

bind <C-c> quit
bind <C-l> redraw
bind <C-s> sort
bind <Up> cursor up
bind <Down> cursor down
bind k cursor up
bind j cursor down
bind <PgUp> cursor pgup
bind <PgDn> cursor pgdn
bind <Home> cursor home
bind <End> cursor end
bind gg cursor home
bind gc cursor current
bind G cursor end
bind t list next
bind T list previous
bind <C-d> list duplicate
bind <C-g> list remove
bind <C-j> isolate artist
bind <C-t> isolate albumartist,album
bind : inputmode input
bind / inputmode search
bind v inputmode visual
bind V inputmode visual
bind <F3> inputmode search
bind m select toggle
bind a add
bind <Delete> remove
bind x remove
bind <Enter> play cursor
bind <Space> pause
bind s stop
bind h previous
bind l next
bind + volume +2
bind - volume -2
bind M volume mute
`

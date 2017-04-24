package options

// AddDefaultOptions adds internal options that can be set by the user through
// the command-line interface.
func (o *Options) AddDefaultOptions() {
	o.Add(NewStringOption("columns", "artist,track,title,album,year,time"))
}

// Defaults is the default, internal configuration file.
const Defaults string = `
bind <C-c> quit
bind <C-l> redraw
bind <Up> cursor up
bind <Down> cursor down
bind k cursor up
bind j cursor down
bind <PgUp> cursor pgup
bind <PgDn> cursor pgdn
bind <Home> cursor home
bind <End> cursor end
bind t list next
bind T list previous
bind : inputmode input
bind / inputmode search
bind <F3> inputmode search
bind <Enter> play cursor
bind s stop
`

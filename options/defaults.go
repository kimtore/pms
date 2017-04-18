package options

// AddDefaultOptions adds internal options that can be set by the user through
// the command-line interface.
func (o *Options) AddDefaultOptions() {
	o.Add(NewStringOption("columns", "artist,track,title,album,year,time"))
}

// Defaults is the default, internal configuration file.
const Defaults string = `
bind <C-c> quit
`

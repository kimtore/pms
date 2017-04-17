package options

func (o *Options) AddDefaultOptions() {
	o.Add(NewStringOption("columns", "artist,track,title,album,year,time"))
}

const Defaults string = `
`

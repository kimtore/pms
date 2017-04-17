package options

func (o *Options) AddDefaultOptions() {
	o.Add(NewStringOption("pms", "Practical Music Search"))
	o.Add(NewStringOption("columns", "artist,track,title,album,year,time,comment"))
}

const Defaults string = `
`

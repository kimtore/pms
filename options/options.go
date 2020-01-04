package options

import (
	"fmt"
	"github.com/spf13/viper"
)

// Option names.
const (
	Center       = "center"
	Columns      = "columns"
	Limit        = "limit"
	Sort         = "sort"
	Topbar       = "topbar"
	PollInterval = "pollinterval"

	SpotifyClientID     = "spotifyclientid"
	SpotifyClientSecret = "spotifyclientsecret"

	LogFile      = "logfile"
	LogOverwrite = "logoverwrite"
)

// Option types.
const (
	boolType   = false
	intType    = 0
	stringType = ""
)

// Initialize option types.
// Default values must be defined in the Defaults string.
func init() {
	viper.Set(Center, boolType)
	viper.Set(Columns, stringType)
	viper.Set(Limit, intType)
	viper.Set(LogFile, stringType)
	viper.Set(LogOverwrite, boolType)
	viper.Set(PollInterval, intType)
	viper.Set(Sort, stringType)
	viper.Set(SpotifyClientID, stringType)
	viper.Set(SpotifyClientSecret, stringType)
	viper.Set(Topbar, stringType)
}

// Return a human-readable representation of an option.
// This string can be used in a config file.
func Print(key string, opt interface{}) string {
	switch v := opt.(type) {
	case string:
		return fmt.Sprintf("%s=\"%s\"", key, v)
	case int:
		return fmt.Sprintf("%s=%d", key, v)
	case bool:
		if !v {
			return fmt.Sprintf("no%s", key)
		}
		return fmt.Sprintf("%s", key)
	default:
		return fmt.Sprintf("%s=%v", key, v)
	}
}

// Default configuration file.
const Defaults string = `
# Global options
set nocenter
set columns=artist,title,track,album,year,time,popularity
set sort=track,disc,album,year,albumArtist
set topbar="|$shortname $version||;${tag|artist} - ${tag|title}||\\#${tag|track} ${tag|album};$volume $mode $elapsed ${state} $time;|[${list|index}/${list|total}] ${list|title}||;;"
set limit=50
set pollinterval=10

# Logging
set nologoverwrite
set logfile=

# Song tag styles
style album teal
style albumArtist teal
style artist yellow dim
style date green
style disc default
style popularity dim
style time darkmagenta
style title white
style track default
style year default
style _id gray

# Tracklist styles
style currentSong black yellow
style cursor black white
style header teal bold
style selection white blue

# Topbar styles
style elapsedTime green
style elapsedPercentage green
style listIndex darkblue
style listTitle blue bold
style listTotal darkblue
style mute red
style shortName bold
style state default
style switches teal
style tagMissing red
style topbar darkgray
style version gray
style volume green

# Other styles
style commandText default
style errorText black red
style logLevel white
style logMessage gray
style readout default
style searchText white bold
style sequenceText teal
style statusbar default
style timestamp teal
style visualText teal

# Keyboard bindings: cursor and viewport movement
bind list <Up> cursor up
bind list k cursor up
bind list <Down> cursor down
bind list j cursor down
bind list <PgUp> viewport pgup
bind list <PgDn> viewport pgdn
bind list <C-b> viewport pgup
bind list <C-f> viewport pgdn
bind list <C-u> viewport halfpgup
bind list <C-d> viewport halfpgdn
bind list <C-y> viewport up
bind list <C-e> viewport down
bind list <Home> cursor home
bind list gg cursor home
bind list <End> cursor end
bind list G cursor end
bind list gc cursor current
bind list R cursor random
bind list H cursor high
bind list M cursor middle
bind list L cursor low
bind list zb viewport high
bind list z- viewport high
bind list zz viewport middle
bind list z. viewport middle
bind list zt viewport low
bind list z<Enter> viewport low

# Tracklist specifics
bind tracklist b cursor prevOf album
bind tracklist e cursor nextOf album

# Keyboard bindings: input mode
bind global : inputmode input
bind global / inputmode search
bind global <F3> inputmode search
bind global v select visual
bind global V select visual

# Keyboard bindings: player and mixer
bind tracklist <Enter> play selection
bind global <Space> pause
bind global s stop
bind global h previous
bind global l next
bind global + volume +2
bind global - volume -2
bind global <left> seek -5
bind global <right> seek +5
bind global <Alt-M> volume mute
bind global S single

# Keyboard bindings: other
bind global <C-c> quit
bind global <C-l> redraw
bind list <C-s> sort
bind tracklist i print file
bind global gt list next
bind global gT list previous
bind global t list next
bind global T list previous
bind list <C-w>d list duplicate
bind list <C-g> list remove
bind tracklist <C-j> isolate artist
bind tracklist <C-t> isolate albumArtist album
bind tracklist & select nearby albumArtist album
bind list m select toggle
bind tracklist a add
bind tracklist <Delete> cut
bind tracklist x cut
bind tracklist y yank
bind tracklist p paste after
bind tracklist P paste before
`

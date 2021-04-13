package options

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Option names.
const (
	Center            = "center"
	Columns           = "columns"
	Limit             = "limit"
	Sort              = "sort"
	Topbar            = "topbar"
	PollInterval      = "pollinterval"
	SpotifyAuthServer = "spotifyauthserver"
	LogFile           = "logfile"
	LogOverwrite      = "logoverwrite"
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
	viper.Set(SpotifyAuthServer, stringType)
	viper.Set(Topbar, stringType)
}

// Methods for getting options from Viper.
var (
	Get       = viper.Get
	GetString = viper.Get
	GetInt    = viper.GetInt
	GetBool   = viper.GetBool
)

// Split a string option into a comma-delimited list.
func GetList(key string) []string {
	return strings.Split(viper.GetString(key), ",")
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
set topbar="${tag|artist} - ${tag|title}|$shortname $version|$elapsed $state $time;\\#${tag|track} ${tag|album}|${list|title} [${list|index}/${list|total}]|$device $mode $volume;;"
set spotifyauthserver="http://localhost:59999"
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
style deviceName teal
style deviceType teal
style elapsedTime green
style elapsedPercentage green
style listIndex teal
style listTitle teal bold
style listTotal teal
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
style currentDevice white green
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
bind global <Up> cursor up
bind global k cursor up
bind global <Down> cursor down
bind global j cursor down
bind global <PgUp> viewport pgup
bind global <PgDn> viewport pgdn
bind global <C-b> viewport pgup
bind global <C-f> viewport pgdn
bind global <C-u> viewport halfpgup
bind global <C-d> viewport halfpgdn
bind global <C-y> viewport up
bind global <C-e> viewport down
bind global <Home> cursor home
bind global gg cursor home
bind global <End> cursor end
bind global G cursor end
bind global gc cursor current
bind global R cursor random
bind global H cursor high
bind global M cursor middle
bind global L cursor low
bind global zb viewport high
bind global z- viewport high
bind global zz viewport middle
bind global z. viewport middle
bind global zt viewport low
bind global z<Enter> viewport low

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
bind tracklist a add
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

# Special windows
bind global c show library
bind global w show windows
bind windows <Enter> show selected
bind library <Enter> show selected
bind devices <Enter> device activate
bind playlists <Enter> list open

# Keyboard bindings: other
bind global <C-l> redraw
bind global <C-s> sort
bind tracklist i print file
bind global gt list next
bind global gT list previous
bind global t list next
bind global T list previous
bind global <C-w>d list duplicate
bind global <C-g> list remove
bind tracklist <C-j> isolate artist
bind tracklist <C-t> isolate albumArtist album
bind tracklist & select nearby albumArtist album
bind global m select toggle
bind global <C-a> select all
bind global <C-c> select none
bind tracklist <Delete> cut
bind tracklist x cut
bind tracklist y yank
bind tracklist p paste after
bind tracklist P paste before
`

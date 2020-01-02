// Package config holds all of visp's configuration options.

package config

import (
	"github.com/ambientsound/pms/xdg"
	"github.com/mitchellh/mapstructure"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Spotify struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

type Log struct {
	File      string
	Overwrite bool
	Level string
}

type Config struct {
	Log     Log
	Spotify Spotify
}

func decoderHook(dc *mapstructure.DecoderConfig) {
	dc.TagName = "json"
}

func init() {
	viper.SetConfigName("visp")
	viper.SetConfigType("yaml")

	for _, dir := range xdg.ConfigDirectories() {
		viper.AddConfigPath(dir)
	}

	flag.String("spotify.clientid", "", "Spotify app client ID")
	flag.String("spotify.clientsecret", "", "Spotify app client secret")
	flag.String("spotify.accesstoken", "", "Spotify access token")
	flag.String("spotify.refreshtoken", "", "Spotify refresh token")

	flag.String("log.file", "/dev/null", "Write program log to this file")
	flag.String("log.level", "info", "Log level (error, warning, info, debug, trace)")
	flag.Bool("log.overwrite", false, "Overwrite log file instead of appending")
}

func Configuration() (*Config, error) {
	var err error
	var cfg Config

	err = viper.ReadInConfig()
	switch err.(type) {
	default:
		return nil, err
	case nil, viper.ConfigFileNotFoundError:
	}

	flag.Parse()

	err = viper.BindPFlags(flag.CommandLine)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg, decoderHook)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

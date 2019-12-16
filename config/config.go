// Package config holds all of visp's configuration options.

package config

import (
	"github.com/mitchellh/mapstructure"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Spotify struct {
	Username string
	Password string
}

type Config struct {
	Spotify Spotify
}

func decoderHook(dc *mapstructure.DecoderConfig) {
	dc.TagName = "json"
	dc.ErrorUnused = true
}

func init() {
	viper.SetConfigName("visp")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME/.config/visp")

	flag.String("spotify.username", "", "Spotify username")
	flag.String("spotify.password", "", "Spotify password")
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

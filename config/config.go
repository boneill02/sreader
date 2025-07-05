package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Constants for configuration paths
const Confdir string = "/.config/sreader"
const Datadir string = "/.local/share/sreader"
const Confpath string = Confdir + "/config.toml"
const Urlspath string = Confdir + "/urls"

type Config struct {
	BG              string
	FG              string
	TitleBG         string
	TitleFG         string
	SelectedTitleBG string
	SelectedTitleFG string
	DescFG          string
	DescBG          string
	SelectedDescFG  string
	SelectedDescBG  string
	Player          string
	Browser         string
}

func LoadConfig(path string) *Config {
	// Define default configuration values
	const defaultPlayer = "mpv"
	const defaultBrowser = "firefox"

	conf := &Config{
		FG:              "#000000",
		BG:              "#FAF6F6",
		TitleBG:         "#ffffff",
		TitleFG:         "#191923",
		SelectedTitleBG: "#7FB685",
		SelectedTitleFG: "#191923",
		DescFG:          "#191923",
		DescBG:          "#FAF6F6",
		SelectedDescFG:  "#191923",
		SelectedDescBG:  "#7FB685",
		Player:          defaultPlayer,
		Browser:         defaultBrowser,
	}

	// Load environment variables
	if envPlayer := os.Getenv("PLAYER"); envPlayer != "" {
		conf.Player = envPlayer
	}
	if envBrowser := os.Getenv("BROWSER"); envBrowser != "" {
		conf.Browser = envBrowser
	}

	// Load config file
	if path != "" {
		file, err := os.Open(path)
		if err == nil {
			defer file.Close()
			toml.NewDecoder(file).Decode(conf)
		}
	}

	return conf
}

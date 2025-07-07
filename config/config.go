package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type SreaderConfig struct {
	ConfDir         string
	DataDir         string
	ConfFile        string
	URLs            []*string
	DBFile          string
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

const (
	defaultBG              string = "#000000"
	defaultFG              string = "#FFFFFF"
	defaultTitleBG         string = "#FFFFFF"
	defaultTitleFG         string = "#000000"
	defaultSelectedTitleBG string = "#7FB685"
	defaultSelectedTitleFG string = "#000000"
	defaultDescFG          string = "#000000"
	defaultDescBG          string = "#FFFFFF"
	defaultSelectedDescFG  string = "#000000"
	defaultSelectedDescBG  string = "#7FB685"
	defaultPlayer          string = "mpv"
	defaultBrowser         string = "firefox"
)

// Defaults
var (
	defaultConfDir  string = os.Getenv("HOME") + "/.config/sreader"
	defaultDataDir  string = os.Getenv("HOME") + "/.local/share/sreader"
	defaultConfPath string = defaultConfDir + "/config.toml"
	defaultDBFile   string = defaultDataDir + "/sreader.db"
	Config  *SreaderConfig = &SreaderConfig{
		ConfDir:         defaultConfDir,  // Not usable in config file
		ConfFile:        defaultConfPath, // Not usable in config file
		DataDir:         defaultDataDir,
		URLs:            nil,
		DBFile:          defaultDBFile,
		BG:              defaultBG,
		FG:              defaultFG,
		TitleBG:         defaultTitleBG,
		TitleFG:         defaultTitleFG,
		SelectedTitleBG: defaultSelectedTitleBG,
		SelectedTitleFG: defaultSelectedTitleFG,
		DescFG:          defaultDescFG,
		DescBG:          defaultDescBG,
		SelectedDescFG:  defaultSelectedDescFG,
		SelectedDescBG:  defaultSelectedDescBG,
		Player:          defaultPlayer,
		Browser:         defaultBrowser,
	}
)

func LoadConfig(path string) {
	// Make directories if non-existent
	os.MkdirAll(Config.DataDir, os.ModePerm)
	os.MkdirAll(Config.ConfDir, os.ModePerm)

	// Define default configuration
	if envPlayer := os.Getenv("PLAYER"); envPlayer != "" {
		Config.Player = envPlayer
	}
	if envBrowser := os.Getenv("BROWSER"); envBrowser != "" {
		Config.Browser = envBrowser
	}

	// Load config file
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalln("Failed to open configuration file:", err.Error())
		}
		defer file.Close()
		_, err = toml.NewDecoder(file).Decode(Config)
		if err != nil {
			log.Fatalln("Failed to parse configuration file:", err.Error())
		}
	}

	if Config.URLs == nil {
		log.Fatalln("No URLs in configuration.")
	}

}

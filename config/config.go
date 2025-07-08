package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type SreaderConfig struct {
	URLs []*string

	// Paths
	ConfFile string
	DBFile   string
	LogFile  string
	TmpDir   string

	// Colors
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

	// Keys
	UpKey      string
	DownKey    string
	LeftKey    string
	RightKey   string
	QuitKey    string
	SyncKey    string
	BrowserKey string
	PlayerKey  string

	// External applications
	Player  string
	Browser string
}

const (
	// Default paths
	defaultConfFile string = "~/.config/sreader/config.toml"
	defaultDBFile   string = "~/.local/share/sreader/sreader.db"
	defaultLogFile  string = "~/.local/share/sreader/sreader.log"
	defaultTmpDir   string = "~/.local/share/sreader"

	// Default colors
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

	// Default keys
	defaultUpKey      string = "k"
	defaultDownKey    string = "j"
	defaultLeftKey    string = "h"
	defaultRightKey   string = "l"
	defaultQuitKey    string = "q"
	defaultSyncKey    string = "r"
	defaultBrowserKey string = "o"
	defaultPlayerKey  string = "v"

	// Default external applications
	defaultPlayer  string = "mpv"
	defaultBrowser string = "firefox"
)

// Defaults
var (
	Config *SreaderConfig = &SreaderConfig{
		URLs: nil,

		// Paths
		ConfFile: defaultConfFile, // Not usable in config file
		DBFile:   defaultDBFile,
		LogFile:  defaultLogFile,
		TmpDir:   defaultTmpDir,

		// Colors
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

		// Keys
		UpKey:      defaultUpKey,
		DownKey:    defaultDownKey,
		LeftKey:    defaultLeftKey,
		RightKey:   defaultRightKey,
		QuitKey:    defaultQuitKey,
		SyncKey:    defaultSyncKey,
		BrowserKey: defaultBrowserKey,
		PlayerKey:  defaultPlayerKey,

		// External applications
		Player:  defaultPlayer,
		Browser: defaultBrowser,
	}
)

func ExpandHome(path string) string {
	if path == "" {
		return ""
	}
	if path[0] == '~' {
		return os.Getenv("HOME") + path[1:]
	}
	return path
}

func LoadConfig(path string) {
	// Load environment variables if set
	if envPlayer := os.Getenv("PLAYER"); envPlayer != "" {
		Config.Player = envPlayer
	}
	if envBrowser := os.Getenv("BROWSER"); envBrowser != "" {
		Config.Browser = envBrowser
	}
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		Config.DBFile = dataHome + "/sreader/sreader.db"
		Config.LogFile = dataHome + "/sreader/sreader.log"
		Config.TmpDir = dataHome + "/sreader"
	}

	// Load config file
	if path != "" {
		path = ExpandHome(path)
		file, err := os.Open(path)
		if err != nil {
			WriteDefaultConfig(path)
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

	// Make directories if non-existent
	dbDir := getDirectoryOfFile(Config.DBFile)
	tmpDir := ExpandHome(Config.TmpDir)
	logDir := getDirectoryOfFile(Config.LogFile)
	os.MkdirAll(tmpDir, 0700)
	os.MkdirAll(dbDir, 0700)
	os.MkdirAll(logDir, 0700)

	log.Println("Configuration loaded successfully.")
}

func WriteDefaultConfig(path string) {
	file, err := os.Create(ExpandHome(path))
	if err != nil {
		log.Fatalln("Failed to create configuration file:", err.Error())
	}
	defer file.Close()

	// Write default configuration to file
	if err := toml.NewEncoder(file).Encode(Config); err != nil {
		log.Fatalln("Failed to write default configuration:", err.Error())
	}

	log.Println("Default configuration written to", path)
	log.Println("Please edit the configuration file to set your RSS feed URLs.")
	os.Exit(0)
}

func getDirectoryOfFile(path string) string {
	if path == "" {
		return ""
	}

	path = ExpandHome(path)
	dir := path
	if idx := len(dir) - 1; idx >= 0 {
		for i := len(dir) - 1; i >= 0; i-- {
			if dir[i] == '/' {
				dir = dir[:i]
				break
			}
		}
	}

	return dir
}

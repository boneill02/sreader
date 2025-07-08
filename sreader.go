package main

import (
	"flag"
	"log"
	"os"

	"github.com/bmoneill/sreader/config"
	"github.com/bmoneill/sreader/feed"
	"github.com/bmoneill/sreader/ui"
)

func main() {
	// Modify default config file path if XDG_CONFIG_HOME is set
	confPath := config.ExpandHome(config.Config.ConfFile)
	if xdgConfDir := os.Getenv("XDG_CONFIG_HOME"); xdgConfDir != "" {
		confPath = xdgConfDir + "/sreader/config.toml"
	}

	// Parse command line flags
	confFlag := flag.String("c", confPath, "Path to the configuration file")
	syncFlag := flag.Bool("s", false, "Sync feeds and exit")
	flag.Parse()

	config.LoadConfig(*confFlag)
	feed.InitDB()

	// sync and quit if called with "-s" flag
	if *syncFlag {
		feed.Sync()
		return
	}

	writer, err := os.Create(config.ExpandHome(config.Config.LogFile))
	if err != nil {
		log.Fatalln("Failed to create log file:", err.Error())
	}

	log.SetOutput(writer)

	feeds := feed.GetFeeds()
	ui := ui.Init(feeds)
	if _, err := ui.Run(); err != nil {
		log.Fatalln("Failed to start UI", err.Error())
	}
}

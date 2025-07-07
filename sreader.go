package main

import (
	"log"
	"os"

	"github.com/bmoneill/sreader/config"
	"github.com/bmoneill/sreader/feed"
	"github.com/bmoneill/sreader/ui"
)

func main() {
	config.LoadConfig(os.Getenv("HOME") + config.Config.ConfFile)

	feed.Init()

	/* sync and quit if called with the arg "sync" */
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "sync":
			feed.Sync()
			return
		}
	}

	writer, err := os.Create(config.Config.DataDir + "/sreader.log")
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

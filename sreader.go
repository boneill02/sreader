package main

import (
	"log"
	"os"

	"github.com/boneill02/sreader/config"
	"github.com/boneill02/sreader/feed"
	"github.com/boneill02/sreader/ui"
)

func main() {
	feed.Init()

	/* sync and quit if called with the arg "sync" */
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "sync":
			feed.Sync()
			return
		}
	}

	conf := config.LoadConfig(os.Getenv("HOME") + config.Confpath)
	feeds := feed.LoadFeeds()
	ui := ui.Init(feeds, conf)
	if _, err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

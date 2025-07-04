package main

import (
	"os"

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

	feeds := feed.LoadFeeds()
	tui := ui.Init(feeds)
	err := tui.Run()
	if err != nil {
		panic(err)
	}
}

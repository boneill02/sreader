/*
	sreader.go: A simple RSS reader written in Go. See LICENSE for
	copyright details.
	Author: Ben O'Neill <ben@oneill.sh>

	This uses hex-encoded SHA1 hashes of the desired feeds' URLs to
	store them. It's much better than putting the URL or feed name as
	a filename because that would look ugly.

	TODO show when entry has been read*
	TODO view for all entries from all feeds at once (most recent first)
	TODO filter out read entries*
	TODO searching through entries via regex
	TODO nickname feeds (so this is better than newsboat)*
	TODO add config file for keybindings and default browser/player
	TODO show refreshing status instead of freezing

	*: we might need to use a database because of these which sucks,
	maybe i will make it JSON, maybe plaintext, but I will try to
	avoid using a database as much as I can.
*/

package main

import (
	"os"
	"github.com/boneill02/sreader/feed"
	"github.com/boneill02/sreader/ui"
)

func main() {
	/* sync and quit if called with the arg "sync" */
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "sync":
			feed.Sync()
			return
		}
	}

	feeds := feed.Init()
	tui := ui.Init(feeds)
	err := tui.Run()
	if err != nil {
		panic(err)
	}
}

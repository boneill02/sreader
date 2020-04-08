/*
	sreader.go: A simple RSS reader written in Go. See LICENSE for copyright details.
	Author: Ben O'Neill <ben@benoneill.xyz>

	This uses hex-encoded SHA1 hashes of the desired feeds' URLs to store them. It's
	much better than putting the URL or feed name as a filename because that would
	look ugly.

	TODO show when entry has been read*
	TODO view for all entries from all feeds at once (most recent first)
	TODO filter out read entries*
	TODO searching through entries via regex
	TODO nickname feeds (so this is better than newsboat)*
	TODO add config file for keybindings and default browser/player
	TODO show refreshing status instead of freezing

	*: we might need to use a database because of these which sucks, maybe i will make
	it JSON, maybe plaintext, but I will try to avoid using a database as much as I can.
*/

package main

import (
	"os"
	"os/exec"
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
	"github.com/mmcdole/gofeed"
	"jaytaylor.com/html2text"
	"io"
	"io/ioutil"
	"strings"
	"net/http"
	"crypto/sha1"
	"encoding/hex"
)

/* UI stuff */
var ui tui.UI
var title *tui.Label
var maintable *tui.Table
var mainview *tui.Box
var feedtable *tui.Table
var feedview *tui.Box
var content *tui.Label
var contentarea *tui.ScrollArea
var entryview *tui.Box
var theme *tui.Theme
var view int // keep track of current view: mainview=0,feedview=1,entryview=2

var titlestr string
var confdir string // config directory
var datadir string // data directory (xml files)
var urls []string

/* sync all feeds (download files) */
func sync() {
	title.SetText(titlestr + "syncing...")
	for _, url := range urls {
		if len(url) < 1 {
			continue
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		urlsum := sha1.Sum([]byte(url))
		filename := datadir + "/" + hex.EncodeToString(urlsum[:])
		out, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	}
}

/* parse feed from data directory */
func get_feed(url string) *gofeed.Feed {
	urlsum := sha1.Sum([]byte(url))
	file, err := os.Open(datadir + "/" + hex.EncodeToString(urlsum[:]))

	if err != nil {
		panic(err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)

	if err != nil {
		panic(err)
	}

	return feed
}

/* open feed in video player */
func open_in_player(url string) {
	player := os.Getenv("PLAYER")

	if player == "" {
		player = "mpv" // default player
	}

	cmd := exec.Command("setsid", "nohup", player, url)
	cmd.Start()
}

/* open feed in default browser */
func open_in_browser(url string) {
	browser := os.Getenv("BROWSER")
	if browser != "" {
		cmd := exec.Command(browser, url)
		cmd.Start()
	}
}

/* initialize the main view (first thing you see, view 0) */
func init_mainview(feeds []*gofeed.Feed) {
	maintable = tui.NewTable(0, 0)
	maintable.SetFocused(true)

	for _, feed := range feeds {
		maintable.AppendRow(tui.NewLabel(feed.Title))
	}

	mainpadding := tui.NewLabel("")
	mainpadding.SetSizePolicy(tui.Preferred, tui.Expanding)
	mainview = tui.NewVBox(title, maintable, mainpadding)
}

/* initialize the feed view (list of entries in feed, view 1) */
func init_feedview() {
	feedtable = tui.NewTable(0, 0)
	feedpadding := tui.NewLabel("")
	feedpadding.SetSizePolicy(tui.Preferred, tui.Expanding)
	feedview = tui.NewVBox(title, feedtable, feedpadding)
}

/* update feed view for new feed */
func update_feedview(feed *gofeed.Feed) {
	items := feed.Items
	feedtable.RemoveRows()
	if len(items) > 80 {
		items = items[:80]
	}
	for _, item := range items {
		feedtable.AppendRow(tui.NewLabel(item.Title))
	}

	feedtable.SetFocused(true)
	feedtable.Select(0)
}

/* initialize the entry view (entry content, view 2) */
func init_entryview() {
	content = tui.NewLabel("")
	content.SetSizePolicy(tui.Preferred, tui.Expanding)

	contentarea = tui.NewScrollArea(content)
	entryview = tui.NewVBox(title, contentarea)
}

/* update entryview when a different one is opened */
func update_entryview(feed *gofeed.Feed, item *gofeed.Item) {
	metatext := "Feed: " + feed.Title + "\nTitle: " + item.Title + "\nDate: " + item.Published + "\nLink: " + item.Link
	feedtext, err := html2text.FromString(item.Description + "\n" + item.Content, html2text.Options{PrettyTables: true})
	if err != nil {
		panic(err)
	}

	content.SetText(metatext + "\n\n\n" + wordwrap.WrapString(feedtext, 80))
}

/* create a ui based on feeds */
func build_ui(feeds []*gofeed.Feed) tui.UI {
	title = tui.NewLabel(titlestr)

	init_mainview(feeds)
	init_feedview()
	update_feedview(feeds[0])
	init_entryview()
	update_entryview(feeds[0], feeds[0].Items[0])

	root := tui.NewVBox(mainview, feedview, entryview)
	ui, err := tui.New(root)

	ui.SetWidget(mainview)

	if err != nil {
		panic(err)
	}

	ui.SetKeybinding("h", func() {
		switch view {
		case 0:
			ui.Quit()
			break
		case 1:
			ui.SetWidget(mainview)
			view = 0
			break
		case 2:
			ui.SetWidget(feedview)
			view = 1
		}
	})

	ui.SetKeybinding("j", func() {
		if (view == 2) {
			contentarea.Scroll(0, 1)
		}
	})
	ui.SetKeybinding("k", func() {
		if (view == 2) {
			contentarea.Scroll(0, -1)
		}
	})

	ui.SetKeybinding("l", func() {
		switch view {
		case 0:
			update_feedview(feeds[maintable.Selected()])
			ui.SetWidget(feedview)
			view = 1
			break
		case 1:
			update_entryview(feeds[maintable.Selected()], feeds[maintable.Selected()].Items[feedtable.Selected()])
			ui.SetWidget(entryview)
			contentarea.ScrollToTop()
			view = 2
			break
		}
	})

	ui.SetKeybinding("r", func() {
		sync()
	})

	ui.SetKeybinding("o", func() {
		if view != 0 {
			open_in_browser(feeds[maintable.Selected()].Items[feedtable.Selected()].Link)
		}
	})

	ui.SetKeybinding("v", func() {
		if view != 0 {
			open_in_player(feeds[maintable.Selected()].Items[feedtable.Selected()].Link)
		}
	})

	ui.SetKeybinding("q", func() { ui.Quit() })

	return ui
}

func main() {
	/* set configuration stuff */
	confdir = os.Getenv("HOME") + "/.config/sreader"
	datadir = os.Getenv("HOME") + "/.local/share/sreader"
	urlsfile := confdir + "/urls"
	titlestr = "sreader: "

	/* this won't do anything if the files exist already */
	os.MkdirAll(confdir, os.ModePerm)
	os.MkdirAll(datadir, os.ModePerm)

	_, err := os.Stat(urlsfile)
    if os.IsNotExist(err) {
		file, err := os.Create(urlsfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	dat, err := ioutil.ReadFile(urlsfile)

	if err != nil {
		panic(err)
	}

	urls = strings.Split(string(dat), "\n")

	/* sync and quit if called with the arg "sync" */
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "sync":
			sync()
			return
		}
	}

	var feeds []*gofeed.Feed;

	for _, url := range urls {
		if len(url) > 0 {
			feeds = append(feeds, get_feed(url))
		}
	}

	ui = build_ui(feeds)

	err = ui.Run()
	if err != nil {
		panic(err)
	}
}

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
)

var ui tui.UI
var maintable *tui.Table
var mainview *tui.Box
var feedtable *tui.Table
var feedview *tui.Box
var content *tui.Label
var contentarea *tui.ScrollArea
var entryview *tui.Box
var view int

func sync(urls []string) {
	basedir := os.Getenv("HOME") + "/.local/share/sreader"
	for _, url := range urls {
		if len(url) < 1 {
			continue
		}
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		urlsum := sha1.Sum([]byte(url))
		filename := basedir + "/" + string(urlsum[:])
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

func get_feed(url string) *gofeed.Feed {
	basedir := os.Getenv("HOME") + "/.local/share/sreader"
	urlsum := sha1.Sum([]byte(url))
	file, err := os.Open(basedir + "/" + string(urlsum[:]))

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
		player = "mpv" // default
	}

	cmd := exec.Command(player, url)
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

func init_mainview(feeds []*gofeed.Feed) {
	maintable = tui.NewTable(0, 0)
	maintable.SetFocused(true)

	for _, feed := range feeds {
		maintable.AppendRow(tui.NewLabel(feed.Title))
	}

	mainpadding := tui.NewLabel("")
	mainpadding.SetSizePolicy(tui.Preferred, tui.Expanding)
	mainview = tui.NewVBox(maintable, mainpadding)
}

func init_feedview() {
	feedtable = tui.NewTable(0, 0)
	feedpadding := tui.NewLabel("")
	feedpadding.SetSizePolicy(tui.Preferred, tui.Expanding)
	feedview = tui.NewVBox(feedtable, feedpadding)
}

func update_feedview(feed *gofeed.Feed) {
	items := feed.Items
	feedtable.RemoveRows()
	for _, item := range items {
		feedtable.AppendRow(tui.NewLabel(item.Title))
	}

	feedtable.SetFocused(true)
	feedtable.Select(0)
}

func init_entryview() {
	content = tui.NewLabel("")
	content.SetSizePolicy(tui.Preferred, tui.Expanding)

	contentarea = tui.NewScrollArea(content)
	entryview = tui.NewVBox(contentarea)
}

func update_entryview(feed *gofeed.Feed, item *gofeed.Item) {

	metatext := "Feed: " + feed.Title + "\nTitle: " + item.Title + "\nDate: " + item.Published + "\nLink: " + item.Link
	feedtext, err := html2text.FromString(item.Description + "\n" + item.Content, html2text.Options{PrettyTables: true})
	if err != nil {
		panic(err)
	}

	content.SetText(metatext + "\n\n\n" + wordwrap.WrapString(feedtext, 80))
}

func build_ui(feeds []*gofeed.Feed) tui.UI {
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
			view = 2
			break
		}
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

	return ui
}

func main() {
	dat, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/sreader/urls")
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "sync":
			sync(strings.Split(string(dat), "\n"))
			return
		}
	}

	var feeds []*gofeed.Feed;

	for _, url := range strings.Split(string(dat), "\n") {
		if len(url) > 0 {
			feeds = append(feeds, get_feed(url))
		}
	}


	ui = build_ui(feeds)

	ui.SetKeybinding("q", func() { ui.Quit() })

	err = ui.Run()
	if err != nil {
		panic(err)
	}
}

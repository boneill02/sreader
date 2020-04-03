package main

import (
	"os"
	"os/exec"
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
	"github.com/mmcdole/gofeed"
	"jaytaylor.com/html2text"
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

func get_feed(url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)

	if err != nil {
		panic(err)
	}

	return feed
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

	return ui
}

func main() {
	feedurls := os.Args[1:]

	if len(feedurls) == 0 {
		panic("no feed url specified")
	}

	var feeds []*gofeed.Feed;

	for _, url := range feedurls {
		feeds = append(feeds, get_feed(url))
	}

	ui = build_ui(feeds)

	ui.SetKeybinding("q", func() { ui.Quit() })
	if err := ui.Run(); err != nil {
		panic(err)
	}
}

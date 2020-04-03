package main

import (
	"github.com/marcusolsson/tui-go"
	"github.com/SlyMarbo/rss"
	"jaytaylor.com/html2text"
)


var items []*rss.Item;

func get_feed(url string) {
	feed, err := rss.Fetch("http://benoneill.xyz/posts/index.xml")

	if err != nil {
		panic(err)
	}

	items = feed.Items

}

func main() {
	feed := "https://benoneill.xyz/posts/index.xml"

	get_feed(feed)

	feedbox := tui.NewTable(0, 0) // create table for feeds
	feedbox.SetFocused(true)

	for _, item := range items {
		feedbox.AppendRow(
			tui.NewLabel(item.Title),
			tui.NewLabel(item.Date.String()),
			tui.NewLabel(item.Link),
		)
	}

	title := tui.NewLabel("")
	date := tui.NewLabel("")
	url := tui.NewLabel("")

	content := tui.NewLabel("")
	content.SetSizePolicy(tui.Preferred, tui.Expanding)

	feedbox.OnSelectionChanged(func(t *tui.Table) {
		e := items[t.Selected()]
		title.SetText(e.Title)
		date.SetText(e.Date.String())
		url.SetText(e.Link)

		feedtext, err := html2text.FromString(e.Summary, html2text.Options{PrettyTables: true})
		if err != nil {
			panic(err)
		}
		content.SetText(feedtext)
	})

	feedbox.Select(0)

	root := tui.NewVBox(
		feedbox,
		tui.NewPadder(1, 0, tui.NewLabel("")),
		tui.NewSpacer(),
		content,
	)

	ui, err := tui.New(root)

	if err != nil {
		panic(err)
	}

	ui.SetKeybinding("q", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

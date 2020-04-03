package main

import (
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
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

	feedbox := tui.NewTable(0, 0) // create table for feed names
	feedbox.SetFocused(true)

	for _, item := range items {
		feedbox.AppendRow(
			tui.NewLabel(item.Title),
		)
	}

	content := tui.NewLabel("")
	content.SetSizePolicy(tui.Preferred, tui.Expanding)

	contentarea := tui.NewScrollArea(content)
	contentview := tui.NewVBox(contentarea)

	info := tui.NewLabel("")

	feedbox.OnSelectionChanged(func(t *tui.Table) {
		e := items[t.Selected()]
		info.SetText(e.Title)

		feedtext, err := html2text.FromString(e.Summary, html2text.Options{PrettyTables: true})
		if err != nil {
			panic(err)
		}
		content.SetText(wordwrap.WrapString(feedtext, 80))
	})

	feedbox.Select(0)

	feedpad := tui.NewLabel("")
	feedpad.SetSizePolicy(tui.Preferred, tui.Expanding)

	feedview := tui.NewVBox(feedbox, feedpad)

	root := tui.NewVBox(feedview)

	ui, err := tui.New(root)


	view := 0

	if err != nil {
		panic(err)
	}

	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("h", func() {
		if view == 0 {
			ui.Quit()
		} else if view == 1 {
			ui.SetWidget(feedview)
			contentarea.ScrollToTop() // Autoscroll back
			view = 0
		}
	})
	ui.SetKeybinding("j", func() { contentarea.Scroll(0, 1) })
	ui.SetKeybinding("k", func() { contentarea.Scroll(0, -1) })
	ui.SetKeybinding("l", func() {
		ui.SetWidget(contentview)
		view = 1
	})

	ui.SetKeybinding("Esc", func() {
		if view == 0 {
			ui.Quit()
		} else if view == 1 {
			ui.SetWidget(feedview)
			contentarea.ScrollToTop() // Autoscroll back
			view = 0
		}
	})
	ui.SetKeybinding("Enter", func() {
		ui.SetWidget(contentview)
		view = 1
	})
	ui.SetKeybinding("Up", func() { contentarea.Scroll(0, -1) })
	ui.SetKeybinding("Down", func() { contentarea.Scroll(0, 1) })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

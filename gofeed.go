package main

import (
	"os"
	"os/exec"
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
	"github.com/SlyMarbo/rss"
	"jaytaylor.com/html2text"
)

func get_feed(url string) *rss.Feed {
	feed, err := rss.Fetch(url)

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

func main() {
	url := "https://www.fsf.org/static/fsforg/rss/blogs.xml"

	feed := get_feed(url)

	feedbox := tui.NewTable(0, 0) // create table for feed names
	feedbox.SetFocused(true)

	for _, item := range feed.Items {
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
		e := feed.Items[t.Selected()]
		info.SetText(e.Title)

		metatext := "Feed: " + feed.Title + "\nTitle: " + e.Title + "\nDate: " + e.Date.String() + "\nLink: " + e.Link
		feedtext, err := html2text.FromString(e.Summary + e.Content, html2text.Options{PrettyTables: true})
		if err != nil {
			panic(err)
		}
		content.SetText(metatext + "\n\n\n" + wordwrap.WrapString(feedtext, 80))
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
			view = 0
		}
	})
	ui.SetKeybinding("j", func() { contentarea.Scroll(0, 1) })
	ui.SetKeybinding("k", func() { contentarea.Scroll(0, -1) })
	ui.SetKeybinding("l", func() {
		ui.SetWidget(contentview)
		contentarea.ScrollToTop() // Autoscroll back
		view = 1
	})

	ui.SetKeybinding("o", func() {
		open_in_browser(feed.Items[feedbox.Selected()].Link)
	})

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

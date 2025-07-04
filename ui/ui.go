package ui

import (
	"github.com/boneill02/sreader/feed"
	"github.com/marcusolsson/tui-go"
	"github.com/mmcdole/gofeed"
)

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

const titlestr string = "sreader: "

/**
 * Create a ui based on the given feeds.
 */
func Init(feeds []*gofeed.Feed) tui.UI {
	title = tui.NewLabel(titlestr)

	InitIndexView(feeds)
	InitFeedView()
	UpdateFeedView(feeds[0])
	InitEntryView()
	UpdateEntryView(feeds[0], feeds[0].Items[0])

	root := tui.NewVBox(mainview, feedview, entryview)
	ui, err := tui.New(root)

	ui.SetWidget(mainview)

	if err != nil {
		panic(err)
	}

	SetKeybindings(ui, feeds)

	return ui
}

/**
 * Initialize keybindings for the UI
 */
func SetKeybindings(ui tui.UI, feeds []*gofeed.Feed) {
	// Go to previous view or quit if in indexview
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

	// Scroll down in content area
	ui.SetKeybinding("j", func() {
		if view == 2 {
			contentarea.Scroll(0, 1)
		}
	})

	// Scroll up in content area
	ui.SetKeybinding("k", func() {
		if view == 2 {
			contentarea.Scroll(0, -1)
		}
	})

	// Enter feedview for selected feed or entryview for selected entry
	ui.SetKeybinding("l", func() {
		switch view {
		case 0:
			UpdateFeedView(feeds[maintable.Selected()])
			ui.SetWidget(feedview)
			view = 1
			break
		case 1:
			UpdateEntryView(feeds[maintable.Selected()], feeds[maintable.Selected()].Items[feedtable.Selected()])
			ui.SetWidget(entryview)
			contentarea.ScrollToTop()
			view = 2
			break
		}
	})

	// Sync feeds
	ui.SetKeybinding("r", func() {
		feed.Sync()
	})

	// Open in browser
	ui.SetKeybinding("o", func() {
		if view != 0 {
			feed.OpenInBrowser(feeds[maintable.Selected()].Items[feedtable.Selected()].Link)
		}
	})

	// Open in player
	ui.SetKeybinding("v", func() {
		if view != 0 {
			feed.OpenInPlayer(feeds[maintable.Selected()].Items[feedtable.Selected()].Link)
		}
	})

	// Quit
	ui.SetKeybinding("q", func() { ui.Quit() })
}

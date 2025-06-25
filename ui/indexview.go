package ui

import (
	"github.com/mmcdole/gofeed"
	"github.com/marcusolsson/tui-go"
)

/**
 * initialize the index view (first thing you see, view 0)
 */
func InitIndexView(feeds []*gofeed.Feed) {
	maintable = tui.NewTable(0, 0)
	maintable.SetFocused(true)

	for _, feed := range feeds {
		maintable.AppendRow(tui.NewLabel(feed.Title))
	}

	mainpadding := tui.NewLabel("")
	mainpadding.SetSizePolicy(tui.Preferred, tui.Expanding)
	mainview = tui.NewVBox(title, maintable, mainpadding)
}

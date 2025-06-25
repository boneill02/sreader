package ui

import (
	"github.com/mmcdole/gofeed"
	"github.com/marcusolsson/tui-go"
)

/**
 * Initialize the feed view (list of entries in feed, view 1)
 */
func InitFeedView() {
	feedtable = tui.NewTable(0, 0)
	feedpadding := tui.NewLabel("")
	feedpadding.SetSizePolicy(tui.Preferred, tui.Expanding)
	feedview = tui.NewVBox(title, feedtable, feedpadding)
}

/**
 * Update feed view for new feed
 */
func UpdateFeedView(feed *gofeed.Feed) {
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

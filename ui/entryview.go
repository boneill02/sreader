package ui

import (
	"github.com/k3a/html2text"
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
	"github.com/mmcdole/gofeed"
)

/**
 * initialize the entry view (entry content, view 2)
 */
func InitEntryView() {
	content = tui.NewLabel("")
	contentarea = tui.NewScrollArea(content)
	contentarea.SetSizePolicy(tui.Preferred, tui.Expanding)
	entryview = tui.NewVBox(title, contentarea)
}

/**
 * update entryview when a different entry is opened
 */
func UpdateEntryView(feed *gofeed.Feed, item *gofeed.Item) {
	metatext := "Feed: " + feed.Title + "\nTitle: " + item.Title + "\nDate: " + item.Published + "\nLink: " + item.Link
	feedtext := html2text.HTML2Text(item.Description + "\n" + item.Content)
	content.SetText(metatext + "\n\n\n" + wordwrap.WrapString(feedtext, 80))
}

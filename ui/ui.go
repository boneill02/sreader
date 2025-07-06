package ui

import (
	"fmt"

	"github.com/boneill02/sreader/config"
	"github.com/boneill02/sreader/feed"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/k3a/html2text"
	"github.com/mmcdole/gofeed"
)

const titlestr = "sreader: "

type viewState int

// View states
const (
	mainView viewState = iota
	feedView
	entryView
)

var (
	listDelegate list.DefaultDelegate
	appStyle     lipgloss.Style
)

type feedItem struct {
	title string
	desc  string
	link  string
}

func (f feedItem) Title() string       { return f.title }
func (f feedItem) Description() string { return f.desc }
func (f feedItem) FilterValue() string { return f.title }

type model struct {
	feeds     []*gofeed.Feed
	config    *config.Config
	view      viewState
	feedList  list.Model
	entryList list.Model
	entryView viewport.Model
	currFeed  int
	currEntry int
	width     int
	height    int
}

func NewModel(feeds []*gofeed.Feed, conf *config.Config, width, height int) model {
	feedItems := make([]list.Item, len(feeds))
	for i, f := range feeds {
		feedItems[i] = feedItem{title: f.Title, desc: f.Description, link: ""}
	}
	feedList := list.New(feedItems, list.NewDefaultDelegate(), width, height)
	feedList.Title = "Feeds"
	feedList.SetShowHelp(false)

	entryItems := []list.Item{}
	if len(feeds) > 0 {
		for _, item := range feeds[0].Items {
			entryItems = append(entryItems, feedItem{title: item.Title, desc: item.Description, link: item.Link})
		}
	}
	entryList := list.New(entryItems, list.NewDefaultDelegate(), width, height)
	entryList.Title = "Entries"
	entryList.SetShowHelp(false)

	vp := viewport.New(width, height)
	if len(feeds) > 0 && len(feeds[0].Items) > 0 {
		vp.SetContent(feeds[0].Items[0].Content)
	}

	return model{
		feeds:     feeds,
		config:    conf,
		view:      mainView,
		feedList:  feedList,
		entryList: entryList,
		entryView: vp,
		currFeed:  0,
		currEntry: 0,
		width:     width,
		height:    height,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.feedList.SetSize(msg.Width, msg.Height)
		m.entryList.SetSize(msg.Width, msg.Height)
		m.entryView.Width = msg.Width
		m.entryView.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "h":
			switch m.view {
			case mainView:
				return m, tea.Quit
			case feedView:
				m.view = mainView
			case entryView:
				m.view = feedView
			}
		case "l":
			switch m.view {
			case mainView:
				m.currFeed = m.feedList.Index()
				m.updateEntryList()
				m.view = feedView
			case feedView:
				m.currEntry = m.entryList.Index()
				m.updateViewport()
				m.view = entryView
			}
		case "j":
			switch m.view {
			case mainView:
				m.feedList, _ = m.feedList.Update(msg)
			case feedView:
				m.entryList, _ = m.entryList.Update(msg)
			case entryView:
				m.entryView.ScrollDown(1)
			}
		case "k":
			switch m.view {
			case mainView:
				m.feedList, _ = m.feedList.Update(msg)
			case feedView:
				m.entryList, _ = m.entryList.Update(msg)
			case entryView:
				m.entryView.ScrollUp(1)
			}
		case "r":
			feed.Sync()
		case "o":
			if m.view == feedView || m.view == entryView {
				link := m.feeds[m.currFeed].Items[m.currEntry].Link
				feed.OpenInBrowser(link, m.config.Browser)
			}
		case "v":
			if m.view == feedView || m.view == entryView {
				link := m.feeds[m.currFeed].Items[m.currEntry].Link
				feed.OpenInPlayer(link, m.config.Player)
			}
		default:
			switch m.view {
			case mainView:
				m.feedList, _ = m.feedList.Update(msg)
			case feedView:
				m.entryList, _ = m.entryList.Update(msg)
			case entryView:
				m.entryView, _ = m.entryView.Update(msg)
			}
		}
	}
	return m, nil
}

func (m *model) updateEntryList() {
	entryItems := []list.Item{}
	if m.currFeed < len(m.feeds) {
		for _, item := range m.feeds[m.currFeed].Items {
			entryItems = append(entryItems, feedItem{
				title: item.Title,
				link:  item.Link,
			})
		}
	}
	m.entryList.SetItems(entryItems)
	m.entryList.SetDelegate(listDelegate)
	m.entryList.Select(0)
	m.currEntry = 0
}

/**
 * Converts HTML to plain text and wraps lines at the specified width.
 */
func htmlTruncate(html string, width int) string {
	s := html2text.HTML2Text(html)
	var result []rune
	lineLen := 0
	isLink := false
	i := 0
	for i < len(s) {
		ch := s[i]
		// Check for start of URL
		if ch == 'h' && i+3 < len(s) && s[i:i+4] == "http" {
			isLink = true
		}

		// Check if end of URL
		if isLink && (ch == ' ' || ch == '\n' || ch == '\t') {
			isLink = false
		}
		result = append(result, rune(ch))
		if ch == '\n' {
			lineLen = 0
		} else {
			lineLen++
		}
		if lineLen >= width && !isLink {
			// Wrap word if it exceeds width
			// Find the last space before the width limit
			for j := len(result) - 1; j >= 0; j-- {
				if result[j] == ' ' || result[j] == '\n' {
					// Replace the space with a newline
					result[j] = '\n'
					lineLen = len(result) - j - 1 // Reset line length after newline
					break
				}
			}
		}
		i++
	}
	return string(result)
}

func (m *model) updateViewport() {
	if m.currFeed < len(m.feeds) && m.currEntry < len(m.feeds[m.currFeed].Items) {
		// Set the content to the selected entry's content
		content := "\nDate: " + m.feeds[m.currFeed].Items[m.currEntry].PublishedParsed.String()
		content += "\nLink: " + m.feeds[m.currFeed].Items[m.currEntry].Link
		content += "\nDescription:\n" + htmlTruncate(m.feeds[m.currFeed].Items[m.currEntry].Description, m.width-2)
		content += "\n\nContent:\n" + htmlTruncate(m.feeds[m.currFeed].Items[m.currEntry].Content, m.width-2)
		m.entryView.SetContent(content)
		m.entryView.GotoTop()
	}
}

func (m model) View() string {
	s := fmt.Sprintf("%s\n\n", titlestr)
	switch m.view {
	case mainView:
		s += m.feedList.View()
	case feedView:
		s += m.entryList.View()
	case entryView:
		s += m.entryView.View()
	}
	s += "\n[h] back [l] enter [j/k] move [q] quit [r] sync [o] open [v] play"

	// Render the entire UI with the app style
	return appStyle.Render(lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, s))
}

func Init(feeds []*gofeed.Feed, conf *config.Config) *tea.Program {
	width, height := 500, 24 // width set to 500, hopefully enough for most screens

	// Styles
	bg := lipgloss.Color(conf.BG)
	fg := lipgloss.Color(conf.FG)
	selectedTitleFG := lipgloss.Color(conf.SelectedTitleFG)
	selectedTitleBG := lipgloss.Color(conf.SelectedTitleBG)
	selectedDescFG := lipgloss.Color(conf.SelectedDescFG)
	selectedDescBG := lipgloss.Color(conf.SelectedDescBG)
	titleFG := lipgloss.Color(conf.TitleFG)
	titleBG := lipgloss.Color(conf.TitleBG)
	descFG := lipgloss.Color(conf.DescFG)
	descBG := lipgloss.Color(conf.DescBG)

	// Load list delegate with styles
	listDelegate = list.NewDefaultDelegate()
	listDelegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(titleFG).
		Background(titleBG).
		Width(width)
	listDelegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(selectedTitleFG).
		Background(selectedTitleBG).
		Width(width)
	listDelegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(descFG).
		Background(descBG).
		Width(width)
	listDelegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(selectedDescFG).
		Background(selectedDescBG).
		Width(width)

	m := NewModel(feeds, conf, width, height)
	m.feedList.SetDelegate(listDelegate)
	appStyle = lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Width(width)
	return tea.NewProgram(m, tea.WithAltScreen())
}

package ui

import (
	"fmt"

	html2markdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/bmoneill/sreader/config"
	"github.com/bmoneill/sreader/feed"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	feeds     []*feed.Feed
	config    *config.SreaderConfig
	view      viewState
	feedList  list.Model
	entryList list.Model
	entryView viewport.Model
	currFeed  int
	currEntry int
	width     int
	height    int
}

/**
 * Initializes the UI with the given feeds and configuration.
 */
func Init(feeds []*feed.Feed) *tea.Program {
	width, height := 500, 24 // width set to 500, hopefully enough for most screens

	// Styles
	bg := lipgloss.Color(config.Config.BG)
	fg := lipgloss.Color(config.Config.FG)
	selectedTitleFG := lipgloss.Color(config.Config.SelectedTitleFG)
	selectedTitleBG := lipgloss.Color(config.Config.SelectedTitleBG)
	selectedDescFG := lipgloss.Color(config.Config.SelectedDescFG)
	selectedDescBG := lipgloss.Color(config.Config.SelectedDescBG)
	titleFG := lipgloss.Color(config.Config.TitleFG)
	titleBG := lipgloss.Color(config.Config.TitleBG)
	descFG := lipgloss.Color(config.Config.DescFG)
	descBG := lipgloss.Color(config.Config.DescBG)

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

	m := newModel(feeds, width, height)
	m.feedList.SetDelegate(listDelegate)
	appStyle = lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Width(width)
	return tea.NewProgram(m, tea.WithAltScreen())
}

func (m model) Init() tea.Cmd {
	return nil
}

/**
 * Handles user input and updates the model accordingly.
 */
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
				m.currFeed = 0
			case entryView:
				m.view = feedView
				m.currEntry = 0
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
			m.feeds = feed.GetFeeds()
		case "o":
			if m.view == feedView || m.view == entryView {
				link := m.feeds[m.currFeed].Entries[m.currEntry].URL
				feed.OpenInBrowser(link, m.config.Browser)
			}
		case "v":
			if m.view == feedView || m.view == entryView {
				link := m.feeds[m.currFeed].Entries[m.currEntry].URL
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

/**
 * In entryList, updates the list of entries based on the currently selected feed.
 */
func (m *model) updateEntryList() {
	entryItems := []list.Item{}
	if m.currFeed < len(m.feeds) {
		for _, item := range m.feeds[m.currFeed].Entries {
			entryItems = append(entryItems, feedItem{
				title: item.Title,
				link:  item.URL,
			})
		}
	}
	m.entryList.SetItems(entryItems)
	m.entryList.SetDelegate(listDelegate)
	m.entryList.Select(0)
	m.currEntry = 0
}

/**
 * In entryView, updates the viewport with the content of the currently selected entry.
 */
func (m *model) updateViewport() {
	if m.currFeed < len(m.feeds) && m.currEntry < len(m.feeds[m.currFeed].Entries) {
		// Set the content to the selected entry's content
		content := "\nDate: " + m.feeds[m.currFeed].Entries[m.currEntry].DatePublished
		content += "\nLink: " + m.feeds[m.currFeed].Entries[m.currEntry].URL
		content += "\n\n" + htmlTruncate(m.feeds[m.currFeed].Entries[m.currEntry].Description, m.width-2)
		content += "\n\n" + htmlTruncate(m.feeds[m.currFeed].Entries[m.currEntry].Content, m.width-2)
		m.entryView.SetContent(content)
		m.entryView.GotoTop()
	}
}

/**
 * Renders the current view of the model.
 */
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

/**
 * Converts HTML to plain text and wraps lines at the specified width.
 */
func htmlTruncate(content string, width int) string {
	s, _ := html2markdown.ConvertString(content)
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

func newModel(feeds []*feed.Feed, width, height int) model {
	feedItems := make([]list.Item, len(feeds))
	for i, f := range feeds {
		feedItems[i] = feedItem{title: f.Title, desc: f.Description, link: f.URL}
	}
	feedList := list.New(feedItems, list.NewDefaultDelegate(), width, height)
	feedList.Title = "Feeds"
	feedList.SetShowHelp(false)

	entryItems := []list.Item{}
	entryList := list.New(entryItems, list.NewDefaultDelegate(), width, height)
	entryList.Title = "Entries"
	entryList.SetShowHelp(false)

	vp := viewport.New(width, height)
	if len(feeds) > 0 && len(feeds[0].Entries) > 0 {
		vp.SetContent(feeds[0].Entries[0].Content)
	}

	return model{
		feeds:     feeds,
		config:    config.Config,
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

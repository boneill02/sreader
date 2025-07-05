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
	viewport  viewport.Model
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

	entryItems := []list.Item{}
	if len(feeds) > 0 {
		for _, item := range feeds[0].Items {
			entryItems = append(entryItems, feedItem{title: item.Title, desc: item.Description, link: item.Link})
		}
	}
	entryList := list.New(entryItems, list.NewDefaultDelegate(), width, height)
	entryList.Title = "Entries"

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
		viewport:  vp,
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
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
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
				m.viewport.ScrollDown(1)
			}
		case "k":
			switch m.view {
			case mainView:
				m.feedList, _ = m.feedList.Update(msg)
			case feedView:
				m.entryList, _ = m.entryList.Update(msg)
			case entryView:
				m.viewport.ScrollUp(1)
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
				m.viewport, _ = m.viewport.Update(msg)
			}
		}
	}
	return m, nil
}

func (m *model) updateEntryList() {
	entryItems := []list.Item{}
	if m.currFeed < len(m.feeds) {
		for _, item := range m.feeds[m.currFeed].Items {
			entryItems = append(entryItems, feedItem{title: item.Title, desc: item.Description, link: item.Link})
		}
	}
	m.entryList.SetItems(entryItems)
	m.entryList.SetDelegate(listDelegate)
	m.entryList.Select(0)
	m.currEntry = 0
}

func (m *model) updateViewport() {
	if m.currFeed < len(m.feeds) && m.currEntry < len(m.feeds[m.currFeed].Items) {
		// Set the content to the selected entry's content
		content := "\nDate: " + m.feeds[m.currFeed].Items[m.currEntry].PublishedParsed.String()
		content += "\nLink: " + m.feeds[m.currFeed].Items[m.currEntry].Link
		content += "\nDescription:\n" + html2text.HTML2Text(m.feeds[m.currFeed].Items[m.currEntry].Description)
		content += "\n\nContent:\n" + html2text.HTML2Text(m.feeds[m.currFeed].Items[m.currEntry].Content)
		m.viewport.SetContent(content)
		m.viewport.GotoTop()
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
		s += m.viewport.View()
	}
	s += "\n[h] back [l] enter [j/k] move [q] quit [r] sync [o] open [v] play"
	return appStyle.Render(s)
}

func Init(feeds []*gofeed.Feed, conf *config.Config) *tea.Program {
	width, height := 80, 24 // default
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

	// List styles
	listDelegate = list.NewDefaultDelegate()
	listDelegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(titleFG).
		Background(titleBG)
	listDelegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(selectedTitleFG).
		Background(selectedTitleBG)
	listDelegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(descFG).
		Background(descBG)
	listDelegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(selectedDescFG).
		Background(selectedDescBG)

	m := NewModel(feeds, conf, width, height)
	m.feedList.SetDelegate(listDelegate)

	appStyle = lipgloss.NewStyle().
		Foreground(fg).
		Background(bg)
	return tea.NewProgram(m, tea.WithAltScreen())
}

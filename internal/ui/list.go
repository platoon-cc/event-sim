package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("270"))
	currentItemStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	defaultHeight     = 14
	defaultWidth      = 20
)

type ItemSelectFn func(ListItem)

type ListItem struct {
	Key    string
	Value  string
	Active bool
}

func (i ListItem) FilterValue() string { return "" }

type model struct {
	list     list.Model
	onSelect ItemSelectFn
	choice   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(ListItem)
			if ok {
				m.choice = i.Value
				if m.onSelect != nil {
					m.onSelect(i)
				}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Value)
	if i.Active {
		str += currentItemStyle.Render("(current)")
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type List struct {
	title string
	items []list.Item
}

func NewList(title string) List {
	return List{
		title: title,
	}
}

func (l *List) AddItem(key string, value string, isActive bool) {
	l.items = append(l.items, ListItem{
		Key:    key,
		Value:  value,
		Active: isActive,
	})
}

func (l *List) Run(onSelect ItemSelectFn) error {
	m := model{}
	m.list = list.New(l.items, itemDelegate{}, defaultWidth, defaultHeight)
	m.list.Title = l.title
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	m.list.Styles.Title = titleStyle
	m.list.Styles.PaginationStyle = paginationStyle
	m.list.Styles.HelpStyle = helpStyle

	m.onSelect = onSelect

	_, err := tea.NewProgram(m).Run()
	return err
}

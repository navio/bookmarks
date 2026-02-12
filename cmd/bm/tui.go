package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/navio/bookmarks/internal/bookmarks"
)

type bookmarkItem struct {
	b bookmarks.Bookmark
}

func (i bookmarkItem) Title() string       { return i.b.Name }
func (i bookmarkItem) Description() string { return i.b.Path }
func (i bookmarkItem) FilterValue() string {
	return i.b.Name + " " + i.b.Path + " " + strings.Join(i.b.Tags, ",")
}

// ----------------
// FIND (list)
// ----------------

type findModel struct {
	list     list.Model
	selected string
	tags     []string
}

func newFindModel(items []list.Item, title string, tags []string) findModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Bold(true)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Bold(true)
	delegate.Styles.DimmedTitle = delegate.Styles.DimmedTitle.Bold(true)

	lm := list.New(items, delegate, 0, 0)
	lm.Title = title
	lm.SetShowStatusBar(true)
	lm.SetFilteringEnabled(true)
	lm.KeyMap.Quit.SetEnabled(true)
	return findModel{list: lm, tags: tags}
}

func (m findModel) Init() tea.Cmd { return nil }

func (m findModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if it, ok := m.list.SelectedItem().(bookmarkItem); ok {
				m.selected = it.b.Path
				return m, tea.Quit
			}
		case "c":
			if it, ok := m.list.SelectedItem().(bookmarkItem); ok {
				if err := clipboard.WriteAll(it.b.Path); err != nil {
					m.list.NewStatusMessage(statusStyle.Render("copy failed: " + err.Error()))
					return m, nil
				}
				m.list.NewStatusMessage(statusStyle.Render("copied: " + it.b.Path))
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

var tagBannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("212")).
	Bold(true)

func (m findModel) View() string {
	var b strings.Builder
	if len(m.tags) > 0 {
		b.WriteString(tagBannerStyle.Render("filters: "+strings.Join(m.tags, ", ")) + "\n")
	}
	b.WriteString(m.list.View())
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("enter: print  •  c: copy  •  /: filter  •  q: quit")
	b.WriteString("\n" + help)
	return b.String()
}

// ----------------
// TABLE
// ----------------

type tableModel struct {
	table    table.Model
	selected string
}

var statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

func newTableModel(rows []table.Row, title string) tableModel {
	columns := []table.Column{
		{Title: "Name", Width: 18},
		{Title: "Path", Width: 48},
		{Title: "Tags", Width: 20},
		{Title: "Created", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true)
	styles.Selected = styles.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(styles)

	_ = title // shown in View
	return tableModel{table: t}
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			row := m.table.SelectedRow()
			if len(row) >= 2 {
				m.selected = row[1]
			}
			return m, tea.Quit
		case "c":
			row := m.table.SelectedRow()
			if len(row) >= 2 {
				_ = clipboard.WriteAll(row[1])
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		// keep the table responsive to terminal size
		m.table.SetHeight(max(5, msg.Height-3))
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Render("bm table")
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("enter: print  •  c: copy  •  q: quit")
	return header + "\n" + m.table.View() + "\n" + help
}

// Helpers

func buildTableRows(entries []bookmarks.Bookmark) []table.Row {
	rows := make([]table.Row, 0, len(entries))
	for _, e := range entries {
		created := ""
		if !e.CreatedAt.IsZero() {
			created = e.CreatedAt.In(time.Local).Format("2006-01-02")
		}
		rows = append(rows, table.Row{
			e.Name,
			e.Path,
			strings.Join(e.Tags, ","),
			created,
		})
	}
	return rows
}

func runFindTUI(entries []bookmarks.Bookmark, title string, tags []string) (string, error) {
	items := make([]list.Item, 0, len(entries))
	for _, e := range entries {
		items = append(items, bookmarkItem{b: e})
	}
	m := newFindModel(items, title, tags)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	fm, ok := final.(findModel)
	if !ok {
		return "", fmt.Errorf("unexpected model")
	}
	return fm.selected, nil
}

func runTableTUI(entries []bookmarks.Bookmark, title string) (string, error) {
	rows := buildTableRows(entries)
	m := newTableModel(rows, title)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	tm, ok := final.(tableModel)
	if !ok {
		return "", fmt.Errorf("unexpected model")
	}
	return tm.selected, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

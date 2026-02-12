package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/navio/bookmarks/internal/bookmarks"
)

// ----------------
// FIND (filepicker)
// ----------------

type findModel struct {
	filepicker filepicker.Model
	selected   string
	quitting   bool
	title      string
	status     string
}

func newFindModel(startDir string, title string) findModel {
	fp := filepicker.New()
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.ShowHidden = false
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.CurrentDirectory = startDir
	return findModel{filepicker: fp, title: title}
}

func (m findModel) Init() tea.Cmd { return m.filepicker.Init() }

func (m findModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "c":
			if m.filepicker.Path != "" {
				if err := clipboard.WriteAll(m.filepicker.Path); err != nil {
					m.status = statusStyle.Render("copy failed: " + err.Error())
				} else {
					m.status = statusStyle.Render("copied: " + m.filepicker.Path)
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selected = path
		return m, tea.Quit
	}

	return m, cmd
}

func (m findModel) View() string {
	if m.quitting {
		return ""
	}
	header := lipgloss.NewStyle().Bold(true).Render(m.title)
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("enter: select  •  c: copy  •  h/←: back  •  l/→: open  •  q: quit")
	view := header + "\n" + m.filepicker.View()
	if m.status != "" {
		view += "\n" + m.status
	}
	return view + "\n" + help
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

func runFindTUI(startDir string, title string) (string, error) {
	m := newFindModel(startDir, title)
	p := tea.NewProgram(m, tea.WithAltScreen())
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
	p := tea.NewProgram(m, tea.WithAltScreen())
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

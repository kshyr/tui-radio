package client

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kshyr/tui-radio/internal/audio"
	"github.com/kshyr/tui-radio/internal/radio"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	AlignHorizontal(lipgloss.Center)

type model struct {
	table table.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			go audio.Play(m.table.SelectedRow()[2])
			return m, tea.Batch(
				tea.Printf("Listening to %s!\n", m.table.SelectedRow()[0]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		baseStyle.Render(m.table.View()),
		baseStyle.Render(m.table.View()),
	)
}

func (m model) Init() tea.Cmd {
	return nil
}

func DefaultClient() *tea.Program {

	columns := []table.Column{
		{Title: "Station Name", Width: 20},
		{Title: "Language", Width: 15},
		{Title: "URL", Width: 40},
	}

	rows := []table.Row{}

	stations := radio.GetStations()
	for _, station := range stations {
		newRow := make([]string, 3)
		newRow[0] = station.Name
		newRow[1] = station.Language
		newRow[2] = station.URL
		rows = append(rows, newRow)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	return tea.NewProgram(m)
}

func NewClient(table table.Model) *tea.Program {
	m := model{table}
	return tea.NewProgram(m)
}

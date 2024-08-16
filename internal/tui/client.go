package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kshyr/tui-radio/internal/audio"
	"gitlab.com/AgentNemo/goradios"
)

type Client struct{}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

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
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func DefaultClient() *tea.Program {

	columns := []table.Column{
		{Title: "Station Name", Width: 20},
		{Title: "Language", Width: 15},
		{Title: "URL", Width: 40},
	}

	rows := []table.Row{}

	stations := goradios.FetchStations(goradios.StationsByCountry, "Germany")
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
		table.WithHeight(7),
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

func Listen(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		foo := make(map[string]interface{})
		err := json.Unmarshal([]byte(scanner.Text()), &foo)
		if err != nil {
			return "", err
		}
		fmt.Println(foo)
	}
	return scanner.Text(), nil
}

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count int
}

func intialState() model {
	return model{count: 0}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit // Quit program
		case "+":
			m.count++ // Increment on spacebar
		case "-":
			m.count--
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"Press SPACE to count: %d\n\n(ctrl+c or q to quit)\n",
		m.count,
	)
}

func main() {
	p := tea.NewProgram(intialState())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
}

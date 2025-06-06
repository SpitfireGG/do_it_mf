package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type (
	errMsg error
)

type model struct {
	count     int
	textInput textinput.Model
	err       errMsg
}

func initialModel() model {

	t1 := textinput.New()
	t1.Placeholder = "placeholder"
	t1.Focus()
	t1.CharLimit = 256
	t1.Width = 40

	return model{
		count:     0,
		textInput: t1,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl-c":
			return m, tea.Quit
		case "+":
			m.count++
		case "-":
			m.count--
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Press SPACE to count: %d\n%s\n(ctrl+c or q to quit)\n",
		m.count, m.textInput.View(),
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: show the string "no tasks submitted and added only for few seconds"

type (
	errMsg error
)

type model struct {
	count     int
	textInput textinput.Model
	err       errMsg
	tasks     []string
}

func initialModel() model {

	t1 := textinput.New()
	t1.Placeholder = "enter something..."
	t1.Focus()
	t1.CharLimit = 256
	t1.Width = 40

	return model{
		count:     0,
		textInput: t1,
		err:       nil,
		tasks:     []string{},
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

		case "enter":
			submittedText := strings.TrimSpace(m.textInput.Value())
			if submittedText != "" {
				m.tasks = append(m.tasks, submittedText)
				m.count++
				m.textInput.SetValue("")
			}

		case "delete":
			if len(m.tasks) > 0 {
				m.tasks = m.tasks[:len(m.tasks)-1] // remove the last input from the slice
				m.count--
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd

		}
	case errMsg:
		m.err = msg
		return m, nil

	}
	return m, cmd
}

func (m model) View() string {

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("current count : %d\n\n", m.count))
	s.WriteString("submitted task!\n")

	if len(m.tasks) == 0 {
		s.WriteString("no tasks submitted yet!!!")
	} else {
		for i, task := range m.tasks {
			s.WriteString(fmt.Sprintf("%d: %s\n", i+1, task))
		}
	}
	s.WriteString("\n")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	s.WriteString("press Enter to submit, backspace/delete to remove last\n")
	s.WriteString("press q or ctrl+c to quit\n")

	// handle the error finally
	if m.err != nil {
		s.WriteString(fmt.Sprintf("%s\n", m.err))
	}
	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
}

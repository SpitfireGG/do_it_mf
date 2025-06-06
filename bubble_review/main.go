package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: show the string "no tasks submitted and added only for few seconds"

type statusMsgClear struct{}

type (
	errMsg error
)

var (
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(0, 1).
				Align(lipgloss.Center)

	border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1F456E")).
		Padding(1, 2)
)

type model struct {
	count           int
	textInput       textinput.Model
	err             errMsg
	tasks           []string
	statusMsg       string
	statusMsgExpire time.Time
	height          int
	width           int
}

func initialModel() model {

	t1 := textinput.New()
	t1.Placeholder = "enter something..."
	t1.Focus()
	t1.CharLimit = 256
	t1.Width = 40
	initialStatus := "no tasks submitted yet! start typing." // default message
	initialExpire := time.Now().Add(3 * time.Second)         // default expiry time

	return model{
		count:           0,
		textInput:       t1,
		err:             nil,
		tasks:           []string{},
		statusMsg:       initialStatus,
		statusMsgExpire: initialExpire,
		height:          0,
		width:           0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.createSatusMsgCmd())
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

				m.statusMsg = fmt.Sprintf("added: %s\n", submittedText)
				m.statusMsgExpire = time.Now().Add(2 * time.Second) // display the message for 2 seconds
				return m, m.createSatusMsgCmd()
			}

		case "delete":
			if len(m.tasks) > 0 {
				// keep track of what was removed
				removedTask := m.tasks[:len(m.tasks)-1]
				m.tasks = m.tasks[:len(m.tasks)-1] // remove the last input from the slice
				m.count--

				if len(m.tasks) == 0 {
					m.statusMsg = "no tasks added yet!"
					m.statusMsgExpire = time.Now().Add(3 * time.Second)
					return m, m.createSatusMsgCmd()
				} else {
					m.statusMsg = fmt.Sprintf("Removed: \"%s\"", removedTask)
					m.statusMsgExpire = time.Now().Add(2 * time.Second)
				}
				return m, m.createSatusMsgCmd()
			}
			/* default:
			m.textInput, cmd = m.textInput.Updte(msg)
			return m, cmd */
		}
	case statusMsgClear:
		if time.Now().After(m.statusMsgExpire) {
			m.statusMsg = ""
			m.statusMsgExpire = time.Time{}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = msg.Width / 2
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd

	case errMsg:
		m.err = msg
		return m, nil

	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) createSatusMsgCmd() tea.Cmd {
	duration := m.statusMsgExpire.Sub(time.Now())
	if duration <= 0 {
		return nil
	}
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return statusMsgClear{}
	})

}

func (m model) View() string {

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("current count : %d\n\n", m.count))

	if len(m.tasks) == 0 {

	} else {
		for i, task := range m.tasks {
			s.WriteString(fmt.Sprintf("%d: %s\n", i+1, task))
		}
	}
	s.WriteString("\n")

	if m.statusMsg != "" && time.Now().Before(m.statusMsgExpire) {
		s.WriteString(statusMessageStyle.Render(m.statusMsg))
		s.WriteString("\n\n")
	}

	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	s.WriteString("press Enter to submit, backspace/delete to remove last\n")
	s.WriteString("press q or ctrl+c to quit\n")

	// handle the error finally
	if m.err != nil {
		s.WriteString(fmt.Sprintf("%s\n", m.err))
	}
	return border.Width(m.width - 2).Render(s.String())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
}

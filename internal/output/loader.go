package output

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type loader struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func InitLoader() loader {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return loader{spinner: s}
}

func (l loader) Init() tea.Cmd {
	return l.spinner.Tick
}

func (l loader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			l.quitting = true
			return l, tea.Quit
		default:
			return l, nil
		}

	case errMsg:
		l.err = msg
		return l, nil

	default:
		var cmd tea.Cmd
		l.spinner, cmd = l.spinner.Update(msg)
		return l, cmd
	}
}

func (l loader) View() string {
	if l.err != nil {
		return l.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s\n\n", l.spinner.View())
	if l.quitting {
		return str + "\n"
	}
	return str
}

// func main() {
// 	p := tea.NewProgram(InitLoader())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

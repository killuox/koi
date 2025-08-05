package output

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/killuox/koi/internal/api"
)

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
)

func ShowResponse(r api.Result) {
	colorCode := getColorForStatus(r.Status)

	title := fmt.Sprintf("%s%v%s • %s %s • %vms",
		colorCode, r.Status, ColorReset, r.Method, r.Url, r.Duration.Milliseconds())
	p := tea.NewProgram(
		Pager{content: string(r.Body), title: title},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

func getColorForStatus(status int) string {
	if status >= 200 && status <= 299 {
		return ColorGreen
	} else if status >= 500 && status <= 599 {
		return ColorRed
	}
	return ColorYellow // Default for other statuses (e.g., 3xx, 4xx)
}

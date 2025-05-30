package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinser/pacmanai/internal/app"
)

func main() {
	p := tea.NewProgram(app.NewModel())
	if _, err := p.Run(); err != nil {
		println("Error:", err)
		os.Exit(1)
	}
}

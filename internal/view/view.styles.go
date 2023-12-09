package view

import "github.com/charmbracelet/lipgloss"

const (
	Padding  = 1
	MaxWidth = 50
)

// Text
var TextOpacity = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var TextPrimary = lipgloss.NewStyle().Foreground(lipgloss.Color("#9286D9")).Render

// Http method
var HttpPostStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#2EC6EC")).
	Background(lipgloss.Color("#0A2536")).
	PaddingLeft(1).PaddingRight(1).Render

var HttpGetStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#21BD9C")).
	Background(lipgloss.Color("#072A26")).
	PaddingLeft(1).PaddingRight(1).Render

// Colors
var coral1 = lipgloss.Color("#C385B6")
var coral2 = lipgloss.Color("#D685A6")
var coral3 = lipgloss.Color("#E0859E")
var purple1 = lipgloss.Color("#B585BF")
var purple2 = lipgloss.Color("#A585C9")
var purple3 = lipgloss.Color("#9286D9")

func (v *View) GetColumnColor(col int) lipgloss.Color {
	switch col {
	case 1:
		return coral1
	case 2:
		return coral2
	case 3:
		return coral3
	case 4:
		return purple1
	case 5:
		return purple2
	case 6:
		return purple3
	}

	return lipgloss.Color("#ffffff")
}

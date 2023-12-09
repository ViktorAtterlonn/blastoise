package commands

import (
	"blastoise/internal/services"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	padding      = 1
	maxWidth     = 50
	tickInterval = 1
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var primaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9286D9")).Render

type tickMsg time.Time

type Model struct {
	ctx      Ctx
	progress progress.Model
}

func NewModel(ctx Ctx) Model {
	return Model{
		ctx:      ctx,
		progress: progress.New(progress.WithGradient("#9286D9", "#EF7F9E")),
	}
}

func (m Model) Start() {
	go m.ctx.service.Execute(m.ctx.url, m.ctx.method, m.ctx.rps, m.ctx.duration, m.ctx.requestsChan)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

	results := <-m.ctx.requestsChan

	pad := strings.Repeat(" ", padding)

	fmt.Println("\n" +
		pad + renderHttpMethod("GET") + " " + helpStyle(m.ctx.url+fmt.Sprintf(" (%d requests/second)", m.ctx.rps)) + "\n\n" +
		m.renderResults(results) + "\n\n")
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		percent := tickInterval / float64(m.ctx.duration)

		cmd := m.progress.IncrPercent(percent)

		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m Model) View() string {
	pad := strings.Repeat(" ", padding)

	return "\n" +
		pad + renderHttpMethod("GET") + " " + helpStyle(m.ctx.url) + "\n\n" +
		pad + helpStyle(fmt.Sprintf("Sending %s rps for %d seconds", primaryStyle(fmt.Sprintf("%d", m.ctx.rps)), m.ctx.duration)) + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press ctrl + c to quit")

}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

var httpPostStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#2EC6EC")).
	Background(lipgloss.Color("#0A2536")).
	PaddingLeft(1).PaddingRight(1).Render
var httpGetStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#21BD9C")).
	Background(lipgloss.Color("#072A26")).
	PaddingLeft(1).PaddingRight(1).Render

func renderHttpMethod(method string) string {
	switch method {
	case "POST":
		return httpPostStyle(method)
	case "GET":
		return httpGetStyle(method)
	}

	return ""
}

var coral1 = lipgloss.Color("#C385B6")
var coral2 = lipgloss.Color("#D685A6")
var coral3 = lipgloss.Color("#E0859E")
var purple1 = lipgloss.Color("#B585BF")
var purple2 = lipgloss.Color("#A585C9")
var purple3 = lipgloss.Color("#9286D9")

func getColumnColor(col int) lipgloss.Color {
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

func (m Model) renderResults(requests []*services.RequestResult) string {
	responseTimes := m.ctx.service.SummarizeReponseTimesInPercentiles(requests)
	statusCodes := m.ctx.service.SummarizeStatusCodes(requests)

	rows := [][]string{
		{"Response time",
			fmt.Sprintf("%d ms", responseTimes["P10"]),
			fmt.Sprintf("%d ms", responseTimes["P20"]),
			fmt.Sprintf("%d ms", responseTimes["P50"]),
			fmt.Sprintf("%d ms", responseTimes["P75"]),
			fmt.Sprintf("%d ms", responseTimes["P90"]),
			fmt.Sprintf("%d ms", responseTimes["P99"]),
		},
	}

	baseStyle := lipgloss.NewStyle().Padding(0, 1)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Foreground(getColumnColor(col)).Bold((true)).Padding(0, 1)
			case row%2 == 0:
				return baseStyle
			default:
				return baseStyle
			}
		}).
		Headers("Stat", "P10", "P20", "P50", "P75", "P90", "P99").
		Rows(rows...)

	rows2 := [][]string{
		{"Count",
			fmt.Sprintf("%d", statusCodes[200]),
			fmt.Sprintf("%d", statusCodes[400]),
			fmt.Sprintf("%d", statusCodes[500]),
		},
	}

	t2 := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold((true)).Padding(0, 1)
			case row%2 == 0:
				return baseStyle
			default:
				return baseStyle
			}
		}).
		Headers("Status", "2xx", "4xx", "5xx").
		Rows(rows2...)

	return t.Render() + "\n\n" + t2.Render()
}

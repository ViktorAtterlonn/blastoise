package view

import (
	"blastoise/internal/structs"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	tickInterval = 1
)

type tickMsg time.Time

type View struct {
	ctx      *structs.Ctx
	progress progress.Model
}

func NewView(ctx *structs.Ctx) View {
	return View{
		ctx:      ctx,
		progress: progress.New(progress.WithGradient("#9286D9", "#EF7F9E")),
	}
}

func (v View) Init() tea.Cmd {
	return tickCmd()
}

func (v View) Start() {

	if _, err := tea.NewProgram(v).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

	results := <-v.ctx.ResultChan

	pad := strings.Repeat(" ", Padding)

	fmt.Println("\n\n" +
		pad + v.renderHttpMethod() + " " + Text(v.ctx.Url) + TextOpacity(fmt.Sprintf(" (%d requests/second)", v.ctx.Rps)) + "\n\n" +
		v.renderResults(results) + "\n\n")
}

func (v View) View() string {
	pad := strings.Repeat(" ", Padding)

	return "\n\n" +
		pad + v.renderHttpMethod() + " " + Text(v.ctx.Url) + "\n\n" +
		pad + TextOpacity(fmt.Sprintf("Sending %d requests/second for %d seconds", v.ctx.Rps, v.ctx.Duration)) + "\n\n" +
		pad + v.progress.View() + "\n\n" +
		pad + TextOpacity("Press ctrl + c to quit")
}

func (v View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return v, tea.Quit
		}

		return v, nil

	case tea.WindowSizeMsg:
		v.progress.Width = msg.Width - Padding*2 - 4
		if v.progress.Width > MaxWidth {
			v.progress.Width = MaxWidth
		}
		return v, nil

	case tickMsg:
		if v.progress.Percent() == 1.0 {
			return v, tea.Quit
		}

		percent := tickInterval / float64(v.ctx.Duration)

		cmd := v.progress.IncrPercent(percent)

		return v, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := v.progress.Update(msg)
		v.progress = progressModel.(progress.Model)
		return v, cmd

	default:
		return v, nil
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (v View) renderHttpMethod() string {
	switch v.ctx.Method {
	case "POST":
		return HttpPostStyle(v.ctx.Method)
	case "GET":
		return HttpGetStyle(v.ctx.Method)
	}

	return ""
}

func (v View) renderResults(results []*structs.RequestResult) string {

	responseTimes := v.SummarizeReponseTimesInPercentiles(results)
	statusCodes := v.SummarizeStatusCodes(results)

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
				return lipgloss.NewStyle().Foreground(v.GetColumnColor(col)).Bold((true)).Padding(0, 1)
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

func (s *View) SummarizeStatusCodes(results []*structs.RequestResult) map[int]int {
	statusCodes := make(map[int]int)

	for _, result := range results {
		statusCodes[result.StatusCode]++
	}

	return statusCodes
}

// P10, P20, P50, P75, P90, P95, P99
func (s *View) SummarizeReponseTimesInPercentiles(results []*structs.RequestResult) map[string]int {
	percentiles := make(map[string]int)

	// Step 1: Sort the results by duration
	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration < results[j].Duration
	})

	// Step 2: Calculate percentiles
	totalResults := len(results)
	percentiles["P10"] = calculatePercentile(results, 10, totalResults)
	percentiles["P20"] = calculatePercentile(results, 20, totalResults)
	percentiles["P50"] = calculatePercentile(results, 50, totalResults)
	percentiles["P75"] = calculatePercentile(results, 75, totalResults)
	percentiles["P90"] = calculatePercentile(results, 90, totalResults)
	percentiles["P95"] = calculatePercentile(results, 95, totalResults)
	percentiles["P99"] = calculatePercentile(results, 99, totalResults)

	// Step 3: Return the percentiles
	return percentiles
}

func calculatePercentile(results []*structs.RequestResult, percentile int, totalResults int) int {
	if totalResults == 0 {
		return 0
	}

	index := int(math.Ceil(float64(percentile)/100*float64(totalResults))) - 1
	if index < 0 {
		index = 0
	} else if index >= totalResults {
		index = totalResults - 1
	}

	return results[index].Duration
}

package progress

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressMsg represents a progress update for the custom progress bar.
type ProgressMsg struct {
	Current int64
	Total   int64
	Done    bool
	Error   error
}

// model holds the state for the custom progress bar UI.
type model struct {
	current      int64
	total        int64
	err          error
	done         bool
	startTime    time.Time
	width        int // Terminal width for responsive bar
	quitting     bool
	progressChan chan ProgressMsg // Channel to receive updates from the runner
}

// Progress bar styles (copied from your old application)
var (
	progressBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("205")).
				Foreground(lipgloss.Color("230"))

	progressEmptyStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("240"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)

// NewModel initializes a new custom progress bar model for Bubble Tea.
// It returns the model and the channel to which progress updates should be sent.
func NewModel(total int64) (model, chan ProgressMsg) {
	progressCh := make(chan ProgressMsg, 100) // Buffered channel for updates
	m := model{
		total:        total,
		current:      0,
		done:         false,
		startTime:    time.Now(),
		progressChan: progressCh,
		width:        80, // Default width, will be updated by WindowSizeMsg
	}
	return m, progressCh
}

// Init is the first command that is run when the program starts.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,  // Enter alternative screen buffer for a cleaner TUI
		m.waitForProgress(), // Start listening for progress updates immediately
	)
}

// Update handles messages sent to the Bubble Tea program.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		// Handle window resizing to make the progress bar responsive.
		m.width = msg.Width
		if m.width > 80 { // Cap the width at 80 characters for readability
			m.width = 80
		}
		return m, nil

	case ProgressMsg:
		// Update the model's state based on the received progress message.
		m.current = msg.Current
		m.total = msg.Total
		m.done = msg.Done
		m.err = msg.Error

		if m.done || m.err != nil {
			m.quitting = true
			return m, tea.Quit // Quit the program when done or an error occurs
		}
		// Continue waiting for the next progress update.
		return m, m.waitForProgress()

	default:
		return m, nil
	}
}

// View renders the UI.
func (m model) View() string {
	if m.quitting {
		if m.err != nil {
			return errorStyle.Render(fmt.Sprintf("\n❌ Error: %v\n\n", m.err))
		}
		return "" // Clear the screen if quitting without error
	}

	progressVal := float64(m.current) / float64(m.total)
	if progressVal > 1 {
		progressVal = 1 // Cap at 100%
	}

	elapsed := time.Since(m.startTime)
	var eta time.Duration
	if m.current > 0 && progressVal > 0 {
		eta = time.Duration(float64(elapsed) * (1.0 - progressVal) / progressVal)
	}

	// Calculate bar width dynamically based on terminal width
	barWidth := m.width - 30 // Allocate space for text around the bar
	if barWidth < 10 {       // Ensure minimum bar width
		barWidth = 10
	}

	filled := int(progressVal * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	var progressBar string
	if filled > 0 {
		progressBar = progressBarStyle.Render(strings.Repeat("█", filled))
	}
	if filled < barWidth {
		progressBar += progressEmptyStyle.Render(strings.Repeat("░", barWidth-filled))
	}

	percentage := progressVal * 100
	status := fmt.Sprintf("Progress: %s/%s (%.1f%%)",
		formatNumber(m.current), formatNumber(m.total), percentage)

	timing := fmt.Sprintf("Elapsed: %v", elapsed.Round(time.Second))
	if eta > 0 && !m.done {
		timing += fmt.Sprintf(" | ETA: %v", eta.Round(time.Second))
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n%s\n\n%s\n",
		statusStyle.Render("Generating data..."),
		progressBar,
		status,
		timing,
		lipgloss.NewStyle().Faint(true).Render("Press q, Ctrl+C, or Esc to quit"),
	)
}

// waitForProgress is a command that waits for the next ProgressMsg on the channel.
func (m model) waitForProgress() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.progressChan
		if !ok {
			// Channel closed, signal completion
			return ProgressMsg{Done: true}
		}
		return msg
	}
}

// formatNumber formats a number with comma separators for better readability.
func formatNumber(n int64) string {
	str := strconv.FormatInt(n, 10)
	if len(str) <= 3 {
		return str
	}
	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}
	return result.String()
}

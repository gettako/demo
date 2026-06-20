// history.go
// Description: Interactive history table with Bubbletea native integration.

package components

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"gettako.dev/tako/contracts"
)

type History struct {
	app   contracts.Application
	table table.Model
}

func NewHistory(app contracts.Application) *History {
	return &History{app: app}
}

func (c *History) ID() string { return "history" }

func (c *History) Init(ctx contracts.Context) {
	columns := []table.Column{
		{Title: "ID", Width: 15},
		{Title: "Date", Width: 12},
		{Title: "Dest", Width: 15},
		{Title: "Amount", Width: 10},
		{Title: "Status", Width: 10},
	}

	var history []Transaction
	_ = ctx.KV().Get("history.all", &history)

	var rows []table.Row
	for _, tx := range history {
		rows = append(rows, table.Row{
			tx.ID,
			tx.Timestamp.Format("Jan 02"),
			tx.Dest,
			fmt.Sprintf("%d", tx.Amount),
			tx.Status,
		})
	}

	if len(rows) == 0 {
		rows = append(rows, table.Row{"-", "-", "No history", "-", "-"})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(PrimaryColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#ffffff")).
		Background(PrimaryColor).
		Bold(false)
	t.SetStyles(s)

	c.table = t
}

func (c *History) Reload() {
	var history []Transaction
	_ = c.app.KV().Get("history.all", &history)

	var rows []table.Row
	for _, tx := range history {
		rows = append(rows, table.Row{
			tx.ID,
			tx.Timestamp.Format("Jan 02"),
			tx.Dest,
			fmt.Sprintf("%d", tx.Amount),
			tx.Status,
		})
	}
	if len(rows) == 0 {
		rows = append(rows, table.Row{"-", "-", "No history", "-", "-"})
	}
	c.table.SetRows(rows)
}

func (c *History) Keys(km contracts.KeyManager) {}

func (c *History) Update(msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	c.table, cmd = c.table.Update(msg)
	c.app.UI().RequestRender()

	// Only return true if we handled a navigation event
	if _, ok := msg.(tea.KeyMsg); ok {
		return cmd, true
	}

	// Check for storage updates
	if _, ok := msg.(contracts.Event); ok {
		c.Reload()
		return cmd, true
	}

	return cmd, false
}

func (c *History) Render() string {
	width, _ := c.app.UI().Dimensions()

	// Calculate dynamic column widths (accounting for borders and padding)
	boxWidth := width - 14
	boxWidth = max(50, boxWidth) // Minimum width

	// Bubbletea tables add padding to each cell (Padding(0, 1) = 2 chars per column).
	// For 5 columns, that's 10 extra characters.
	availableForCols := boxWidth - 10
	availableForCols = max(20, availableForCols)

	colID := int(float32(availableForCols) * 0.20)
	colDate := int(float32(availableForCols) * 0.15)
	colDest := int(float32(availableForCols) * 0.30)
	colAmt := int(float32(availableForCols) * 0.20)
	colStat := availableForCols - colID - colDate - colDest - colAmt // Remaining

	c.table.SetColumns([]table.Column{
		{Title: "ID", Width: colID},
		{Title: "Date", Width: colDate},
		{Title: "Dest", Width: colDest},
		{Title: "Amount", Width: colAmt},
		{Title: "Status", Width: colStat},
	})
	c.table.SetWidth(boxWidth)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 4).
		MarginBottom(1).
		Width(width - 4)

	title := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("Transaction History")

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", c.table.View()))
}

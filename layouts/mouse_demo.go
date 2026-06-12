// Package main provides an demo demonstrating the Mouse Router feature.
//
// This demo shows how to:
//   - Register clickable zones with hitboxes
//   - Handle left/right click and scroll events
//   - Implement drag-and-drop between two panels
//
// Package layouts provides UI layout implementations for the Tako demo app.
package layouts

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/gettako/tako/internal/router"
	takoCtx "github.com/gettako/tako/internal/tako"
)

// ── App State ─────────────────────────────────────────────────────────────────

// MouseState holds the state for the Mouse Router demo panels.
var MouseState = struct {
	Items     []string
	Dropped   []string
	Selected  int
	LastEvent string
}{
	Items:   []string{"🍎 Apple", "🍌 Banana", "🍇 Grapes", "🍊 Orange", "🍓 Strawberry"},
	Dropped: []string{},
}

// ── Mouse Demo Layout ─────────────────────────────────────────────────────────

// MouseDemoLayout renders two side-by-side panels and updates hitboxes on every
// render so the Mouse Router always has accurate positions.
type MouseDemoLayout struct{}

func (l *MouseDemoLayout) View(ctx *takoCtx.Context, r *router.Router) tea.View {
	const panelW = 30

	panelH := len(MouseState.Items) + 5
	if panelH < 10 {
		panelH = 10
	}

	// ── Left Panel ───────────────────────────────────────────────────────────
	leftLines := []string{" LEFT PANEL (click/scroll/drag) ", ""}
	for i, item := range MouseState.Items {
		cursor := "  "
		if i == MouseState.Selected {
			cursor = "▶ "
		}
		leftLines = append(leftLines, fmt.Sprintf(" %s%d. %s", cursor, i+1, item))
	}
	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(panelW).Height(panelH).
		Render(strings.Join(leftLines, "\n"))

	// Hitbox for left panel starts at col=0, row=2 (after header + gap)
	ctx.Mouse().UpdateHitbox("left-panel", 0, 2, panelW+2, panelH+2)

	// ── Right Panel ──────────────────────────────────────────────────────────
	rightLines := []string{" RIGHT PANEL (drop items here) ", ""}
	if len(MouseState.Dropped) == 0 {
		rightLines = append(rightLines, "  (empty — drag from left panel)")
	} else {
		for _, item := range MouseState.Dropped {
			rightLines = append(rightLines, "  ✓ "+item)
		}
	}
	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#10B981")).
		Width(panelW).Height(panelH).
		Render(strings.Join(rightLines, "\n"))

	rightX := panelW + 4 // left panel width + border (2) + gap (2)
	ctx.Mouse().UpdateHitbox("right-panel", rightX, 2, panelW+2, panelH+2)

	// ── Header / Status ───────────────────────────────────────────────────────
	header := lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).Padding(0, 1).
		Render(" 🖱️  TAKO MOUSE ROUTER DEMO ")

	status := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(
		"Last: %s  |  Left: %d items  |  Right: %d items  |  [ctrl+c] Quit",
		MouseState.LastEvent, len(MouseState.Items), len(MouseState.Dropped),
	))

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, "  ", rightBox)
	full := lipgloss.JoinVertical(lipgloss.Left, header, "", panels, "", status)

	v := tea.NewView(full)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

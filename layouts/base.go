package layouts

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/gettako/tako/internal/router"
	internal_tako "github.com/gettako/tako/internal/tako"
)

// AppLayout defines the visual UI layout for the miniapp.
type AppLayout struct{}

func (l *AppLayout) View(ctx *internal_tako.Context, r *router.Router) tea.View {
	// 1. Header
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)
	header := headerStyle.Render(" 🐙 TAKO FRAMEWORK - MINI APP DEMO ")

	// 2. Sidebar Widgets
	var sidebarWidgets []string
	if ctx != nil && ctx.Hooks() != nil {
		raw := ctx.Hooks().All("app.sidebar")
		for _, w := range raw {
			if s, ok := w.(string); ok {
				sidebarWidgets = append(sidebarWidgets, s)
			}
		}
	}
	sidebar := lipgloss.JoinVertical(lipgloss.Left, sidebarWidgets...)

	// 3. Main Content
	var mainContent strings.Builder
	mainContent.WriteString("Welcome to the Tako Base Layer!\n")
	mainContent.WriteString("This area represents the main workspace zone.\n\n")

	if r != nil && r.Stack() != nil {
		top := r.Stack().Top()
		if top != "base" {
			// Try to render overlay from hooks
			hookName := "tako.overlay." + top
			overlay := ctx.Hooks().Get(hookName)
			
			// Fallback to legacy app.overlay if not found
			if overlay == nil {
				hookName = "app.overlay." + top
				overlay = ctx.Hooks().Get(hookName)
			}
			
			if s, ok := overlay.(string); ok {
				if strings.HasPrefix(top, "tako.dialog.") {
					// Apply dialog styling
					dialogStyle := lipgloss.NewStyle().
						Border(lipgloss.ThickBorder()).
						BorderForeground(lipgloss.Color("#EF4444")).
						Padding(1, 2).
						Align(lipgloss.Center)
					mainContent.WriteString(dialogStyle.Render("⚠️ DIALOG\n\n" + s))
				} else {
					mainContent.WriteString(s)
				}
			} else {
				fmt.Fprintf(&mainContent, " [OVERLAY ACTIVE: %s] \n No hook provider for %s found.", top, hookName)
			}
		} else {
			mainContent.WriteString(" [IDLE] No active overlays. \n Try pressing [ctrl+f] or [?].")
		}
	}

	mainPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Width(60).
		Height(15).
		Render(mainContent.String())

	// 4. Footer
	var footerWidgets []string
	if ctx != nil && ctx.Hooks() != nil {
		raw := ctx.Hooks().All("app.footer")
		for _, w := range raw {
			if s, ok := w.(string); ok {
				footerWidgets = append(footerWidgets, s)
			}
		}
	}
	footerLine := lipgloss.NewStyle().
		Faint(true).
		Render("Keys: [ctrl+f] Search | [ctrl+d] Dash | [?] Help | [ctrl+s] Status | [ctrl+c] Exit ")
	footer := lipgloss.JoinHorizontal(lipgloss.Bottom, footerLine, lipgloss.JoinHorizontal(lipgloss.Left, footerWidgets...))

	// Combine Everything
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, "  ", mainPanel)
	fullView := lipgloss.JoinVertical(lipgloss.Left, header, "", body, "", footer)

	view := tea.NewView(fullView)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

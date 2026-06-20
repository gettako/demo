// layout.go
// Description: The root layout component of TakoBank.

package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type Layout struct {
	app contracts.Application
}

func NewLayout(app contracts.Application) *Layout {
	return &Layout{app: app}
}

func (c *Layout) ID() string { return "layout" }

func (c *Layout) Init(ctx contracts.Context) {
	// Initialize default state if not present
	if !ctx.State().Has("balance") {
		ctx.State().Mutate("balance").Value(15000000.0).Broadcast()
	}
	if !ctx.State().Has("pin_entered") {
		ctx.State().Mutate("pin_entered").Value(false).Broadcast()
	}
}

func (c *Layout) Keys(km contracts.KeyManager) {
	// Global Keybindings
	km.Bind("ctrl+q", func() {
		c.app.Shutdown()
	})

	km.Bind("ctrl+t", func() {
		// Toggle Theme
		current := c.app.Config().Get("app.theme", "dark").String()
		next := "light"
		switch current {
		case "light":
			next = "cyberpunk"
		case "cyberpunk":
			next = "dark"
		}
		c.app.Config().Set("app.theme", next)
		c.app.Events().Dispatch(contracts.EventThemeChanged, next)
	})

	km.Bind("ctrl+l", func() {
		// Toggle Language
		lang := c.app.Lang()
		current := lang.Active()
		if current == "en" {
			lang.SetLocale("id")
		} else {
			lang.SetLocale("en")
		}
	})

	// Add navigation global shortcuts
	km.Bind("d", func() { c.app.Router().GoTo("dashboard") })
	km.Bind("t", func() { c.app.Router().GoTo("transfer") })
	km.Bind("h", func() { c.app.Router().GoTo("history") })
	km.Bind("p", func() { c.app.Router().GoTo("profile") })
	km.Bind("s", func() { c.app.Router().GoTo("settings") })

	// Arrow keys for tab navigation
	routes := []string{"dashboard", "transfer", "history", "profile", "settings"}
	navigate := func(offset int) {
		current := c.app.Router().Current().Name
		idx := -1
		for i, r := range routes {
			if r == current {
				idx = i
				break
			}
		}
		if idx != -1 {
			newIdx := (idx + offset + len(routes)) % len(routes)
			c.app.Router().GoTo(routes[newIdx])
		}
	}
	km.Bind("left", func() { navigate(-1) })
	km.Bind("right", func() { navigate(1) })

	// Number shortcuts (Alt+1,2,3,4) and F-keys
	km.Bind("alt+1 1 f1", func() { c.app.Router().GoTo("dashboard") })
	km.Bind("alt+2 2 f2", func() { c.app.Router().GoTo("transfer") })
	km.Bind("alt+3 3 f3", func() { c.app.Router().GoTo("history") })
	km.Bind("alt+4 4 f4", func() { c.app.Router().GoTo("profile") })
	km.Bind("alt+5 5 f5", func() { c.app.Router().GoTo("settings") })
}

//nolint:funlen
func (c *Layout) Render() string {
	width, height := c.app.UI().Dimensions()
	lang := c.app.Lang()

	// Layout usable inner width (Terminal width - 4 to account for app border and padding)
	innerWidth := width - 4
	innerWidth = max(10, innerWidth)

	// HEADER
	title := lang.T("app.title")
	titleText := fmt.Sprintf("🏦 %s 🏦", title)
	padLen := (innerWidth - lipgloss.Width(titleText)) / 2
	padLen = max(0, padLen)
	leftPad := strings.Repeat(" ", padLen)
	rightPad := innerWidth - lipgloss.Width(titleText) - padLen
	rightPad = max(0, rightPad)
	headerText := leftPad + titleText + strings.Repeat(" ", rightPad)
	header := lipgloss.NewStyle().
		Background(PrimaryColor).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Render(headerText)

	// NAV BAR
	currentRoute := c.app.Router().Current().Path
	navItems := []string{
		c.renderNavItem("[F1] "+lang.T("nav.dashboard"), "/dashboard", currentRoute),
		c.renderNavItem("[F2] "+lang.T("nav.transfer"), "/transfer", currentRoute),
		c.renderNavItem("[F3] "+lang.T("nav.history"), "/history", currentRoute),
		c.renderNavItem("[F4] Profile", "/profile", currentRoute),
		c.renderNavItem("[F5] Settings", "/settings", currentRoute),
	}
	navBar := lipgloss.JoinHorizontal(lipgloss.Center, navItems...)
	navBar = lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Center).PaddingTop(1).PaddingBottom(1).Render(navBar)

	// BODY (Active Route View)
	bodyContent := c.app.UI().RenderView()

	// FOOTER HOOKS
	footerHooks := c.renderFooterHooks(innerWidth)

	// Layout Middle section
	middle := lipgloss.NewStyle().Width(innerWidth).Render(bodyContent)

	// Calculate remaining height to push footer down
	usedHeight := lipgloss.Height(header) + lipgloss.Height(navBar) + lipgloss.Height(middle)
	if footerHooks != "" {
		usedHeight += lipgloss.Height(footerHooks)
	}

	padding := height - usedHeight - 1 // 1 for base footer
	padding = max(0, padding)
	middle = lipgloss.JoinVertical(lipgloss.Left, middle, strings.Repeat("\n", padding))

	if footerHooks != "" {
		middle = lipgloss.JoinVertical(lipgloss.Left, middle, lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Center).Render(footerHooks))
	}

	// FOOTER
	activeLang := strings.ToUpper(lang.Active())
	footerText := fmt.Sprintf("global: [ctrl+q] quit | [ctrl+t] theme | [ctrl+l] lang (%s) | nav: [← →] change tab", activeLang)
	footer := lipgloss.NewStyle().
		Foreground(MutedColor).
		Width(innerWidth).
		Align(lipgloss.Center).
		Render(footerText)

	// Combine components
	content := lipgloss.JoinVertical(lipgloss.Left, header, navBar, middle, footer)

	// Apply an app border frame
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Width(width - 2).
		Height(height - 2).
		Render(content)
}

func (c *Layout) renderNavItem(label, route, currentRoute string) string {
	style := lipgloss.NewStyle().Padding(0, 2)
	if route == currentRoute {
		style = style.Foreground(PrimaryColor).Bold(true).Underline(true)
	} else {
		style = style.Foreground(lipgloss.Color("#888888"))
	}
	return style.Render(label)
}

func (c *Layout) renderFooterHooks(innerWidth int) string {
	hooks := c.app.Hooks().All("app.footer")
	var blocks []string
	for _, hook := range hooks {
		if content, ok := hook.(string); ok {
			blocks = append(blocks, content)
		}
	}
	if len(blocks) > 0 {
		// Use lipgloss.JoinHorizontal to properly align multiline blocks
		separator := lipgloss.NewStyle().Foreground(PrimaryColor).Render("│")

		// Distribute width evenly
		blockWidth := (innerWidth - lipgloss.Width(separator)*(len(blocks)-1)) / len(blocks)

		var joinedBlocks []string
		for i, block := range blocks {
			if i > 0 {
				joinedBlocks = append(joinedBlocks, separator)
			}
			// Wrap each block with even width
			wrappedBlock := lipgloss.NewStyle().Width(blockWidth).Align(lipgloss.Center).Render(block)
			joinedBlocks = append(joinedBlocks, wrappedBlock)
		}

		content := lipgloss.JoinHorizontal(lipgloss.Center, joinedBlocks...)
		return lipgloss.NewStyle().
			Padding(1, 0).
			Render(content)
	}
	return ""
}

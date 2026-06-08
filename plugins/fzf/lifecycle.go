package fzf

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/takoterm/tako/contracts"
	"github.com/takoterm/tako/internal/plugin"
	"github.com/takoterm/tako/internal/tako"
)

var fzfStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#7D56F4")).
	Padding(0, 1).
	Width(45)

type Plugin struct {
	plugin.NoopLifecycle
	ctx    *tako.Context
	logger contracts.Logger
}

func (p *Plugin) OnInit(ctx *tako.Context) error {
	p.ctx = ctx

	// Resolve logger from container safely
	var l contracts.Logger
	if err := ctx.Container().Make(&l); err == nil {
		p.logger = l
	}

	if p.logger != nil {
		p.logger.Info("FZF Plugin Initialized inside Mini-App")
	}

	// Register sidebar widget
	ctx.Hooks().Add("app.sidebar", func() any {
		var content string
		if ctx.Overlay().Top() == "fzf" {
			content = "🔍 [FZF SEARCH ACTIVE]\nQuery: (type to filter)\nPress [esc] to dismiss"
		} else {
			content = "🔍 [FZF] Press [ctrl+f] to Search"
		}
		return fzfStyle.Render(content)
	})

	// ── New API: single Show() call replaces Stack.Push + Router.Focus + Hooks.Set ──

	// Open FZF search overlay
	ctx.Keys().Bind("ctrl+f", func() {
		if ctx.Overlay().IsActive() {
			return // already showing something — ignore to prevent double-push
		}
		if p.logger != nil {
			p.logger.Info("FZF Shortcut ctrl+f captured: opening search layer")
		}
		ctx.Overlay().Show("fzf", p.render)
		ctx.Emit("fzf:opened", "FZF search layer activated")
	})

	return nil
}

func (p *Plugin) OnActivate(ctx *tako.Context) error {
	if p.logger != nil {
		p.logger.Info("FZF Plugin activated")
	}
	return nil
}

// render is the FZF overlay view function — returns any (UI-agnostic).
func (p *Plugin) render() any {
	var sb strings.Builder
	sb.WriteString("  FZF SEARCH OVERLAY  \n")
	sb.WriteString("======================\n\n")
	sb.WriteString(" [ ] app.go\n")
	sb.WriteString(" [ ] main.go\n")
	sb.WriteString(" [ ] plugin.go\n")
	sb.WriteString(" [ ] README.md\n\n")
	sb.WriteString("Type your search query...")
	return fzfStyle.Render(sb.String())
}

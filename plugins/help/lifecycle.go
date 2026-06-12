package help

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/gettako/tako/contracts"
	"github.com/gettako/tako/internal/plugin"
	"github.com/gettako/tako/internal/tako"
)

var helpStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("#FFD700")).
	Padding(1, 2).
	Width(60)

type Plugin struct {
	plugin.NoopLifecycle
	ctx      *tako.Context
	isActive bool
}

func (p *Plugin) OnInit(ctx *tako.Context) error {
	p.ctx = ctx

	// Register footer shortcut hint
	ctx.Hooks().Add("app.footer", func() any {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render(" [?] Help ")
	})

	// Register overlay view hook
	ctx.Hooks().Set("app.overlay.help", func() any {
		if !p.isActive {
			return nil
		}
		var sb strings.Builder
		sb.WriteString("  TAKO MINI-APP HELP  \n")
		sb.WriteString("======================\n\n")
		sb.WriteString("Global Shortcuts:\n")
		sb.WriteString(" [ctrl+c]  Exit Application\n")
		sb.WriteString(" [ctrl+f]  Open Search (FZF)\n")
		sb.WriteString(" [?]       Toggle this Help\n")
		sb.WriteString(" [esc]     Close current overlay\n\n")
		sb.WriteString("Navigation:\n")
		sb.WriteString(" [h]       Show base logger message\n")
		sb.WriteString(" [ctrl+s]  Trigger diagnostic ping\n\n")
		sb.WriteString("Press any key to close this help.")
		return helpStyle.Render(sb.String())
	})

	// Register global keybinding to open help layer
	ctx.Keys().Bind("?", func() {
		if p.isActive {
			return
		}
		p.isActive = true
		ctx.Stack().Push("help")
		ctx.Emit("help:opened", nil)
	})

	// Subscribe to close event
	ctx.On("help:closed", func(e contracts.Event) {
		p.isActive = false
	})

	return nil
}

func (p *Plugin) OnActivate(_ *tako.Context) error {
	return nil
}

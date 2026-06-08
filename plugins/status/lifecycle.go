package status

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/takoterm/tako/contracts"
	"github.com/takoterm/tako/internal/plugin"
	"github.com/takoterm/tako/internal/tako"
)

var statusStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#04B575")).
	Padding(0, 1).
	Width(45)

type Plugin struct {
	plugin.NoopLifecycle
	ctx          *tako.Context
	logger       contracts.Logger
	config       contracts.Config
	fzfCount     int
	currentState string
	lastEvents   []string
	subCancel    func()
}

func (p *Plugin) OnInit(ctx *tako.Context) error {
	p.ctx = ctx
	p.currentState = "Idle"

	// Resolve Logger from container safely
	var l contracts.Logger
	if err := ctx.Container().Make(&l); err == nil {
		p.logger = l
	}

	// Resolve Config from container safely
	var cfg contracts.Config
	if err := ctx.Container().Make(&cfg); err == nil {
		p.config = cfg
	}

	if p.logger != nil {
		p.logger.Info("Status Plugin Initialized inside Mini-App")
	}

	// Register sidebar widget
	ctx.Hooks().Add("app.sidebar", func() any {
		var sb strings.Builder
		sb.WriteString("📊  [SYSTEM MONITOR]\n")
		fmt.Fprintf(&sb, "State: %s\n", p.currentState)
		fmt.Fprintf(&sb, "FZF Opens: %d\n", p.fzfCount)

		debugVal := "false"
		if p.config != nil {
			debugVal = fmt.Sprintf("%v", p.config.Bool("app.debug"))
		}
		fmt.Fprintf(&sb, "App Debug Mode: %s\n\n", debugVal)

		sb.WriteString("Recent Events:\n")
		if len(p.lastEvents) == 0 {
			sb.WriteString(" - None")
		} else {
			for i, ev := range p.lastEvents {
				fmt.Fprintf(&sb, " %d. %s\n", i+1, ev)
			}
		}
		return statusStyle.Render(sb.String())
	})

	// Register footer widget
	ctx.Hooks().Add("app.footer", func() any {
		color := "#04B575"
		if p.currentState == "Searching" {
			color = "#7D56F4"
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color(color)).
			Padding(0, 1).
			Render(" STATE: " + p.currentState + " ")
	})

	// Register global keybinding for diagnostic ping
	ctx.Keys().Bind("ctrl+s", func() {
		if p.logger != nil {
			p.logger.Info("Diagnostic shortcut ctrl+s captured!")
		}
		p.addEvent("Diagnostic Ping Triggered")
	})

	return nil
}

func (p *Plugin) OnActivate(ctx *tako.Context) error {
	if p.logger != nil {
		p.logger.Info("Status Plugin activated")
	}

	// Subscribe to events — use ctx as the scoped context so subscriptions
	// are automatically pruned when the plugin context is cancelled.
	cancelFZFOpened := ctx.Subscribe(ctx, "fzf:opened", func(e contracts.Event) {
		p.currentState = "Searching"
		p.fzfCount++
		p.addEvent("FZF Layer Opened")
	})

	cancelFZFClosed := ctx.Subscribe(ctx, "fzf:closed", func(e contracts.Event) {
		p.currentState = "Idle"
		p.addEvent("FZF Layer Closed")
	})

	cancelHelpOpened := ctx.Subscribe(ctx, "help:opened", func(e contracts.Event) {
		p.currentState = "Reading Help"
		p.addEvent("Help Layer Opened")
	})

	cancelHelpClosed := ctx.Subscribe(ctx, "help:closed", func(e contracts.Event) {
		p.currentState = "Idle"
		p.addEvent("Help Layer Closed")
	})

	p.subCancel = func() {
		cancelFZFOpened()
		cancelFZFClosed()
		cancelHelpOpened()
		cancelHelpClosed()
	}

	return nil
}

func (p *Plugin) OnDeactivate(_ *tako.Context) error {
	if p.subCancel != nil {
		p.subCancel()
	}
	return nil
}

func (p *Plugin) addEvent(msg string) {
	p.lastEvents = append(p.lastEvents, msg)
	if len(p.lastEvents) > 3 {
		p.lastEvents = p.lastEvents[len(p.lastEvents)-3:]
	}
}

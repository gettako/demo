package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/gettako/tako/contracts"
)

// CounterComponent is a self-contained UI widget implementing contracts.Component.
// It manages its own state, rendering, and keybindings.
type CounterComponent struct {
	count int
	style lipgloss.Style
}

// NewCounterComponent creates a new CounterComponent.
func NewCounterComponent() *CounterComponent {
	return &CounterComponent{
		count: 0,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5F87")).
			Padding(1, 4).
			Align(lipgloss.Center),
	}
}

// ID returns the unique identifier for this component and its key zone.
func (c *CounterComponent) ID() string {
	return "counter_widget"
}

// Render returns the component's visual output as a formatted string.
func (c *CounterComponent) Render() any {
	var sb strings.Builder
	sb.WriteString("👾 Custom Component API 👾\n\n")
	sb.WriteString(fmt.Sprintf("Count: %d\n\n", c.count))
	sb.WriteString(" [Up] increment | [Down] decrement | [Esc] close ")

	return c.style.Render(sb.String())
}

// RegisterKeys wires up the inputs for this specific component.
// These bindings are only active when this component is focused in the router.
func (c *CounterComponent) RegisterKeys(keys contracts.KeyManager) {
	zone := keys.Zone(c.ID())

	zone.Bind("up", func() {
		c.count++
	})

	zone.Bind("down", func() {
		c.count--
	})
	
	// Note: Esc closing is usually handled globally in the main app layout,
	// but components can define their own specific handlers if needed.
}

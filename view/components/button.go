package components

import (
	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type Button struct {
	ctx   contracts.Context
	slots map[string]contracts.Slot
}

func NewButton(app contracts.Application) *Button {
	return &Button{ctx: app.Context()}
}

func (c *Button) ID() string                   { return "button" }
func (c *Button) Init(ctx contracts.Context)   {}
func (c *Button) Keys(km contracts.KeyManager) {}

func (c *Button) SetSlots(slots map[string]contracts.Slot) {
	c.slots = slots
}

func (c *Button) Render() string {
	content := ""
	if defaultSlot, ok := c.slots["default"]; ok {
		content += defaultSlot.Render(c.ctx)
	}

	return lipgloss.NewStyle().
		Background(lipgloss.Color("#cb6843")).
		Foreground(lipgloss.Color("#ffffff")).
		Padding(0, 1).
		MarginRight(1).
		Render(content)
}

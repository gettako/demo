package components

import (
	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type Card struct {
	ctx   contracts.Context
	slots map[string]contracts.Slot
}

func NewCard(app contracts.Application) *Card {
	return &Card{ctx: app.Context()}
}

func (c *Card) ID() string                   { return "card" }
func (c *Card) Init(ctx contracts.Context)   {}
func (c *Card) Keys(km contracts.KeyManager) {}

func (c *Card) SetSlots(slots map[string]contracts.Slot) {
	c.slots = slots
}

func (c *Card) Render() string {
	content := ""
	if headerSlot, ok := c.slots["header"]; ok {
		content += headerSlot.Render(c.ctx) + "\n"
	}
	if defaultSlot, ok := c.slots["default"]; ok {
		content += defaultSlot.Render(c.ctx)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4368cb")).
		Padding(1, 2).
		MarginBottom(1).
		Render(content)
}

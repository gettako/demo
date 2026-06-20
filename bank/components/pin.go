// pin.go
// Description: PIN Entry component.

package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type Pin struct {
	app contracts.Application
	pin string
}

func NewPin(app contracts.Application) *Pin {
	return &Pin{app: app}
}

func (c *Pin) ID() string { return "pin" }

func (c *Pin) Init(ctx contracts.Context) {
	c.pin = ""
}

func (c *Pin) Keys(km contracts.KeyManager) {
	zone := km.Zone("pin")

	zone.Bind("enter", func() {
		// Hardcoded correct PIN is "123456" for demo. Also allow empty PIN bypass.
		if c.pin == "123456" || c.pin == "" {
			c.app.State().Mutate("pin_entered").Value(true).Broadcast()

			// Go to intended URL
			intended := "/dashboard"
			if s := c.app.State().Get("intended_url").String(); s != "" {
				intended = s
			}
			c.app.Router().GoTo(intended)
		} else {
			// Wrong pin, clear it
			c.pin = ""
		}
	})

	zone.Bind("backspace", func() {
		if len(c.pin) > 0 {
			c.pin = c.pin[:len(c.pin)-1]
		}
	})

	// Bind numbers 0-9
	for i := range 10 {
		num := fmt.Sprintf("%d", i)
		zone.Bind(num, func() {
			if len(c.pin) < 6 {
				c.pin += num
			}
		})
	}

	// Focus the zone after it has been created by the KeyManager
	c.app.Input().Focus(0, "pin")
}

func (c *Pin) Render() string {
	lang := c.app.Lang()

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(2, 8).
		Align(lipgloss.Center)

	logo := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render(LogoAscii)

	title := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(lang.T("transfer.pin_required"))

	// Display bullets for PIN
	displayPin := strings.Repeat("•", len(c.pin)) + strings.Repeat("", 6-len(c.pin))
	displayPin = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(displayPin)

	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		title,
		"",
		displayPin,
		"",
		lipgloss.NewStyle().Foreground(MutedColor).Render("[enter] submit"),
	)

	box := boxStyle.Render(content)
	width, height := c.app.UI().Dimensions()
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

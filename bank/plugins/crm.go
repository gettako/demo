// crm.go
// Description: Simulated third-party plugin injecting hooks.

package plugins

import (
	"banking/components"

	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type CRMProvider struct{}

func NewCRMProvider() *CRMProvider {
	return &CRMProvider{}
}

func (p *CRMProvider) Register(app contracts.Application) error {
	return nil
}

func (p *CRMProvider) Boot(app contracts.Application) error {
	app.Hooks().Add("app.footer", func() any {
		return lipgloss.NewStyle().
			Foreground(components.PrimaryColor).
			Align(lipgloss.Center).
			Bold(true).
			Render("✨ SPECIAL OFFER ✨\n\nGet a TakoBank\nPlatinum Card today\nwith 0% Interest!")
	})

	app.Hooks().Add("app.footer", func() any {
		return lipgloss.NewStyle().
			Foreground(components.PrimaryColor).
			Align(lipgloss.Center).
			Render("Need Help?\nCall 14000")
	})

	return nil
}

package providers

import (
	"viewdemo/components"

	"gettako.dev/tako/contracts"
)

type ViewServiceProvider struct{}

func (p *ViewServiceProvider) Register(app contracts.Application) error { return nil }

func (p *ViewServiceProvider) Boot(app contracts.Application) error {
	engine := app.ViewEngine()

	// Mendaftarkan custom tag HTML ke Go struct factory
	engine.RegisterComponent("x-card", func() contracts.Component {
		return components.NewCard(app)
	})
	engine.RegisterComponent("x-button", func() contracts.Component {
		return components.NewButton(app)
	})
	return nil
}

// dashboard.go
// Description: Dashboard component displaying balance and rates.

package components

import (
	"encoding/json"
	"fmt"

	"banking/utils"

	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type Dashboard struct {
	app contracts.Application
}

func NewDashboard(app contracts.Application) *Dashboard {
	return &Dashboard{app: app}
}

func (c *Dashboard) ID() string { return "dashboard" }

func (c *Dashboard) Init(ctx contracts.Context) {
	// Re-render automatically when balance changes
	ctx.State().Watch("balance").OnUpdate(func(oldVal, newVal any) {
		ctx.UI().RequestRender()
	}).Subscribe(ctx.StdCtx())

	// Re-render automatically when exchange rates refresh
	ctx.Events().On("exchange_rates_updated", func(event contracts.Event) {
		ctx.UI().RequestRender()
	})
}

func (c *Dashboard) Keys(km contracts.KeyManager) {}

func (c *Dashboard) Render() string {
	lang := c.app.Lang()
	balFloat := c.app.State().Get("balance", 0.0).Float64()

	balStr := utils.FormatIDR(balFloat)

	width, _ := c.app.UI().Dimensions()
	innerWidth := width - 4
	innerWidth = max(50, innerWidth)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 4).
		MarginBottom(1).
		Width(innerWidth)

	cardASCII := `
 ╭───────────────────────────────────╮
 │  TakoBank                   [  ]  │
 │                                   │
 │  1234  5678  9012  3456           │
 │                                   │
 │  VALID THRU 12/28                 │
 │  SPONGEBOB SQUAREPANTS            │
 ╰───────────────────────────────────╯`

	cardStyle := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(cardASCII)

	balTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(lang.T("dashboard.balance"))
	balValue := lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Render(balStr)

	balanceBox := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, cardStyle, "", balTitle, balValue))

	// Get Exchange Rates from Cache
	if cache := c.app.Cache(); cache != nil {
		var ratesStr string
		if cache.Get("exchange_rates", &ratesStr) {
			var rates map[string]float64
			_ = json.Unmarshal([]byte(ratesStr), &rates)

			ratesList := ""
			for k, v := range rates {
				ratesList += fmt.Sprintf("%s: %s\n", k, utils.FormatIDR(v))
			}

			rateTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(lang.T("dashboard.rates"))
			ratesBox := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rateTitle, ratesList))

			return lipgloss.JoinVertical(lipgloss.Left, balanceBox, ratesBox)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, balanceBox, lipgloss.NewStyle().Foreground(MutedColor).Render("loading exchange rates..."))
}

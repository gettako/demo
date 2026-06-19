package components

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"gettako.dev/tako/contracts"
)

type UserProfile struct {
	Name    string
	Level   string
	Branch  string
	Address string
}

type Profile struct {
	app     contracts.Application
	profile UserProfile
	loading bool
}

func NewProfile(app contracts.Application) *Profile {
	return &Profile{app: app}
}

func (c *Profile) ID() string { return "profile" }

func (c *Profile) Init(ctx contracts.Context) {
	c.loading = true

	// Simulate slow API call in a goroutine
	ctx.Spawn(func(goCtx context.Context) {
		cache := ctx.Cache()
		if cache == nil {
			c.loading = false
			c.app.UI().RequestRender()
			return
		}

		// Use cache with a 10 second TTL.
		// If cached data is not found, Remember runs the closure, stores it, and returns it.
		val, _ := cache.Remember("user_profile", 10*time.Second, func() (any, error) {
			// Simulate slow API
			time.Sleep(2 * time.Second)
			return UserProfile{
				Name:    "SpongeBob SquarePants",
				Level:   "Platinum Priority",
				Branch:  "Bikini Bottom",
				Address: "124 Conch Street",
			}, nil
		})

		if p, ok := val.(UserProfile); ok {
			c.profile = p
		}
		c.loading = false
		c.app.UI().RequestRender()
	})
}

func (c *Profile) Keys(km contracts.KeyManager) {}

func (c *Profile) Update(msg tea.Msg) (tea.Cmd, bool) {
	return nil, false
}

func (c *Profile) Render() string {
	width, _ := c.app.UI().Dimensions()

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(2, 4).
		Width(width - 4).
		Align(lipgloss.Center)

	title := lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("User Profile")

	if c.loading {
		return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			lipgloss.NewStyle().Foreground(WarningColor).Render("Fetching profile data securely..."),
			lipgloss.NewStyle().Foreground(MutedColor).Render("please wait a moment"),
		))
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"==========================================",
		"",
		fmt.Sprintf("Name    : %s", lipgloss.NewStyle().Foreground(SuccessColor).Render(c.profile.Name)),
		fmt.Sprintf("Level   : %s", lipgloss.NewStyle().Foreground(SuccessColor).Render(c.profile.Level)),
		fmt.Sprintf("Branch  : %s", c.profile.Branch),
		fmt.Sprintf("Address : %s", c.profile.Address),
		"",
		"==========================================",
		"\nCRM Details:",
		lipgloss.NewStyle().Foreground(MutedColor).Render("this data is cached for 10 seconds."),
	)

	return lipgloss.Place(0, 0, lipgloss.Center, lipgloss.Center, boxStyle.Render(content))
}

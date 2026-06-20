package components

import (
	"gettako.dev/tako/contracts"
)

type App struct {
	app contracts.Application

	// Public fields for reactivity
	Title string
	Count int
}

func NewApp(app contracts.Application) *App {
	return &App{
		app:   app,
		Title: "Tako View Engine Demo",
		Count: 0,
	}
}

func (c *App) ID() string { return "app" }

func (c *App) Init(ctx contracts.Context) {}

func (c *App) Keys(km contracts.KeyManager) {
	km.Bind("enter", c.Increment)
	km.Bind("r", c.Reset)
}

func (c *App) Increment() {
	c.Count++
}

func (c *App) Reset() {
	c.Count = 0
}

func (c *App) Render() string {
	// Let the view engine render the HTML based on this struct instance
	return c.app.ViewEngine().Render(c, "components/app.html")
}

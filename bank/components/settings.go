package components

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"gettako.dev/tako/contracts"
)

type Settings struct {
	app       contracts.Application
	node      contracts.ViewNode
	demoInput *DemoInput
}

// DemoInput component
type DemoInput struct {
	input textinput.Model
}

func NewDemoInput() *DemoInput {
	i := textinput.New()
	i.Placeholder = "Enter your settings here..."
	i.SetWidth(60)
	i.Focus()
	return &DemoInput{input: i}
}
func (c *DemoInput) ID() string                   { return "demo-input" }
func (c *DemoInput) Init(ctx contracts.Context)   {}
func (c *DemoInput) Keys(km contracts.KeyManager) {}
func (c *DemoInput) Render() string               { return c.input.View() }

// DemoTable component
type DemoTable struct {
	content string
}

func NewDemoTable(ctx contracts.Context) *DemoTable {
	t := table.New().
		Headers("Setting", "Value").
		Rows(
			[]string{"Theme", "Dark Mode"},
			[]string{"Language", "English"},
			[]string{"Notifications", "Enabled"},
		).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#cb6843"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle().Padding(0, 1)
			if row == 0 {
				return s.Bold(true).Foreground(lipgloss.Color("#cb6843"))
			}
			return s
		})

	return &DemoTable{content: t.Render()}
}
func (c *DemoTable) ID() string                   { return "demo-table" }
func (c *DemoTable) Init(ctx contracts.Context)   {}
func (c *DemoTable) Keys(km contracts.KeyManager) {}
func (c *DemoTable) Render() string               { return c.content }

// DemoCard component with Slot support
type DemoCard struct {
	ctx   contracts.Context
	slots map[string]contracts.Slot
}

func NewDemoCard(ctx contracts.Context) *DemoCard {
	return &DemoCard{ctx: ctx}
}
func (c *DemoCard) ID() string                   { return "demo-card" }
func (c *DemoCard) Init(ctx contracts.Context)   {}
func (c *DemoCard) Keys(km contracts.KeyManager) {}
func (c *DemoCard) SetSlots(slots map[string]contracts.Slot) {
	c.slots = slots
}
func (c *DemoCard) Render() string {
	content := ""
	if slot, ok := c.slots["default"]; ok {
		content = slot.Render(c.ctx)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cb6843")).
		Padding(1, 2).
		MarginBottom(1).
		Render(content)
}

func NewSettings(app contracts.Application) *Settings {
	s := &Settings{app: app}
	s.demoInput = NewDemoInput()

	// Register components dynamically for the view engine
	app.ViewEngine().Register("demo-input", func() contracts.Component {
		return s.demoInput
	})
	app.ViewEngine().Register("demo-table", func() contracts.Component {
		return NewDemoTable(app.Context())
	})
	app.ViewEngine().Register("demo-card", func() contracts.Component {
		return NewDemoCard(app.Context())
	})

	if !app.Context().State().Has("SearchQuery") {
		app.Context().State().Mutate("SearchQuery").Value("").Broadcast()
	}

	// Load the view using the new XML View Engine
	// Assuming working directory is demo/
	node, err := app.ViewEngine().Load("components/settings.html")
	if err != nil {
		app.Context().Logger().Error("Failed to load settings view", "error", err)
	} else {
		s.node = node
	}

	return s
}

func (c *Settings) ID() string { return "settings" }

func (c *Settings) Init(ctx contracts.Context) {
	if c.demoInput != nil {
		c.demoInput.input.Focus()
	}
}

func (c *Settings) Keys(km contracts.KeyManager) {}

// UpdateNative intercepts BubbleTea messages and routes them to our nested input component
func (c *Settings) UpdateNative(msg any) (any, bool) {
	if c.demoInput == nil {
		return nil, false
	}

	if teaMsg, ok := msg.(tea.Msg); ok {
		var cmd tea.Cmd
		c.demoInput.input, cmd = c.demoInput.input.Update(teaMsg)

		// Sync input value to global state so it's reactive
		c.app.Context().State().Mutate("SearchQuery").Value(c.demoInput.input.Value()).Broadcast()

		// Check if it's a keypress to consume
		if keyMsg, isKey := msg.(tea.KeyMsg); isKey {
			switch keyMsg.String() {
			case "tab", "shift+tab", "enter", "up", "down", "left", "right", "f1", "f2", "f3", "f4", "f5", "alt+1", "alt+2", "alt+3", "alt+4", "alt+5":
				// Don't consume navigation keys
				return cmd, false
			default:
				// Consume typing keys
				c.app.UI().RequestRender()
				return cmd, true
			}
		}

		// Request a render unconditionally so blinking works
		c.app.UI().RequestRender()
		return cmd, false
	}

	return nil, false
}

func (c *Settings) Render() string {
	if c.node == nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render("Error: XML view not loaded.")
	}

	// Render the persistent AST node using the context to evaluate bindings
	// To add vertical padding:
	out := c.node.RenderTree(c.app.Context())

	// Center the output in the screen for presentation
	width, _ := c.app.UI().Dimensions()
	return lipgloss.NewStyle().Width(width - 4).Align(lipgloss.Center).PaddingTop(2).Render(out)
}

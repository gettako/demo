package foo

import (
	"github.com/gettako/tako/internal/plugin"
	"github.com/gettako/tako/internal/tako"
)

func init() {
	manifest := plugin.Manifest{
		ID:          "com.example.foo",
		Name:        "foo",
		Version:     "1.0.0",
		Description: "A new Tako plugin",
		Type:        plugin.TypeHeadless,
	}

	plugin.Register(manifest, func() plugin.Lifecycle {
		return &Plugin{}
	})
}

// Plugin is the main lifecycle controller for the foo plugin.
type Plugin struct {
	plugin.NoopLifecycle
	ctx *tako.Context
}

func (p *Plugin) OnInit(ctx *tako.Context) error {
	p.ctx = ctx
	// Setup your plugin here
	return nil
}

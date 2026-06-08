package help

import "github.com/takoterm/tako/internal/plugin"

func init() {
	plugin.Register(Manifest, func() plugin.Lifecycle {
		return &Plugin{}
	})
}

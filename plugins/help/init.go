package help

import "github.com/gettako/tako/internal/plugin"

func init() {
	plugin.Register(Manifest, func() plugin.Lifecycle {
		return &Plugin{}
	})
}

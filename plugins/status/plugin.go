package status

import "github.com/takoterm/tako/internal/plugin"

// Manifest declares the identity and capabilities of the Status plugin.
var Manifest = plugin.Manifest{
	ID:          "status",
	Name:        "System Status Monitor",
	Version:     "1.0.0",
	Description: "Shows events and statuses inside the Tako Terminal.",
	Author:      "Tako Team",
}

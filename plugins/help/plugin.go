package help

import "github.com/takoterm/tako/internal/plugin"

// Manifest declares the identity and capabilities of the Help plugin.
var Manifest = plugin.Manifest{
	ID:          "help",
	Name:        "Help System",
	Version:     "1.0.0",
	Description: "Provides an interactive help overlay.",
	Author:      "Tako Team",
}

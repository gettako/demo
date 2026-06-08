package fzf

import "github.com/takoterm/tako/internal/plugin"

// Manifest declares the identity and capabilities of the FZF plugin.
var Manifest = plugin.Manifest{
	ID:          "fzf",
	Name:        "Fuzzy Finder",
	Version:     "1.0.0",
	Description: "Provides fuzzy file search capabilities.",
	Author:      "Tako Team",
}

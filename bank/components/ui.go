package components

import "charm.land/lipgloss/v2"

var (
	PrimaryColor = lipgloss.Color("#cb6843") // Tako Orange
	MutedColor   = lipgloss.Color("#9ca3af") // Gray-400
	DangerColor  = lipgloss.Color("#dc2626") // Red
	SuccessColor = lipgloss.Color("#16a34a") // Emerald
	WarningColor = lipgloss.Color("#d97706") // Amber
)

// Parse is a simple mock of ui.Parse to strip tags or ignore them.
func Parse(text string) string {
	// In a real app we'd parse <dim> etc., here we just return the text
	// after removing basic tags.
	// We'll just leave it or rely on lipgloss for styling in the demo.
	return text
}

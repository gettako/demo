package main

import (
	"fmt"
	lipgloss "charm.land/lipgloss/v2"
)

func main() {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderTitle(" Label ").
		BorderForeground(lipgloss.Color("#874833")).
		BorderTitleForeground(lipgloss.Color("#e65200")). // brand
		Width(40).
		Height(5).
		Padding(0, 1)

	fmt.Println(style.Render("Hello world\nLine 2"))
}

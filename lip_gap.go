package main

import (
	"fmt"
	"charm.land/lipgloss/v2"
)

func main() {
	box1 := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render("Box 1")
	box2 := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render("Box 2")

	fmt.Println("=== WITH \\n ===")
	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, box1, "\n", box2))
	
	fmt.Println("=== WITHOUT \\n ===")
	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, box1, box2))
}

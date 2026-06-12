package main

import (
	"fmt"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

func main() {
	headers := []string{"ID", "Name"}
	rows := [][]string{
		{"1", "Tako"},
		{"2", "Puro"},
	}

	colWidths := []int{5, 10}

	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).BorderColumn(false).BorderRow(false).
		BorderHeader(true).
		Width(40)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		s := lipgloss.NewStyle().Padding(0, 1)
		if col < len(colWidths)-1 {
			s = s.Width(colWidths[col])
		}
		
		if row == 1 { // tako
			s = s.Background(lipgloss.Color("#123456")).Foreground(lipgloss.Color("#ffffff"))
		}
		return s
	})

	fmt.Println(t.Render())
}

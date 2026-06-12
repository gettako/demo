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
		{"3", "Octo"},
	}

	colWidths := []int{10, 20}

	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).BorderColumn(false).BorderRow(false).
		BorderHeader(true)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		s := lipgloss.NewStyle().Width(colWidths[col]).Padding(0, 1)
		if row == 0 {
			s = s.Bold(true)
		} else if row == 1 { // tako
			s = s.Background(lipgloss.Color("#123456")).Foreground(lipgloss.Color("#ffffff"))
		}
		return s
	})

	fmt.Println(t.Render())
}

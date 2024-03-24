package main

import "github.com/charmbracelet/lipgloss"

var breadcrumbDividerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#333333"))

var breadcrumbTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#888888"))

var breadcrumbContainerStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("#444444")).
	BorderBottom(true).
	BorderStyle(lipgloss.NormalBorder()).
	Width(80)

func renderBreadcrumb(m model) string {
	output := ""
	for i, item := range m.breadcrumb {
		if i == 0 {
			output += "🐦"
		} else {
			output += breadcrumbDividerStyle.Render(" > ") +
				addSpaceToEmoji(item.group.Emoji) +
				breadcrumbTextStyle.Render(item.group.Name)
		}
	}
	return breadcrumbContainerStyle.Render(output)
}

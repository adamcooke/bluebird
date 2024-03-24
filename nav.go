package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var globalStyle = lipgloss.NewStyle().Padding(1, 2)
var dimmedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
var whiteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))

func renderNav(m model) string {
	output := m.textInput.View() +
		whiteStyle.Render("↑↓") + dimmedStyle.Render(" navigate") +
		" • " +
		whiteStyle.Render("enter") + dimmedStyle.Render(" select") +
		" • " +
		whiteStyle.Render("esc") + dimmedStyle.Render(" back") +
		" • " +
		whiteStyle.Render("?") + dimmedStyle.Render(" help")

	if m.showHelp {
		output += "\n"
		output += strings.Repeat(" ", 21)
		output += whiteStyle.Render("ctrl+l") + dimmedStyle.Render(" show commands") +
			" • " +
			whiteStyle.Render("ctrl+c") + dimmedStyle.Render(" quit")

	}
	return output
}

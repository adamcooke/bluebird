package main

import (
	"github.com/charmbracelet/lipgloss"
)

var rowStyleInactive = lipgloss.NewStyle().
	Padding(0, 0, 0, 2).
	Border(lipgloss.Border{Left: " "}).
	BorderLeft(true).BorderBottom(false).BorderTop(false).BorderRight(false)

var rowStyleActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB454")).
	BorderLeft(true).
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("#FFB454")).
	Padding(0, 0, 0, 2)

var descriptionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#777777")).
	Bold(false).
	Width(70)

var welcomeText = lipgloss.NewStyle().
	Padding(0, 0, 0, 2).
	Width(70)

func renderList(m model) string {
	output := ""
	items := m.FilteredItems()

	if len(items) == 0 && len(m.breadcrumb) == 1 && len(m.textInput.Value()) == 0 {
		output := "\n"
		output += addSpaceToEmoji("👋")
		output += lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD869")).Bold(true).Render("Hey. You don't seem to have anything configured yet.")
		output += " To get started you'll need to add some values to your configuration file. You can do this by looking in "
		output += lipgloss.NewStyle().Foreground(lipgloss.Color("#169FFF")).Render(m.config.path)

		return welcomeText.Render(output)
	} else if len(items) == 0 {
		return rowStyleInactive.Render(
			"\n" + addSpaceToEmoji("🫙") + descriptionStyle.Render("Nothing to see here\n"),
		)
	}

	addExtraSpace := false
	for i, item := range items {
		style := rowStyleInactive
		if i == m.cursor {
			style = rowStyleActive
		}

		rowOutput := ""
		rowOutput += item.Emoji()
		rowOutput += item.Name

		if item.Description != "" {
			rowOutput += "\n"
			rowOutput += descriptionStyle.Render(item.Description)
			addExtraSpace = true
		}

		if m.showCommands && item.IsCommand() && item.Command.Command != "" {
			rowOutput += "\n"
			rowOutput += descriptionStyle.Render("$ " + item.Command.Command)
			addExtraSpace = true
		}

		output += style.Render(rowOutput)

		if addExtraSpace {
			output += "\n"
			addExtraSpace = false
		}

		if i < len(items)-1 {
			output += "\n"
		}
	}

	return output
}

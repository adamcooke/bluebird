package main

import "github.com/charmbracelet/lipgloss"

var buttonStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#000000")).
	Background(lipgloss.Color("#ffffff")).
	Padding(0, 3).
	MarginTop(1)

var dialogBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#ff0000")).
	Padding(1, 1).
	BorderTop(true).
	BorderLeft(true).
	BorderRight(true).
	BorderBottom(true)

func renderDialog(text string, buttonText string) string {
	okButton := buttonStyle.Render(buttonText)
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	question := lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render(text)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	return lipgloss.Place(80, 9,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("猫咪"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}

func renderRunConfirmation(_ model) string {
	okButton := buttonStyle.Render("Press ENTER to continue")

	question := lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render("Are you sure you wish to run this?")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
	cancelText := lipgloss.NewStyle().Width(40).
		Align(lipgloss.Center).
		PaddingTop(1).
		Foreground(lipgloss.Color("#777777")).
		Render("Press escape to cancel")
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons, cancelText)

	return lipgloss.Place(80, 13,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("猫咪"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#5C353E")),
	)
}

package styles

import "github.com/charmbracelet/lipgloss"

var (
	GlobalStyle = lipgloss.NewStyle().Padding(1, 2)

	SubtleText = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	WhiteText = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))

	BoldText = lipgloss.NewStyle().Bold(true)

	ErrorText = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	BreadcrumbStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#444444")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			MarginBottom(1)

	ActiveListItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB454")).
				BorderLeft(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#FFB454")).
				Padding(0, 0, 0, 2)

	LoadingListItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555")).
				BorderLeft(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#555555")).
				Padding(0, 0, 0, 2)

	InactiveListItemStyle = lipgloss.NewStyle().
				Padding(0, 0, 0, 2).
				Border(lipgloss.Border{Left: " "}).
				BorderLeft(true).BorderBottom(false).BorderTop(false).BorderRight(false)
)

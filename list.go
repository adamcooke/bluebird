package main

import (
	"fmt"
	"io"

	"github.com/adamcooke/bluebird/config"
	"github.com/adamcooke/bluebird/list"
	"github.com/adamcooke/bluebird/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listItemDelegate struct{}

func (d listItemDelegate) Height() int {
	return 2
}

func (d listItemDelegate) Spacing() int {
	return 1
}

func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(config.Item)
	if !ok {
		return
	}

	showSpinner := false

	var style lipgloss.Style
	if index == m.Index() {
		style = styles.ActiveListItemStyle
		if m.ShowSpinner() {
			showSpinner = true
			style = styles.LoadingListItemStyle
		}
	} else {
		style = styles.InactiveListItemStyle
	}

	output := ""
	output += addSpaceToEmoji(item.EmojiWithDefault()) + " "
	output += styles.BoldText.Render(item.Name)
	if showSpinner {
		output += " "
		output += m.Spinner().View()
	}
	output += "\n"

	description := item.Description
	if item.Description == "" {
		description = item.Backend().Description()
	}

	output += styles.SubtleText.Render(description)

	fmt.Fprint(w, style.Render(output))
}

func initialList() list.Model {
	// This starts very small but it will be updated as soon as know
	// the size of the terminal.
	list := list.New([]list.Item{}, listItemDelegate{}, 1, 1)

	list.SetSpinner(spinner.Points)
	list.Styles.Spinner.Foreground(lipgloss.Color("#555555"))

	// Disable all features of the list. The only thing we want here
	// is the pagination support.
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowFilter(false)
	list.SetShowPagination(false)

	return list
}

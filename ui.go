package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type commandFinishedMessage struct{ err error }

type model struct {
	config          config
	breadcrumb      []breadcrumbItem
	cursor          int
	textInput       textinput.Model
	errorMessage    string
	quitting        bool
	showCommands    bool
	showHelp        bool
	proposedCommand *command
}

type breadcrumbItem struct {
	group  *group
	cursor int
	filter string
}

func (m model) ActiveGroup() *group {
	return m.breadcrumb[len(m.breadcrumb)-1].group
}

func (m model) FilteredItems() []item {
	items := m.ActiveGroup().Items()

	searchValue := m.textInput.Value()
	if searchValue == "" {
		return items
	}

	filteredItems := []item{}
	for _, g := range items {
		if strings.Contains(strings.ToLower(g.Name), strings.ToLower(searchValue)) {
			filteredItems = append(filteredItems, g)
		}
	}

	return filteredItems
}

func initialModel(config config) model {
	ti := textinput.New()
	ti.Placeholder = "type to filter"
	ti.Focus()
	ti.CharLimit = 20
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	ti.Prompt = ""
	ti.Cursor.SetMode(cursor.CursorHide)
	ti.Width = 20

	return model{
		config:     config,
		textInput:  ti,
		breadcrumb: []breadcrumbItem{{group: &config.RootGroup}},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commandFinishedMessage:
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+l":
			m.showCommands = !m.showCommands
			return m, nil
		case "?":
			m.showHelp = !m.showHelp
			return m, nil

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down":
			if len(m.FilteredItems())-1 > m.cursor {
				m.cursor++
			}
			return m, nil
		case "esc", "left":
			if m.errorMessage != "" {
				m.errorMessage = ""
				return m, nil
			}

			if m.proposedCommand != nil {
				m.proposedCommand = nil
				return m, nil
			}

			// Clear any filtering at this point if there is any before actually
			// changing levels.
			if m.textInput.Value() != "" {
				m.textInput.SetValue("")
				return m, nil
			}

			// If we have a breadcrumb stack, we'll start popping things off
			// that stack now until it's only got a single item left.
			if len(m.breadcrumb) > 1 {
				previousItem := m.breadcrumb[len(m.breadcrumb)-1]
				m.breadcrumb = m.breadcrumb[:len(m.breadcrumb)-1]
				m.cursor = previousItem.cursor
				m.textInput.SetValue(previousItem.filter)
			}

			return m, nil

		case "enter", "right":
			if m.errorMessage != "" {
				m.errorMessage = ""
				return m, nil
			}

			if m.proposedCommand != nil {
				return runCommand(m, m.proposedCommand)
			}

			// If there are no items on display, just return nil because pressing
			// enter at this stage should do nothing.
			items := m.FilteredItems()
			if len(items) == 0 {
				return m, nil
			}

			// Grab the current cursor and the item which has actually been
			// selected from the table.
			currentCursor := m.cursor
			selectedItem := items[currentCursor]

			// If this is a command, we're going to display a warning if
			// appropriate and then get on with rendering it.
			if selectedItem.IsCommand() {
				if selectedItem.Command.Command == "" {
					m.errorMessage = fmt.Sprintf("No command available for %s", selectedItem.Command.Name)
					return m, nil
				}

				if selectedItem.Command.Prompt {
					// If the command has a prompt, we should render that prompt.
					m.proposedCommand = selectedItem.Command
					return m, nil
				} else {
					// If the command does not have a prompt enabled, we can just launch it immediately.
					return runCommand(m, selectedItem.Command)
				}
			}

			// Otherwise, we're changing group here... We need to add the
			// currently visible group to the breadcrumb stack here and
			// update the table.
			m.breadcrumb = append(m.breadcrumb, breadcrumbItem{
				group:  selectedItem.Group,
				cursor: currentCursor,
				filter: m.textInput.Value(),
			})

			m.textInput.SetValue("")
			m.cursor = 0
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.cursor = 0
	return m, cmd
}

func runCommand(m model, command *command) (tea.Model, tea.Cmd) {
	m.quitting = true
	c := exec.Command("/bin/sh", "-c", command.Command) //nolint:gosec
	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		return commandFinishedMessage{err}
	})
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	output := ""
	output += renderBreadcrumb(m) + "\n"
	if m.errorMessage != "" {
		output += renderDialog(m.errorMessage, "Back")
	} else if m.proposedCommand != nil {
		output += renderRunConfirmation(m)
	} else {
		output += renderList(m)
		output += "\n\n"
		output += renderNav(m)
	}

	return globalStyle.Render(output) + "\n"
}

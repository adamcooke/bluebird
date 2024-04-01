package main

import (
	"os/exec"
	"strings"

	"github.com/adamcooke/bluebird/config"
	"github.com/adamcooke/bluebird/list"
	"github.com/adamcooke/bluebird/styles"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type listItemsAvailableNowMsg config.Item
type listItemsLoadingErrorMsg error
type commandFinishedMessage struct{ err error }

type model struct {
	config         config.Config
	list           list.Model
	breadcrumb     []breadcrumbItem
	textInput      textinput.Model
	quitting       bool
	terminalWidth  int
	terminalHeight int
	disableInput   bool
	errorMessage   string
}

type breadcrumbItem struct {
	item   *config.Item
	cursor int
	page   int
	filter string
}

func (m model) ActiveItems() []config.Item {
	if len(m.breadcrumb) == 0 {
		return m.config.Items
	} else {
		return m.breadcrumb[len(m.breadcrumb)-1].item.Backend().Items()
	}
}

func (m model) FilterItems(items []config.Item) []config.Item {
	searchValue := m.textInput.Value()
	if searchValue == "" {
		return items
	}

	filteredItems := []config.Item{}
	for _, g := range items {
		if strings.Contains(strings.ToLower(g.Name), strings.ToLower(searchValue)) {
			filteredItems = append(filteredItems, g)
		}
	}

	return filteredItems
}

func (m model) Init() tea.Cmd {
	return nil
}

func itemsForList(items []config.Item) []list.Item {
	l := make([]list.Item, len(items))
	for i, item := range items {
		l[i] = item
	}
	return l
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commandFinishedMessage:
		m.quitting = true
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// Whenever we get a new window size message, we update the
		// model so it knows what size terminal it has to work with.
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height

		// We'll need to update the size of the list at this point to
		// make it appropriatelysized.
		m.list.SetWidth(m.terminalWidth)
		m.list.SetHeight(m.terminalHeight - 7)

		// We should now load the initial items into the list, we
		// don't want to do this prematurely so we'll do on window size
		// message.
		if len(m.breadcrumb) == 0 {
			m.list.SetItems(itemsForList(m.FilterItems(m.config.Items)))
		}

		return m, nil
	case spinner.TickMsg:
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case listItemsAvailableNowMsg:
		// We have list items available to display, they should now be
		// shown
		item := config.Item(msg)

		// Add the new item to the breadcrumb
		m.breadcrumb = append(m.breadcrumb, breadcrumbItem{
			item:   &item,
			cursor: m.list.Cursor(),
			page:   m.list.Paginator.Page,
			filter: m.textInput.Value(),
		})

		// Reset filtering
		m.textInput.SetValue("")

		// Update items
		m.list.SetItems(itemsForList(m.ActiveItems()))

		// Set the cursor back to zero
		m.list.Paginator.Page = 0
		m.list.SetCursor(0)

		// Stop the spinner when loading a new list item
		m.list.StopSpinner()
		m.disableInput = false

		return m, nil

	case listItemsLoadingErrorMsg:
		m.errorMessage = msg.Error()
		m.disableInput = false
		m.list.StopSpinner()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down":
			if m.disableInput || m.errorMessage != "" {
				return m, nil
			}

			m.list, cmd = m.list.Update(msg)
			return m, cmd
		case "ctrl+c":
			// When ctrl+c is pressed this is our signal to get out of here
			// and we'll just send the quit command and be on our way.
			m.quitting = true
			return m, tea.Quit
		case "enter", "right":
			if m.disableInput {
				return m, nil
			}

			if m.errorMessage != "" {
				m.errorMessage = ""
				return m, nil
			}

			// When pressing enter, we'll need to load this up
			// and display the next set of items.
			item, ok := m.list.SelectedItem().(config.Item)
			if !ok {
				return m, nil
			}

			if item.Backend().IsList() {
				// If we have items available, we'll skip the spinner
				if !item.Backend().Loaded() {
					m.disableInput = true
					return m, tea.Batch(
						m.list.StartSpinner(),
						loadList(item),
					)
				} else {
					return m, loadList(item)
				}

			} else {
				if item.Backend().Command() == "" {
					m.errorMessage = "No command available for this item."
					return m, nil
				}

				m.quitting = true
				c := exec.Command("/bin/bash", "-c", item.Backend().Command())
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return commandFinishedMessage{err}
				})
			}

		case "esc", "left":
			if m.disableInput {
				return m, nil
			}

			if m.errorMessage != "" {
				m.errorMessage = ""
				return m, nil
			}

			// If the text input is nil, clear it when pressing escape.
			if m.textInput.Value() != "" {
				m.textInput.SetValue("")
				m.list.SetItems(itemsForList(m.FilterItems(m.ActiveItems())))
				m.list.SetCursor(0)
				m.list.Paginator.Page = 0
				return m, cmd
			}

			breadcrumbSize := len(m.breadcrumb)
			if breadcrumbSize == 0 {
				// If we're at the root, we don't want to do anything here.
				return m, nil
			}

			currentItem := m.breadcrumb[len(m.breadcrumb)-1]

			var newItems []config.Item
			if breadcrumbSize > 1 {
				previousItem := m.breadcrumb[len(m.breadcrumb)-2]
				m.breadcrumb = m.breadcrumb[:len(m.breadcrumb)-1]
				newItems = previousItem.item.Backend().Items()
			} else {
				// If we're at the end of the line, we'll load in the
				// root level items again now and clear the breadcrumb
				// entirely.
				m.breadcrumb = []breadcrumbItem{}
				newItems = m.config.Items
			}
			m.textInput.SetValue(currentItem.filter)
			m.list.SetItems(itemsForList(m.FilterItems(newItems)))
			m.list.Paginator.Page = currentItem.page
			m.list.SetCursor(currentItem.cursor)

			return m, nil
		}
	}

	if m.disableInput {
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.list.SetItems(itemsForList(m.FilterItems(m.ActiveItems())))
	m.list.SetCursor(0)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "\n"
	}

	breadcrumb := "🐦"

	for _, item := range m.breadcrumb {
		breadcrumb += styles.SubtleText.Render(" > ") +
			addSpaceToEmoji(item.item.Emoji) +
			item.item.Name
	}

	output := styles.BreadcrumbStyle.Copy().Width(m.terminalWidth - 4).Render(breadcrumb)
	output += "\n"

	if m.errorMessage != "" {
		output += styles.ErrorText.Render("Oh no. An error occurred while doing that.")
		output += "\n\n>> "
		output += m.errorMessage
		output += "\n\n"
		output += styles.SubtleText.Render("Press enter or esc to go back.")
		return styles.GlobalStyle.Render(output) + "\n"
	}

	output += m.list.View()
	output += "\n"

	output += m.textInput.View() +
		styles.WhiteText.Render("↑↓") + styles.SubtleText.Render(" navigate") +
		" • " +
		styles.WhiteText.Render("enter") + styles.SubtleText.Render(" select") +
		" • " +
		styles.WhiteText.Render("esc") + styles.SubtleText.Render(" back") +
		" • " +
		styles.WhiteText.Render("?") + styles.SubtleText.Render(" help")

	return styles.GlobalStyle.Render(output) + "\n"
}

func initialModel(config config.Config) model {
	// The list is one of the core components of this so we're going to
	// adding that to our initial model.
	list := initialList()

	// Text input
	ti := textinput.New()
	ti.Placeholder = "type to filter"
	ti.Focus()
	ti.CharLimit = 20
	ti.PlaceholderStyle = styles.SubtleText
	ti.Prompt = ""
	ti.Cursor.SetMode(cursor.CursorHide)
	ti.Width = 20

	// Create a model and return it
	return model{
		config:     config,
		textInput:  ti,
		list:       list,
		breadcrumb: []breadcrumbItem{},
	}
}

// A tea command which will begin the process of loading a list.
// When complete, it will
func loadList(item config.Item) tea.Cmd {
	return func() tea.Msg {
		// If the items are already loaded, we don't need to do anything
		// here and we can just say that items are available immediately.
		if item.Backend().Loaded() {
			return listItemsAvailableNowMsg(item)
		}

		// Load the items at this point.
		err := item.Backend().Load()
		if err == nil {
			return listItemsAvailableNowMsg(item)
		} else {
			return listItemsLoadingErrorMsg(err)
		}
	}
}

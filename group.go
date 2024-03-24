package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/mattn/go-runewidth"
)

type group struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Emoji       string    `yaml:"emoji"`
	Groups      []group   `yaml:"groups"`
	Commands    []command `yaml:"commands"`
}

type item struct {
	Name        string
	Description string
	emoji       string
	Group       *group
	Command     *command
}

// Is this item a group?
func (i item) IsGroup() bool {
	return i.Group != nil
}

// Is this item a command?
func (i item) IsCommand() bool {
	return i.Command != nil
}

func (i item) Emoji() string {
	emoji := ""
	if len(i.emoji) > 0 {
		emoji = i.emoji
	}

	if emoji == "" {
		if i.IsCommand() {
			emoji = "🪓"
		} else {
			emoji = "🗂️"
		}
	}

	emojiWidth := runewidth.StringWidth(emoji)
	spaces := strings.Repeat(" ", 4-emojiWidth) // adjust the number of spaces based on the width of the emoji
	return fmt.Sprintf("%s%s", emoji, spaces)
}

// Returns all items which should be displayed for this group.
// In order to actually display the table the items must all
// respond with an interface that is suitable.
func (g group) Items() []item {
	slice := []item{}

	for _, group := range g.Groups {
		slice = append(slice, item{
			Name:        group.Name,
			Description: group.Description,
			emoji:       group.Emoji,
			Group:       &group,
		})
	}

	for _, command := range g.Commands {
		slice = append(slice, item{
			Name:        command.Name,
			Description: command.Description,
			emoji:       command.Emoji,
			Command:     &command,
		})
	}

	return slice
}

func (g group) TableRows() []table.Row {
	rows := []table.Row{}

	for _, item := range g.Items() {
		rows = append(rows, table.Row{
			item.Emoji(),
			item.Name,
			item.Description,
		})
	}

	return rows
}

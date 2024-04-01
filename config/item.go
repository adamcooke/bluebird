package config

import "fmt"

type ItemType int

const (
	UnknownItemType ItemType = iota
	ListItemType
	CommandItemType
)

// Define mapping of string to item type constant
var itemTypeMap = map[string]ItemType{
	"unknown": UnknownItemType,
	"list":    ListItemType,
	"command": CommandItemType,
}

type Item struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Emoji       string   `yaml:"emoji"`
	Type        ItemType `yaml:"type"`

	ListOptions    ListOptions    `yaml:"listOptions"`
	CommandOptions CommandOptions `yaml:"commandOptions"`

	backend Backend
}

type ListOptions struct {
	Items []Item `yaml:"items"`
}

type CommandOptions struct {
	Command string `yaml:"command"`
}

// Return the backend for this item
func (i Item) Backend() Backend {
	if i.backend != nil {
		return i.backend
	}

	switch i.Type {
	case ListItemType:
		i.backend = ListBackend{item: i}
	case CommandItemType:
		i.backend = CommandBackend{item: i}
	default:
		i.backend = EmptyBackend{}
	}

	return i.backend
}

func (i Item) FilterValue() string {
	return i.Name
}

func (i Item) EmojiWithDefault() string {
	if len(i.Emoji) != 0 {
		return i.Emoji
	}

	return "🐦"
}

// Handle mapping the type value to the correct type constant
func (it *ItemType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	var ok bool
	*it, ok = itemTypeMap[str]
	if !ok {
		return fmt.Errorf("unknown item type: %s", str)
	}

	return nil
}

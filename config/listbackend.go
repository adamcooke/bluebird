package config

import (
	"fmt"
)

type ListBackend struct {
	EmptyBackend
	item Item
}

func (b ListBackend) Items() []Item {
	return b.item.ListOptions.Items
}

func (b ListBackend) Description() string {
	count := len(b.item.ListOptions.Items)
	word := "items"
	if count == 1 {
		word = "item"
	}
	return fmt.Sprintf("%d %s", count, word)
}

func (b ListBackend) IsList() bool {
	return true
}

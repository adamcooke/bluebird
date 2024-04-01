package config

type CommandBackend struct {
	EmptyBackend
	item Item
}

func (b CommandBackend) Description() string {
	return b.item.CommandOptions.Command
}

func (b CommandBackend) Command() string {
	return b.item.CommandOptions.Command
}

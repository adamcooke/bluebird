package config

// A backend is a generic
type Backend interface {
	Items() []Item

	Load() error

	// Have the items been loaded?
	Loaded() bool

	// Return a description to display for this item
	Description() string

	// Is this item a list. If a list this will return
	// true otherwise it will be false.
	IsList() bool

	// Return the command to execute if this is a command backend
	Command() string
}

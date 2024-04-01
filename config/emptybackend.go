package config

type EmptyBackend struct{}

// Load all the items into this backend
func (b EmptyBackend) Load() error {
	return nil
}

// Returns all the items in this backend
func (b EmptyBackend) Items() []Item {
	return []Item{}
}

// Have items been loaded?
func (b EmptyBackend) Loaded() bool {
	return true
}

func (b EmptyBackend) Description() string {
	return "Invalid type"
}

func (b EmptyBackend) IsList() bool {
	return false
}

func (b EmptyBackend) Command() string {
	return ""
}

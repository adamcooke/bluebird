package config

type Config struct {
	Path  string
	Items []Item `yaml:"items"`
}

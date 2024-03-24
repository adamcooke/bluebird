package main

type command struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Emoji       string `yaml:"emoji"`
	Prompt      bool   `yaml:"prompt"`
	Command     string `yaml:"command"`
}

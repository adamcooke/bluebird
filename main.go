package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	version = "dev"
	commit  = "none"
)

var rootCmd = &cobra.Command{
	Use:   "bb",
	Short: "Bluebird is a alias tool.",
	Long:  `Bluebird is a TUI for accessing various aliased commands. Ideal for accessing other services over SSH or through kubectl.`,
	Run: func(cmd *cobra.Command, args []string) {

		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("bluebird %s (commit: %s)\n", version, commit)
			os.Exit(0)
		}

		configFilePath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		c := config{path: configFilePath}

		data, err := os.ReadFile(configFilePath)
		if err == nil {
			err = yaml.Unmarshal([]byte(data), &c)
			if err != nil {
				fmt.Printf("Could not parse config file at %s\n", configFilePath)
				fmt.Printf("%s", err)
				os.Exit(1)
			}
		}

		p := tea.NewProgram(initialModel(c))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func main() {
	homeDir := os.Getenv("HOME")
	rootCmd.PersistentFlags().StringP("config", "c", fmt.Sprintf("%s/.bluebird.yaml", homeDir), "Path to the config file")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

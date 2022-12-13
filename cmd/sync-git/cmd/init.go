package cmd

import (
	"sync-git/internal/embeds"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Write a config file to the current directory",
	Long: `Write a config file to the current directory
To get started with a sync project.`,
	Run: func(cmd *cobra.Command, args []string) {
		embeds.WriteConfigFile()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

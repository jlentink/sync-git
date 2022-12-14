/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync-git")
		fmt.Printf(" - Version: %s\n", version)
		if verbose {
			fmt.Printf(" - Date %s\n", date)
			fmt.Printf(" - Commit %s\n", commit)
			fmt.Printf(" - Built by %s\n", builtBy)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

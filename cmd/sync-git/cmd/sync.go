package cmd

import (
	"github.com/spf13/cobra"
	"github.com/subchen/go-log"
	"sync-git/internal/cnf"
	"sync-git/internal/gitActions"
	"time"
)

var polling bool
var verbose bool
var alwaysPush bool

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the git repositories",
	Long:  `Sync source git repository to the destination repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.Default.Level = log.DEBUG
		}
		gitActions.Clone()
		gitActions.AddRemotes()
		for {
			gitActions.Fetch()
			for _, b := range gitActions.ListRemoteBranches("origin") {
				if gitActions.CompareRemoteChanged(b) || alwaysPush {
					gitActions.Checkout(b)
					gitActions.Pull(b)
					for _, d := range cnf.GetDestinations() {
						gitActions.Push(d, b)
					}
				}
			}
			for _, d := range cnf.GetDestinations() {
				gitActions.DeleteArtifactBranches(d)
			}

			if !polling {
				break
			}
			sleep := time.Duration(cnf.GetInt("git.sleep"))
			time.Sleep(sleep * time.Second)
		}

	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	syncCmd.Flags().BoolVarP(&polling, "poll", "a", false, "Keep polling")
	syncCmd.Flags().BoolVarP(&alwaysPush, "always-push", "p", false, "Always push the branch to remote")
	syncCmd.Flags().StringVarP(&cnf.ConfigLocation, "config", "c", "", "config file (default is /etc/sync-git.toml or ./sync-git.toml)")
}

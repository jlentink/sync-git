package gitActions

import (
	"github.com/go-git/go-git/v5"
	"github.com/subchen/go-log"
	"os"
	"strings"
	"sync-git/internal/cnf"
)

func getBranches(remote bool) []string {
	options := []string{"branch", "--format='%(refname:short)'"}
	if remote {
		options = append(options, "-r")
	}
	output, err := execute(options...)
	branches := make([]string, 0)
	if err != nil {
		log.Errorf("Error listing branches %s (%t)", err, remote)
		os.Exit(1)
	}
	for _, branch := range output {
		branch = strings.TrimSpace(branch)
		branch = strings.Replace(branch, "'", "", -1)
		if len(branch) > 0 {
			branches = append(branches, branch)
		}
	}
	return branches

}

func ListRemoteBranches(remote string) []string {
	branches := getBranches(true)
	if len(remote) == 0 {
		return branches
	}
	filtered := make([]string, 0)
	for _, branch := range branches {
		if strings.HasPrefix(branch, remote) {
			filtered = append(filtered, branch)
		}
	}
	return filtered

}
func ListLocalBranches() []string {
	return getBranches(false)
}

func BranchExists(branch string) bool {
	branches := ListLocalBranches()
	for _, b := range branches {
		if b == branch {
			return true
		}
	}
	return false
}

func Checkout(branch string) {
	log.Debugf("Checking out %s", branch)
	if BranchExists(branch) {
		execute("checkout", "--track", branch)
	} else {
		localBranch := strings.Replace(branch, "origin/", "", -1)
		execute("checkout", localBranch)
	}

}

func Pull(branch string) {
	localBranch := strings.Replace(branch, "origin/", "", -1)
	log.Debugf("Pulling %s", localBranch)
	execute("pull", "origin", localBranch)
}

func Push(destination *cnf.GitDestination, branch string) {
	repo := Open()
	err := repo.Push(&git.PushOptions{
		RemoteURL: destination.CloneUrl(),
		Force:     true,
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Debugf("Branch %s is up to date", branch)
		} else {
			log.Errorf("Error pushing %s", err)
			os.Exit(1)
		}
	}
}

func CompareRemoteChanged(branch string) bool {
	localBranch := strings.Replace(branch, "origin/", "", -1)
	if !BranchExists(localBranch) {
		return true
	}
	remoteHash := BranchHash(branch)
	localHash := BranchHash(localBranch)
	if remoteHash == localHash {
		log.Infof("Branch %s (%s) is up to date", branch, remoteHash)
		return false
	}
	return true
}

func BranchHash(branch string) string {
	log.Debugf("Getting hash for " + branch)
	output, err := execute("rev-parse", branch)
	if err != nil {
		log.Errorf("Error getting branch hash %s", err)
		os.Exit(1)
	}
	return output[0]
}
